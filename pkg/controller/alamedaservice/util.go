package alamedaservice

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/containers-ai/federatorai-operator/pkg/component"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
)

func newDefautlImageConfig() component.ImageConfig {
	return component.NewDefautlImageConfig()
}

var (
	relatedImageEnvList = []string{
		"RELATED_IMAGE_ADMISSION_CONTROLLER",
		"RELATED_IMAGE_AI_DISPATCHER",
		"RELATED_IMAGE_AI_ENGINE",
		"RELATED_IMAGE_ANALYZER",
		"RELATED_IMAGE_DATAHUB",
		"RELATED_IMAGE_EVICTIONER",
		"RELATED_IMAGE_EXECUTOR",
		"RELATED_IMAGE_GRAFANA",
		"RELATED_IMAGE_INFLUXDB",
		"RELATED_IMAGE_NOTIFIER",
		"RELATED_IMAGE_OPERATOR",
		"RELATED_IMAGE_RABBITMQ",
		"RELATED_IMAGE_RECOMMENDER",
		"RELATED_IMAGE_FEDEMETER_API",
		"RELATED_IMAGE_FEDERATORAI_AGENT_GPU",
		"RELATED_IMAGE_FEDERATORAI_AGENT_PRELOADER",
		"RELATED_IMAGE_FEDERATORAI_AGENT",
		"RELATED_IMAGE_FEDERATORAI_AGENT_APP",
		"RELATED_IMAGE_FEDERATORAI_RESTAPI",
		"RELATED_IMAGE_FEDEMETER_INFLUXDB",
		"RELATED_IMAGE_FEDERATORAI_DASHBOARD_FRONTEND",
		"RELATED_IMAGE_FEDERATORAI_DASHBOARD_BACKEND",
	}
	imageLocationImmutableFields = map[string]struct{}{
		"InfluxDB": struct{}{},
		"Alpine":   struct{}{},
	}
	imageTagImmutableFields = map[string]struct{}{
		"InfluxDB": struct{}{},
		"Alpine":   struct{}{},
	}
)

// setImageConfigWithImageLocation replace each image's location under imageConfig with imageLocation
// Ignore fields exist in map "imageLocationImmutableFields"
// If image does not contain "ㄥ", insert (location + "ㄥ") into it.
// If image is empty, ignore it.
func setImageConfigWithImageLocation(imageConfig component.ImageConfig, imageLocation string) component.ImageConfig {

	rf := reflect.TypeOf(imageConfig)
	rv := reflect.ValueOf(&imageConfig).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if _, exist := imageLocationImmutableFields[rf.Field(i).Name]; exist {
			continue
		}

		v := rv.Field(i).String()
		// skip mutating image location since the image is empty
		if v == "" {
			continue
		}

		indexLastSlash := strings.LastIndex(v, "/")
		if indexLastSlash != -1 {
			v = fmt.Sprintf(`%s%s`, imageLocation, v[indexLastSlash:])
		} else {
			v = fmt.Sprintf(`%s/%s`, imageLocation, v)
		}
		rv.Field(i).SetString(v)
	}

	return imageConfig
}

// setImageConfigWithImageTag append or replace each image's tag under imageConfig with imageTag
// Ignore fields exist in map "imageTagImmutableFields"
// If imageTag is empty, remove tag from each image.
// If image is empty, ignore it.
func setImageConfigWithImageTag(imageConfig component.ImageConfig, imageTag string) component.ImageConfig {

	rf := reflect.TypeOf(imageConfig)
	rv := reflect.ValueOf(&imageConfig).Elem()
	for i := 0; i < rv.NumField(); i++ {
		if _, exist := imageTagImmutableFields[rf.Field(i).Name]; exist {
			continue
		}

		v := rv.Field(i).String()
		// skip mutating image tag since the image is empty
		if v == "" {
			continue
		}

		indexLastColon := strings.LastIndex(v, ":")
		if indexLastColon != -1 {
			if imageTag != "" {
				v = fmt.Sprintf(`%s:%s`, v[:indexLastColon], imageTag)
			} else {
				v = fmt.Sprintf(`%s`, v[:indexLastColon])
			}
		} else {
			if imageTag != "" {
				v = fmt.Sprintf(`%s:%s`, v, imageTag)
			} else {
				v = fmt.Sprintf(`%s`, v)
			}
		}
		rv.Field(i).SetString(v)
	}
	return imageConfig
}

