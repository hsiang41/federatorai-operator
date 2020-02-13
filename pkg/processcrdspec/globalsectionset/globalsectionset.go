package globalsectionset

import (
	"strings"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/updateenvvar"
	"github.com/containers-ai/federatorai-operator/pkg/util"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func GlobalSectionSetParamterToStatefulset(ss *appsv1.StatefulSet, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	switch ss.Name {
	case util.FedemeterInflixDBSSN:
	}
}

func GlobalSectionSetParamterToDeployment(dep *appsv1.Deployment, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	switch dep.Name {
	case util.AlamedaaiDPN:
		{
			//Global section set persistentVolumeClaim to mountPath
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-ai-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaaiCTN, "alameda-ai-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaoperatorDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-operator-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaoperatorCTN, "alameda-operator-type-storage", util.AlamedaGroup)
		}
	case util.AlamedadatahubDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-datahub-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedadatahubCTN, "alameda-datahub-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaevictionerDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-evictioner-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaevictionerCTN, "alameda-evictioner-type-storage", util.AlamedaGroup)
		}
	case util.AdmissioncontrollerDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "admission-controller-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AdmissioncontrollerCTN, "admission-controller-type-storage", util.AlamedaGroup)
		}
	case util.AlamedarecommenderDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-recommender-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedarecommenderCTN, "alameda-recommender-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaexecutorDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-executor-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaexecutorCTN, "alameda-executor-type-storage", util.AlamedaGroup)
		}
	case util.AlamedadispatcherDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-dispatcher-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedadispatcherCTN, "alameda-dispatcher-type-storage", util.AlamedaGroup)
		}
	case util.AlamedaRabbitMQDPN:
	case util.AlamedaanalyzerDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-analyzer-type.pvc", util.AlamedaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaanalyzerCTN, "alameda-analyzer-type-storage", util.AlamedaGroup)
		}
	case util.FedemeterDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "fedemeter-type.pvc", util.FedemeterGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.FedemeterCTN, "fedemeter-type-storage", util.FedemeterGroup)
		}
	case util.InfluxdbDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "my-alameda.influxdb-type.pvc", util.InfluxDBGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.InfluxdbCTN, "influxdb-type-storage", util.InfluxDBGroup)
		}
	case util.GrafanaDPN:
		{
			util.SetStorageToVolumeSource(dep, asp.Storages, "my-alameda.grafana-type.pvc", util.GrafanaGroup)
			util.SetStorageToMountPath(dep, asp.Storages, util.GrafanaCTN, "grafana-type-storage", util.GrafanaGroup)
		}
	case util.AlamedaweavescopeDPN:
		{
			util.SetImagePullPolicy(dep, util.AlamedaweavescopeCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
		}
	case util.AlamedaNotifierDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "alameda-notifier-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.AlamedaNofitierCTN, "alameda-notifier-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-agent-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiAgentCTN, "federatorai-agent-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentAppDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-agent-app-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiAgentAppCTN, "federatorai-agent-app-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentGPUDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-agent-gpu-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiAgentGPUCTN, "federatorai-agent-gpu-type-storage", util.AlamedaGroup)
	case util.FederatoraiRestDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-rest-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiRestCTN, "federatorai-rest-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentPreloaderDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-agent-preloader-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiAgentPreloaderCTN, "federatorai-agent-preloader-type-storage", util.AlamedaGroup)
	case util.FederatoraiFrontendDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-frontend-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiFrontendCTN, "federatorai-frontend-type-storage", util.AlamedaGroup)
	case util.FederatoraiBackendDPN:
		util.SetStorageToVolumeSource(dep, asp.Storages, "federatorai-backend-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.Storages, util.FederatoraiBackendCTN, "federatorai-backend-type-storage", util.AlamedaGroup)
	}

	envVars := getEnvVarsToUpdateByDeployment(dep.Name, asp)
	updateenvvar.UpdateEnvVarsToDeployment(dep, envVars)
}

func GlobalSectionSetParamterToDaemonSet(ds *appsv1.DaemonSet, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	switch ds.Name {
	case util.AlamedaweavescopeAgentDS:
		util.SetDaemonSetImagePullPolicy(ds, util.AlamedaweavescopeAgentCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
	}
}

func GlobalSectionSetParamterToPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	for _, pvcusage := range v1alpha1.PvcUsage {
		if strings.Contains(pvc.Name, string(pvcusage)) {
			util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.Storages, pvcusage)
		}
	}
}

func getEnvVarsToUpdateByDeployment(deploymentName string, asp *alamedaserviceparamter.AlamedaServiceParamter) []corev1.EnvVar {

	var envVars []corev1.EnvVar

	switch deploymentName {
	case util.AlamedaaiDPN:
		envVars = getAlamedaAIEnvVarsToUpdate(asp)
	default:
	}

	return envVars
}

func getAlamedaAIEnvVarsToUpdate(asp *alamedaserviceparamter.AlamedaServiceParamter) []corev1.EnvVar {

	envVars := make([]corev1.EnvVar, 0)

	switch asp.EnableDispatcher {
	case true:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "PREDICT_QUEUE_ENABLED",
			Value: "true",
		})
	case false:
		envVars = append(envVars, corev1.EnvVar{
			Name:  "PREDICT_QUEUE_ENABLED",
			Value: "false",
		})
	}

	return envVars
}
