package recipe_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/eirini/eirinifakes"
	. "code.cloudfoundry.org/eirini/recipe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("PackageInstaller", func() {

	var (
		err         error
		appID       string
		targetDir   string
		zipFilePath string
		installer   Installer
		server      *ghttp.Server
		extractor   *eirinifakes.FakeExtractor
	)

	BeforeEach(func() {
		appID = "guid"
		targetDir = "testdata"
		zipFilePath = filepath.Join(targetDir, appID) + ".zip"
		extractor = new(eirinifakes.FakeExtractor)
		server = ghttp.NewServer()
		serverURL, urlErr := url.Parse(server.URL())
		Expect(urlErr).ToNot(HaveOccurred())
		installer = &PackageInstaller{ServerURL: serverURL, Client: &http.Client{}, Extractor: extractor}
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/v2/apps/"+appID+"/download"),
				ghttp.RespondWith(http.StatusOK, "appbits"),
			),
		)
	})

	JustBeforeEach(func() {
		err = installer.Install(appID, targetDir)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("Install", func() {
		assertNoInteractionsWithExtractor := func() {
			It("shoud not interact with the extractor", func() {
				Expect(extractor.Invocations()).To(BeEmpty())
			})
		}

		assertExtractorInteractions := func() {
			It("should use the extractor to extract the zip file", func() {
				src, actualTargetDir := extractor.ExtractArgsForCall(0)
				Expect(extractor.ExtractCallCount()).To(Equal(1))
				Expect(src).To(Equal(zipFilePath))
				Expect(actualTargetDir).To(Equal(targetDir))
			})
		}

		Context("When an empty appID is provided", func() {
			BeforeEach(func() {
				appID = ""
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("empty appID provided")))
			})
			assertNoInteractionsWithExtractor()
		})

		Context("When an empty targetDir is provided", func() {
			BeforeEach(func() {
				targetDir = ""
			})

			It("should return an error", func() {
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("empty targetDir provided")))
			})
			assertNoInteractionsWithExtractor()
		})

		FContext("When package is installed successfully", func() {
			AfterEach(func() {
				osError := os.Remove(zipFilePath)
				Expect(osError).ToNot(HaveOccurred())
			})

			It("writes the downloaded content to the given file", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(zipFilePath).Should(BeAnExistingFile())

				file, readErr := ioutil.ReadFile(filepath.Clean(zipFilePath))
				Expect(readErr).ToNot(HaveOccurred())
				Expect(string(file)).To(Equal("appbits"))

			})

			assertExtractorInteractions()
		})

		Context("When the download fails", func() {
			Context("When the cfclient returns an error", func() {
				BeforeEach(func() {
				})

				It("should error with an corresponding error message", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to perform request")))
				})

				assertNoInteractionsWithExtractor()
			})

			Context("When the cfclient does not return OK HTTP status", func() {
				BeforeEach(func() {
				})

				It("should return an meaningful err message", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("Download failed. Status Code")))
				})
			})

			Context("When the extractor returns an error", func() {
				var expectedErrorMessage string

				BeforeEach(func() {
					expectedErrorMessage = "failed to extract zip"
					extractor.ExtractReturns(errors.New(expectedErrorMessage))
				})

				AfterEach(func() {
					osError := os.Remove(zipFilePath)
					Expect(osError).ToNot(HaveOccurred())
				})

				assertExtractorInteractions()

				It("should return an error", func() {
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring(expectedErrorMessage)))
				})
			})

		})
	})
})
