package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

//Create returns self encoded token
func Create(message string, secret string, digestType string) (string, error) {

	signedToken, err := signToken(message, secret, digestType)
	return signedToken, err

}

//Verify extracts payload from the given token post verification
func Verify(token string, secret string) (string, error) {

	payLoad, err := verifyToken(token, secret)
	return payLoad, err
}

func signToken(payload string, key string, digestType string) (string, error) {

	var (
		protected, plainText, signature string
		err                             error
	)

	switch digestType {
	case "sha256":
		protected = "HS256"
	case "sha384":
		protected = "HS384"
	case "sha512":
		protected = "HS512"
	}

	plainText = base64.StdEncoding.EncodeToString([]byte(protected)) +
		"." +
		base64.StdEncoding.EncodeToString([]byte(payload))

	signature, err = hmacHash(digestType, key, plainText)

	return plainText + "." + base64.StdEncoding.EncodeToString([]byte(signature)), err
}

func verifyToken(token string, key string) (string, error) {
	var (
		digestType, challenge string
	)

	protected, payload, plainText, signature, err := decodeToken(token)

	switch protected {
	case "HS256":
		digestType = "sha256"
	case "HS384":
		digestType = "sha384"
	case "HS512":
		digestType = "sha512"
	}

	challenge, err = hmacHash(digestType, key, plainText)

	if challenge == signature {
		return payload, err
	} else {
		return "", err
	}
}

func decodeToken(token string) (string, string, string, string, error) {

	//todo check return for absence of "." and handle error approp
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 3 {
		//todo return approp error
		return "", "", "", "", nil
	}

	plainText := tokenParts[0] + "." + tokenParts[1]
	protected, err := base64.StdEncoding.DecodeString(tokenParts[0])
	if err != nil {
		//todo return approp error
		return "", "", "", "", nil
	}
	payload, err := base64.StdEncoding.DecodeString(tokenParts[1])
	if err != nil {
		//todo return approp error
		return "", "", "", "", nil
	}
	signature, err := base64.StdEncoding.DecodeString(tokenParts[2])
	if err != nil {
		//todo return approp error
		return "", "", "", "", nil
	}

	return string(protected), string(payload), plainText, string(signature), err
}

func hmacHash(digestType string, key string, plainText string) (string, error) {

	switch digestType {
	case "sha256":
		return string(hmac.New(sha256.New, []byte(key)).Sum([]byte(plainText))), nil
	default:
		//todo return appropriate error
		return "", nil
	}
}
