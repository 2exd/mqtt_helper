package utils

import (
	"fmt"
	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
	"os"
)

var ClipBoardChan = make(chan string, 10)

func GetClipBoard() {
	fmt.Println("--- Please press q+w+e to send clipboard ---")
	hook.Register(hook.KeyDown, []string{"q", "w", "e"}, func(e hook.Event) {
		fmt.Println("q+w+e")
		ClipBoardFunc()
	})

	fmt.Println("--- Please press ctrl + shift + w to stop clipboard ---")
	hook.Register(hook.KeyDown, []string{"w", "ctrl", "shift"}, func(e hook.Event) {
		fmt.Println("ctrl-shift-w")
		hook.End()
	})
}

func ClipBoardFunc() {
	// 尝试从剪切板获取文本
	clipboardText, err := clipboard.ReadAll()
	if err != nil {
		fmt.Println("无法读取剪切板内容:", err)
		return
	}
	ClipBoardChan <- clipboardText
}

func SaveClipBoard(s string, baseName string) (string, error) {
	fileName := fmt.Sprintf("%s.txt", baseName)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return fileName, err
	}
	defer file.Close()
	// 将内容写入文件
	_, err = file.WriteString(s)
	if err != nil {
		return fileName, err
	}
	return fileName, nil
}
