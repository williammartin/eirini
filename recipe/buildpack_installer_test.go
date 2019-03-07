package recipe_test

import (
	. "code.cloudfoundry.org/eirini/recipe"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"net/http"
)

var _ = Describe("BuildpackInstaller", func() {

	var (
		buildpack Buildpack
		server *ghttp.Server
		responseContent string
		downloadURL string
		buildpackInstaller BuildpackInstaller
		actualBytes []byte
		expectedBytes []byte
		err error
	)

	Context("when a buildpack URL is given", func() {

		BeforeEach(func() {
			buildpackInstaller = BuildpackInstaller{
				Client: http.DefaultClient,
			}
			server = ghttp.NewServer()
		})

		JustBeforeEach(func() {
			expectedBytes = []byte(responseContent)
			actualBytes, err = buildpackInstaller.OpenUrl(&buildpack)
		})

		Context("and it is a valid URL", func() {
			BeforeEach(func() {

				responseContent = "the content"

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/buildpack"),
						ghttp.RespondWith(http.StatusOK, responseContent),
					),
				)
				downloadURL = server.URL() + "/buildpack"

				buildpack = Buildpack{
					Name: "custom",
					Key: "some_key",
					Url: downloadURL,
					SkipDetect: true,
				}
			})


			It("should not fail", func() {
				Expect(err).To(BeNil())
			})

			It("it downloads the buildpack contents", func() {
				Expect(actualBytes).To(Equal(expectedBytes))
			})
		})

		Context("and it is NOT a valid url", func() {

			BeforeEach(func() {
				buildpack = Buildpack{
					Name: "custom",
					Key: "some_key",
					Url: "___terrible::::__url",
					SkipDetect: true,
				}
			})

			It("should fail", func() {
				Expect(err).ToNot(BeNil())
			})
		})

		Context("and it is an unresponsive url", func() {

			BeforeEach(func() {

				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/buildpack"),
						ghttp.RespondWith(http.StatusNotFound, responseContent),
					),
				)

				buildpack = Buildpack{
					Name: "custom",
					Key: "some_key",
					Url: server.URL() + "/buildpack",
					SkipDetect: true,
				}
			})

			It("should fail", func() {
				Expect(err).ToNot(BeNil())
			})
		})
	})
})
