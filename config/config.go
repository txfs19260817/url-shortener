package config

import (
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"github.com/txfs19260817/url-shortener/database"
	"github.com/txfs19260817/url-shortener/service"
)

type Config struct {
	Server  Server                 `yaml:"server"`
	Mongodb database.MongoDBConfig `yaml:"mongodb"`
	Service service.ServicesConfig `yaml:"service"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func (cfg *Config) ReadConfig(fullpath string) error {
	dir, file := filepath.Split(fullpath)
	ext := filepath.Ext(file)
	v := viper.New()
	v.SetConfigName(strings.TrimSuffix(file, ext)) // name of config file (without extension)
	v.SetConfigType(ext[1:])                       // REQUIRED if the config file does not have the extension in the name
	v.AddConfigPath(dir)                           // path to look for the config file in
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	if err := v.Unmarshal(cfg, func(d *mapstructure.DecoderConfig) { d.TagName = "yaml" }); err != nil {
		return err
	}
	return nil
}
