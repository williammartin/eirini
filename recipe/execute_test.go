package recipe_test

import (
	"io/ioutil"
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	bap "code.cloudfoundry.org/buildpackapplifecycle"
	. "code.cloudfoundry.org/eirini/recipe"
	"code.cloudfoundry.org/eirini/recipe/recipefakes"
)

var _ = Describe("PacksExecutor", func() {

	var (
		executor       Executor
		commander      *recipefakes.FakeCommander
		resultModifier *recipefakes.FakeStagingResultModifier
		tmpfile        *os.File
		resultContents string
	)

	createTmpFile := func() {
		var err error

		tmpfile, err = ioutil.TempFile("", "metadata_result")
		Expect(err).ToNot(HaveOccurred())

		_, err = tmpfile.Write([]byte(resultContents))
		Expect(err).ToNot(HaveOccurred())

		err = tmpfile.Close()
		Expect(err).ToNot(HaveOccurred())
	}

	BeforeEach(func() {
		commander = new(recipefakes.FakeCommander)
		resultModifier = new(recipefakes.FakeStagingResultModifier)

		resultModifier.ModifyStub = func(result bap.StagingResult) (bap.StagingResult, error) {
			return result, nil
		}

		resultContents = `{"lifecycle_type":"no-type", "execution_metadata":"data"}`
	})

	JustBeforeEach(func() {
		createTmpFile()
		packsConf := PacksBuilderConf{
			PacksBuilderPath:          "/packs/builder",
			BuildDir:                  "some-dir",
			BuildpacksDir:             "/var/lib/buildpacks",
			OutputDropletLocation:     "/out/droplet.tgz",
			OutputBuildArtifactsCache: "/cache/cache.tgz",
			OutputMetadataLocation:    tmpfile.Name(),
		}

		executor = &PacksExecutor{
			Conf:      packsConf,
			Commander: commander,
		}

	})

	AfterEach(func() {
		Expect(os.Remove(tmpfile.Name())).To(Succeed())
	})

	Context("When executing a recipe", func() {

		var (
			err    error
			server *ghttp.Server
		)

		BeforeEach(func() {
			server = ghttp.NewServer()
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
			err = executor.ExecuteRecipe()
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			server.Close()
		})

		It("should not return an error", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("should run the packs builder", func() {
			Expect(commander.ExecCallCount()).To(Equal(1))

			cmd, args := commander.ExecArgsForCall(0)
			Expect(cmd).To(Equal("/packs/builder"))
			Expect(args).To(ConsistOf(
				"-buildDir", "some-dir",
				"-buildpacksDir", "/var/lib/buildpacks",
				"-outputDroplet", "/out/droplet.tgz",
				"-outputBuildArtifactsCache", "/cache/cache.tgz",
				"-outputMetadata", tmpfile.Name(),
			))
		})

	})

})
