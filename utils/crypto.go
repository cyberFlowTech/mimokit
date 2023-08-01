package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

func ToMD5(str string) string {
	hash := md5.New()
	hash.Write([]byte(str))
	return hex.EncodeToString(hash.Sum(nil))
}

func ToSha256(rawPassword string) (string, error) {
	h := sha256.New()
	_, err := h.Write([]byte(rawPassword))
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func aesencrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func aesdecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

// AESEncrpypt 外部调用
func AESEncrpypt(origData string, key []byte) (string, error) {
	resultByte, err := aesencrypt([]byte(origData), key)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(resultByte), nil
}

// AESDecrypt 外部调用
func AESDecrypt(origData string, key []byte) (string, error) {
	decodeBase64Byte, _ := base64.StdEncoding.DecodeString(origData)
	resultByte, err := aesdecrypt(decodeBase64Byte, key)
	if err != nil {
		return "", err
	}
	return string(resultByte), nil
}

// ReversibleEncrypt 可逆加密
func reversibleEncrypt(data, key string) string {
	keyBytes := []byte(ToMD5(key))
	keyLen := len(keyBytes)
	x := 0
	dataLen := len(data)
	char := make([]byte, 0, dataLen)
	str := make([]byte, 0, dataLen)
	for i := 0; i < dataLen; i++ {
		if x == keyLen {
			x = 0
		}
		char = append(char, keyBytes[x])
		x++
	}

	for i := 0; i < dataLen; i++ {
		str = append(str, byte(int(data[i])+(int(char[i])%256)))
	}
	return base64.StdEncoding.EncodeToString(str)
}

// ReversibleDecrypt 解密
func reversibleDecrypt(encrypted, key string) (string, error) {
	encryptedBytes, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	keyBytes := []byte(ToMD5(key))
	keyLen := len(keyBytes)
	x := 0
	encryptedLen := len(encryptedBytes)
	char := make([]byte, 0, encryptedLen)
	str := make([]byte, 0, encryptedLen)

	for i := 0; i < encryptedLen; i++ {
		if x == keyLen {
			x = 0
		}
		char = append(char, keyBytes[x])
		x++
	}

	for i := 0; i < encryptedLen; i++ {
		str = append(str, byte(int(encryptedBytes[i])-(int(char[i])%256)))
	}

	return string(str), nil
}

// IdEncrypt id加密
func IdEncrypt(data, userId string) string {
	if len(userId) < 0 {
		userId = "0"
	}
	str := userId + "|" + data
	return reversibleEncrypt(str, ClubEncryptKey)
}

// IdDecrypt id解密
func IdDecrypt(encrypted string) (string, error) {
	str, err := reversibleDecrypt(encrypted, ClubEncryptKey)
	if err != nil {
		return "", err
	}
	arr := strings.Split(str, "|")

	return arr[1], nil
}
