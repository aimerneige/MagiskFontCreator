package internal

import (
	"archive/zip"
	"embed"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"

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

//go:embed xml/simple.xml
var simpleXMLTemplate string

//go:embed simple
var simple embed.FS

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

// AllBindCollection all bind collection
func AllBindCollection(w webview.WebView) {
	w.Bind("github", openURL("https://github.com/aimerneige/MagiskFontCreator"))
	w.Bind("issue", openURL("https://github.com/aimerneige/MagiskFontCreator/issue"))
	w.Bind("aboutPage", switchPage(w, aboutHTML))
	w.Bind("cjkPage", switchPage(w, cjkHTML))
	w.Bind("indexPage", switchPage(w, indexHTML))
	w.Bind("multiPage", switchPage(w, multiHTML))
	w.Bind("singlePage", switchPage(w, singleHTML))
	w.Bind("generateSingle", generateSingleFontPack())
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

// generateSingleFontPack 生成单字体字体模块
func generateSingleFontPack() interface{} {
	return func(fontBase64, fontType, moduleID, moduleName, version, versionCode, author, description string) error {
		// 将 base64 字符串解码为字节数组
		fontBytes, err := b64.StdEncoding.DecodeString(fontBase64)
		if err != nil {
			return err
		}
		// 替换所有回车符为空格
		moduleID = strings.Replace(moduleID, "\n", " ", -1)
		moduleName = strings.Replace(moduleName, "\n", " ", -1)
		version = strings.Replace(version, "\n", " ", -1)
		versionCode = strings.Replace(versionCode, "\n", " ", -1)
		author = strings.Replace(author, "\n", " ", -1)
		description = strings.Replace(description, "\n", " ", -1)
		// module.prop 字符串
		moduleProp := fmt.Sprintf("id=%s\nname=%s\nversion=%s\nversionCode=%s\nauthor=%s\ndescription=%s", moduleID, moduleName, version, versionCode, author, description)

		return zipSimple(fontBytes, fontType, moduleProp)
	}
}

func zipSimple(fontData []byte, fontType string, modulePropString string) error {
	if fontType != "ttf" && fontType != "otf" {
		return fmt.Errorf("Invalid font type %s", fontType)
	}
	// buf := new(bytes.Buffer)
	f, err := os.Create("./module.zip")
	if err != nil {
		return err
	}
	defer f.Close()
	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()
	assetsFilePath := []string{
		"META-INF/com/google/android/update-binary",
		"META-INF/com/google/android/updater-script",
		"system/fonts/EmptyFont-Black.ttf",
		"system/fonts/EmptyFont-BlackItalic.ttf",
		"system/fonts/EmptyFont-Bold.ttf",
		"system/fonts/EmptyFont-BoldItalic.ttf",
		"system/fonts/EmptyFont-Italic.ttf",
		"system/fonts/EmptyFont-Light.ttf",
		"system/fonts/EmptyFont-LightItalic.ttf",
		"system/fonts/EmptyFont-Medium.ttf",
		"system/fonts/EmptyFont-MediumItalic.ttf",
		"system/fonts/EmptyFont-Regular.ttf",
		"system/fonts/EmptyFont-Thin.ttf",
		"system/fonts/EmptyFont-ThinItalic.ttf",
		"system/product/fonts/GoogleSans-Bold.ttf",
		"system/product/fonts/GoogleSans-BoldItalic.ttf",
		"system/product/fonts/GoogleSans-Italic.ttf",
		"system/product/fonts/GoogleSans-Medium.ttf",
		"system/product/fonts/GoogleSans-MediumItalic.ttf",
		"system/product/fonts/GoogleSans-Regular.ttf",
		"customize.sh",
		"post-fs-data.sh",
		"README.md",
		"sepolicy.rule",
		"service.sh",
		"system.prop",
		"uninstall.sh",
	}
	for _, path := range assetsFilePath {
		assetFile, err := simple.Open("simple/" + path)
		if err != nil {
			return err
		}
		defer assetFile.Close()
		assetZip, err := zipWriter.Create(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(assetZip, assetFile)
		if err != nil {
			return err
		}
	}
	fontFileName := "customized." + fontType
	writeFileToZip(fontData, "system/fonts/"+fontFileName, zipWriter)
	writeFileToZip([]byte(modulePropString), "module.prop", zipWriter)
	fontXML := getSimpleFontXML(fontFileName, fontFileName, fontFileName, fontFileName, fontFileName, fontFileName, fontFileName, fontFileName, fontFileName)
	writeFileToZip([]byte(fontXML), "system/etc/fonts.xml", zipWriter)

	return nil
}

func writeFileToZip(data []byte, path string, zipWriter *zip.Writer) error {
	f, err := zipWriter.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func getSimpleFontXML(w1, w2, w3, w4, w5, w6, w7, w8, w9 string) string {
	ret := simpleXMLTemplate
	ret = strings.Replace(ret, "fontw1.ttf", w1, -1)
	ret = strings.Replace(ret, "fontw2.ttf", w2, -1)
	ret = strings.Replace(ret, "fontw3.ttf", w3, -1)
	ret = strings.Replace(ret, "fontw4.ttf", w4, -1)
	ret = strings.Replace(ret, "fontw5.ttf", w5, -1)
	ret = strings.Replace(ret, "fontw6.ttf", w6, -1)
	ret = strings.Replace(ret, "fontw7.ttf", w7, -1)
	ret = strings.Replace(ret, "fontw8.ttf", w8, -1)
	ret = strings.Replace(ret, "fontw9.ttf", w9, -1)

	return ret
}
