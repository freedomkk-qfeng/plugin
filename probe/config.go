package main

import (
	"bytes"
	"log"

	"github.com/spf13/viper"
	"github.com/toolkits/pkg/file"
)

const Version = "0.1.0"

type ConfYaml struct {
	Probe probeSection `yaml:"probe"`
	Ping  []string     `yaml:"ping"`
	Url   []string     `yaml:"url"`
}

type probeSection struct {
	Timeout int64             `yaml:"timeout"`
	Limit   int64             `yaml:"limit"`
	Headers map[string]string `yaml:"headers"`
}

func Parse(conf string) (config ConfYaml) {
	bs, err := file.ReadBytes(conf)
	if err != nil {
		log.Fatalf("cannot read config.yml: %v", err)
		return
	}

	viper.SetConfigType("yaml")
	err = viper.ReadConfig(bytes.NewBuffer(bs))
	if err != nil {
		log.Fatalf("cannot read config.yml: %v", err)
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unmarshal %v", err)
		return
	}

	return
}
