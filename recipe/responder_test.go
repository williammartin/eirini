package recipe_test

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/runtimeschema/cc_messages"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Responder", func() {
	Context("when responding to cc-uploader", func() {
		var (
			err       error
			server    *ghttp.Server
			responder recipe.Responder
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
		})

		JustBeforeEach(func() {
			cfg := recipe.Config{
				AppID:              "app-id",
				StagingGUID:        "staging-guid",
				CompletionCallback: "completion-call-me-back",
				EiriniAddr:         server.URL(),
				DropletUploadURL:   "droplet.eu/upload",
				PackageDownloadURL: server.URL() + "app-id",
			}

			responder = recipe.NewResponder(cfg)
		})

		AfterEach(func() {
			server.Close()
		})

		Context("when there is an error", func() {
			BeforeEach(func() {
				server.RouteToHandler("PUT", "/stage/staging-guid/completed",
					ghttp.VerifyJSON(`{
						"task_guid": "staging-guid",
						"failed": true,
						"failure_reason": "sploded!",
						"result": "",
						"annotation": "{\"lifecycle\":\"\",\"completion_callback\":\"completion-call-me-back\"}",
						"created_at": 0
					}`),
				)
			})

			It("should respond with failure", func() {
				err = errors.New("sploded!")
				responder.RespondWithFailure(err)
			})
		})

		Context("when the response is success", func() {

			var (
				resp            *models.TaskCallbackResponse
				resultsFilePath string
				resultContents  string
				buildpacks      []byte
			)

			Context("when preparing the response results", func() {
				var err error

				Context("when the results file is missing", func() {
					JustBeforeEach(func() {
						resp, err = responder.PrepareSuccessResponse(resultsFilePath, string(buildpacks))
					})

					It("should error with missing file msg", func() {
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("failed to read result.json"))
					})
				})

				Context("when the results json file is invalid", func() {
					JustBeforeEach(func() {
						resultsFilePath = resultsFile(resultContents)

						buildpack := cc_messages.Buildpack{}

						buildpacks, err = json.Marshal([]cc_messages.Buildpack{buildpack})
						Expect(err).NotTo(HaveOccurred())

						resp, err = responder.PrepareSuccessResponse(resultsFilePath, string(buildpacks))
					})

					It("should error when unmarhsaling the content", func() {
						Expect(err).To(HaveOccurred())
						Expect(err.Error()).To(ContainSubstring("unexpected end of JSON input"))
					})
				})

			})

			Context("when response preparation is successful", func() {

				BeforeEach(func() {
					resultContents = `{"lifecycle_type":"no-type", "execution_metadata":"data"}`

					server.RouteToHandler("PUT", "/stage/staging-guid/completed",
						ghttp.VerifyJSON(`{
						"task_guid": "staging-guid",
						"failed": false,
						"failure_reason": "",
						"result": "{\"lifecycle_metadata\":{\"detected_buildpack\":\"\",\"buildpacks\":null},\"process_types\":null,\"execution_metadata\":\"data\",\"lifecycle_type\":\"no-type\"}",
						"annotation": "{\"lifecycle\":\"\",\"completion_callback\":\"completion-call-me-back\"}",
						"created_at": 0
					}`),
					)
				})

				JustBeforeEach(func() {
					resultsFilePath = resultsFile(resultContents)

					buildpack := cc_messages.Buildpack{}

					buildpacks, err = json.Marshal([]cc_messages.Buildpack{buildpack})
					Expect(err).NotTo(HaveOccurred())
					resp, err = responder.PrepareSuccessResponse(resultsFilePath, string(buildpacks))
					Expect(err).NotTo(HaveOccurred())
				})

				AfterEach(func() {
					Expect(os.Remove(resultsFilePath)).To(Succeed())
				})

				It("should respond with failure", func() {
					err = responder.RespondWithSuccess(resp)
					Expect(err).NotTo(HaveOccurred())
				})
			})
		})

	})
})

func resultsFile(content string) string {
	var err error

	tmpfile, err := ioutil.TempFile("", "metadata_result")
	Expect(err).ToNot(HaveOccurred())

	_, err = tmpfile.Write([]byte(content))
	Expect(err).ToNot(HaveOccurred())

	err = tmpfile.Close()
	Expect(err).ToNot(HaveOccurred())

	return tmpfile.Name()
}
