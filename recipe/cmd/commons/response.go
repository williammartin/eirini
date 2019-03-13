package commons

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"code.cloudfoundry.org/bbs/models"
	bap "code.cloudfoundry.org/buildpackapplifecycle"
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/runtimeschema/cc_messages"
	"github.com/pkg/errors"
)

func RespondWithFailure(failure error) {
	recipeConfig := RecipeConfig()
	cbResponse := createFailureResponse(failure, recipeConfig.StagingGUID, recipeConfig.CompletionCallback)

	if completeErr := sendCompleteResponse(recipeConfig.EiriniAddr, cbResponse); completeErr != nil {
		fmt.Println("Error processsing completion callback:", completeErr.Error())
	}
}

func RespondWithSuccess() error {
	recipeConfig := RecipeConfig()

	cbResponse, err := createSuccessResponse(RecipeConfig(), OutputMetadataLocation, BuildpackJson())
	if err != nil {
		return err
	}

	err = sendCompleteResponse(recipeConfig.EiriniAddr, cbResponse)
	if err != nil {
		return err
	}

	return nil
}

func createSuccessResponse(recipeConfig recipe.Config, outputMetadataLocation string, buildpackJson string) (*models.TaskCallbackResponse, error) {
	stagingResult, err := getStagingResult(outputMetadataLocation)
	if err != nil {
		return nil, err
	}

	modifier := &recipe.BuildpacksKeyModifier{CCBuildpacksJSON: buildpackJson}
	stagingResult, err = modifier.Modify(stagingResult)
	if err != nil {
		return nil, err
	}

	result, err := json.Marshal(stagingResult)
	if err != nil {
		return nil, err
	}

	annotation := cc_messages.StagingTaskAnnotation{
		CompletionCallback: recipeConfig.CompletionCallback,
	}

	annotationJSON, err := json.Marshal(annotation)
	if err != nil {
		return nil, err
	}

	return &models.TaskCallbackResponse{
		TaskGuid:   recipeConfig.StagingGUID,
		Result:     string(result),
		Failed:     false,
		Annotation: string(annotationJSON),
	}, nil
}

func createFailureResponse(failure error, stagingGUID, completionCallback string) *models.TaskCallbackResponse {
	annotation := cc_messages.StagingTaskAnnotation{
		CompletionCallback: completionCallback,
	}

	annotationJSON, err := json.Marshal(annotation)
	if err != nil {
		panic(err)
	}

	return &models.TaskCallbackResponse{
		TaskGuid:      stagingGUID,
		Failed:        true,
		FailureReason: failure.Error(),
		Annotation:    string(annotationJSON),
	}
}

func getStagingResult(path string) (bap.StagingResult, error) {
	contents, err := ioutil.ReadFile(filepath.Clean(path))
	if err != nil {
		return bap.StagingResult{}, errors.Wrap(err, "failed to read result.json")
	}
	var stagingResult bap.StagingResult
	err = json.Unmarshal(contents, &stagingResult)
	if err != nil {
		return bap.StagingResult{}, err
	}
	return stagingResult, nil
}

func sendCompleteResponse(eiriniAddress string, response *models.TaskCallbackResponse) error {
	responseJSON, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	uri := fmt.Sprintf("%s/stage/%s/completed", eiriniAddress, response.TaskGuid)
	req, err := http.NewRequest("PUT", uri, bytes.NewBuffer(responseJSON))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}

	if resp.StatusCode >= 400 {
		return errors.New("Request not successful")
	}

	return nil
}
