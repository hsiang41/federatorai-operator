package influxdb

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	federatoraiv1alpha1 "github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"

	"github.com/containers-ai/federatorai-operator/pkg/influxdb"
	"github.com/containers-ai/federatorai-operator/pkg/util"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	k8sscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

var (
	UpgradeInfluxDBSchemaCMD = &cobra.Command{
		Use:   "influxdb",
		Short: "upgrade influxdb schema from Kilo to Lima",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := parseFlag()
			if err != nil {
				return errors.Wrap(err, "parse flag failed")
			}

			err = initK8SClient()
			if err != nil {
				return errors.Wrap(err, "init k8s client failed")
			}
			return upgradeFromKiloToLima()
		},
	}

	k8sClient           client.Client
	kiloDatabasesToDrop = []string{"alameda_cluster_status", "alameda_planning", "alameda_prediction"}
	// limaNewMeasurements is a map from database name to list of measurements that are created since Lima
	limaNewMeasurements = map[string][]string{
		"alameda_cluster_status": []string{
			"pod",
		},
	}

	timeout        = "300s"
	upgradeTimeout = 300 * time.Second

	logger = log.Log.WithName("upgrader-influxdb")
)

func init() {
	UpgradeInfluxDBSchemaCMD.Flags().StringVar(&timeout, "timeout", "300s", "Timeout limit to execute update process.")
}

func parseFlag() error {
	var err error

	upgradeTimeout, err = time.ParseDuration(timeout)
	if err != nil {
		return errors.Wrap(err, `parse flag "timeout" failed`)
	}

	return nil
}

func initK8SClient() error {

	cfg, err := config.GetConfig()
	if err != nil {
		return errors.Wrap(err, "get kubernetes client config failed")
	}

	scheme := k8sscheme.Scheme
	federatoraiv1alpha1.SchemeBuilder.AddToScheme(scheme)
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		return errors.Wrap(err, "new kubernetes in-cluster client failed")
	}

	return nil
}

func upgradeFromKiloToLima() error {

	logger.Info("Start checking upgradation requirement of influxdb from Kilo to Lima.")

	ctx, cancel := context.WithTimeout(context.Background(), upgradeTimeout)
	defer cancel()
	alamedaServices, err := listAlamedaServices(ctx)
	if err != nil {
		return errors.Wrap(err, "list alamedaservice failed")
	}

	for _, alamedaService := range alamedaServices {
		logger = logger.WithValues("AlamedaService.Namespace", alamedaService.GetNamespace(), "AlamedaService.Name", alamedaService.GetName())

		resourceFactory, err := newResourceFactoryByAlamedaService(alamedaService)
		if err != nil {
			return errors.Wrap(err, "get resources template by alamedaservice failed")
		}

		service := resourceFactory.getAlamedaInfluxdbService()
		exist, err := isAlamedaInfluxdbServiceExist(ctx, service)
		if err != nil {
			return errors.Wrap(err, "check if AlamedaInfluxdb service exists failed")
		} else if !exist {
			logger.Info("AlamedaInfluxdb Service not exist, skip upgradtion.")
			continue
		}

		influxdbCfg := resourceFactory.getAlamedaInfluxdbConfig()
		influxdbClient, err := influxdb.NewClient(influxdbCfg)
		if err != nil {
			return errors.Wrap(err, "new influxdb client failed")
		}

		logger.Info("Check if Lima measurements exist in databases.")
		exist, err = isLimaMeasurementsExist(influxdbClient)
		if err != nil {
			return errors.Wrap(err, "check if Lima measurements exist failed")
		} else if exist {
			logger.Info("Lima measurements exist in databases, skip upgradtion.")
			continue
		}
		logger.Info("Lima measurements not exist in databases, start influxdb upgradtion.")

		controllers := resourceFactory.listWorkloadControllersWithoutAlamedaInfluxdb()
		logger.Info("Scale down and wait replicase of workload controllers to 0.", "controllers", controllers)
		if err := scaleAndWaitWorkloadControllers(ctx, controllers, 0); err != nil {
			return errors.Wrap(err, "scale workload controller failed")
		}

		logger.Info("Drop Kilo databases from influxdb.", "databases", kiloDatabasesToDrop)
		if err := dropKiloDatabases(influxdbClient); err != nil {
			return errors.Wrap(err, "drop Kilo influxdb databases failed")
		}
	}

	logger.Info("Upgradion done.")
	return nil
}

func listAlamedaServices(ctx context.Context) ([]federatoraiv1alpha1.AlamedaService, error) {
	alamedaserviceList := federatoraiv1alpha1.AlamedaServiceList{}
	if err := k8sClient.List(ctx, &alamedaserviceList); err != nil {
		return nil, errors.Wrap(err, "list alamedaservices failed")
	}
	return alamedaserviceList.Items, nil
}

func isAlamedaInfluxdbServiceExist(ctx context.Context, service corev1.Service) (bool, error) {
	instance := corev1.Service{}
	err := k8sClient.Get(ctx, client.ObjectKey{Namespace: service.GetNamespace(), Name: service.GetName()}, &instance)
	if err != nil && !k8serrors.IsNotFound(err) {
		return false, errors.Wrap(err, "get service failed")
	}

	exist := true
	if k8serrors.IsNotFound(err) {
		exist = false
	}
	return exist, nil
}

