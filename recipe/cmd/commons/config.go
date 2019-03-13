package commons

import (
	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
	"os"
)

const (
	BuildpacksDir             = "/var/lib/buildpacks"
	OutputDropletLocation     = "/out/droplet.tgz"
	OutputBuildArtifactsCache = "/cache/cache.tgz"
	OutputMetadataLocation    = "/out/result.json"
)

func RecipeConfig() recipe.Config {
	appID := os.Getenv(eirini.EnvAppID)
	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)
	appBitsDownloadURL := os.Getenv(eirini.EnvDownloadURL)
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)

	return recipe.Config{
		AppID:              appID,
		StagingGUID:        stagingGUID,
		CompletionCallback: completionCallback,
		EiriniAddr:         eiriniAddress,
		DropletUploadURL:   dropletUploadURL,
		PackageDownloadURL: appBitsDownloadURL,
	}
}

func BuildpackJson() string {
	return os.Getenv(eirini.EnvBuildpacks)
}

func PacksConfig() recipe.PacksBuilderConf {
	return recipe.PacksBuilderConf{
		BuildpacksDir:             BuildpacksDir,
		OutputDropletLocation:     OutputDropletLocation,
		OutputBuildArtifactsCache: OutputBuildArtifactsCache,
		OutputMetadataLocation:    OutputMetadataLocation,
	}
}
