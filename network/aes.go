package network

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func EncryptData(key, data []byte) (string, error) {
	// 创建AES加密块，使用指定的密钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 创建随机生成的初始化向量
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", fmt.Errorf("failed to generate random IV: %w", err)
	}

	// 创建AES加密模式，使用CBC模式，并提供初始化向量
	mode := cipher.NewCBCEncrypter(block, iv)

	// 对数据进行填充，确保其长度是块大小的倍数
	paddedData := pkcs7Padding(data, aes.BlockSize)

	// 创建加密缓冲区
	encrypted := make([]byte, len(paddedData))

	// 加密数据
	mode.CryptBlocks(encrypted, paddedData)

	// 将初始化向量和加密后的数据进行组合
	ciphertext := append(iv, encrypted...)

	// 对加密后的数据进行Base64编码，以便于传输和存储
	encoded := base64.StdEncoding.EncodeToString(ciphertext)

	return encoded, nil
}

func DecryptData(key []byte, encoded string) ([]byte, error) {
	// 对Base64编码的密文进行解码
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// 创建AES解密块，使用指定的密钥
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create AES cipher: %w", err)
	}

	// 检查密文长度是否合法
	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext is too short")
	}

	// 提取初始化向量
	iv := ciphertext[:aes.BlockSize]

	// 创建AES解密模式，使用CBC模式，并提供初始化向量
	mode := cipher.NewCBCDecrypter(block, iv)

	// 创建解密缓冲区
	decrypted := make([]byte, len(ciphertext)-aes.BlockSize)

	// 解密数据
	mode.CryptBlocks(decrypted, ciphertext[aes.BlockSize:])

	// 对解密后的数据进行去填充操作
	unpaddedData := pkcs7UnPadding(decrypted)

	return unpaddedData, nil
}

// pkcs7Padding 使用PKCS7填充对齐方式填充数据
func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padData := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padData...)
}

// pkcs7UnPadding 去除PKCS7填充的数据
func pkcs7UnPadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}
