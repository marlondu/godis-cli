package core

// created: 2018/9/6
import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

const (
	KEY = "YXILOVEYOUSOMUCH"
	IV  = "0102030405060708"
)

func Encrypt(text string) string {
	if text == "" {
		return ""
	}
	return aesEncrypt([]byte(KEY), []byte(text))
}

func Decrypt(text string) string {
	if text == "" {
		return ""
	}
	decoded, _ := base64.StdEncoding.DecodeString(text)
	return aesDecrypt([]byte(KEY), decoded)
}

func aesEncrypt(key []byte, src []byte) string {
	iv := []byte(IV)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Panic(err)
	}
	src = PKCS5Padding(src, block.BlockSize())
	encrypter := cipher.NewCBCEncrypter(block, iv)
	var dst = make([]byte, len(src))
	encrypter.CryptBlocks(dst, src)
	encoded := base64.StdEncoding.EncodeToString(dst)
	return encoded
}

func aesDecrypt(key []byte, src []byte) string {
	iv := []byte(IV)
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	decrypter := cipher.NewCBCDecrypter(block, iv)
	var dst = make([]byte, len(src))
	decrypter.CryptBlocks(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	dst = UnPKCS5Padding(dst)
	return string(dst)
}

func PKCS5Padding(src []byte, blockSize int) []byte {
	paddingSize := blockSize - len(src)%blockSize
	paddingBytes := bytes.Repeat([]byte{byte(paddingSize)}, paddingSize)
	return append(src, paddingBytes...)
}

func UnPKCS5Padding(src []byte) []byte {
	length := len(src)
	unPadding := int(src[length-1])
	return src[:(length - unPadding)]
}
