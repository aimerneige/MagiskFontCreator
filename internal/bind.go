package internal

import (
	_ "embed" // import embedded assets
	"fmt"
	"os/exec"
	"runtime"

	b64 "encoding/base64"

	"github.com/webview/webview"
)

//go:embed html/about.html
var aboutHTML string

//go:embed html/cjk.html
var cjkHTML string

//go:embed html/index.html
var indexHTML string

//go:embed html/multi.html
var multiHTML string

//go:embed html/single.html
var singleHTML string

// AllBindCollection all bind collection
func AllBindCollection(w webview.WebView) {
	w.Bind("github", openURL("https://github.com/aimerneige/MagiskFontCreator"))
	w.Bind("issue", openURL("https://github.com/aimerneige/MagiskFontCreator/issue"))
	w.Bind("aboutPage", switchPage(w, aboutHTML))
	w.Bind("cjkPage", switchPage(w, cjkHTML))
	w.Bind("indexPage", switchPage(w, indexHTML))
	w.Bind("multiPage", switchPage(w, multiHTML))
	w.Bind("singlePage", switchPage(w, singleHTML))
}

// openURL open url in browser
func openURL(url string) interface{} {
	return func() error {
		var err error
		switch runtime.GOOS {
		case "linux":
			err = exec.Command("xdg-open", url).Start()
		case "windows":
			err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		case "darwin":
			err = exec.Command("open", url).Start()
		default:
			err = fmt.Errorf("unsupported platform")
		}
		return err
	}
}

// switchPage switch html page
func switchPage(w webview.WebView, html string) interface{} {
	return func() error {
		sEnc := b64.StdEncoding.EncodeToString([]byte(html))
		w.Navigate("data:text/html;base64," + sEnc)
		return nil
	}
}
