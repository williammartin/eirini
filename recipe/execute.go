package recipe

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"code.cloudfoundry.org/buildpackapplifecycle"
	"code.cloudfoundry.org/buildpackapplifecycle/buildpackrunner"
)

type PacksExecutor struct {
	BuildDir                  string
	BuildpacksDir             string
	OutputDropletLocation     string
	OutputBuildArtifactsCache string
	OutputMetadataLocation    string
}

func (e *PacksExecutor) ExecuteRecipe() error {

	conf, err := e.newBuilderConfig()
	if err != nil {
		return err
	}
	runner := buildpackrunner.New(conf)

	infoFilePath, err := runner.Run()
	fmt.Println(infoFilePath)
	fmt.Println(err.Error())
	fmt.Println("COOOOL")
	return err
}

func (e *PacksExecutor) newBuilderConfig() (*buildpackapplifecycle.LifecycleBuilderConfig, error) {

	flagSet := flag.NewFlagSet("builder", flag.ExitOnError)

	flagSet.String(
		"buildDir",
		e.BuildDir,
		"directory containing raw app bits",
	)

	flagSet.String(
		"outputDroplet",
		e.OutputDropletLocation,
		"file where compressed droplet should be written",
	)

	flagSet.String(
		"outputMetadata",
		e.OutputMetadataLocation,
		"directory in which to write the app metadata",
	)

	flagSet.String(
		e.OutputBuildArtifactsCache,
		"/tmp/output-cache",
		"file where compressed contents of new cached build artifacts should be written",
	)

	flagSet.String(
		"buildpacksDir",
		e.BuildpacksDir,
		"directory containing the buildpacks to try",
	)

	flagSet.String(
		"buildArtifactsCacheDir",
		"/tmp/cache",
		"directory where previous cached build artifacts should be extracted",
	)

	buildpacks, err := reduceJSON(path.Join(e.BuildpacksDir, "config.json"), "name")
	if err != nil {
		return nil, err
	}

	flagSet.String(
		"buildpackOrder",
		strings.Join(buildpacks, ","),
		"comma-separated list of buildpacks, to be tried in order",
	)

	/*
		flagSet.String(
			"buildpacksDownloadDir",
			"/tmp/buildpackdownloads",
			"directory to download buildpacks to",
		)

	*/

	flagSet.Bool(
		"skipDetect",
		len(buildpacks) == 1,
		"skip buildpack detect",
	)

	flagSet.Bool(
		"skipCertVerify",
		false,
		"skip SSL certificate verification",
	)

	return &buildpackapplifecycle.LifecycleBuilderConfig{
		FlagSet:        flagSet,
		ExecutablePath: "/tmp/lifecycle/builder",
	}, nil

}

func reduceJSON(path string, key string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	var list []map[string]string
	if err := json.NewDecoder(f).Decode(&list); err != nil {
		return nil, err
	}

	var out []string
	for _, m := range list {
		out = append(out, m[key])
	}
	return out, nil
}
