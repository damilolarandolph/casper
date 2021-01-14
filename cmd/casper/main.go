package main

import (
	"os"
	"runtime"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

func main() {

	app := widgets.NewQApplication(len(os.Args), os.Args)

	window := gui.NewQWindow(gui.QGuiApplication_PrimaryScreen())
	window.Resize2(250, 150)
	window.SetTitle("Simple exampl")
	window.Show()
	app.Exec()
	return
}

func init() {
	runtime.LockOSThread()
}
