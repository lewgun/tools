package main

import (
	"fmt"
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
			provider.PATH:       "./testlog.log",
			provider.PREFIX:     "[lewgun]",
			provider.IS_ROLLING: true,
			provider.LOG_SIZE:   1,
			provider.FLAG:       provider.Lshortfile | provider.LstdFlags,
		},
	}

	manager.Open(provider.CONSOLE, manager.LEVEL_DEBUG, provider.Lshortfile|provider.LstdFlags, &arg)
	manager.Run()
	defer manager.Stop()

	for i := 0; i < 10000; i++ {

		manager.Error("hello world")
		//manager.Errorf("With format%s", "hello world")
		//	println(i)
	}

	logger := manager.Current()
	logger.(manager.StdLogger).Panic("need crash")
	fmt.Println("\n\n\n")
	//	manager.Change(provider.CONSOLE, manager.LEVEL_INFO, provider.Lshortfile|provider.LstdFlags)
	manager.Error("console hello world")
	manager.Debugf("console With format %s can't be output ", "hello world")
	//
	for {

		time.Sleep(time.Second * 3)
	}

}
