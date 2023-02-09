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
	return GetDataDir("app")
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
 * get the app data in cwd/.magicspace
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
