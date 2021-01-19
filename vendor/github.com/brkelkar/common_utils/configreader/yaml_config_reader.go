package configreader

import (
	"errors"
	"os"
	"strconv"

	gc "github.com/brkelkar/common_utils/gcsbucketclient"
	"github.com/brkelkar/common_utils/logger"
	"gopkg.in/yaml.v2"
)

//Config provides option get Config variables from .yml file or
//from enviroment variable
//Incase of both eviroment variables overwrites .yml file variables
type Config struct {
	Server struct {
		Port int    `yaml:"port", envconfig:"SERVER_PORT"`
		Host string `yaml:"host", envconfig:"SERVER_HOST"`
	} `yaml:"server"`

	DatabaseName struct {
		AwacsDBName         string `yaml:"acawsDBName", envconfig:"AWACS_DB"`
		AwacsSmartDBName    string `yaml:"acawsSmartDBName", envconfig:"AWACS_SMART_DB"`
		SmartStockistDBName string `yaml:"smartStockistDBName", envconfig:"AWACS_SMART_STOCKIST_DB"`
	} `yaml:"databaseName"`

	DatabaseConfig struct {
		Port     int    `yaml:"port", envconfig:"DB_PORT"`
		Host     string `yaml:"host", envconfig:"DB_HOST"`
		Username string `yaml:"user", envconfig:"DB_USERNAME"`
		Password string `yaml:"pass", envconfig:"DB_PASSWORD"`
	} `yaml:"databaseConfig"`

	GrpcConfig struct {
		Port int    `yaml:"port", envconfig:"GRPC_PORT"`
		Host string `yaml:"host", envconfig:"GCPC_HOST"`
	} `yaml:"grpcConfig"`
}

var m = make(map[string]string)

func processError(err error) {
	logger.Error("Error while reading config", err)
	os.Exit(2)
}

//ReadFile reads given config file
//and loads into config object
func (cfg *Config) ReadFile(fileName string) {
	f, err := os.Open(fileName)
	if err != nil {
		processError(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

//ReadGcsFile reads given config file
//and loads into config object
func (cfg *Config) ReadGcsFile(filePath string) {
	bucketName, fileName, err := gc.GetBucketAndFileName(filePath)
	if err != nil {
		processError(err)
	}
	var gcsObj gc.GcsBucketClient
	gcsClient := gcsObj.InitClient().SetBucketName(bucketName).SetNewReader(fileName)
	if !gcsClient.GetLastStatus() {
		processError(errors.New("Error while reading file from GCS for filepath=" + filePath))
	}

	decoder := yaml.NewDecoder(gcsClient.GetReader())
	err = decoder.Decode(cfg)
	if err != nil {
		processError(err)
	}
}

//ReadEnv reads enviroment variables
//and loads into config object
func (cfg *Config) ReadEnv(envObj []string) map[string]string {
	for _, element := range envObj {
		val, present := os.LookupEnv(element)
		if present == true {
			m[element] = val
		}
	}

	return m
}

//MapEnv enviroment variables
//and loads into config object
func (cfg *Config) MapEnv(m map[string]string) {
	for key, val := range m {
		switch key {
		case "SERVER_HOST":
			cfg.Server.Host = val
		case "SERVER_PORT":
			port, err := strconv.Atoi(val)
			if err == nil {
				cfg.Server.Port = port
			}
		case "AWACS_DB":
			cfg.DatabaseName.AwacsDBName = val
		case "AWACS_SMART_DB":
			cfg.DatabaseName.AwacsSmartDBName = val
		case "AWACS_SMART_STOCKIST_DB":
			cfg.DatabaseName.SmartStockistDBName = val
		case "DB_PORT":
			port, err := strconv.Atoi(val)
			if err == nil {
				cfg.DatabaseConfig.Port = port
			}
		case "DB_HOST":
			cfg.DatabaseConfig.Host = val
		case "DB_USERNAME":
			cfg.DatabaseConfig.Username = val
		case "DB_PASSWORD":
			cfg.DatabaseConfig.Password = val

		case "GRPC_HOST":
			cfg.GrpcConfig.Host = val
		case "GRPC_PORT":
			port, err := strconv.Atoi(val)
			if err == nil {
				cfg.GrpcConfig.Port = port
			}

		}
	}
}