func setImageConfigWithAlamedaServiceParameterGlobalConfiguration(imageConfig component.ImageConfig, asp alamedaserviceparamter.AlamedaServiceParamter) component.ImageConfig {
	if asp.ImageLocation != "" {
		imageConfig = setImageConfigWithImageLocation(imageConfig, asp.ImageLocation)
	}

	if imageTag := asp.Version; imageTag != "" {
		imageConfig = setImageConfigWithImageTag(imageConfig, imageTag)
	}
	return imageConfig
}

func setImageConfigWithAlamedaServiceParameter(imageConfig component.ImageConfig, asp alamedaserviceparamter.AlamedaServiceParamter) component.ImageConfig {
	if image := asp.InfluxdbSectionSet.Image; image != "" {
		imageConfig.SetInfluxdb(image)
	}

	if image := asp.GrafanaSectionSet.Image; image != "" {
		imageConfig.SetGrafana(image)
	}

	if image := asp.AlamedaAISectionSet.Image; image != "" {
		imageConfig.SetAIEngine(image)
	}

	if image := asp.AlamedaOperatorSectionSet.Image; image != "" {
		imageConfig.SetOperator(image)
	}

	if image := asp.AlamedaDatahubSectionSet.Image; image != "" {
		imageConfig.SetDatahub(image)
	}

	if image := asp.AlamedaEvictionerSectionSet.Image; image != "" {
		imageConfig.SetEvictioner(image)
	}

	if image := asp.AdmissionControllerSectionSet.Image; image != "" {
		imageConfig.SetAdmissionController(image)
	}

	if image := asp.AlamedaRecommenderSectionSet.Image; image != "" {
		imageConfig.SetRecommender(image)
	}

	if image := asp.AlamedaExecutorSectionSet.Image; image != "" {
		imageConfig.SetExecutor(image)
	}

	if image := asp.AlamedaDispatcherSectionSet.Image; image != "" {
		imageConfig.SetAIDispatcher(image)
	}

	if image := asp.AlamedaFedemeterSectionSet.Image; image != "" {
		imageConfig.SetFedemeterAPI(image)
	}

	if image := asp.AlamedaFedemeterInfluxdbSectionSet.Image; image != "" {
		imageConfig.SetFedemeterInfluxdb(image)
	}

	if image := asp.AlamedaAnalyzerSectionSet.Image; image != "" {
		imageConfig.SetAnalyzer(image)
	}

	if image := asp.AlamedaNotifierSectionSet.Image; image != "" {
		imageConfig.SetNotifier(image)
	}

	if image := asp.AlamedaRabbitMQSectionSet.Image; image != "" {
		imageConfig.SetRabbitMQ(image)
	}

	if image := asp.FederatoraiAgentSectionSet.Image; image != "" {
		imageConfig.SetFederatoraiAgent(image)
	}

	if image := asp.FederatoraiAgentGPUSectionSet.Image; image != "" {
		imageConfig.SetFederatoraiAgentGPU(image)
	}

	if image := asp.FederatoraiAgentAppSectionSet.Image; image != "" {
		imageConfig.SetFederatoraiAgentApp(image)
	}

	if image := asp.FederatoraiRestSectionSet.Image; image != "" {
		imageConfig.SetFederatoraiRestAPI(image)
	}

	if image := asp.FederatoraiAgentPreloaderSectionSet.Image; image != "" {
		imageConfig.SetFederatoraiAgentPreloader(image)
	}
	if image := asp.FederatoraiFrontendSectionSet.Image; image != "" {
		imageConfig.SetFrontend(image)
	}
	if image := asp.FederatoraiBackendSectionSet.Image; image != "" {
		imageConfig.SetBackend(image)
	}
	return imageConfig
}

