package utils

import (
	"fmt"
	"github.com/atotto/clipboard"
	hook "github.com/robotn/gohook"
	"os"
)

var ClipBoardChan = make(chan string, 10)

func GetClipBoard() {
	fmt.Println("--- 快捷键 q+w+e 发送剪切板内容给 server 端 ---")
	hook.Register(hook.KeyDown, []string{"q", "w", "e"}, func(e hook.Event) {
		ClipBoardFunc()
		fmt.Printf("%s, 剪切板内容已发送~\n", GetHhmmss())

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
	fmt.Printf("%s 收到新的内容\n", fileName)
	return fileName, nil
}
