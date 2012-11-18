/*
This driver is about use plain text to save log. 

*/
package provider

import (
	//"encoding/binary"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"

//	"time"
)

const (
	FILE     = "file"
	SPLITTER = "."
)

func init() {
	register(FILE, newFile)
}

// create a logger wich was drvied by plain text
func newFile(arg *Arg, prefix string) interface{} {
	if arg == nil {
		return nil
	}

	f := &file{}

	if err := f.init(arg, prefix); err != nil {
		return nil
	}
	return f

}

type file struct {
	//the global unique name for this logger
	uname string

	//the body destination, required
	originPath string

	//the body prefix, option
	prefix string

	//the body flag, option
	flag int

	//write type, ption
	typ int

	//the max value in MB, option
	logSize int64

	logCount int

	f *os.File

	logger *log.Logger
	locker sync.Mutex
}

func (f *file) init(arg *Arg, prefix string) error {

	extras := arg.Extras

	if _, isExists := extras[PATH]; !isExists {
		return fmt.Errorf("The required parameter is missing\n")
	}

	f.uname = arg.Driver
	f.prefix = prefix

	for k, v := range extras {
		switch k {
		case PATH:
			f.originPath = v.(string)

		case PREFIX:
			f.prefix = v.(string)

		case FLAG:
			f.flag = v.(int)

		case LOG_SIZE:
			f.logSize = int64(v.(int) << 20) // MB ->byte

		case LOG_COUNT:
			f.logCount = v.(int)

		case WRITE_TYPE:
			f.typ = v.(int)
			if f.typ != SINGLE_APPEND && f.typ != MULTI_APPEND /*&& f.typ != ROLLING*/ {
				f.typ = SINGLE_APPEND
			}

		default:
			//placeholder
		}
	}

	return f.openLogger(true)

}

func (f *file) openLogger(isInit bool) error {

	var (
		err   error
		name2 string

		baseOrigPath = filepath.Base(f.originPath)
	)

	dir, err := os.Open(filepath.Dir(f.originPath))
	if err != nil {
		return err
	}

	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	//找取文件中最近被更新的日志文件
	var logFiles []string
	for _, fi := range fileInfos {
		name2 = fi.Name()

		//均为日志文件
		if strings.Contains(strings.ToLower(name2), baseOrigPath) {
			logFiles = append(logFiles, name2)
		}
	}

	sort.Sort(sort.StringSlice(logFiles))
	if !isInit {
		//取得最老一个日志文件的后缀
		oldestLogFile := logFiles[len(logFiles)-1]

		var postfix int
		if oldestLogFile != baseOrigPath {
			postfix, err = strconv.Atoi(oldestLogFile[len(baseOrigPath)+1:])
			if err != nil {
				return err
			}
		}

		//删除最后一个文件
		if postfix >= f.logCount {
			os.Remove(oldestLogFile)
			logFiles = logFiles[:len(logFiles)-1]
		}

		//文件重命名
		for i := postfix + 1; i >= 1; i-- {
			//for i := postfix; i >= 1; i-- {

			var oldName string
			if i-1 == 0 {
				oldName = f.originPath
			} else {
				oldName = fmt.Sprintf("%s.%d", f.originPath, i-1)
			}

			newName := fmt.Sprintf("%s.%d", f.originPath, i)
			os.Rename(oldName, newName)
		}

	}

	f.f, err = os.OpenFile(f.originPath, os.O_CREATE|os.O_APPEND, os.ModePerm)

	if err != nil {
		//panic(err)
		return err
	}

	f.logger = log.New(f.f, "", 0)

	return nil

}

func (f *file) logMultiAppend(data string) {
	stats, err := f.f.Stat()
	if err != nil {
		return
	}
	size := stats.Size()

	if size+int64(len(data)) >= f.logSize {

		//关闭旧文件
		f.f.Close()
		f.openLogger(false)
	}

	//常规操作
	f.logAppend(data)

}

func (f *file) logAppend(data string) {

	//f.f.Write(all)
	f.logger.Print(data)
}

func (f *file) logRolling(data string) {
	//placeholder
}

func (f *file) logHelper(all []byte) {

	switch f.typ {
	case SINGLE_APPEND:
		f.logAppend(string(all))

	case MULTI_APPEND:
		f.logMultiAppend(string(all))

	//case ROLLING:
	//	f.logAppend(string(all))
	default:
		f.logAppend(string(all))
	}

}

