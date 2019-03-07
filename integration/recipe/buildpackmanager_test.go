package recipe_test

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"code.cloudfoundry.org/eirini/recipe"
	. "code.cloudfoundry.org/eirini/recipe/cmd/buildpack"
)

var _ = Describe("Buildpackmanager", func() {

	var (
		client           *http.Client
		buildpackDir     string
		buildpackManager *BuildpackManager
		buildpacks       []recipe.Buildpack
		server           *ghttp.Server
		responseContent  []byte
		err              error
	)

	BeforeEach(func() {
		client = http.DefaultClient

		buildpackDir, err = ioutil.TempDir("", "buildpacks")
		Expect(err).ToNot(HaveOccurred())

		buildpackManager = New(client, buildpackDir)

		responseContent, err = makeZippedPackage()
		Expect(err).ToNot(HaveOccurred())

		server = ghttp.NewServer()
		server.AppendHandlers(
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/my-buildpack"),
				ghttp.RespondWith(http.StatusOK, responseContent),
			),
			ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/your-buildpack"),
				ghttp.RespondWith(http.StatusOK, responseContent),
			),
		)

		buildpacks = []recipe.Buildpack{
			{
				Name: "my_buildpack",
				Key:  "my-key",
				Url:  fmt.Sprintf("%s/my-buildpack", server.URL()),
			},
			{
				Name: "your_buildpack",
				Key:  "your-key",
				Url:  fmt.Sprintf("%s/your-buildpack", server.URL()),
			},
		}
	})

	Context("When a list of Buildpacks is provided", func() {
		Context("and the buildpacks need be installed", func() {

			JustBeforeEach(func() {
				err = buildpackManager.Install(buildpacks)
			})

			It("should not fail", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should download all buildpacks to the given directory", func() {
				myMd5Dir := fmt.Sprintf("%x", md5.Sum([]byte("my_buildpack")))
				yourMd5Dir := fmt.Sprintf("%x", md5.Sum([]byte("your_buildpack")))
				Expect(filepath.Join(buildpackDir, myMd5Dir)).To(BeADirectory())
				Expect(filepath.Join(buildpackDir, yourMd5Dir)).To(BeADirectory())
			})

			//TODO Write Test that zip was extracted
		})
	})
})

func makeZippedPackage() ([]byte, error) {
	buf := bytes.Buffer{}
	w := zip.NewWriter(&buf)

	//TODO: Add Content to Zip

	err := w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
