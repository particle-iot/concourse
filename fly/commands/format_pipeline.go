package commands

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/fly/commands/internal/displayhelpers"
	yamlpatch "github.com/krishicks/yaml-patch"
)

type FormatPipelineCommand struct {
	Config atc.PathFlag `short:"c" long:"config" required:"true" description:"Pipeline configuration file"`
	Write  bool         `short:"w" long:"write" description:"Do not print to stdout; overwrite the file in place"`
}

func (command *FormatPipelineCommand) Execute(args []string) error {
	configPath := string(command.Config)
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		displayhelpers.FailWithErrorf("could not read config file", err)
	}

	placeholderWrapper := yamlpatch.NewPlaceholderWrapper("{{", "}}")
	wrappedConfigBytes := placeholderWrapper.Wrap(configBytes)

	var config atc.Config
	err = yaml.Unmarshal(wrappedConfigBytes, &config)
	if err != nil {
		displayhelpers.FailWithErrorf("could not unmarshal config", err)
	}

	formattedBytes, err := yaml.Marshal(config)
	if err != nil {
		displayhelpers.FailWithErrorf("could not marshal config", err)
	}

	unwrappedConfigBytes := placeholderWrapper.Unwrap(formattedBytes)

	if command.Write {
		fi, err := os.Stat(configPath)
		if err != nil {
			displayhelpers.FailWithErrorf("could not stat config file", err)
		}

		err = os.WriteFile(configPath, unwrappedConfigBytes, fi.Mode())
		if err != nil {
			displayhelpers.FailWithErrorf("could not write formatted config to %s", err, command.Config)
		}
	} else {
		_, err = fmt.Fprint(os.Stdout, string(unwrappedConfigBytes))
		if err != nil {
			displayhelpers.FailWithErrorf("could not write formatted config to stdout", err)
		}
	}

	return nil
}