func (f *file) genStdLogHelper(format string, params ...interface{}) []interface{} {

	_, file, line, _ := runtime.Caller(2)
	if format != "" {
		format = f.prefix + " " + format

		params = append([]interface{}{file, line}, params...)
	} else {
		params = append([]interface{}{file, line, f.prefix}, params...)
	}
	return params

}
func (f *file) Tracef(format string, params ...interface{}) {

	all := outputHelper(format, false, f.flag, f.prefix, params...)
	f.logHelper(all)
}

func (f *file) Debugf(format string, params ...interface{}) {
	all := outputHelper(format, false, f.flag, f.prefix, params...)
	f.logHelper(all)
}
func (f *file) Infof(format string, params ...interface{}) {
	all := outputHelper(format, false, f.flag, f.prefix, params...)
	f.logHelper(all)
}
func (f *file) Warnf(format string, params ...interface{}) {
	all := outputHelper(format, false, f.flag, f.prefix, params...)
	f.logHelper(all)
}
func (f *file) Errorf(format string, params ...interface{}) {
	all := outputHelper(format, false, f.flag, f.prefix, params...)
	f.logHelper(all)
}
func (f *file) Criticalf(format string, params ...interface{}) {
	all := outputHelper(format, false, f.flag, f.prefix, params...)
	f.logHelper(all)
}

func (f *file) Release() {
	if f.f != nil {
		f.f.Close()
		f.f = nil
	}
}

//standard log methods
func (f *file) Fatal(v ...interface{}) {
	params := f.genStdLogHelper("", v...)

	all := outputHelper("", false, f.flag, PREFIX_STD_FATAL, params...)
	f.logHelper(all)
	os.Exit(1)

}
func (f *file) Fatalf(format string, v ...interface{}) {
	params := f.genStdLogHelper(format, v...)
	all := outputHelper(format, false, f.flag, PREFIX_STD_FATAL, params...)
	f.logHelper(all)
	os.Exit(1)
}
func (f *file) Fatalln(v ...interface{}) {
	params := f.genStdLogHelper("", v...)
	all := outputHelper("", true, f.flag, PREFIX_STD_FATAL, params...)
	f.logHelper(all)
	os.Exit(1)
}

func (f *file) Output(calldepth int, s string) error {
	return f.logger.Output(calldepth+1, s)

}
func (f *file) Panic(v ...interface{}) {
	params := f.genStdLogHelper("", v...)
	all := outputHelper("", false, f.flag, PREFIX_STD_PANIC, params...)
	panic(string(all))

}
func (f *file) Panicf(format string, v ...interface{}) {
	params := f.genStdLogHelper(format, v...)
	all := outputHelper(format, true, f.flag, PREFIX_STD_PANIC, params...)
	panic(string(all))
}
func (f *file) Panicln(v ...interface{}) {
	params := f.genStdLogHelper("", v...)
	all := outputHelper("", true, f.flag, PREFIX_STD_PANIC, params...)
	panic(string(all))
}

func (f *file) Print(v ...interface{}) {
	params := f.genStdLogHelper("", v...)
	all := outputHelper("", false, f.flag, PREFIX_STD_PRINT, params...)
	f.logHelper(all)

}
func (f *file) Printf(format string, v ...interface{}) {
	params := f.genStdLogHelper(format, v...)
	all := outputHelper(format, false, f.flag, PREFIX_STD_PRINT, params...)
	f.logHelper(all)
}
func (f *file) Println(v ...interface{}) {
	params := f.genStdLogHelper("", v...)
	all := outputHelper("", true, f.flag, PREFIX_STD_PRINT, params...)
	f.logHelper(all)
}

func (f *file) Prefix() string {
	f.locker.Lock()
	defer f.locker.Unlock()
	return f.prefix
}

func (f *file) Flags() int {
	f.locker.Lock()
	defer f.locker.Unlock()
	return f.flag

}

func (f *file) SetFlags(flag int) {
	f.locker.Lock()
	defer f.locker.Unlock()
	f.flag = flag
}
func (f *file) SetPrefix(prefix string) {
	f.locker.Lock()
	defer f.locker.Unlock()
	f.prefix = prefix

}
