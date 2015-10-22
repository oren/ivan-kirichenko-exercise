package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/Sirupsen/logrus"
	app "github.com/seesawlabs/ivan-kirichenko-exercise/application"
	yaml "gopkg.in/yaml.v2"
)

func main() {
	logger := logrus.New() // TODO: configure logger more if needed

	configPath := flag.String("config", "config.yaml", "mandatory path to config file")
	migrate := flag.Bool("migrate", false, "path to directiry with migration scripts. If provided, runs database migrations and exits")

	if configPath == nil {
		logger.Fatal("config file must be provided")
	}

	config, err := readConfig(*configPath)
	if err != nil {
		logger.Fatal(err.Error())
	}
	logger.Infof("read config %+v", config)

	if migrate != nil && *migrate {
		runMigrations(config, logger)
		return
	}

	runApplication(config, logger)
}

func readConfig(path string) (*app.Config, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %s", err.Error())
	}
	config := &app.Config{}
	if err := yaml.Unmarshal(content, config); err != nil {
		return nil, fmt.Errorf("could not parse yaml config file '%s': %s", path, err.Error())
	}

	return config, nil
}

func runApplication(config *app.Config, logger *logrus.Logger) {
	application, err := app.NewApp(config, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}

	logger.Infoln("starting the application...")
	application.Run()
}

func runMigrations(config *app.Config, logger *logrus.Logger) {
	logger.Infof("running migrations in '%s'...", config.DbFile)
	migrator, err := app.NewMigratable(config, logger)
	if err != nil {
		logger.Fatal(err.Error())
	}
	err = migrator.Migrate()
	if err != nil {
		logger.Fatal(err.Error())
	}
}
