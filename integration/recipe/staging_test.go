package recipe_test

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/urljoiner"

	"code.cloudfoundry.org/eirini"
	"code.cloudfoundry.org/tlsconfig"
	"github.com/onsi/gomega/gexec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = FDescribe("StagingText", func() {

	const (
		stagingGUID        = "5b00de6b-d8f4-476b-b070-303367b46cef"
		completionCallback = "url://some_endpoint"
	)

	var (
		err            error
		server         *ghttp.Server
		appbitBytes    []byte
		buildpackBytes []byte
		session        *gexec.Session
		buildpacks     []recipe.Buildpack
		buildpacksDir  string
		workspaceDir   string
	)

	BeforeEach(func() {

		workspaceDir, err = ioutil.TempDir("", "workspace")
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv(eirini.EnvWorkspaceDir, workspaceDir)
		Expect(err).NotTo(HaveOccurred())

		buildpacksDir, err = ioutil.TempDir("", "buildpacks")
		Expect(err).NotTo(HaveOccurred())
		err = os.Setenv(eirini.EnvBuildpacksDir, buildpacksDir)
		Expect(err).NotTo(HaveOccurred())

		err = os.Setenv(eirini.EnvStagingGUID, stagingGUID)
		Expect(err).NotTo(HaveOccurred())

		err = os.Setenv(eirini.EnvCompletionCallback, completionCallback)
		Expect(err).NotTo(HaveOccurred())

		appbitBytes, err = ioutil.ReadFile("testdata/catnip.zip")
		Expect(err).NotTo(HaveOccurred())

		buildpackBytes, err = ioutil.ReadFile("testdata/binary-buildpack-cflinuxfs2-v1.0.31.zip")
		Expect(err).NotTo(HaveOccurred())

		certsPath, err := filepath.Abs("testdata/certs")
		Expect(err).NotTo(HaveOccurred())

		certPath := filepath.Join(certsPath, "cc-server-crt")
		keyPath := filepath.Join(certsPath, "cc-server-crt-key")
		caCertPath := filepath.Join(certsPath, "internal-ca-cert")

		tlsConfig, err := tlsconfig.Build(
			tlsconfig.WithInternalServiceDefaults(),
			tlsconfig.WithIdentityFromFile(certPath, keyPath),
		).Server(
			tlsconfig.WithClientAuthenticationFromFile(caCertPath),
		)
		Expect(err).NotTo(HaveOccurred())

		server = ghttp.NewUnstartedServer()
		server.HTTPTestServer.TLS = tlsConfig

		server.AppendHandlers(

			// Downloader
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/my-buildpack"),
				ghttp.RespondWith(http.StatusOK, buildpackBytes),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/my-app-bits"),
				ghttp.RespondWith(http.StatusOK, appbitBytes),
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

		server.Start()

		err = os.Setenv(eirini.EnvDownloadURL, urljoiner.Join(server.URL(), "my-app-bits"))
		Expect(err).ToNot(HaveOccurred())

		buildpacks = []recipe.Buildpack{
			{
				Name: "binary_buildpack",
				Key:  "binary_buildpack",
				URL:  urljoiner.Join(server.URL(), "/my-buildpack"),
			},
		}

		buildpackJSON, err := json.Marshal(buildpacks)
		Expect(err).ToNot(HaveOccurred())

		err = os.Setenv(eirini.EnvBuildpacks, string(buildpackJSON))
		Expect(err).ToNot(HaveOccurred())

		err = os.Setenv(eirini.EnvCertsPath, certsPath)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err = os.RemoveAll(buildpacksDir)
		Expect(err).ToNot(HaveOccurred())

		err = os.Unsetenv(eirini.EnvCertsPath)
		Expect(err).ToNot(HaveOccurred())
		err = os.Unsetenv(eirini.EnvStagingGUID)
		Expect(err).NotTo(HaveOccurred())
		err = os.Unsetenv(eirini.EnvCompletionCallback)
		Expect(err).NotTo(HaveOccurred())
		err = os.Unsetenv(eirini.EnvBuildpacksDir)
		Expect(err).NotTo(HaveOccurred())
		err = os.Unsetenv(eirini.EnvBuildpacks)
		Expect(err).ToNot(HaveOccurred())
		err = os.Unsetenv(eirini.EnvDownloadURL)
		Expect(err).ToNot(HaveOccurred())
		err = os.Unsetenv(eirini.EnvWorkspaceDir)
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when a droplet needs building...", func() {

		Context("downloads the app and buildpacks", func() {

			JustBeforeEach(func() {
				cmd := exec.Command(binaries.DownloaderPath)
				session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
				Eventually(session).Should(gexec.Exit())
			})

			It("runs successfully", func() {
				Expect(session.ExitCode()).To(BeZero())
			})

			It("installs the buildpack json", func() {

				expectedFile := filepath.Join(buildpacksDir, "config.json")
				Expect(expectedFile).To(BeARegularFile())

				actualBytes, err := ioutil.ReadFile(expectedFile)
				Expect(err).ToNot(HaveOccurred())

				var actualBuildpacks []recipe.Buildpack
				err = json.Unmarshal(actualBytes, &actualBuildpacks)
				Expect(err).ToNot(HaveOccurred())

				Expect(actualBuildpacks).To(Equal(buildpacks))
			})

			It("installs the binary buildpack", func() {
				md5Hash := fmt.Sprintf("%x", md5.Sum([]byte("binary_buildpack")))
				expectedBuildpackPath := path.Join(buildpacksDir, md5Hash)
				Expect(expectedBuildpackPath).To(BeADirectory())
			})

			It("places the app bits in the workspace", func() {
				actualBytes, err := ioutil.ReadFile(path.Join(workspaceDir, "catnip"))
				Expect(err).NotTo(HaveOccurred())
				expectedBytes, err := ioutil.ReadFile("testdata/catnip")
				Expect(err).NotTo(HaveOccurred())
				Expect(actualBytes).To(Equal(expectedBytes))
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
