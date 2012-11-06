package provider

import (
	"fmt"
	"log"
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

func newConsole(arg *Arg) interface{} {
	if arg == nil {
		return nil
	}

	c := &console{}

	if err := c.init(arg); err != nil {
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

func (c *console) init(arg *Arg) error {

	extras := arg.Extras

	c.uname = arg.Driver

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

func (c *console) Tracef(format string, params ...interface{}) {
	c.logHelper(format, params...)
}

func (c *console) Debugf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}
func (c *console) Infof(format string, params ...interface{}) {
	c.logHelper(format, params...)
}
func (c *console) Warnf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}
func (c *console) Errorf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}

func (c *console) Criticalf(format string, params ...interface{}) {
	c.logHelper(format, params...)
}

func (c *console) logHelper(format string, params ...interface{}) {

	var (
		body string
		all  []byte
	)

	if format != "" {
		body = fmt.Sprintf(format, params[2:]...)

	} else {
		body = fmt.Sprint(params[2:]...)
	}
	c.locker.Lock()
	defer c.locker.Unlock()
	all = formatHelper(params[0].(string), params[1].(int), c.flag, c.prefix, body)
	if all[len(all)-1] != '\n' {
		all = append(all, '\n')
	}
	c.logger.Print(string(all))
	//c.writer.Write(all)

}

func (c *console) Release() {

}

func (c *console) Fatal(v ...interface{}) {
	c.logger.Fatal(v...)

}
func (c *console) Fatalf(format string, v ...interface{}) {
	c.logger.Fatalf(format, v...)
}
func (c *console) Fatalln(v ...interface{}) {
	c.logger.Fatalln(v...)
}
func (c *console) Flags() int {
	c.locker.Lock()
	defer c.locker.Unlock()
	return c.flag

}
func (c *console) Output(calldepth int, s string) error {
	return c.logger.Output(calldepth+1, s)

}
func (c *console) Panic(v ...interface{}) {
	c.logger.Panic(v...)

}
func (c *console) Panicf(format string, v ...interface{}) {
	c.logger.Panicf(format, v...)
}
func (c *console) Panicln(v ...interface{}) {
	c.logger.Panicln(v...)
}
func (c *console) Prefix() string {
	c.locker.Lock()
	defer c.locker.Unlock()
	return c.prefix
}
func (c *console) Print(v ...interface{}) {
	c.logger.Print(v...)
}
func (c *console) Printf(format string, v ...interface{}) {
	c.logger.Printf(format, v...)
}
func (c *console) Println(v ...interface{}) {
	c.logger.Println(v...)
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
