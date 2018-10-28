package main

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config for connection to btc wallet
type S3Conf struct {
	Id     string `yaml:"id"`
	Key    string `yaml:"key"`
	Token  string `yaml:"token"`
	Region string `yaml:"region"`
	Bucket string `yaml:"bucket"`
}

func (s3c S3Conf) GetId() string {
	return s3c.Id
}

func (s3c S3Conf) GetKey() string {
	return s3c.Key
}

func (s3c S3Conf) GetToken() string {
	return s3c.Token
}
func (s3c S3Conf) GetRegion() string {
	return s3c.Region
}
func (s3c S3Conf) GetBucket() string {
	return s3c.Bucket
}

type ConsumerConf struct {
	Addr       string `yaml:"addr"`
	NotifTopic string `yaml:"notif_topic"`
	ClientId   string `yaml:"client_id"`
}

func (pc ConsumerConf) GetAddr() string {
	return pc.Addr
}

func (pc ConsumerConf) GetTopic() string {
	return pc.NotifTopic
}

func (pc ConsumerConf) GetClientId() string {
	return pc.ClientId
}

// Struct for all configs
type Config struct {
	S3conf *S3Conf       `yaml:"s3"`
	ConsumerConf *ConsumerConf `yaml:"consumer"`
}

// Config global object
var config *Config = &Config{}

// Init all configs from config.yml
func InitConfig() {
	absPath, _ := filepath.Abs("./config.yml")

	log.WithFields(log.Fields{
		"path": absPath,
	}).Info("Start inition config")

	yamlFile, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.WithFields(log.Fields{
			"path": absPath,
			"err":  err,
		}).Fatal("Initialization config error")
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		log.WithFields(log.Fields{
			"path": absPath,
			"err":  err,
		}).Fatal("Initialization config error")
	}

	log.WithFields(log.Fields{
		"success": true,
	}).Info("Inition config")
}

func GetConsumerConfig() *ConsumerConf {
	return config.ConsumerConf
}

func GetS3Conf() *S3Conf {
	return config.S3conf
}
