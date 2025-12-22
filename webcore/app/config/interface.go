package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var InstanceViper map[string]*viper.Viper = make(map[string]*viper.Viper)

type Configurable interface {
	SetDefaults() map[string]any
	SetEnvBindings() map[string]string
}

func LoadDefaultConfig[T Configurable](c T, ismodule bool) error {
	return LoadConfig("", c, "config", "yaml", []string{}, ismodule)
}

func LoadDefaultConfigModule[T Configurable](name string, c T) error {
	return LoadConfig(name, c, "config", "yaml", []string{}, true)
}

func LoadConfig[T Configurable](prefix string, c T, file string, ext string, path []string, ismodule bool) error {
	var v *viper.Viper
	replacer := strings.NewReplacer(".", "_")
	name := file + "." + ext

	if InstanceViper[name] == nil {
		v = viper.New()

		v.SetConfigName(file)
		v.SetConfigType(ext)
		if len(path) == 0 {
			v.AddConfigPath(".")
		} else {
			for _, p := range path {
				v.AddConfigPath(p)
			}
		}

		// Override with environment variables
		v.AutomaticEnv()

		if err := v.ReadInConfig(); err != nil {
			// If config file is not found, use defaults and environment variables
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return err
			}
		}

		log.Printf("Muat konfigurasi dari environment variables")

		// Replace dots with underscores for environment variable keys
		v.SetEnvKeyReplacer(replacer)
		InstanceViper[name] = v
	} else {
		v = InstanceViper[name]
	}

	// Set defaults with priority to environment variables
	setPriorityDefaults(c, v, replacer, prefix, ismodule)

	if err := v.Unmarshal(c); err != nil {
		return err
	}

	return nil
}

func setPriorityDefaults(c Configurable, v *viper.Viper, replacer *strings.Replacer, prefix string, ismodule bool) {
	var modPrefix string
	if prefix != "" {
		if ismodule {
			modPrefix = "module." + prefix + "."
		} else {
			modPrefix = prefix + "."
		}
	}

	// Force binding of specific environment variables
	bindings := c.SetEnvBindings()
	for runtimeKey, envKey := range bindings {
		v.BindEnv(runtimeKey, envKey)
	}

	defaults := c.SetDefaults()

	log.Printf("Scan Values:")
	for _, runtimeKey := range v.AllKeys() {
		runtimeValue := v.Get(runtimeKey)
		envFilekey := replacer.Replace(runtimeKey)

		cut := false
		runtimeKeyCut := runtimeKey
		if prefix != "" && strings.HasPrefix(runtimeKey, modPrefix) {
			runtimeKeyCut = strings.TrimPrefix(runtimeKey, modPrefix)
			cut = true
		}

		if runtimeKey != envFilekey {
			if runtimeValue == nil {
				envFileValue := v.Get(envFilekey)
				if envFileValue != nil {
					log.Printf(" %s = %v -> [%s]", runtimeKey, envFileValue, envFilekey)
					v.SetDefault(runtimeKey, envFileValue)
					if cut {
						log.Printf(" ~%s = %v -> [%s-CUT-PREFIX]", runtimeKeyCut, envFileValue, envFilekey)
						v.SetDefault(runtimeKeyCut, envFileValue)
					}
				} else if defValue, ok := defaults[runtimeKey]; ok {
					log.Printf(" %s = %v -> [DEFAULTS]", runtimeKey, defValue)
					v.SetDefault(runtimeKey, defValue)
					if cut {
						log.Printf(" ~%s = %v -> [DEFAULTS-CUT-PREFIX]", runtimeKeyCut, defValue)
						v.SetDefault(runtimeKeyCut, defValue)
					}
				}
			} else {
				log.Printf(" %s = %v -> [RUNTIME]", runtimeKey, runtimeValue)
				if cut {
					log.Printf(" ~%s = %v -> [RUNTIME-CUT-PREFIX]", runtimeKeyCut, runtimeValue)
					v.SetDefault(runtimeKeyCut, runtimeValue)
				}
			}
		}
	}
}
