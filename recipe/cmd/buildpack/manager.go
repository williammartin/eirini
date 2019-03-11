package buildpack

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"code.cloudfoundry.org/eirini/recipe"
)

const configFileName = "config.json"

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
		if err := b.install(buildpack); err != nil {
			return err
		}
	}

	return b.writeBuildpackJson(buildpacks)
}

func (b *BuildpackManager) install(buildpack recipe.Buildpack) (err error) {

	var bytes []byte
	bytes, err = recipe.OpenBuildpackUrl(&buildpack, b.client)
	if err != nil {
		return err
	}

	var tempDirName string
	tempDirName, err = ioutil.TempDir("", "buildpacks")
	if err != nil {
		return err
	}

	fileName := filepath.Join(tempDirName, fmt.Sprintf("buildback-%d-.zip", time.Now().Nanosecond()))
	defer func() {
		err = os.Remove(fileName)
	}()

	err = ioutil.WriteFile(fileName, bytes, 0777)
	if err != nil {
		return err
	}

	buildpackName := fmt.Sprintf("%x", md5.Sum([]byte(buildpack.Name)))
	buildpackPath := filepath.Join(b.buildpackDir, buildpackName)
	err = os.MkdirAll(buildpackPath, 0777)
	if err != nil {
		return err
	}

	err = b.unzipper.Extract(fileName, buildpackPath)
	if err != nil {
		return err
	}

	return err
}

func (b *BuildpackManager) writeBuildpackJson(buildpacks []recipe.Buildpack) error {
	bytes, err := json.Marshal(buildpacks)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(b.buildpackDir, configFileName), bytes, 0644)
	if err != nil {
		return err
	}

	return nil
}
