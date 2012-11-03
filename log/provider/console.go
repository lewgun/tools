package provider

import (
	"fmt"
)

const (
	CONSOLE = "console"
)

func init() {
	register(CONSOLE, newConsole)
}

func newConsole(arg *Arg) interface{} {
	return &Console{}

}

//日志接口
type Console struct {
}

func (c *Console) Tracef(format string, params ...interface{}) {
	c.logHelper(format, params...)
}

func (c *Console) Debugf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}
func (c *Console) Infof(format string, params ...interface{}) {
	c.logHelper(format, params...)
}
func (c *Console) Warnf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}
func (c *Console) Errorf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}

func (c *Console) Criticalf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}

func (c *Console) logHelper(format string, params ...interface{}) {

	if format != "" {
		if format[len(format)-1:] != "\n" {
			format += "\n"
		}
		fmt.Printf(format, params...)

	} else {
		fmt.Println(params...)
	}

}

func (c *Console) Release() {

}
