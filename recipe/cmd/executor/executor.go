package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
)

func main() {

	stagingGUID := os.Getenv(eirini.EnvStagingGUID)
	completionCallback := os.Getenv(eirini.EnvCompletionCallback)
	eiriniAddress := os.Getenv(eirini.EnvEiriniAddress)
	buildpacksDir, ok := os.LookupEnv(eirini.EnvBuildpacksDir)
	if !ok {
		buildpacksDir = eirini.RecipeBuildPacksDir
	}

	outputDropletLocation, ok := os.LookupEnv(eirini.EnvOutputDropletLocation)
	if !ok {
		outputDropletLocation = eirini.RecipeOutputDropletLocation
	}

	outputBuildArtifactsCache, ok := os.LookupEnv(eirini.EnvOutputBuildArtifactsCache)
	if !ok {
		outputBuildArtifactsCache = eirini.RecipeOutputBuildArtifactsCache
	}

	outputMetadataLocation, ok := os.LookupEnv(eirini.EnvOutputMetadataLocation)
	if !ok {
		outputMetadataLocation = eirini.RecipeOutputMetadataLocation
	}

	workspaceDir, ok := os.LookupEnv(eirini.EnvWorkspaceDir)
	if !ok {
		workspaceDir = eirini.RecipeWorkspaceDir
	}

	responder := recipe.NewResponder(stagingGUID, completionCallback, eiriniAddress)

	executor := &recipe.PacksExecutor{
		BuildDir:                  workspaceDir,
		BuildpacksDir:             buildpacksDir,
		OutputDropletLocation:     outputDropletLocation,
		OutputBuildArtifactsCache: outputBuildArtifactsCache,
		OutputMetadataLocation:    outputMetadataLocation,
	}

	err := executor.ExecuteRecipe()
	if err != nil {
		responder.RespondWithFailure(err)
		os.Exit(1)
	}

	fmt.Println("Recipe Execution completed")
}
