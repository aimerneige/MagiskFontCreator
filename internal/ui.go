package internal

import (
	_ "embed" // import embedded assets

	"github.com/webview/webview"
)

const (
	title  = "Magisk 字体模块生成器"
	width  = 480
	height = 720
)

// ShowMainWindow shows the main window of the application.
func ShowMainWindow(debug bool) {
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle(title)
	w.SetSize(width, height, webview.HintNone)
	AllBindCollection(w)
	w.SetHtml(indexHTML)
	w.Run()
}
