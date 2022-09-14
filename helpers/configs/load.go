package configs

import (
	"fmt"
	"path"
	"strings"

	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func LoadConfig[TGlobal any, TLocal any, TGroup any](
	name string,
	cmd *cobra.Command,
	groupFn func(*TGlobal, *TLocal) *TGroup,
) (*TGroup, error) {
	global, err := loadGlobalConfig[TGlobal](cmd)
	if err != nil {
		return nil, err
	}
	override, err := loadServiceConfig[TGlobal](name, cmd)
	if err != nil {
		return nil, err
	}
	globalOverwritten := new(TGlobal)
	mergo.Merge(globalOverwritten, override)
	mergo.Merge(globalOverwritten, global)

	local, err := loadServiceConfig[TLocal](name, cmd)
	if err != nil {
		return nil, err
	}
	return groupFn(globalOverwritten, local), nil
}

func loadGlobalConfig[T any](cmd *cobra.Command) (*T, error) {
	return createConfig[T](cmd, path.Join(".", "config.json"))
}

func loadServiceConfig[T any](name string, cmd *cobra.Command) (*T, error) {
	return createConfig[T](cmd, path.Join(".", name, "config.json"))
}

func createConfig[T any](cmd *cobra.Command, filename string) (*T, error) {
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	viper.SetEnvPrefix("irbe")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	configFile := viper.GetString("config")
	if configFile == "" {
		configFile = filename
		if mode := viper.GetString("mode"); mode != "" {
			configFile = fmt.Sprintf("%s.%s.json", strings.TrimRight(filename, ".json"), mode)
		}
	}

	viper.SetConfigFile(configFile)

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return populateConfig(new(T))
}
