package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/blankdots/minimal-kube-app/internal/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config is a parent object for all the different configuration parts
type Config struct {
	Database database.DBConfig
	API      APIConf
	CronJob  CronJobConfig
}

type APIConf struct {
	Host        string
	Port        int
	StaticToken string
}

type CronJobConfig struct {
	APIBase   string   // e.g. https://registry.npmjs.org
	Packages  []string // packages to fetch, e.g. express, lodash
	Interval  time.Duration
}

func App(app string) (*Config, error) {

	// look for config.yaml in root path
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	// logging as JSON by default
	log.SetFormatter(&log.JSONFormatter{})

	log.Infoln("reading config")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Infoln("No config file found, using ENVs only")
		} else {
			log.Errorln("ReadInConfig Error")
		}
	}

	log.SetLevel(log.DebugLevel)

	if viper.IsSet("log.level") {
		stringLevel := viper.GetString("log.level")
		intLevel, err := log.ParseLevel(stringLevel)
		if err != nil {
			log.Infof("Log level: '%s' not supported, setting to 'debug'", stringLevel)
			intLevel = log.DebugLevel
		}
		log.SetLevel(intLevel)
		log.Infof("Setting log level to '%s'", stringLevel)
	}

	c := &Config{}

	switch app {
	case "cronjob":
		log.Debug("Configuring cronjob")

		c.cronJob()

		c.configDatabase()

	case "api":
		log.Debug("Configuring REST API")

		c.configAPI()

		c.configDatabase()

	default:
		return nil, fmt.Errorf("cannot recognize app '%s' to configure", app)
	}

	return c, nil
}

func (c *Config) configAPI() {
	viper.SetDefault("api.host", "0.0.0.0")
	viper.SetDefault("api.port", 5005)
	viper.SetDefault("api.token", "test")

	api := APIConf{}

	api.Host = viper.GetString("api.host")
	api.Port = viper.GetInt("api.port")

	// configure a static token
	api.StaticToken = viper.GetString("api.token")

	c.API = api
}

func (c *Config) cronJob() {

	// npm Registry API (no auth, package dependencies)
	viper.SetDefault("cronjob.apibase", "https://registry.npmjs.org")
	viper.SetDefault("cronjob.packages", []string{"express", "lodash", "react"})

	cj := CronJobConfig{}

	cj.APIBase = viper.GetString("cronjob.apibase")
	cj.Packages = viper.GetStringSlice("cronjob.packages")
	// Env CRONJOB_PACKAGES="express,lodash,react" often comes as a single slice element
	if len(cj.Packages) == 1 && strings.Contains(cj.Packages[0], ",") {
		cj.Packages = strings.Split(cj.Packages[0], ",")
		for i := range cj.Packages {
			cj.Packages[i] = strings.TrimSpace(cj.Packages[i])
		}
	}
	if len(cj.Packages) == 0 {
		if s := viper.GetString("cronjob.packages"); s != "" {
			cj.Packages = strings.Split(s, ",")
			for i := range cj.Packages {
				cj.Packages[i] = strings.TrimSpace(cj.Packages[i])
			}
		} else {
			cj.Packages = []string{"express"}
		}
	}

	c.CronJob = cj

}

func (c *Config) configDatabase() {
	db := database.DBConfig{}

	// splitting into a more traditional way of setting the environment variables
	host := viper.GetString("db.host")
	port := viper.GetInt("db.port")
	user := viper.GetString("db.user")
	password := viper.GetString("db.password")
	database := viper.GetString("db.database")

	// postgres://username:password@localhost:5432/database_name
	db.URL = "postgres://" + user + ":" + password + "@" + host + ":" + fmt.Sprint(port) + "/" + database

	c.Database = db
}
