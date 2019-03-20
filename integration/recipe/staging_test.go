package recipe_test

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"

	"code.cloudfoundry.org/eirini"

	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = FDescribe("StagingText", func() {

	var (
		err            error
		server         *ghttp.Server
		appbitBytes    []byte
		buildpackBytes []byte
		session        *gexec.Session
	)

	BeforeEach(func() {

		appbitBytes, err = ioutil.ReadFile("testdata/catnip")
		Expect(err).NotTo(HaveOccurred())

		buildpackBytes, err = ioutil.ReadFile("testdata/binary-buildpack-cflinuxfs2-v1.0.31.zip")
		Expect(err).NotTo(HaveOccurred())

		server = ghttp.NewServer()
		server.AppendHandlers(

			// Downloader
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/my-app-bits"),
				ghttp.RespondWith(http.StatusOK, appbitBytes),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/my-buildpack"),
				ghttp.RespondWith(http.StatusOK, buildpackBytes),
			),

			// Uploader
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/my-droplet"),
				ghttp.RespondWith(http.StatusOK, ""),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/cc-success"),
				ghttp.RespondWith(http.StatusOK, ""),
			),
		)

		err = os.Setenv(eirini.EnvCertsPath, tempCertsPath)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err = os.Unsetenv(eirini.EnvCertsPath)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when a droplet needs building...", func() {

		Context("downloads the app and buildpacks", func() {

			JustBeforeEach(func() {
				cmd := exec.Command(binaries.DownloaderPath)
				session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			})

			FIt("runs successfully", func() {
				Expect(session.ExitCode()).To(BeZero())
			})

			It("installs the buildpack json", func() {

			})

			It("installs the binary buildpack", func() {

			})

			It("places the app bits in the workspace", func() {

			})

			Context("creates the droplet", func() {

				JustBeforeEach(func() {
					cmd := exec.Command(binaries.ExecutorPath)
					session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				})

				Context("uploads the droplet", func() {

					JustBeforeEach(func() {
						cmd := exec.Command(binaries.UploaderPath)
						session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
					})

				})
			})
		})
	})
})
