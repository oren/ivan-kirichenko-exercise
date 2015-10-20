package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	core "github.com/seesawlabs/ivan-kirichenko-exercise"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	logger := logrus.New() // TODO: configure logger more if needed

	configPath := flag.String("config", "config.yaml", "mandatory path to config file")
	if configPath == nil {
		logger.Fatal("config file must be provided")
	}

	config, err := readConfig(*configPath)
	if err != nil {
		logger.Fatal(err.Error())
	}

	app, err := core.NewApp(config, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	app.Run()
}

func readConfig(path string) (*core.Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %s", err.Error())
	}
	config := &core.Config{}
	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, fmt.Errorf("could not parse yaml config file '%s': %s", path, err.Error())
	}

	return config, nil
}
