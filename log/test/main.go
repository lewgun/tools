package main

import (
	"log/manager"
	"log/provider"
	"time"
)

const (
	PATH       = "path"
	PREFIX     = "prefix"
	IS_ROLLING = "is_rolling"
	LOG_SIZE   = "log_size"
)

func main() {

	arg := provider.Arg{
		Driver: "hello test",
		Extras: map[string]interface{}{
			provider.PATH:   "./testlog.log",
			provider.PREFIX: "[lewgun]",
		},
	}

	manager.Open(provider.FILE, manager.LEVEL_DEBUG, manager.Lshortfile, &arg)
	manager.Run()
	defer manager.Stop()

	manager.Error("hello world")
	manager.Errorf("With format%s", "hello world")
	manager.Change(provider.CONSOLE, manager.LEVEL_ERROR, manager.Lshortfile)
	manager.Error("console hello world")
	manager.Debugf("console With format %s can't be output ", "hello world")

	time.Sleep(time.Second * 3)

}
