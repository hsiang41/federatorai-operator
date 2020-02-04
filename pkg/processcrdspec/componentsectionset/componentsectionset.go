package componentsectionset

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
	"github.com/containers-ai/federatorai-operator/pkg/util"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

func SectionSetParamterToDeployment(dep *appsv1.Deployment, asp *alamedaserviceparamter.AlamedaServiceParamter) {

	envVars := []corev1.EnvVar{}
	switch dep.Name {
	case util.AlamedaaiDPN:
		envVars = asp.AlamedaAISectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaaiCTN, asp.AlamedaAISectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaAISectionSet.Storages, "alameda-ai-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaAISectionSet.Storages, util.AlamedaaiCTN, "alameda-ai-type-storage", util.AlamedaGroup)
	case util.AlamedaoperatorDPN:
		envVars = asp.AlamedaOperatorSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaoperatorCTN, asp.AlamedaOperatorSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaOperatorSectionSet.Storages, "alameda-operator-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaOperatorSectionSet.Storages, util.AlamedaoperatorCTN, "alameda-operator-type-storage", util.AlamedaGroup)
	case util.AlamedadatahubDPN:
		envVars = asp.AlamedaDatahubSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedadatahubCTN, asp.AlamedaDatahubSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaDatahubSectionSet.Storages, "alameda-datahub-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaDatahubSectionSet.Storages, util.AlamedadatahubCTN, "alameda-datahub-type-storage", util.AlamedaGroup)
	case util.AlamedaevictionerDPN:
		envVars = asp.AlamedaEvictionerSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaevictionerCTN, asp.AlamedaEvictionerSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaEvictionerSectionSet.Storages, "alameda-evictioner-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaEvictionerSectionSet.Storages, util.AlamedaevictionerCTN, "alameda-evictioner-type-storage", util.AlamedaGroup)
	case util.AdmissioncontrollerDPN:
		envVars = asp.AdmissionControllerSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AdmissioncontrollerCTN, asp.AdmissionControllerSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AdmissionControllerSectionSet.Storages, "admission-controller-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AdmissionControllerSectionSet.Storages, util.AdmissioncontrollerCTN, "admission-controller-type-storage", util.AlamedaGroup)
	case util.AlamedarecommenderDPN:
		envVars = asp.AlamedaRecommenderSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedarecommenderCTN, asp.AlamedaRecommenderSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaRecommenderSectionSet.Storages, "alameda-recommender-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaRecommenderSectionSet.Storages, util.AlamedarecommenderCTN, "alameda-recommender-type.pvc", util.AlamedaGroup)
	case util.AlamedaexecutorDPN:
		envVars = asp.AlamedaExecutorSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaexecutorCTN, asp.AlamedaExecutorSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaExecutorSectionSet.Storages, "alameda-executor-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaExecutorSectionSet.Storages, util.AlamedaexecutorCTN, "alameda-executor-type.pvc", util.AlamedaGroup)
	case util.AlamedadispatcherDPN:
		envVars = asp.AlamedaDispatcherSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedadispatcherCTN, asp.AlamedaDispatcherSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaDispatcherSectionSet.Storages, "alameda-dispatcher-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaDispatcherSectionSet.Storages, util.AlamedadispatcherCTN, "alameda-dispatcher-type.pvc", util.AlamedaGroup)
	case util.AlamedaRabbitMQDPN:
		envVars = asp.AlamedaRabbitMQSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaRabbitMQCTN, asp.AlamedaRabbitMQSectionSet.ImagePullPolicy)
	case util.AlamedaanalyzerDPN:
		envVars = asp.AlamedaAnalyzerSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaanalyzerCTN, asp.AlamedaAnalyzerSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaAnalyzerSectionSet.Storages, "alameda-analyzer-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaAnalyzerSectionSet.Storages, util.AlamedaanalyzerCTN, "alameda-analyzer-type.pvc", util.AlamedaGroup)
	case util.FedemeterDPN:
		envVars = asp.AlamedaFedemeterSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.FedemeterCTN, asp.AlamedaFedemeterSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaFedemeterSectionSet.Storages, "fedemeter-type.pvc", util.FedemeterGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaFedemeterSectionSet.Storages, util.FedemeterCTN, "fedemeter-type.pvc", util.FedemeterGroup)
	case util.InfluxdbDPN:
		envVars = asp.InfluxdbSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.InfluxdbCTN, asp.InfluxdbSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.InfluxdbSectionSet.Storages, "my-alameda.influxdb-type.pvc", util.InfluxDBGroup)
		util.SetStorageToMountPath(dep, asp.InfluxdbSectionSet.Storages, util.InfluxdbCTN, "influxdb-type-storage", util.InfluxDBGroup)
	case util.GrafanaDPN:
		envVars = asp.GrafanaSectionSet.EnvVars
		util.SetBootStrapImageStruct(dep, asp.GrafanaSectionSet, util.GetTokenCTN)
		util.SetImagePullPolicy(dep, util.GrafanaCTN, asp.GrafanaSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.GrafanaSectionSet.Storages, "my-alameda.grafana-type.pvc", util.GrafanaGroup)
		util.SetStorageToMountPath(dep, asp.GrafanaSectionSet.Storages, util.GrafanaCTN, "grafana-type-storage", util.GrafanaGroup)
	case util.AlamedaweavescopeDPN:
		envVars = asp.AlamedaWeavescopeSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaweavescopeCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
	case util.AlamedaweavescopeProbeDPN:
		envVars = asp.AlamedaWeavescopeSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaweavescopeProbeCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
	case util.AlamedaNotifierDPN:
		envVars = asp.AlamedaNotifierSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.AlamedaNofitierCTN, asp.AlamedaNotifierSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.AlamedaNotifierSectionSet.Storages, "alameda-notifier-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.AlamedaNotifierSectionSet.Storages, util.AlamedaNofitierCTN, "alameda-notifier-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentDPN:
		envVars = asp.FederatoraiAgentSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.FederatoraiAgentCTN, asp.FederatoraiAgentSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.FederatoraiAgentSectionSet.Storages, "federatorai-agent-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.FederatoraiAgentSectionSet.Storages, util.FederatoraiAgentCTN, "federatorai-agent-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentGPUDPN:
		envVars = asp.FederatoraiAgentGPUSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.FederatoraiAgentGPUCTN, asp.FederatoraiAgentGPUSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.FederatoraiAgentGPUSectionSet.Storages, "federatorai-agent-gpu-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.FederatoraiAgentGPUSectionSet.Storages, util.FederatoraiAgentGPUCTN, "federatorai-agent-gpu-type-storage", util.AlamedaGroup)
	case util.FederatoraiRestDPN:
		envVars = asp.FederatoraiRestSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.FederatoraiRestCTN, asp.FederatoraiRestSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.FederatoraiRestSectionSet.Storages, "federatorai-rest-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.FederatoraiRestSectionSet.Storages, util.FederatoraiRestCTN, "federatorai-rest-type-storage", util.AlamedaGroup)
	case util.FederatoraiAgentPreloaderDPN:
		envVars = asp.FederatoraiAgentPreloaderSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.FederatoraiAgentPreloaderCTN, asp.FederatoraiAgentPreloaderSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.FederatoraiAgentPreloaderSectionSet.Storages, "federatorai-agent-preloader-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.FederatoraiAgentPreloaderSectionSet.Storages, util.FederatoraiAgentPreloaderCTN, "federatorai-agent-preloader-type-storage", util.AlamedaGroup)
	case util.FederatoraiFrontendDPN:
		envVars = asp.FederatoraiFrontendSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.FederatoraiFrontendCTN, asp.FederatoraiFrontendSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.FederatoraiFrontendSectionSet.Storages, "federatorai-frontend-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.FederatoraiFrontendSectionSet.Storages, util.FederatoraiFrontendCTN, "federatorai-frontend-type-storage", util.AlamedaGroup)
	case util.FederatoraiBackendDPN:
		envVars = asp.FederatoraiBackendSectionSet.EnvVars
		util.SetImagePullPolicy(dep, util.FederatoraiBackendCTN, asp.FederatoraiBackendSectionSet.ImagePullPolicy)
		util.SetStorageToVolumeSource(dep, asp.FederatoraiBackendSectionSet.Storages, "federatorai-backend-type.pvc", util.AlamedaGroup)
		util.SetStorageToMountPath(dep, asp.FederatoraiBackendSectionSet.Storages, util.FederatoraiBackendCTN, "federatorai-backend-type-storage", util.AlamedaGroup)
	}

	for i, container := range dep.Spec.Template.Spec.Containers {
		newEnv := replaceOrAppendEnvVar(container.Env, envVars)
		dep.Spec.Template.Spec.Containers[i].Env = newEnv
	}
}

