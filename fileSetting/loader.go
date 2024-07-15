package fileSetting

import (
	"embed"
	"fmt"
	"io/fs"
)

type Option struct {
	Module            string //加载器 loader
	SrcFile           string //bin文件路径
	ShellcodeEncode   string //加密模式 ecb
	Separate          string //分离选项
	ShellcodeLocation string //分离后的shellcode路径
	GoBuild           string //go编译参数
	OtherOpt          bool   //garble
}

//go:embed module/*
var contentFS embed.FS

// GetLoaderNames 获取模板目录下的文件名
func GetLoaderNames() (loaderNames []string) {
	// 通过 ReadDir 读取目录中的文件
	files, err := fs.ReadDir(contentFS, "module")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}
	// 遍历文件并将内容保存到数组中
	for _, file := range files {
		loaderNames = append(loaderNames, file.Name())
	}
	return loaderNames
}
