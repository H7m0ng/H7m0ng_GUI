package fileSetting

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	globalRand *rand.Rand
)

// GoFileExists 此函数的目的是确保result目录下的go文件都被编译
func GoFileExists() bool {
	dir := "./result/"
	// 使用filepath.WalkDir遍历目录下的所有文件和子目录
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // 如果有错误，返回错误
		}

		// 检查文件名是否以.go结尾
		if strings.HasSuffix(d.Name(), ".go") {
			// 找到.go文件，直接返回nil表示成功，不再继续遍历
			return nil
		}

		// 如果不是.go文件，继续遍历
		return nil
	})
	return false
}

// GenerateGoFile 此函数的目的是将指定目录下的文件拷贝到指定目录，根据前缀确定写入内容
func GenerateGoFile(options Option) {
	destDir := "./result"
	//将指定模版的文件拷贝到本地文件并重命名
	if options.Module == "ALL" { // 判断是否为全部模块加载
		prefix := "ALL"
		files, _ := LoadFiles(prefix) // 如何prefix == all ，必须读取除ALL.txt
		// 写入文件内容到目标目录，
		if err := WriteFiles(files, destDir); err != nil {
			fmt.Println("Error writing files:", err)
			return
		}
		fmt.Println("Files copied successfully.")
	} else { // 复制指定loader
		prefix := options.Module
		files, err := LoadFiles(prefix)
		if err != nil {
			fmt.Println("Error loading files:", err)
			return
		}
		if len(files) == 0 {
			fmt.Println("No files found with the prefix:", prefix)
			return
		}
		// 写入文件内容到目标目录
		if err := WriteFiles(files, destDir); err != nil {
			fmt.Println("Error writing files:", err)
			return
		}
		fmt.Println("Files copied successfully.")
	}
}

type FileContent struct {
	Name    string
	Content string
}

// LoadFiles 根据文件名前缀加载指定（嵌入的）目录下的所有文件内容到切片中
func LoadFiles(prefix string) ([]FileContent, error) {
	var files []FileContent
	err := fs.WalkDir(contentFS, "module", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // 返回遍历错误
		}
		if d.IsDir() {
			return nil // 忽略目录
		}
		// 如果 prefix 是 "ALL"，或者文件名以 prefix 开头（忽略大小写）
		if prefix == "ALL" || strings.ToLower(d.Name()) == strings.ToLower(prefix+".txt") {
			data, err := fs.ReadFile(contentFS, path)
			if err != nil {
				return err // 返回读取文件错误
			}
			files = append(files, FileContent{
				Name:    filepath.Base(path), // 获取文件名（不带路径）
				Content: string(data),
			})
		}
		return nil
	})
	return files, err
}

// WriteFiles 将文件内容写入到目标目录
func WriteFiles(files []FileContent, destDir string) error {
	// 确保目标目录存在
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}
	for _, file := range files {
		// 计算目标文件路径
		rel, err := filepath.Rel(filepath.Dir(file.Name), file.Name)
		if err != nil {
			return err
		}
		if rel == "ALL.txt" {
			continue
		}
		rel = strings.TrimSuffix(rel, ".txt")
		rel = rel + RandomString() + ".go"
		targetPath := filepath.Join(destDir, rel)
		println(rel)
		// 创建目标目录（如果不存在）
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}
		// 写入文件内容
		content := []byte(file.Content)                                     // 或者从其他地方获取内容
		if err := ioutil.WriteFile(targetPath, content, 0644); err != nil { // 或者使用 os.WriteFile
			return err
		}
	}
	return nil
}

// RandomString 生成一个随机字符串
func RandomString() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	globalRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	var r *rand.Rand
	if globalRand != nil {
		r = globalRand
	} else {
		println("Error")
	}
	b := make([]byte, 8)
	for i := range b {
		// 使用全局的 Rand 实例或 rand 包的全局函数来生成随机数
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// DelFile 删除文件功能
func DelFile(filePath string) {
	os.Remove(filePath)
}
