package utility

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var WPORT string

func envFilePath() string {
	dir, _ := os.Getwd()
	var envPath = filepath.Join(dir, ".env")
	_, err := os.Stat(envPath)
	if err != nil {
		//envPath = filepath.Join(dir, "../.env")
		log.Println(err)
	}
	return envPath
}

func init() {

	err := godotenv.Load(envFilePath())
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	WPORT = os.Getenv("WPORT")

}
