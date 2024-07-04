package encode

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes)[:length], nil
}

func AesEncryptByECB(data []byte, key string) ([]byte, error) {
	keyLenMap := map[int]struct{}{16: {}, 24: {}, 32: {}}
	if _, ok := keyLenMap[len(key)]; !ok {
		return nil, fmt.Errorf("invalid key length")
	}
	keyByte := []byte(key)
	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	originByte := PKCS7Padding(data, blockSize)
	encryptResult := make([]byte, len(originByte))
	for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
		block.Encrypt(encryptResult[bs:be], originByte[bs:be])
	}
	return encryptResult, nil
}

// 补码
func PKCS7Padding(originByte []byte, blockSize int) []byte {
	padding := blockSize - len(originByte)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(originByte, padText...)
}

// AesEncryptByCBC 使用AES加密算法在CBC模式下加密数据
func AesEncryptByCBC(plaintext []byte, randomKey string) ([]byte, string) {
	keyBytes, _ := hex.DecodeString(randomKey)
	finalKey := make([]uint8, len(keyBytes))
	for i, b := range keyBytes {
		finalKey[i] = uint8(b)
	}
	block, _ := aes.NewCipher(finalKey)
	blockSize := block.BlockSize()
	plaintext = PKCS7Padding(plaintext, blockSize)
	iv := make([]byte, blockSize)
	io.ReadFull(rand.Reader, iv)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	copy(ciphertext[:blockSize], iv)
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[blockSize:], plaintext)
	ivT := make([]uint8, len(iv))
	for i, b := range iv {
		ivT[i] = uint8(b)
	}
	ivStr := hex.EncodeToString(ivT)
	return ciphertext, ivStr
}

// AesEncryptByCFB 使用AES加密算法在OFB模式下加密数据
func AesEncryptByCFB(plaintext []byte, randomKey string) ([]byte, string) {
	keyBytes, _ := hex.DecodeString(randomKey)
	finalKey := make([]uint8, len(keyBytes))
	for i, b := range keyBytes {
		finalKey[i] = uint8(b)
	}
	block, _ := aes.NewCipher(finalKey)
	iv := make([]byte, aes.BlockSize)
	io.ReadFull(rand.Reader, iv)
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	ivT := make([]uint8, len(iv))
	for i, b := range iv {
		ivT[i] = uint8(b)
	}
	ivStr := hex.EncodeToString(ivT)
	return ciphertext, ivStr
}

// AesEncryptByOFB 使用AES加密算法在OFB模式下加密数据
func AesEncryptByOFB(plaintext []byte, randomKey string) ([]byte, string) {
	keyBytes, _ := hex.DecodeString(randomKey)
	finalKey := make([]uint8, len(keyBytes))
	for i, b := range keyBytes {
		finalKey[i] = uint8(b)
	}
	block, _ := aes.NewCipher(finalKey)
	iv := make([]byte, aes.BlockSize)
	io.ReadFull(rand.Reader, iv)
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, plaintext)
	ivT := make([]uint8, len(iv))
	for i, b := range iv {
		ivT[i] = uint8(b)
	}
	ivStr := hex.EncodeToString(ivT)
	return ciphertext, ivStr
}

func XorEn(shellcode []byte) []byte {
	key1 := 13
	key2 := 241
	key3 := 143
	key4 := 221
	key5 := 98
	shellcode = Xor(shellcode, byte(key1))
	shellcode = Xor(shellcode, byte(key2))
	shellcode = Xor(shellcode, byte(key3))
	shellcode = Xor(shellcode, byte(key4))
	shellcode = Xor(shellcode, byte(key5))
	return shellcode
}

func CreateSeprarteLocalFile(encryptShellcode []byte, fileName string) {
	fileSrc := "./result/" + fileName
	file, err := os.Create(fileSrc) // 创建或覆盖文件
	if err != nil {
		panic(err.Error()) // 错误处理
	}
	defer file.Close() // 确保文件在使用后被关闭

	// 将加密后的数据写入文件
	if _, err := file.Write(encryptShellcode); err != nil {
		panic(err.Error()) // 错误处理
	}
}
