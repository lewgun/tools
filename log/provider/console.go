package provider

import (
	//	"fmt"
	"log"
	"runtime"
	//	"io"
	"os"
	"sync"
)

const (
	CONSOLE = "console"
)

func init() {
	register(CONSOLE, newConsole)
}

func newConsole(arg *Arg, prefix string) interface{} {
	if arg == nil {
		return nil
	}

	c := &console{}

	if err := c.init(arg, prefix); err != nil {
		return nil
	}
	return c

}

//日志接口
type console struct {
	logger *log.Logger
	//the body prefix, option
	prefix string

	//the body flag, option
	flag int

	uname  string
	locker sync.Mutex
	//	writer io.Writer
}

func (c *console) init(arg *Arg, prefix string) error {

	extras := arg.Extras

	c.uname = arg.Driver
	c.prefix = prefix

	for k, v := range extras {
		switch k {

		case PREFIX:
			c.prefix = v.(string)

		case FLAG:
			c.flag = v.(int)

		default:
			//placeholder
		}
	}

	c.logger = log.New(os.Stdout, "", 0)
	//c.writer = os.Stdout //log.New(os.Stdout, "", 0)

	return nil
}

func (c *console) genStdLogHelper(format string, params ...interface{}) []interface{} {

	_, file, line, _ := runtime.Caller(2)
	if format != "" {
		format = c.prefix + " " + format

		params = append([]interface{}{file, line}, params...)
	} else {
		params = append([]interface{}{file, line, c.prefix}, params...)
	}
	return params

}
func (c *console) Tracef(format string, params ...interface{}) {

	all := outputHelper(format, false, c.flag, c.prefix, params...)
	c.logHelper(all)
}

func (c *console) Debugf(format string, params ...interface{}) {
	all := outputHelper(format, false, c.flag, c.prefix, params...)
	c.logHelper(all)
}
func (c *console) Infof(format string, params ...interface{}) {
	all := outputHelper(format, false, c.flag, c.prefix, params...)
	c.logHelper(all)
}
func (c *console) Warnf(format string, params ...interface{}) {
	all := outputHelper(format, false, c.flag, c.prefix, params...)
	c.logHelper(all)
}
func (c *console) Errorf(format string, params ...interface{}) {
	all := outputHelper(format, false, c.flag, c.prefix, params...)
	c.logHelper(all)
}
func (c *console) Criticalf(format string, params ...interface{}) {
	all := outputHelper(format, false, c.flag, c.prefix, params...)
	c.logHelper(all)
}

func (c *console) logHelper(all []byte) {

	c.logger.Print(string(all))
	//c.writer.Write(all)

}

func (c *console) Release() {

}

//standard log methods
func (c *console) Fatal(v ...interface{}) {
	params := c.genStdLogHelper("", v...)

	all := outputHelper("", false, c.flag, PREFIX_STD_FATAL, params...)
	c.logHelper(all)
	os.Exit(1)

}
func (c *console) Fatalf(format string, v ...interface{}) {
	params := c.genStdLogHelper(format, v...)
	all := outputHelper(format, false, c.flag, PREFIX_STD_FATAL, params...)
	c.logHelper(all)
	os.Exit(1)
}
func (c *console) Fatalln(v ...interface{}) {
	params := c.genStdLogHelper("", v...)
	all := outputHelper("", true, c.flag, PREFIX_STD_FATAL, params...)
	c.logHelper(all)
	os.Exit(1)
}

func (c *console) Output(calldepth int, s string) error {
	return c.logger.Output(calldepth+1, s)

}
func (c *console) Panic(v ...interface{}) {
	params := c.genStdLogHelper("", v...)
	all := outputHelper("", false, c.flag, PREFIX_STD_PANIC, params...)
	panic(string(all))

}
func (c *console) Panicf(format string, v ...interface{}) {
	params := c.genStdLogHelper(format, v...)
	all := outputHelper(format, true, c.flag, PREFIX_STD_PANIC, params...)
	panic(string(all))
}
func (c *console) Panicln(v ...interface{}) {
	params := c.genStdLogHelper("", v...)
	all := outputHelper("", true, c.flag, PREFIX_STD_PANIC, params...)
	panic(string(all))
}

func (c *console) Print(v ...interface{}) {
	params := c.genStdLogHelper("", v...)
	all := outputHelper("", false, c.flag, PREFIX_STD_PRINT, params...)
	c.logHelper(all)

}
func (c *console) Printf(format string, v ...interface{}) {
	params := c.genStdLogHelper(format, v...)
	all := outputHelper(format, false, c.flag, PREFIX_STD_PRINT, params...)
	c.logHelper(all)
}
func (c *console) Println(v ...interface{}) {
	params := c.genStdLogHelper("", v...)
	all := outputHelper("", true, c.flag, PREFIX_STD_PRINT, params...)
	c.logHelper(all)
}

func (c *console) Prefix() string {
	c.locker.Lock()
	defer c.locker.Unlock()
	return c.prefix
}

func (c *console) Flags() int {
	c.locker.Lock()
	defer c.locker.Unlock()
	return c.flag

}

func (c *console) SetFlags(flag int) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.flag = flag
}
func (c *console) SetPrefix(prefix string) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.prefix = prefix

}
