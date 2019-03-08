package recipe_test

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/json"
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

	Context("When a list of Buildpacks needs be installed", func() {
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

		It("should write a config.json file in the correct location", func() {
			Expect(filepath.Join(buildpackDir, "config.json")).To(BeAnExistingFile())
		})

		It("marshals the provided buildpacks to the config.json", func() {
			var actualBytes []byte
			actualBytes, err = ioutil.ReadFile(filepath.Join(buildpackDir, "config.json"))
			Expect(err).ToNot(HaveOccurred())

			var actualBuildpacks []recipe.Buildpack
			err = json.Unmarshal(actualBytes, &actualBuildpacks)
			Expect(err).ToNot(HaveOccurred())

			Expect(actualBuildpacks).To(Equal(buildpacks))
		})
	})
})

func makeZippedPackage() ([]byte, error) {
	buf := bytes.Buffer{}
	w := zip.NewWriter(&buf)

	err := w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
