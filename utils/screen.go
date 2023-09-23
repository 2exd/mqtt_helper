package utils

import (
	"fmt"
	"github.com/kbinani/screenshot"
	hook "github.com/robotn/gohook"
	"image"
	"image/png"
	"os"
)

func GetScreen() {
	fmt.Println("--- 快捷键 1+2+3 发送屏幕截图给 server 端 ---")
	hook.Register(hook.KeyDown, []string{"1", "2", "3"}, func(e hook.Event) {
		ScreenshotFunc()
		fmt.Printf("%s, 截图已发送~\n", GetHhmmss())
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
	}
}

func SaveScreen(img *image.RGBA, baseName string) {
	fileName := fmt.Sprintf("%s_%s.png", baseName, GetHh_mm_ss())
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
	fmt.Printf("%s 截图已保存\n", fileName)
}