func SectionSetParamterToDaemonSet(ds *appsv1.DaemonSet, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	envVars := []corev1.EnvVar{}
	switch ds.Name {
	case util.AlamedaweavescopeAgentDS:
		envVars = asp.AlamedaWeavescopeSectionSet.EnvVars
		util.SetDaemonSetImagePullPolicy(ds, util.AlamedaweavescopeAgentCTN, asp.AlamedaWeavescopeSectionSet.ImagePullPolicy)
	}

	for i, container := range ds.Spec.Template.Spec.Containers {
		newEnv := replaceOrAppendEnvVar(container.Env, envVars)
		ds.Spec.Template.Spec.Containers[i].Env = newEnv
	}
}

func SectionSetParamterToStatefulSet(ss *appsv1.StatefulSet, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	envVars := []corev1.EnvVar{}
	switch ss.Name {
	case util.FedemeterInflixDBSSN:
		envVars = asp.AlamedaFedemeterInfluxdbSectionSet.EnvVars
		util.SetStatefulSetImagePullPolicy(ss, util.FedemeterInfluxDBCTN, asp.InfluxdbSectionSet.ImagePullPolicy)
	}

	for i, container := range ss.Spec.Template.Spec.Containers {
		newEnv := replaceOrAppendEnvVar(container.Env, envVars)
		ss.Spec.Template.Spec.Containers[i].Env = newEnv
	}
}