func setImageConfigWithEnv(imageConfig component.ImageConfig) component.ImageConfig {

	for _, env := range relatedImageEnvList {
		switch env {
		case "RELATED_IMAGE_ADMISSION_CONTROLLER":
			if value := os.Getenv("RELATED_IMAGE_ADMISSION_CONTROLLER"); value != "" {
				imageConfig.SetAdmissionController(value)
			}
		case "RELATED_IMAGE_AI_DISPATCHER":
			if value := os.Getenv("RELATED_IMAGE_AI_DISPATCHER"); value != "" {
				imageConfig.SetAIDispatcher(value)
			}
		case "RELATED_IMAGE_AI_ENGINE":
			if value := os.Getenv("RELATED_IMAGE_AI_ENGINE"); value != "" {
				imageConfig.SetAIEngine(value)
			}
		case "RELATED_IMAGE_ANALYZER":
			if value := os.Getenv("RELATED_IMAGE_ANALYZER"); value != "" {
				imageConfig.SetAnalyzer(value)
			}
		case "RELATED_IMAGE_DATAHUB":
			if value := os.Getenv("RELATED_IMAGE_DATAHUB"); value != "" {
				imageConfig.SetDatahub(value)
			}
		case "RELATED_IMAGE_EVICTIONER":
			if value := os.Getenv("RELATED_IMAGE_EVICTIONER"); value != "" {
				imageConfig.SetEvictioner(value)
			}
		case "RELATED_IMAGE_EXECUTOR":
			if value := os.Getenv("RELATED_IMAGE_EXECUTOR"); value != "" {
				imageConfig.SetExecutor(value)
			}
		case "RELATED_IMAGE_GRAFANA":
			if value := os.Getenv("RELATED_IMAGE_GRAFANA"); value != "" {
				imageConfig.SetGrafana(value)
			}
		case "RELATED_IMAGE_INFLUXDB":
			if value := os.Getenv("RELATED_IMAGE_INFLUXDB"); value != "" {
				imageConfig.SetInfluxdb(value)
			}
		case "RELATED_IMAGE_NOTIFIER":
			if value := os.Getenv("RELATED_IMAGE_NOTIFIER"); value != "" {
				imageConfig.SetNotifier(value)
			}
		case "RELATED_IMAGE_OPERATOR":
			if value := os.Getenv("RELATED_IMAGE_OPERATOR"); value != "" {
				imageConfig.SetOperator(value)
			}
		case "RELATED_IMAGE_RABBITMQ":
			if value := os.Getenv("RELATED_IMAGE_RABBITMQ"); value != "" {
				imageConfig.SetRabbitMQ(value)
			}
		case "RELATED_IMAGE_RECOMMENDER":
			if value := os.Getenv("RELATED_IMAGE_RECOMMENDER"); value != "" {
				imageConfig.SetRecommender(value)
			}
		case "RELATED_IMAGE_FEDEMETER_API":
			if value := os.Getenv("RELATED_IMAGE_FEDEMETER_API"); value != "" {
				imageConfig.SetFedemeterAPI(value)
			}
		case "RELATED_IMAGE_FEDERATORAI_AGENT_GPU":
			if value := os.Getenv("RELATED_IMAGE_FEDERATORAI_AGENT_GPU"); value != "" {
				imageConfig.SetFederatoraiAgentGPU(value)
			}
		case "RELATED_IMAGE_FEDERATORAI_AGENT_PRELOADER":
			if value := os.Getenv("RELATED_IMAGE_FEDERATORAI_AGENT_PRELOADER"); value != "" {
				imageConfig.SetFederatoraiAgentPreloader(value)
			}
		case "RELATED_IMAGE_FEDERATORAI_AGENT":
			if value := os.Getenv("RELATED_IMAGE_FEDERATORAI_AGENT"); value != "" {
				imageConfig.SetFederatoraiAgent(value)
			}
		case "RELATED_IMAGE_FEDERATORAI_AGENT_APP":
			if value := os.Getenv("RELATED_IMAGE_FEDERATORAI_AGENT_APP"); value != "" {
				imageConfig.SetFederatoraiAgentApp(value)
			}
		case "RELATED_IMAGE_FEDERATORAI_RESTAPI":
			if value := os.Getenv("RELATED_IMAGE_FEDERATORAI_RESTAPI"); value != "" {
				imageConfig.SetFederatoraiRestAPI(value)
			}
		case "RELATED_IMAGE_FEDEMETER_INFLUXDB":
			if value := os.Getenv("RELATED_IMAGE_FEDEMETER_INFLUXDB"); value != "" {
				imageConfig.SetFedemeterInfluxdb(value)
			}
		case "RELATED_IMAGE_FEDERATORAI_DASHBOARD_FRONTEND":
			if value := os.Getenv("RELATED_IMAGE_FEDERATORAI_DASHBOARD_FRONTEND"); value != "" {
				imageConfig.SetFrontend(value)
			}
		case "RELATED_IMAGE_FEDERATORAI_DASHBOARD_BACKEND":
			if value := os.Getenv("RELATED_IMAGE_FEDERATORAI_DASHBOARD_BACKEND"); value != "" {
				imageConfig.SetBackend(value)
			}
		}

	}

	return imageConfig
}
