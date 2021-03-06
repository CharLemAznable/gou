package pbe

import (
	"strings"

	"github.com/bingoohuang/gou/enc"

	"crypto/cipher"
	"crypto/des" // nolint
	"crypto/md5" // nolint
	"crypto/rand"
)

// Encrypt PrintEncrypt the plainText based on password and iterations with random salt.
// The result contains the first 8 bytes salt before BASE64.
func Encrypt(plainText, password string, iterations int) (string, error) {
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	encText, err := doEncrypt(plainText, password, salt, iterations)
	if err != nil {
		return "", err
	}

	return enc.Base64SafeEncode(append(salt, encText...)), nil
}

// Decrypt PrintDecrypt the cipherText(result of Encrypt) based on password and iterations.
func Decrypt(cipherText, password string, iterations int) (string, error) {
	msgBytes, err := enc.Base64SafeDecode(cipherText)
	if err != nil {
		return "", err
	}

	salt := msgBytes[:8]
	encText := msgBytes[8:]

	return doDecrypt(encText, password, salt, iterations)
}

// EncryptSalt PrintEncrypt the plainText based on password and iterations with fixed salt.
func EncryptSalt(plainText, password, fixedSalt string, iterations int) (string, error) {
	salt := make([]byte, 8)
	copy(salt, fixedSalt)

	encText, err := doEncrypt(plainText, password, salt, iterations)
	if err != nil {
		return "", err
	}

	return enc.Base64SafeEncode(encText), nil
}

// DecryptSalt PrintDecrypt the cipherText(result of EncryptSalt) based on password and iterations.
func DecryptSalt(cipherText, password, fixedSalt string, iterations int) (string, error) {
	msgBytes, err := enc.Base64SafeDecode(cipherText)
	if err != nil {
		return "", err
	}

	salt := make([]byte, 8)
	copy(salt, fixedSalt)

	return doDecrypt(msgBytes, password, salt, iterations)
}

func doEncrypt(plainText, password string, salt []byte, iterations int) ([]byte, error) {
	padNum := byte(8 - len(plainText)%8)
	for i := byte(0); i < padNum; i++ {
		plainText += string(padNum)
	}

	dk, iv := getDerivedKey(password, string(salt), iterations)
	block, err := des.NewCipher(dk) // nolint

	if err != nil {
		return nil, err
	}

	encrypter := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plainText))
	encrypter.CryptBlocks(encrypted, []byte(plainText))

	return encrypted, nil
}

func doDecrypt(encText []byte, password string, salt []byte, iterations int) (string, error) {
	dk, iv := getDerivedKey(password, string(salt), iterations)
	block, err := des.NewCipher(dk) // nolint

	if err != nil {
		return "", err
	}

	decrypter := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encText))
	decrypter.CryptBlocks(decrypted, encText)

	decryptedString := strings.TrimRight(string(decrypted), "\x01\x02\x03\x04\x05\x06\x07\x08")

	return decryptedString, nil
}

func getDerivedKey(password, salt string, iterations int) ([]byte, []byte) {
	key := md5.Sum([]byte(password + salt)) // nolint

	for i := 0; i < iterations-1; i++ {
		key = md5.Sum(key[:]) // nolint
	}

	return key[:8], key[8:]
}
