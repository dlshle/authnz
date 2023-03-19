package yaml

import (
	"go.uber.org/config"
)

// config must be a pointer
func LoadConfig(configFile string, config interface{}) (err error) {
	yaml, err := LoadConfigAsValue(configFile)
	if err != nil {
		return
	}
	err = yaml.Populate(config)
	return
}

func LoadConfigAsMap(configFile string) (res map[string]interface{}, err error) {
	yaml, err := LoadConfigAsValue(configFile)
	if err != nil {
		return
	}
	err = yaml.Populate(&res)
	return
}

func LoadConfigAsValue(configFile string) (val config.Value, err error) {
	yaml, err := config.NewYAML(config.File(configFile))
	if err != nil {
		return
	}
	return yaml.Get(""), nil
}
