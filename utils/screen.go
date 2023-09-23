package utils

import (
	"fmt"
	"github.com/kbinani/screenshot"
	hook "github.com/robotn/gohook"
	"image"
	"image/png"
	"os"
	"time"
)

func GetScreen() {
	fmt.Println("--- Please press 1+2+3 to send screenshot ---")
	hook.Register(hook.KeyDown, []string{"1", "2", "3"}, func(e hook.Event) {
		fmt.Println("1+2+3")
		ScreenshotFunc()
	})

	fmt.Println("--- Please press ctrl + shift + q to stop screenshot ---")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-q")
		hook.End()
	})
}

func ScreenshotFunc() {
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			panic(err)
		}
		ImgChan <- img
		// SaveScreen(img)
	}
}

func SaveScreen(img *image.RGBA, baseName string) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("15_04_05")
	fileName := fmt.Sprintf("%s_%s.png", baseName, formattedTime)
	file, _ := os.Create(fileName)
	err := png.Encode(file, img)
	if err != nil {
		file.Close() // 关闭文件
		panic(err)
	}

	// err = file.Sync() // 强制将文件内容写入磁盘
	// if err != nil {
	// 	file.Close()
	// 	panic(err)
	// }

	err = file.Close() // 保存并关闭文件
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", fileName)
}
