package utils

import (
	"io/ioutil"
	"os"
	"time"
)

// FileChange 文件改动
var FileChange = make(chan []byte, 5)

const (
	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Poll file for changes with this period.
	filePeriod = 10 * time.Second
)

func ReadFileIfModified(lastMod time.Time, filename string) (time.Time, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return lastMod, err
	}
	if !fi.ModTime().After(lastMod) {
		return lastMod, nil
	}
	p, err := ioutil.ReadFile(filename)
	if err != nil {
		return fi.ModTime(), err
	}
	FileChange <- p
	return fi.ModTime(), nil
}
