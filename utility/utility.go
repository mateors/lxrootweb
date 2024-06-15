package utility

import (
	"crypto/aes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
)

const voc string = "abcdfghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numbers string = "0123456789"
const symbols string = "!^*+_-=$#"

var WPORT string
var JWTSECRET, LOGENABLE, EMAILUSER, EMAILPASS, EMAILSERVER, EMAILPORT, STRIPE_SECRETKEY string

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

	rand.New(rand.NewSource(time.Now().UnixNano()))

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
	STRIPE_SECRETKEY = os.Getenv("STRIPE_SECRETKEY")

}

func GeneratePassword(length int, hasNumbers bool, hasSymbols bool) string {
	chars := voc
	if hasNumbers {
		chars = chars + numbers
	}
	if hasSymbols {
		chars = chars + symbols
	}
	return generatePassword(length, chars)
}
func generatePassword(length int, chars string) string {
	password := ""
	for i := 0; i < length; i++ {
		password += string([]rune(chars)[rand.Intn(len(chars))])
	}
	return password
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



func EncryptAES(key []byte, plaintext string) string {

	c, err := aes.NewCipher(key)
	CheckError(err)

	out := make([]byte, len(plaintext))

	c.Encrypt(out, []byte(plaintext))

	return hex.EncodeToString(out)
}

func DecryptAES(key []byte, ct string) {
	ciphertext, _ := hex.DecodeString(ct)

	c, err := aes.NewCipher(key)
	CheckError(err)

	pt := make([]byte, len(ciphertext))
	c.Decrypt(pt, ciphertext)

	s := string(pt[:])
	fmt.Println("DECRYPTED:", s)
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func LicenseEncDec() {

	// cipher key
	key := "thisis32bitlongpassphraseimusing"

	// plaintext
	pt := "-LXROOT LXROOT LXROOT-"

	c := EncryptAES([]byte(key), pt)

	// plaintext
	fmt.Println(pt)

	// ciphertext
	fmt.Println(c)

	// decrypt
	DecryptAES([]byte(key), c)
}
