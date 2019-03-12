package main

import (
	"fmt"
	"os"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/eirini/recipe"
)


func main() {
	dropletUploadURL := os.Getenv(eirini.EnvDropletUploadURL)

	uploader := &recipe.DropletUploader{
		HTTPClient: createUploaderHTTPClient(),
	}

	err := uploader.Upload("/out/droplet.tgz", dropletUploadURL)
	if err != nil {
		fmt.Println("Error while executing staging uploader:", err.Error())
		os.Exit(1)
	}

	// TODO
	/*
	cbResponse, err := createSuccessResponse(recipeConf)
	if err != nil {
		return err
	}

	sendCompleteResponse(recipeConf.EiriniAddr, cbResponse)
	*/

	fmt.Println("Uploading completed")
}