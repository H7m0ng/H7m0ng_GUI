package core

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

//type Option struct {
//	Module            string //模式
//	SrcFile           string //文件路径
//	ShellcodeEncode   string //加密
//	Separate          string //分离
//	ShellcodeLocation string //分离文件路径
//}

var (
	globalRand *rand.Rand
)

// 1.判断文件是否存在
func FileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, fmt.Errorf("文件 %s 不存在", filePath)
	}
	return err == nil, err
}

// 根据前缀进行复制操作
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

// LoadFiles 根据文件名前缀加载指定（嵌入的）目录下的所有文件内容到切片中,总感觉哪里不对
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
		if prefix == "ALL" || strings.HasPrefix(strings.ToLower(d.Name()), strings.ToLower(prefix)) {
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

// WriteFiles 将文件内容写入到目标目录，
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

// 删除加密后的bin文件和生成的go文件
func DelFile(filePath string) {
	os.Remove(filePath)
}