func SectionSetParamterToPersistentVolumeClaim(pvc *corev1.PersistentVolumeClaim, asp *alamedaserviceparamter.AlamedaServiceParamter) {
	for _, pvcusage := range v1alpha1.PvcUsage {
		switch pvc.Name {
		case fmt.Sprintf("alameda-ai-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaAISectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-operator-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaOperatorSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-datahub-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaDatahubSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-evictioner-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaEvictionerSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("admission-controller-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AdmissionControllerSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-recommender-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaRecommenderSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-executor-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaExecutorSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-dispatcher-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaDispatcherSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("alameda-analyzer-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaAnalyzerSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("fedemeter-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.AlamedaFedemeterSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("my-alameda.influxdb-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.InfluxdbSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("my-alameda.grafana-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.GrafanaSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("federatorai-agent-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.FederatoraiAgentSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("federatorai-agent-gpu-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.FederatoraiAgentGPUSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("federatorai-rest-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.FederatoraiRestSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("federatorai-frontend-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.FederatoraiFrontendSectionSet.Storages, pvcusage)
			}
		case fmt.Sprintf("federatorai-backend-%s.pvc", pvcusage):
			{
				util.SetStorageToPersistentVolumeClaimSpec(pvc, asp.FederatoraiBackendSectionSet.Storages, pvcusage)
			}
		}
	}
}

func SectionSetParamterToService(svc *corev1.Service, asp *alamedaserviceparamter.AlamedaServiceParamter) error {

	if asp == nil {
		return errors.New("AlamedaServiceParamter cannnot be nil")
	}

	for _, serviceExposure := range asp.ServiceExposures {
		if svc.Name == serviceExposure.Name {
			if err := processServiceWithServiceExposureSpec(svc, serviceExposure); err != nil {
				return err
			}
		}
	}

	return nil
}

func processServiceWithServiceExposureSpec(svc *corev1.Service, serviceExposure v1alpha1.ServiceExposureSpec) error {

	switch serviceExposure.Type {
	case v1alpha1.ServiceExposureTypeNodePort:
		if serviceExposure.NodePort == nil {
			return errors.New("NodePort cannot be nil")
		}
		if err := processServiceWithNodePortSpec(svc, *serviceExposure.NodePort); err != nil {
			return errors.Wrap(err, "process service with NodePortSpec failed")
		}
	default:
		return errors.Errorf("unsupported ServiceExposureType \"%s\"", serviceExposure.Type)
	}

	return nil
}

func processServiceWithNodePortSpec(svc *corev1.Service, nodePortSpec v1alpha1.NodePortSpec) error {

	svc.Spec.Type = corev1.ServiceTypeNodePort

	for _, portInNodePortSpec := range nodePortSpec.Ports {
		findPort := false
		for i, portInService := range svc.Spec.Ports {
			if portInNodePortSpec.Port == portInService.Port {
				findPort = true
				svc.Spec.Ports[i].NodePort = portInNodePortSpec.NodePort
				break
			}
		}
		if !findPort {
			return errors.Errorf("port %d not exist in service", portInNodePortSpec.Port)
		}
	}

	return nil
}

func replaceOrAppendEnvVar(target, source []corev1.EnvVar) []corev1.EnvVar {

	newEnv := make([]corev1.EnvVar, len(target))
	for i, env := range target {
		newEnv[i] = *env.DeepCopy()
	}

	for _, envVar := range source {
		exist := false
		for i, env := range newEnv {
			if env.Name == envVar.Name {
				exist = true
				newEnv[i] = *envVar.DeepCopy()
				break
			}
		}
		if !exist {
			newEnv = append(newEnv, *envVar.DeepCopy())
		}
	}
	return newEnv
}
