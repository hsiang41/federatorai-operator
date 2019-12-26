package component

const (
	defaultImageAdmissionController       = "quay.io/prophetstor/alameda-admission-ubi:v0.3.8"
	defaultImageAIDispatcher              = "quay.io/prophetstor/alameda-ai-dispatcher:latest"
	defaultImageAIEngine                  = "quay.io/prophetstor/alameda-ai:v0.3.8"
	defaultImageAnalyzer                  = "quay.io/prophetstor/alameda-analyzer-ubi:v0.3.54"
	defaultImageDatahub                   = "quay.io/prophetstor/alameda-datahub-ubi:v0.3.8"
	defaultImageEvictioner                = "quay.io/prophetstor/alameda-evictioner-ubi:v0.3.8"
	defaultImageExecutor                  = "quay.io/prophetstor/alameda-executor-ubi:v0.3.8"
	defaultImageAlpine                    = "alpine"
	defaultImageGrafana                   = "quay.io/prophetstor/alameda-grafana:latest"
	defaultImageInfluxDB                  = "influxdb:1.7-alpine"
	defaultImageNotifier                  = "quay.io/prophetstor/alameda-notifier-ubi:v4.2.259        "
	defaultImageOperator                  = "quay.io/prophetstor/alameda-operator-ubi:v0.3.8"
	defaultImageRabbitMQ                  = "quay.io/prophetstor/alameda-rabbitmq:latest"
	defaultImageRecommender               = "quay.io/prophetstor/alameda-recommender-ubi:v0.3.8"
	defaultImageFedemeterAPI              = "quay.io/prophetstor/fedemeter-api-ubi:v0.3.39"
	defaultImageFederatoraiAgentGPU       = "quay.io/prophetstor/federatorai-agent-gpu:v4.2.300"
	defaultImageFederatoraiAgentPreloader = "quay.io/prophetstor/federatorai-agent-preloader:v4.2.512"
	defaultImageFederatoraiAgent          = "quay.io/prophetstor/federatorai-agent-ubi:v4.2.259"
	defaultImageFederatoraiRestAPI        = "quay.io/prophetstor/federatorai-rest-ubi:v4.2.504"
	defaultImageFedemeterInfluxdb         = "quay.io/prophetstor/fedemeter-influxdb:v0.3.39"
)

type ImageConfig struct {
	AdmissionController       string
	AIDispatcher              string
	AIEngine                  string
	Analyzer                  string
	Datahub                   string
	Evictioner                string
	Executor                  string
	Alpine                    string
	Grafana                   string
	InfluxDB                  string
	Notifier                  string
	Operator                  string
	RabbitMQ                  string
	Recommender               string
	FedemeterAPI              string
	FederatoraiAgentGPU       string
	FederatoraiAgentPreloader string
	FederatoraiAgent          string
	FederatoraiRestAPI        string
	FedemeterInfluxDB         string
}

// NewDefautlImageConfig returns ImageConfig with default value
func NewDefautlImageConfig() ImageConfig {
	return ImageConfig{
		AdmissionController:       defaultImageAdmissionController,
		AIDispatcher:              defaultImageAIDispatcher,
		AIEngine:                  defaultImageAIEngine,
		Analyzer:                  defaultImageAnalyzer,
		Datahub:                   defaultImageDatahub,
		Evictioner:                defaultImageEvictioner,
		Executor:                  defaultImageExecutor,
		Alpine:                    defaultImageAlpine,
		Grafana:                   defaultImageGrafana,
		InfluxDB:                  defaultImageInfluxDB,
		Notifier:                  defaultImageNotifier,
		Operator:                  defaultImageOperator,
		RabbitMQ:                  defaultImageRabbitMQ,
		Recommender:               defaultImageRecommender,
		FedemeterAPI:              defaultImageFedemeterAPI,
		FederatoraiAgentGPU:       defaultImageFederatoraiAgentGPU,
		FederatoraiAgentPreloader: defaultImageFederatoraiAgentPreloader,
		FederatoraiAgent:          defaultImageFederatoraiAgent,
		FederatoraiRestAPI:        defaultImageFederatoraiRestAPI,
		FedemeterInfluxDB:         defaultImageFedemeterInfluxdb,
	}
}

// SetAdmissionController sets image to imageConfig
func (i *ImageConfig) SetAdmissionController(image string) {
	i.AdmissionController = image
}

// SetAIDispatcher sets image to imageConfig
func (i *ImageConfig) SetAIDispatcher(image string) {
	i.AIDispatcher = image
}

// SetAIEngine sets image to imageConfig
func (i *ImageConfig) SetAIEngine(image string) {
	i.AIEngine = image
}

// SetAnalyzer sets image to imageConfig
func (i *ImageConfig) SetAnalyzer(image string) {
	i.Analyzer = image
}

// SetDatahub sets image to imageConfig
func (i *ImageConfig) SetDatahub(image string) {
	i.Datahub = image
}

// SetEvictioner sets image to imageConfig
func (i *ImageConfig) SetEvictioner(image string) {
	i.Evictioner = image
}

// SetExecutor sets image to imageConfig
func (i *ImageConfig) SetExecutor(image string) {
	i.Executor = image
}

// SetAlpine sets image to imageConfig
func (i *ImageConfig) SetAlpine(image string) {
	i.Alpine = image
}

// SetGrafana sets image to imageConfig
func (i *ImageConfig) SetGrafana(image string) {
	i.Grafana = image
}

// SetInfluxdb sets image to imageConfig
func (i *ImageConfig) SetInfluxdb(image string) {
	i.InfluxDB = image
}

// SetNotifier sets image to imageConfig
func (i *ImageConfig) SetNotifier(image string) {
	i.Notifier = image
}

// SetOperator sets image to imageConfig
func (i *ImageConfig) SetOperator(image string) {
	i.Operator = image
}

// SetRabbitMQ sets image to imageConfig
func (i *ImageConfig) SetRabbitMQ(image string) {
	i.RabbitMQ = image
}

// SetRecommender sets image to imageConfig
func (i *ImageConfig) SetRecommender(image string) {
	i.Recommender = image
}

// SetFedemeterAPI sets image to imageConfig
func (i *ImageConfig) SetFedemeterAPI(image string) {
	i.FedemeterAPI = image
}

// SetFederatoraiAgentGPU sets image to imageConfig
func (i *ImageConfig) SetFederatoraiAgentGPU(image string) {
	i.FederatoraiAgentGPU = image
}

// SetFederatoraiAgentPreloader sets image to imageConfig
func (i *ImageConfig) SetFederatoraiAgentPreloader(image string) {
	i.FederatoraiAgentPreloader = image
}

// SetFederatoraiAgent sets image to imageConfig
func (i *ImageConfig) SetFederatoraiAgent(image string) {
	i.FederatoraiAgent = image
}

// SetFederatoraiRestAPI sets image to imageConfig
func (i *ImageConfig) SetFederatoraiRestAPI(image string) {
	i.FederatoraiRestAPI = image
}

// SetFedemeterInfluxdb sets image to imageConfig
func (i *ImageConfig) SetFedemeterInfluxdb(image string) {
	i.FedemeterInfluxDB = image
}
