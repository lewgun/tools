package manager

import (
	"fmt"
	"log/provider"
	//	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (

	//日志记录器容器
	drivers = map[string]Logger{
	//	provider.DEFAULT: provider.Default,
	}
	globalMgr *manager
)

func init() {

	arg := provider.Arg{
		Driver: provider.CONSOLE,
		Extras: map[string]interface{}{
			provider.FLAG: provider.Lshortfile | provider.LstdFlags,
		},
	}

	drv := provider.New(provider.CONSOLE, &arg).(Logger)

	register(provider.CONSOLE, drv)

	//默认管理器
	defMgr := manager{
		defLogger:   drivers[provider.DEFAULT],
		curLogger:   drivers[provider.DEFAULT],
		level:       LEVEL_DEFAULT,                          //默认为>=warn 的消息才输出
		chanLogArgs: make(chan *LogArgs, DEFAULT_CHAN_SIZE), //输出日志所需要的信息
		stopFlag:    make(chan struct{}),
	}

	globalMgr = &defMgr
}

//日志管理器
type manager struct {

	//默认记录器
	defLogger Logger

	//当前记录器
	curLogger Logger

	chanLogArgs chan *LogArgs

	//日志级别
	level int

	//关闭记录器标志
	stopFlag chan struct{}

	locker sync.Mutex
}

func (m *manager) logProcHelper(arg *LogArgs) error {

	level := arg.Level

	//当日志级别异常时，默认为default
	if level < LEVEL_TRACE || level > LEVEL_ALL {
		level = LEVEL_DEFAULT
	}

	var format string

	if arg.Format != "" {
		format = logPrefixs[level] + " " + arg.Format

		arg.Params = append([]interface{}{arg.File, arg.Line}, arg.Params...)
	} else {

		arg.Params = append([]interface{}{arg.File, arg.Line, logPrefixs[level]}, arg.Params...)

	}

	switch level {
	case LEVEL_TRACE:
		m.curLogger.Tracef(format, arg.Params...)

	case LEVEL_DEBUG:
		m.curLogger.Debugf(format, arg.Params...)

	case LEVEL_INFO:
		m.curLogger.Infof(format, arg.Params...)

	case LEVEL_WARN:
		m.curLogger.Warnf(format, arg.Params...)

	case LEVEL_ERROR:
		m.curLogger.Errorf(format, arg.Params...)

	case LEVEL_CRITICAL:
		m.curLogger.Criticalf(format, arg.Params...)
	default:
		panic("never go here")
	}
	return nil
}

func (m *manager) logProc() {
FOOR_FLAG:
	for {
		select {
		case arg, ok := <-m.chanLogArgs:
			//已经清空所有消息,可以退出
			if !ok {
				m.curLogger.Release()
				m.stopFlag <- struct{}{}
				break FOOR_FLAG
			}

			m.logProcHelper(arg)

		case <-m.stopFlag:

			//关闭信道
			close(m.chanLogArgs)

		default:
			//placeholder
			time.Sleep(time.Second)

		} // end of select
	} // end of for 
}

//启动日志记录器
func (m *manager) run() {
	go m.logProc()
}

func (m *manager) stop() {
	m.stopFlag <- struct{}{}
	<-m.stopFlag
}

func (m *manager) current() Logger {
	m.locker.Lock()
	defer m.locker.Unlock()
	return m.curLogger
}
func (m *manager) change(driver string, level int) error {

	m.locker.Lock()
	defer m.locker.Unlock()
	// get the driver which named 'driver' from the driver db.
	logger := get(driver)
	if logger == nil {
		return fmt.Errorf("The driver hasn't been registered")
	}

	// stop the old one
	m.stop()

	// change to the new one 
	m.curLogger = logger
	m.level = level
	m.chanLogArgs = make(chan *LogArgs, DEFAULT_CHAN_SIZE)

	// boot the new one 
	m.run()
	return nil

}

// Register a new driver 
func register(driver string, logger Logger) error {

	lower := strings.ToLower(driver)

	if _, isExists := drivers[lower]; isExists {
		return fmt.Errorf("The driver: %s is exists\n", driver)
	}
	drivers[lower] = logger
	return nil
}

// Get the driver which named 'driver'
func get(driver string) Logger {
	lower := strings.ToLower(driver)
	if logger, isExists := drivers[lower]; isExists {
		return logger
	}
	return nil

}

//打开一个合适的日志记录器
func Open(typ string, level int, arg *provider.Arg) error {

	//设置默认日志级别
	if level < LEVEL_TRACE || level > LEVEL_ALL {
		level = LEVEL_DEFAULT
	}

	if _, isExists := drivers[arg.Driver]; isExists {

		return fmt.Errorf("The driver with the name: %s is existed\n", arg.Driver)
	}

	var drv Logger
	logger := provider.New(typ, arg)
	if logger == nil {
		return fmt.Errorf("create the new driver with name: %s is failed\n", arg.Driver)
	}

	drv = logger.(Logger)
	register(arg.Driver, drv)
	Change(arg.Driver, level)
	return nil
}

//开启日志记录器
func Run() {
	globalMgr.run()
}

//关闭日志记录器
func Stop() {
	globalMgr.stop()
}

//change the current logger to the one which named 'driver'
func Change(driver string, level int) error {
	return globalMgr.change(driver, level)

}

func Current() Logger {
	return globalMgr.current()
}
func logHelper(level int, format string, params ...interface{}) {

	if level&globalMgr.level == 0 {
		return // need not ouput 
	}

	var (
		file string
		line int = -1
	)

	_, file, line, _ = runtime.Caller(2)

	globalMgr.chanLogArgs <- &LogArgs{
		Level:  level,
		File:   file,
		Line:   line,
		Format: format,
		Params: params,
	}
}

//trace 
func Tracef(format string, params ...interface{}) {
	logHelper(LEVEL_TRACE, format, params...)
}
func Trace(params ...interface{}) {
	logHelper(LEVEL_TRACE, "", params...)
}

//debug 
func Debugf(format string, params ...interface{}) {
	logHelper(LEVEL_DEBUG, format, params...)
}
func Debug(params ...interface{}) {
	logHelper(LEVEL_DEBUG, "", params...)
}

//info 
func Infof(format string, params ...interface{}) {
	logHelper(LEVEL_INFO, format, params...)
}
func Info(params ...interface{}) {
	logHelper(LEVEL_INFO, "", params...)
}

//warn 
func Warnf(format string, params ...interface{}) {
	logHelper(LEVEL_WARN, format, params...)
}
func Warn(params ...interface{}) {
	logHelper(LEVEL_WARN, "", params...)
}

//error 
func Errorf(format string, params ...interface{}) {
	logHelper(LEVEL_ERROR, format, params...)
}
func Error(params ...interface{}) {
	logHelper(LEVEL_ERROR, "", params...)
}

//info 
func Criticalf(format string, params ...interface{}) {
	logHelper(LEVEL_CRITICAL, format, params...)
}
func Critical(params ...interface{}) {
	logHelper(LEVEL_CRITICAL, "", params...)
}
