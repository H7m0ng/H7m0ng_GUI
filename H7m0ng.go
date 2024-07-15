package main

import (
	"H7m0ng/building"
	"H7m0ng/gui"
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"github.com/flopp/go-findfont"
	"github.com/gonutz/ide/w32"
	"os"
	"strings"
)

// 初始化要更改一下
func init() {
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "msyh.ttc") {
			os.Setenv("FYNE_FONT", path)
			println(os.Getenv("FYNE_FONT"))
			println("设置中文成功")
			break
		} else if strings.Contains(path, "msyh.ttf") {
			os.Setenv("FYNE_FONT", path)
			println(os.Getenv("FYNE_FONT"))
			println("设置中文成功")
			break
		}
	}
}

func main() {
	a := app.NewWithID("H7m0ng")
	w := a.NewWindow("H7m0ng-GUI(测试版本)") //新建一个窗口
	//  调试窗口，后续可以进行隐藏
	building.CloseWindows(w32.SW_HIDE)
	settingsItem := fyne.NewMenuItem("Settings", func() {
		w := a.NewWindow("Fyne 设置")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(618, 382))
		w.Show()
	})
	mainMenu := fyne.NewMainMenu(
		fyne.NewMenu("FILE", settingsItem),
	)
	tmp := gui.BypassAV(w)
	w.SetContent(tmp)
	w.SetMainMenu(mainMenu)
	w.SetMaster()
	w.Resize(fyne.NewSize(618, 382))
	w.ShowAndRun()
	w.Show()
}
