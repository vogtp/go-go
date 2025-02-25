package cfg

import (
	"fmt"
	"strings"
	"time"

	"log/slog"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const appName = "ragctl"

const ()

func init() {
}

// SetConfigFileName sets the config file name
func SetConfigFileName(n string) {
	viper.Set(CfgFile, n)
}

// Parse parses the config
func Parse() {
	if pflag.Parsed() {
		slog.Debug("pflags already parsed")
		return
	}

	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		slog.Error("cannot bin flags", "error", err)
	}
	capName := strings.ToUpper(appName)
	viper.SetConfigType("yaml")
	viper.SetConfigName(viper.GetString(CfgFile))
	viper.AddConfigPath(".")
	viper.AddConfigPath(fmt.Sprintf("/etc/%s/", appName))
	viper.AddConfigPath(fmt.Sprintf("/%s/", appName))
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s/", appName))
	viper.AddConfigPath(fmt.Sprintf("$%s_HOME/", capName))
	viper.AddConfigPath(fmt.Sprintf("$%s_ROOT/", capName))

	viper.SetEnvPrefix(capName)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	processConfigFile()
}

func processConfigFile() {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			slog.Debug("config not found", "error", err)
		} else {
			slog.Warn("config error", "error", err)
		}
	}
	if viper.GetBool(CfgSave) {
		//autosave the current config
		go func() {
			time.Sleep(10 * time.Second)
			for {
				if !viper.GetBool(CfgSave) {
					slog.Debug("Requested not to save config")
					return
				}
				// should we write it regular and overwrite?

				if err := Save(); err != nil {
					slog.Warn("Could not write config", "error", err)
				}
				time.Sleep(time.Hour)
			}
		}()
	}
}

// Save the config to file
func Save() error {
	slog.Info("Writing config")
	return viper.WriteConfigAs(viper.GetString(CfgFile))
}
