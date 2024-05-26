package utility

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

var WPORT string
var JWTSECRET, LOGENABLE, EMAILUSER, EMAILPASS, EMAILSERVER, EMAILPORT string

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

	JWTSECRET = os.Getenv("JWTSECRET")
	LOGENABLE = os.Getenv("LOGENABLE")
	EMAILUSER = os.Getenv("EMAILUSER")
	EMAILPASS = os.Getenv("EMAILPASS")
	EMAILSERVER = os.Getenv("EMAILSERVER")
	EMAILPORT = os.Getenv("EMAILPORT")

}

func JWTEncode(payload map[string]interface{}, signKey string) (string, error) {

	var mySigningKey = []byte(signKey) //
	token := jwt.New(jwt.SigningMethodHS256)
	//payload|claims
	claims := token.Claims.(jwt.MapClaims)
	for key, val := range payload {
		claims[key] = val
	}
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func JWTDecode(tokenStr, signKey string) (map[string]interface{}, error) {

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("error in token parsing")
		}
		return []byte(signKey), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	claims, isOk := token.Claims.(jwt.MapClaims)
	if !isOk {
		return nil, errors.New("unable to parse")
	}
	return claims, nil
}
