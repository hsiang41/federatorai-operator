package alamedaservice

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/containers-ai/federatorai-operator/pkg/apis/federatorai/v1alpha1"
	"github.com/containers-ai/federatorai-operator/pkg/component"
	"github.com/containers-ai/federatorai-operator/pkg/processcrdspec/alamedaserviceparamter"
)

func TestSetImageConfigWithImageLocation(t *testing.T) {

	type testCaseHave struct {
		imageLocation string
		imageConfig   component.ImageConfig
	}

	type testCase struct {
		have testCaseHave
		want component.ImageConfig
	}

	testCases := []testCase{
		testCase{
			have: testCaseHave{
				imageLocation: "",
				imageConfig:   component.ImageConfig{},
			},
			want: component.ImageConfig{},
		},
		testCase{
			have: testCaseHave{
				imageLocation: "test-registry/test-repository",
				imageConfig: component.ImageConfig{
					AdmissionController:       "test-image-ADMISSION_CONTROLLER",
					AIDispatcher:              "test-image-AI_DISPATCHER",
					AIEngine:                  "test-image-AI_ENGINE",
					Analyzer:                  "test-image-ANALYZER",
					Datahub:                   "test-image-DATAHUB",
					Evictioner:                "test-image-EVICTIONER",
					Executor:                  "test-image-EXECUTOR",
					Alpine:                    "alpine:latest",
					Grafana:                   "test-image-GRAFANA",
					InfluxDB:                  "influxdb:latest",
					Notifier:                  "test-image-NOTIFIER",
					Operator:                  "test-image-OPERATOR",
					RabbitMQ:                  "test-image-RABBITMQ",
					Recommender:               "test-image-RECOMMENDER",
					FedemeterAPI:              "test-image-FEDEMETER_API",
					FederatoraiAgentGPU:       "test-image-FEDERATORAI_AGENT_GPU",
					FederatoraiAgentPreloader: "test-image-FEDERATORAI_AGENT_PRELOADER",
					FederatoraiAgent:          "test-image-FEDERATORAI_AGENT",
					FederatoraiRestAPI:        "test-image-FEDERATORAI_RESTAPI",
					FedemeterInfluxDB:         "test-image-FEDEMETER_INFLUXDB",
				},
			},
			want: component.ImageConfig{
				AdmissionController:       "test-registry/test-repository/test-image-ADMISSION_CONTROLLER",
				AIDispatcher:              "test-registry/test-repository/test-image-AI_DISPATCHER",
				AIEngine:                  "test-registry/test-repository/test-image-AI_ENGINE",
				Analyzer:                  "test-registry/test-repository/test-image-ANALYZER",
				Datahub:                   "test-registry/test-repository/test-image-DATAHUB",
				Evictioner:                "test-registry/test-repository/test-image-EVICTIONER",
				Executor:                  "test-registry/test-repository/test-image-EXECUTOR",
				Alpine:                    "alpine:latest",
				Grafana:                   "test-registry/test-repository/test-image-GRAFANA",
				InfluxDB:                  "influxdb:latest",
				Notifier:                  "test-registry/test-repository/test-image-NOTIFIER",
				Operator:                  "test-registry/test-repository/test-image-OPERATOR",
				RabbitMQ:                  "test-registry/test-repository/test-image-RABBITMQ",
				Recommender:               "test-registry/test-repository/test-image-RECOMMENDER",
				FedemeterAPI:              "test-registry/test-repository/test-image-FEDEMETER_API",
				FederatoraiAgentGPU:       "test-registry/test-repository/test-image-FEDERATORAI_AGENT_GPU",
				FederatoraiAgentPreloader: "test-registry/test-repository/test-image-FEDERATORAI_AGENT_PRELOADER",
				FederatoraiAgent:          "test-registry/test-repository/test-image-FEDERATORAI_AGENT",
				FederatoraiRestAPI:        "test-registry/test-repository/test-image-FEDERATORAI_RESTAPI",
				FedemeterInfluxDB:         "test-registry/test-repository/test-image-FEDEMETER_INFLUXDB",
			},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		actual := setImageConfigWithImageLocation(tc.have.imageConfig, tc.have.imageLocation)
		assert.Equal(tc.want, actual)
	}
}

func TestSetImageConfigWithImageTag(t *testing.T) {

	type testCaseHave struct {
		imageTag    string
		imageConfig component.ImageConfig
	}

	type testCase struct {
		have testCaseHave
		want component.ImageConfig
	}

	testCases := []testCase{
		testCase{
			have: testCaseHave{
				imageTag:    "",
				imageConfig: component.ImageConfig{},
			},
			want: component.ImageConfig{},
		},
		testCase{
			have: testCaseHave{
				imageTag: "test-tag",
				imageConfig: component.ImageConfig{
					AdmissionController:       "test-image-ADMISSION_CONTROLLER",
					AIDispatcher:              "test-image-AI_DISPATCHER",
					AIEngine:                  "test-image-AI_ENGINE",
					Analyzer:                  "test-image-ANALYZER",
					Datahub:                   "test-image-DATAHUB",
					Evictioner:                "test-image-EVICTIONER",
					Executor:                  "test-image-EXECUTOR",
					Alpine:                    "alpine:latest",
					Grafana:                   "test-image-GRAFANA",
					InfluxDB:                  "influxdb:latest",
					Notifier:                  "test-image-NOTIFIER",
					Operator:                  "test-image-OPERATOR",
					RabbitMQ:                  "test-image-RABBITMQ",
					Recommender:               "test-image-RECOMMENDER",
					FedemeterAPI:              "test-image-FEDEMETER_API",
					FederatoraiAgentGPU:       "test-image-FEDERATORAI_AGENT_GPU",
					FederatoraiAgentPreloader: "test-image-FEDERATORAI_AGENT_PRELOADER",
					FederatoraiAgent:          "test-image-FEDERATORAI_AGENT",
					FederatoraiRestAPI:        "test-image-FEDERATORAI_RESTAPI",
					FedemeterInfluxDB:         "test-image-FEDEMETER_INFLUXDB",
				},
			},
			want: component.ImageConfig{
				AdmissionController:       "test-image-ADMISSION_CONTROLLER:test-tag",
				AIDispatcher:              "test-image-AI_DISPATCHER:test-tag",
				AIEngine:                  "test-image-AI_ENGINE:test-tag",
				Analyzer:                  "test-image-ANALYZER:test-tag",
				Datahub:                   "test-image-DATAHUB:test-tag",
				Evictioner:                "test-image-EVICTIONER:test-tag",
				Executor:                  "test-image-EXECUTOR:test-tag",
				Alpine:                    "alpine:latest",
				Grafana:                   "test-image-GRAFANA:test-tag",
				InfluxDB:                  "influxdb:latest",
				Notifier:                  "test-image-NOTIFIER:test-tag",
				Operator:                  "test-image-OPERATOR:test-tag",
				RabbitMQ:                  "test-image-RABBITMQ:test-tag",
				Recommender:               "test-image-RECOMMENDER:test-tag",
				FedemeterAPI:              "test-image-FEDEMETER_API:test-tag",
				FederatoraiAgentGPU:       "test-image-FEDERATORAI_AGENT_GPU:test-tag",
				FederatoraiAgentPreloader: "test-image-FEDERATORAI_AGENT_PRELOADER:test-tag",
				FederatoraiAgent:          "test-image-FEDERATORAI_AGENT:test-tag",
				FederatoraiRestAPI:        "test-image-FEDERATORAI_RESTAPI:test-tag",
				FedemeterInfluxDB:         "test-image-FEDEMETER_INFLUXDB:test-tag",
			},
		},
		testCase{
			have: testCaseHave{
				imageTag: "test-tag",
				imageConfig: component.ImageConfig{
					AdmissionController: "test-image-ADMISSION_CONTROLLER",
				},
			},
			want: component.ImageConfig{
				AdmissionController: "test-image-ADMISSION_CONTROLLER:test-tag",
			},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		actual := setImageConfigWithImageTag(tc.have.imageConfig, tc.have.imageTag)
		assert.Equal(tc.want, actual)
	}
}

func TestSetImageConfigWithEnv(t *testing.T) {

	type setEnvFunc func(key, value string) error

	type testCase struct {
		have map[string]string
		want component.ImageConfig
	}

	testCases := []testCase{
		testCase{
			have: map[string]string{},
			want: component.ImageConfig{},
		},
		testCase{
			have: map[string]string{
				"RELATED_IMAGE_ADMISSION_CONTROLLER":        "test-image-ADMISSION_CONTROLLER",
				"RELATED_IMAGE_AI_DISPATCHER":               "test-image-AI_DISPATCHER",
				"RELATED_IMAGE_AI_ENGINE":                   "test-image-AI_ENGINE",
				"RELATED_IMAGE_ANALYZER":                    "test-image-ANALYZER",
				"RELATED_IMAGE_DATAHUB":                     "test-image-DATAHUB",
				"RELATED_IMAGE_EVICTIONER":                  "test-image-EVICTIONER",
				"RELATED_IMAGE_EXECUTOR":                    "test-image-EXECUTOR",
				"RELATED_IMAGE_ALPINE":                      "test-image-ALPINE",
				"RELATED_IMAGE_GRAFANA":                     "test-image-GRAFANA",
				"RELATED_IMAGE_INFLUXDB":                    "test-image-INFLUXDB",
				"RELATED_IMAGE_NOTIFIER":                    "test-image-NOTIFIER",
				"RELATED_IMAGE_OPERATOR":                    "test-image-OPERATOR",
				"RELATED_IMAGE_RABBITMQ":                    "test-image-RABBITMQ",
				"RELATED_IMAGE_RECOMMENDER":                 "test-image-RECOMMENDER",
				"RELATED_IMAGE_FEDEMETER_API":               "test-image-FEDEMETER_API",
				"RELATED_IMAGE_FEDERATORAI_AGENT_GPU":       "test-image-FEDERATORAI_AGENT_GPU",
				"RELATED_IMAGE_FEDERATORAI_AGENT_PRELOADER": "test-image-FEDERATORAI_AGENT_PRELOADER",
				"RELATED_IMAGE_FEDERATORAI_AGENT":           "test-image-FEDERATORAI_AGENT",
				"RELATED_IMAGE_FEDERATORAI_RESTAPI":         "test-image-FEDERATORAI_RESTAPI",
				"RELATED_IMAGE_FEDEMETER_INFLUXDB":          "test-image-FEDEMETER_INFLUXDB",
			},
			want: component.ImageConfig{
				AdmissionController:       "test-image-ADMISSION_CONTROLLER",
				AIDispatcher:              "test-image-AI_DISPATCHER",
				AIEngine:                  "test-image-AI_ENGINE",
				Analyzer:                  "test-image-ANALYZER",
				Datahub:                   "test-image-DATAHUB",
				Evictioner:                "test-image-EVICTIONER",
				Executor:                  "test-image-EXECUTOR",
				Alpine:                    "test-image-ALPINE",
				Grafana:                   "test-image-GRAFANA",
				InfluxDB:                  "test-image-INFLUXDB",
				Notifier:                  "test-image-NOTIFIER",
				Operator:                  "test-image-OPERATOR",
				RabbitMQ:                  "test-image-RABBITMQ",
				Recommender:               "test-image-RECOMMENDER",
				FedemeterAPI:              "test-image-FEDEMETER_API",
				FederatoraiAgentGPU:       "test-image-FEDERATORAI_AGENT_GPU",
				FederatoraiAgentPreloader: "test-image-FEDERATORAI_AGENT_PRELOADER",
				FederatoraiAgent:          "test-image-FEDERATORAI_AGENT",
				FederatoraiRestAPI:        "test-image-FEDERATORAI_RESTAPI",
				FedemeterInfluxDB:         "test-image-FEDEMETER_INFLUXDB",
			},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {

		for k, v := range tc.have {
			err := os.Setenv(k, v)
			assert.NoError(err)
		}
		actual := setImageConfigWithEnv(component.ImageConfig{})
		assert.Equal(tc.want, actual)
	}
}

func TestSetImageConfigWithAlamedaServiceParameter(t *testing.T) {

	type testCaseHave struct {
		asp         alamedaserviceparamter.AlamedaServiceParamter
		imageConfig component.ImageConfig
	}

	type testCase struct {
		have testCaseHave
		want component.ImageConfig
	}

	testCases := []testCase{
		testCase{
			have: testCaseHave{},
			want: component.ImageConfig{},
		},
		testCase{
			have: testCaseHave{
				asp: alamedaserviceparamter.AlamedaServiceParamter{
					Version:       "vx.x.x",
					ImageLocation: "test-registry/test-repo",
					AlamedaDatahubSectionSet: v1alpha1.AlamedaComponentSpec{
						Image: "datahub-test",
					},
				},
				imageConfig: component.ImageConfig{},
			},
			want: component.ImageConfig{
				Datahub: "datahub-test",
			},
		},
		testCase{
			have: testCaseHave{
				asp: alamedaserviceparamter.AlamedaServiceParamter{
					Version:                  "vx.x.x",
					ImageLocation:            "test-registry/test-repo",
					AlamedaDatahubSectionSet: v1alpha1.AlamedaComponentSpec{},
				},
				imageConfig: component.ImageConfig{
					Datahub: "datahub-test",
				},
			},
			want: component.ImageConfig{
				Datahub: "test-registry/test-repo/datahub-test:vx.x.x",
			},
		},
	}

	assert := assert.New(t)
	for _, tc := range testCases {
		actual := setImageConfigWithAlamedaServiceParameter(tc.have.imageConfig, tc.have.asp)
		assert.Equal(tc.want, actual)
	}
}
