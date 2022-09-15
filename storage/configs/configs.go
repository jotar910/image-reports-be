package configs

import (
	configs_helper "image-reports/helpers/configs"

	"github.com/spf13/cobra"
)

type AppConfig struct {
	configs_helper.GlobalConfig
	ServiceConfig
}

func newAppConfig(global *configs_helper.GlobalConfig, local *ServiceConfig) *AppConfig {
	return &AppConfig{*global, *local}
}

type ServiceConfig struct {
	Path string
}

var Config *AppConfig

func Get() *AppConfig {
	return Config
}

func Initialize(name string) (*AppConfig, error) {
	config, err := configs_helper.LoadConfig(name, createCommand(), newAppConfig)
	if err != nil {
		Config = nil
		return nil, err
	}

	Config = config
	return config, nil
}

func createCommand() *cobra.Command {
	command := &cobra.Command{}
	command.PersistentFlags().StringP("mode", "m", "", "the environment mode (e.g for mode dev, config file must be named config.dev.json)")
	return command
}