func getInfluxdbConfigFromDeploymentAndService(deployment appsv1.Deployment, service corev1.Service) influxdb.Config {
	i := influxdb.Config{}

	enableHTTPS := false
	for _, container := range deployment.Spec.Template.Spec.Containers {
		if container.Name != util.InfluxdbCTN {
			continue
		}
		for _, env := range container.Env {
			envName := env.Name
			switch envName {
			case util.AlamedaInfluxDBAdminUserEnvName:
				i.Username = env.Value
			case util.AlamedaInfluxDBAdminPasswordEnvName:
				i.Password = env.Value
			case util.AlamedaInfluxDBHTTPSEnabledEnvName:
				if env.Value == "true" {
					enableHTTPS = true
				}
			}
		}
	}

	scheme := "http"
	switch enableHTTPS {
	case true:
		scheme = "https"
	case false:
		scheme = "http"
	}
	dns := util.GetServiceDNS(&service)
	i.Address = fmt.Sprintf("%s://%s:%d", scheme, dns, util.AlamedaInfluxDBAPIPort)

	return i
}

func isLimaMeasurementsExist(influxdbClient influxdb.Client) (bool, error) {

	for db, newMeasurements := range limaNewMeasurements {
		measurements, err := influxdbClient.ListMeasurementsInDatabase(db)
		if err != nil {
			return false, errors.Wrap(err, "list measurements in database failed")
		}
		for _, newMeasurement := range newMeasurements {
			for _, measurement := range measurements {
				if newMeasurement == measurement {
					return true, nil
				}
			}
		}
	}

	return false, nil
}

func dropKiloDatabases(influxdbClient influxdb.Client) error {

	if err := influxdbClient.DropDatabases(kiloDatabasesToDrop); err != nil {
		return errors.Wrap(err, "drop kilo databases failed")
	}
	return nil
}

type workloadControllerKind = string

var (
	workloadControllerDeployment  = "deployment"
	workloadControllerStatefulSet = "statefulSet"
)

type workloadController struct {
	Kind      workloadControllerKind
	Namespace string
	Name      string
}

// scaleAndWaitWorkloadControllers scale out the controllers to desired replicas and wait until all replicas ready
func scaleAndWaitWorkloadControllers(ctx context.Context, controllers []workloadController, replicas int32) error {

	wg := errgroup.Group{}

	for _, controller := range controllers {
		copyController := controller
		wg.Go(func() error {

			var instacne runtime.Object
			switch copyController.Kind {
			case workloadControllerDeployment:
				instacne = &appsv1.Deployment{}
			case workloadControllerStatefulSet:
				instacne = &appsv1.StatefulSet{}
			}

			err := k8sClient.Get(ctx, client.ObjectKey{Namespace: copyController.Namespace, Name: copyController.Name}, instacne)
			if err != nil && !k8serrors.IsNotFound(err) {
				return errors.Wrapf(err, "get %s failed", copyController.Kind)
			} else if k8serrors.IsNotFound(err) {
				logger.Info("Workload controller not found, ignore.", "Controller", copyController)
				return nil
			}

			switch copyController.Kind {
			case workloadControllerDeployment:
				deployment, ok := instacne.(*appsv1.Deployment)
				if !ok {
					return errors.Wrapf(err, "convert runtime object to %s failed", copyController.Kind)
				}
				deployment.Spec.Replicas = &replicas
			case workloadControllerStatefulSet:
				statefulSet, ok := instacne.(*appsv1.StatefulSet)
				if !ok {
					return errors.Wrapf(err, "convert runtime object to %s failed", copyController.Kind)
				}
				statefulSet.Spec.Replicas = &replicas
			}

			err = k8sClient.Update(ctx, instacne)
			if err != nil {
				return errors.Wrapf(err, "update %s failed", copyController.Kind)
			}

			pollInterval := 1 * time.Second
			return wait.PollUntil(pollInterval, func() (bool, error) {

				err := k8sClient.Get(ctx, client.ObjectKey{Namespace: copyController.Namespace, Name: copyController.Name}, instacne)
				if err != nil {
					return false, errors.Wrapf(err, "get %s failed", copyController.Kind)
				}

				done := false
				switch copyController.Kind {
				case workloadControllerDeployment:
					deployment, ok := instacne.(*appsv1.Deployment)
					if !ok {
						return false, errors.Wrapf(err, "convert runtime object to %s failed", copyController.Kind)
					}
					done = deployment.Status.ReadyReplicas == replicas
				case workloadControllerStatefulSet:
					statefulSet, ok := instacne.(*appsv1.StatefulSet)
					if !ok {
						return false, errors.Wrapf(err, "convert runtime object to %s failed", copyController.Kind)
					}
					done = statefulSet.Status.ReadyReplicas == replicas
				}
				return done, nil
			}, ctx.Done())

		})
	}

	return wg.Wait()
}
