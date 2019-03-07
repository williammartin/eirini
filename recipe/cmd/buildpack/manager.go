package buildpack

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/eirini/recipe"
)

type BuildpackManager struct {
	client       *http.Client
	buildpackDir string
	unzipper     recipe.Unzipper
}

func New(client *http.Client, buildpackDir string) *BuildpackManager {
	return &BuildpackManager{
		client:       client,
		buildpackDir: buildpackDir,
	}
}

func (b *BuildpackManager) Install(buildpacks []recipe.Buildpack) error {
	for _, buildpack := range buildpacks {
		bytes, err := recipe.OpenBuildpackUrl(&buildpack, b.client)
		if err != nil {
			return err
		}

		tmp, err := ioutil.TempFile("", "buildpack.zip")
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(tmp.Name(), bytes, os.ModePerm)
		if err != nil {
			return err
		}

		buildpackName := fmt.Sprintf("%x", md5.Sum([]byte(buildpack.Name)))
		buildpackPath := filepath.Join(b.buildpackDir, buildpackName)
		err = os.MkdirAll(buildpackPath, os.ModeDir)
		if err != nil {
			return err
		}

		//unpack-zip
		err = b.unzipper.Extract(tmp.Name(), buildpackPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *BuildpackManager) WriteBuildpackJson(buildpacks []recipe.Buildpack) error {
	return nil
}
