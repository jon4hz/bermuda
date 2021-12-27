package config

import (
	_ "embed"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types/filters"
	"github.com/spf13/viper"
)

const DefaultCfgFile = "./bermuda.yml"

type PruneImageConfig struct {
	Until         string   `yaml:"until"`
	ExcludeLabels []string `yaml:"excludeLabels"`
	IncludeLabels []string `yaml:"includeLabels"`
	Active        bool     `yaml:"active"`
	All           bool     `yaml:"all"`
}

type PruneContainerConfig struct {
	Until         string   `yaml:"until"`
	ExcludeLabels []string `yaml:"excludeLabels"`
	IncludeLabels []string `yaml:"includeLabels"`
	Active        bool     `yaml:"active"`
}

type LoggingConfig struct {
	LogFile string `yaml:"logFile"`
	Verbose bool   `yaml:"verbose"`
}

type Config struct {
	Image     *PruneImageConfig     `yaml:"image"`
	Container *PruneContainerConfig `yaml:"container"`
	Logging   *LoggingConfig        `yaml:"logging"`
}

//go:embed default_config.yml
var defaultConfig []byte

var cfg *Config

func LoadConfig(file string) {
	if file == "" {
		file = DefaultCfgFile
	}
	viper.SetConfigFile(file)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		if os.IsNotExist(err) && file == DefaultCfgFile {
			err = deployDefaultConfig()
			if err != nil {
				log.Fatalf("failed to deploy default config: %s", err)
			}
			if err := viper.ReadInConfig(); err != nil {
				log.Fatalf("failed to read config: %s", err)
			}
		} else {
			log.Fatalf("failed to read config: %s", err)
		}
	}
	viper.Unmarshal(&cfg)
}

func Get() *Config {
	return cfg
}

func deployDefaultConfig() error {
	return os.WriteFile(DefaultCfgFile, defaultConfig, 0600)
}

func (p PruneImageConfig) ToFilterArgs() filters.Args {
	f := filters.NewArgs()
	f.Add("dangling", strconv.FormatBool(!p.All))
	if p.Until != "" {
		f.Add("until", p.Until)
	}
	for _, l := range p.ExcludeLabels {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if f.Contains("label!") && f.ExactMatch("label!", parts[1]) {
			continue
		}
		/* if f.Contains("label") && f.ExactMatch("label", parts[1]) {
			continue
		} */
		f.Add("label!", parts[1])
	}
	for _, l := range p.IncludeLabels {
		parts := strings.SplitN(l, "=", 2)
		if len(parts) != 2 {
			continue
		}
		if f.Contains("label!") && f.ExactMatch("label!", parts[1]) {
			continue
		}
		if f.Contains("label") && f.ExactMatch("label", parts[1]) {
			continue
		}
		f.Add("label", parts[1])
	}
	return f
}

func (p PruneContainerConfig) ToFilterArgs() filters.Args {
	f := filters.NewArgs()

	if p.Until != "" {
		f.Add("until", p.Until)
	}
	for _, l := range p.ExcludeLabels {
		if f.Contains("label!") && f.ExactMatch("label!", l) {
			continue
		}
		/* if f.Contains("label") && f.ExactMatch("label", l) {
			continue
		} */
		f.Add("label!", l)
	}

	for _, l := range p.IncludeLabels {
		if f.Contains("label!") && f.ExactMatch("label!", l) {
			continue
		}
		if f.Contains("label") && f.ExactMatch("label", l) {
			continue
		}
		f.Add("label", l)
	}
	return f
}
