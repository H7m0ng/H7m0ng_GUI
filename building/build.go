package building

import (
	"H7m0ng/core"
	"fmt"
	"github.com/gonutz/ide/w32"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func BuildExe(buildOpt string) {
	println("start  build  exe......")
	root := "./result" // 替换为你想要遍历的目录
	goFiles, err := collectGoFiles(root)
	if err != nil {
		fmt.Println("Error collecting Go files:", err)
		return
	}
	// 编译参数的设置
	for _, filePath := range goFiles {
		exeFilePath := strings.Replace(filePath, ".go", ".exe", -1)
		if buildOpt == "" {
			buildOpt = "default"
		}
		if buildOpt == "default" {
			cmd := exec.Command("go", "build", "-o", exeFilePath, filePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			println(cmd.String())
			// 执行命令
			err := cmd.Run()
			println("ExeFilePath:" + exeFilePath)
			if err != nil {
				println("Error")
				return
			}
			core.DelFile(filePath)
		} else {
			if buildOpt == "-ldflags=-w -s -trimpath" {
				buildOpt1 := "-ldflags=-w -s"
				buildOpt2 := "-trimpath"
				cmd := exec.Command("go", "build", "-o", exeFilePath, buildOpt1, buildOpt2, filePath)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				println(cmd.String())
				// 执行命令
				err := cmd.Run()
				println("ExeFilePath:" + exeFilePath)
				if err != nil {
					println("Error")
					return
				}
				core.DelFile(filePath)
			}
			cmd := exec.Command("go", "build", "-o", exeFilePath, buildOpt, filePath)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			println(cmd.String())
			// 执行命令
			err := cmd.Run()
			println("ExeFilePath:" + exeFilePath)
			if err != nil {
				println("Error")
				return
			}
			core.DelFile(filePath)
		}
	}

	println("Build Finally！")
}
func collectGoFiles(dir string) ([]string, error) {
	var goFiles []string
	// 使用filepath.WalkDir遍历目录
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // 返回遍历过程中遇到的错误
		}
		// 检查文件名是否以.go结尾
		if strings.HasSuffix(d.Name(), ".go") {
			// 如果是，将文件名（相对于dir的路径）添加到切片中
			// 注意：filepath.WalkDir提供的path是相对于遍历起始目录的相对路径
			goFiles = append(goFiles, path)
		}
		// 没有错误时返回nil以继续遍历
		return nil
	})
	if err != nil {
		return nil, err // 如果遍历过程中有错误，返回错误
	}

	return goFiles, nil // 返回收集到的.go文件名切片
}

func CloseWindows(commandShow uintptr) {
	console := w32.GetConsoleWindow()
	if console != 0 {
		_, consoleProcID := w32.GetWindowThreadProcessId(console)
		if w32.GetCurrentProcessId() == consoleProcID {
			w32.ShowWindowAsync(console, commandShow)
		}
	}
}
