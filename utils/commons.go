package utils

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"

	color "github.com/fatih/color"
)

func HandleError(err error, text string, _panic bool) {
	if err != nil {
		fmt.Printf("Error: %s\n", text)
		fmt.Printf("%v", err)
		if _panic { panic(err) }
	}
}


func ToHex(item []byte) string {
	return hex.EncodeToString(item)
}

func FromHex(str string) ([]byte, error) {
	return hex.DecodeString(str)
}

func PrintSuccess(str string, vargs... interface{}){
	color.Green(str, vargs...)
}

func PrintInfo(str string, vargs... interface{}){
	color.Cyan(str, vargs...)
}

func PrintError(str string, vargs... interface{}){
	color.Red(str, vargs...)
}

func GetDataDir(subDir string) (string, error) {
	
	curDir,_ := os.Getwd()
	dataDir := path.Join(curDir, ".magicspace", subDir)

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		err = os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			return "", fmt.Errorf("failed to create data dir %s", dataDir)
		}
	}

	return dataDir, nil
}

