package core

import (
	"embed"
	"fmt"
	"io/fs"
)

type Option struct {
	Module            string //加载器 loader
	SrcFile           string
	ShellcodeEncode   string //加密模式 ecb
	Separate          string
	ShellcodeLocation string
	StrSC             *StrSC
	GoBuild           string
	ResourceFile      string
}

//go:embed module/*
var contentFS embed.FS

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

func GetLoaderContent(loaderName string) (loaderContent string) {
	//通过 ReadFile 读取文件内容
	fileData, _ := fs.ReadFile(contentFS, "module/"+loaderName)
	loaderContent = string(fileData)
	return loaderContent
}
