package utils

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

/**
 * get the app data dir
 */
func GetAppDataDir() (string, error) {

	curDir,_ := os.Getwd()
	dataDir := path.Join(curDir, ".data", "app")

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create appdata dir %s", dataDir)
		}
	}

	return dataDir, nil
}

func loadConfig(filename string) (*viper.Viper, error) {

	v := viper.New()

	appDir, err := GetAppDataDir()

	if err != nil {
		return nil,err
	}

	filePath := path.Join(appDir, filename+".toml")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		_, err := os.Create(filePath)

		if err != nil{
			return nil, fmt.Errorf("failed to create appdata file %s, err=%w", filePath, err)
		}
	}

	v.SetConfigName(filename)
	v.SetConfigType("toml")
	v.AddConfigPath(appDir)

	return v, nil
}

/**
 * get the app dava in cwd/.magicspace
 **/
func SaveAppData(filename string, data map[string]interface{}) (*viper.Viper, error) {

	v, err := loadConfig(filename)

	if err != nil {
		return nil,err
	}

	for key, value := range data {
		v.Set(key, value)
	}

	v.WriteConfig()

	return v, nil
}

/**
 * get the app dava in $HOME/.magicspace
 **/
func GetAppData(filename string) (*viper.Viper, error) {

	v,err := loadConfig(filename)

	if err != nil {
		return nil, err
	}

	err = v.ReadInConfig()

	if err != nil {
		return nil, err
	}

	return v, nil
}
