package building

import (
	"H7m0ng/fileSetting"
	"bytes"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type PyContent struct {
	Name    string
	Content string
}

// 修改文件名：将"1.exe"修改为"1-sign.exe"
func modifyFilename(oldPath string) string {
	// 使用filepath.Dir获取目录部分
	dir := filepath.Dir(oldPath)
	// 使用filepath.Base获取文件名部分
	filename := filepath.Base(oldPath)
	ext := filepath.Ext(filename)             // 获取文件扩展名
	base := strings.TrimSuffix(filename, ext) // 移除扩展名，得到文件名前缀
	// 修改文件名前缀，并重新组合文件名
	newFilename := base + "-sign" + ext
	// 重新组合完整路径
	return filepath.Join(dir, newFilename)
}

func GetSign(signSrc string) {
	// 获取指定目录下所有文件的路径
	_ = LoadAndSaveTestPy()
	dir := "./result/"
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		name := d.Name()
		if strings.HasSuffix(strings.ToLower(name), ".exe") {
			fmt.Println(path) // 直接打印由filepath.WalkDir提供的完整路径
			//如果是exe的话就进行签名操作 ，用一点笨方法，没办法，菜是原罪
			signPath := modifyFilename(path) //
			println(signPath)
			// 构建要执行的命令   -i 签名文件  -o 输出路径  -t
			cmd := exec.Command("python", "result/test.py", "-i", signSrc, "-o", signPath, "-t", path)
			println(cmd.String())
			// 创建一个新的bytes.Buffer用于存储命令的输出
			var out bytes.Buffer
			cmd.Stdout = &out // 将命令的标准输出重定向到我们的buffer
			// 执行命令
			err := cmd.Run()
			if err != nil {
				fmt.Println("执行命令出错:", err)
				return nil
			}
			// 输出命令的结果
			fmt.Println("命令执行输出:", out.String())
			fileSetting.DelFile(path)
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

//go:embed py/test.py
var embeddedFiles embed.FS

func LoadAndSaveTestPy() error {
	// 构造 test.py 文件在嵌入文件系统中的路径
	embedPath := "py/test.py"
	println("正在从嵌入的文件系统中加载文件...")
	// 从嵌入的文件系统中打开文件
	file, err := embeddedFiles.Open(embedPath)
	if err != nil {
		return fmt.Errorf("failed to open embedded file: %w", err)
	}
	defer file.Close()
	// 构造要在当前目录下保存的文件路径
	savePath := filepath.Join("./result/", "test.py")
	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(savePath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	// 创建并打开文件以写入
	outputFile, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer outputFile.Close()
	// 将嵌入文件的内容复制到新文件中
	_, err = io.Copy(outputFile, file)
	if err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}
	println("文件已成功保存至", savePath)
	return nil
}
