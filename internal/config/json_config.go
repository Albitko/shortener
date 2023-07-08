package config

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"github.com/Albitko/shortener/internal/entity"
)

func parseJSON(pathToJSONCfg string, cfg *entity.Config) error {
	var jsonCfg entity.JSONConfig
	configFile, _ := os.Open(pathToJSONCfg)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Println("error closing JSON config file ", err)
		}
	}(configFile)
	reader := bufio.NewReader(configFile)
	stat, _ := configFile.Stat()
	appConfigBytes := make([]byte, stat.Size())
	_, err := reader.Read(appConfigBytes)
	if err != nil {
		return err
	}
	err = json.Unmarshal(appConfigBytes, &jsonCfg)
	if err != nil {
		return err
	}

	if jsonCfg.EnableHTTPS {
		cfg.EnableHTTPS = jsonCfg.EnableHTTPS
	}
	if jsonCfg.ServerAddress != "" {
		cfg.ServerAddress = jsonCfg.ServerAddress
	}
	if jsonCfg.FileStoragePath != "" {
		cfg.FileStoragePath = jsonCfg.FileStoragePath
	}
	if jsonCfg.BaseURL != "" {
		cfg.BaseURL = jsonCfg.BaseURL
	}
	if jsonCfg.DatabaseDsn != "" {
		cfg.DatabaseDSN = jsonCfg.DatabaseDsn
	}
	if jsonCfg.TrustedSubnet != "" {
		cfg.TrustedSubnet = jsonCfg.TrustedSubnet
	}

	return nil
}
