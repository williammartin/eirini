package eirini

import (
	"context"
	"fmt"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/eirini/models/cf"
	"code.cloudfoundry.org/eirini/opi"
)

const (
	//Environment Variable Names
	EnvDownloadURL               = "DOWNLOAD_URL"
	EnvBuildpacks                = "BUILDPACKS"
	EnvDropletUploadURL          = "DROPLET_UPLOAD_URL"
	EnvAppID                     = "APP_ID"
	EnvStagingGUID               = "STAGING_GUID"
	EnvCompletionCallback        = "COMPLETION_CALLBACK"
	EnvCfUsername                = "CF_USERNAME"
	EnvCfPassword                = "CF_PASSWORD"
	EnvAPIAddress                = "API_ADDRESS"
	EnvEiriniAddress             = "EIRINI_ADDRESS"
	EnvCertsPath                 = "EIRINI_CERTS_PATH"
	EnvBuildpacksDir             = "EIRINI_BUILDPACKS_DIR"
	EnvWorkspaceDir              = "EIRINI_WORKSPACE_DIR"
	EnvOutputDropletLocation     = "EIRINI_OUTPUT_DROPLET_LOCATION"
	EnvOutputBuildArtifactsCache = "EIRINI_OUTPUT_BUILD_ARTIFACTS_CACHE"
	EnvOutputMetadataLocation    = "EIRINI_OUTPUT_METADATA_LOCATION"
	EnvPacksBuilderPath          = "EIRINI_PACKS_BUILDER_PATH"

	RegisteredRoutes = "routes"

	RecipeBuildPacksDir             = "/var/lib/buildpacks"
	RecipeBuildPacksName            = "recipe_buildpacks"
	RecipeWorkspaceDir              = "/workspace"
	RecipeWorkspaceName             = "recipe_workspace"
	RecipeOutputName                = "staging_output"
	RecipeOutputLocation            = "/out"
	RecipeOutputDropletLocation     = "/out/droplet.tgz"
	RecipeOutputBuildArtifactsCache = "/cache/cache.tgz"
	RecipeOutputMetadataLocation    = "/out/result.json"
	RecipePacksBuilderPath          = "/packs/builder"

	CCUploaderInternalURL = "cc-uploader.service.cf.internal"
	CCCertsMountPath      = "/etc/config/certs"
	CCCertsVolumeName     = "cc-certs-volume"
	CCAPICertName         = "cc-server-crt"
	CCAPIKeyName          = "cc-server-crt-key"
	CCUploaderCertName    = "cc-uploader-crt"
	CCUploaderKeyName     = "cc-uploader-crt-key"
	CCInternalCACertName  = "internal-ca-cert"
)

type Config struct {
	Properties Properties `yaml:"opi"`
}

type Properties struct {
	KubeNamespace     string `yaml:"kube_namespace"`
	NatsPassword      string `yaml:"nats_password"`
	NatsIP            string `yaml:"nats_ip"`
	CcUploaderIP      string `yaml:"cc_uploader_ip"`
	CcInternalAPI     string `yaml:"cc_internal_api"`
	CCCertsSecretName string `yaml:"cc_certs_secret_name"`
	RegistryAddress   string `yaml:"registry_address"`
	EiriniAddress     string `yaml:"eirini_address"`
	DownloaderImage   string `yaml:"downloader_image"`
	UploaderImage     string `yaml:"uploader_image"`
	RunnerImage       string `yaml:"runner_image"`

	MetricsSourceAddress string `yaml:"metrics_source_address"`
	LoggregatorAddress   string `yaml:"loggregator_address"`

	LoggregatorCertPath string `yaml:"loggergator_cert_path"`
	LoggregatorKeyPath  string `yaml:"loggregator_key_path"`
	LoggregatorCAPath   string `yaml:"loggregator_ca_path"`

	CCCertPath string `yaml:"cc_cert_path"`
	CCKeyPath  string `yaml:"cc_key_path"`
	CCCAPath   string `yaml:"cc_ca_path"`
}

//go:generate counterfeiter . Stager
type Stager interface {
	Stage(string, cf.StagingRequest) error
	CompleteStaging(*models.TaskCallbackResponse) error
}

type StagerConfig struct {
	EiriniAddress   string
	DownloaderImage string
	UploaderImage   string
	RunnerImage     string
}

//go:generate counterfeiter . Extractor
type Extractor interface {
	Extract(src, targetDir string) error
}

//go:generate counterfeiter . Bifrost
type Bifrost interface {
	Transfer(ctx context.Context, request cf.DesireLRPRequest) error
	List(ctx context.Context) ([]*models.DesiredLRPSchedulingInfo, error)
	Update(ctx context.Context, update cf.UpdateDesiredLRPRequest) error
	Stop(ctx context.Context, identifier opi.LRPIdentifier) error
	GetApp(ctx context.Context, identifier opi.LRPIdentifier) *models.DesiredLRP
	GetInstances(ctx context.Context, identifier opi.LRPIdentifier) ([]*cf.Instance, error)
}

func GetInternalServiceName(appName string) string {
	//Prefix service as the appName could start with numerical characters, which is not allowed
	return fmt.Sprintf("cf-%s", appName)
}

func GetInternalHeadlessServiceName(appName string) string {
	//Prefix service as the appName could start with numerical characters, which is not allowed
	return fmt.Sprintf("cf-%s-headless", appName)
}
