package main

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"log"
	"net/url"
	"os"

	"github.com/dgryski/dgoogauth"
	"rsc.io/qr"
)

// token=password
func tfaAuthentication(secretBase32, authCode string) error {

	otpc := &dgoogauth.OTPConfig{
		Secret:      secretBase32,
		WindowSize:  3,
		HotpCounter: 0, //Time based
		// UTC:         true,
	}
	val, err := otpc.Authenticate(authCode)
	if err != nil {
		return err
	}
	if !val {
		//fmt.Println("Sorry, Authentication Failed!")
		return errors.New("authentication failed")
	}
	return nil
}

func generateQRcode(account, qrFilename string) (string, error) {

	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	if err != nil {
		log.Println(err)
		return "", err
	}
	secretBase32 := base32.StdEncoding.EncodeToString(secret)
	//fmt.Println(secretBase32)
	//account := "admin@mateors.com"
	issuer := "LXROOT"
	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		log.Println(err)
		return "", err
	}
	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)
	params := url.Values{}
	params.Add("issuer", issuer)
	params.Add("secret", secretBase32)
	URL.RawQuery = params.Encode()
	//fmt.Printf("URL is %s\n", URL.String())
	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		log.Println(err)
		return "", err
	}
	err = os.WriteFile(qrFilename, code.PNG(), 0600)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return secretBase32, nil
}
