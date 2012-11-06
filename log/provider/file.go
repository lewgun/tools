/*
This driver is about use plain text to save log. 

*/
package provider

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"sync"

//	"time"
)

const (
	FILE = "file"
)

func init() {
	register(FILE, newFile)
}

// create a logger wich was drvied by plain text
func newFile(arg *Arg) interface{} {
	if arg == nil {
		return nil
	}

	f := &file{}

	if err := f.init(arg); err != nil {
		return nil
	}
	return f

}

type file struct {
	//the global unique name for this logger
	uname string

	//the body destination, required
	path string

	//the body prefix, option
	prefix string

	//the body flag, option
	flag int

	//rolling or not, option
	isRolling bool

	//the max value in MB, option
	logSize int64

	f *os.File

	logger *log.Logger
	locker sync.Mutex

	index int64
}

func (f *file) initRolling() error {
	var err error
	defer func() {
		if err != nil {
			f.f.Close()
		}
	}()

	//open the logger 
	f.f, err = os.OpenFile(f.path, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			f.f, err = os.Create(f.path)

			//skip the index
			if err == nil {
				_, err = f.f.Seek(INDEX_IN_BYTE, os.SEEK_SET)
				f.index = INDEX_IN_BYTE
			}
		}
		return err
	}

	// seek  to the proper index.
	_, err = f.f.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	var (
		n int64
	)

	err = binary.Read(f.f, binary.LittleEndian, &n)
	if err != nil {
		return err
	}

	f.index = n

	stat, err := f.f.Stat()
	if err != nil {
		return err
	}

	// if the index is beyond the end. seek to 0
	if f.index >= stat.Size() {
		f.index = 0
	}

	//seek to the end for logging
	_, err = f.f.Seek(f.index, os.SEEK_SET)
	if err != nil {
		return err
	}
	return nil
}

func (f *file) init(arg *Arg) error {

	extras := arg.Extras

	if _, isExists := extras[PATH]; !isExists {
		return fmt.Errorf("The required parameter is missing\n")
	}

	f.uname = arg.Driver

	var (
		err error
	)

	for k, v := range extras {
		switch k {
		case PATH:
			f.path = v.(string)

		case PREFIX:
			f.prefix = v.(string)

		case IS_ROLLING:
			f.isRolling = v.(bool)

		case FLAG:
			f.flag = v.(int)

		case LOG_SIZE:
			f.logSize = int64(v.(int) << 20) // MB ->byte

		default:
			//placeholder
		}
	}

	//need rolling
	if f.isRolling {
		err = f.initRolling()
		if err != nil {
			return err
		}

	} else {
		//open the logger 
		f.f, err = os.OpenFile(f.path, os.O_CREATE|os.O_APPEND, os.ModePerm)
		if err != nil {
			return err
		}

	}

	f.logger = log.New(f.f, "", 0)

	return nil
}

func (f *file) logRollingHelper(format string, params ...interface{}) error {

	var (
		body string
		all  []byte
	)
	if format != "" {
		body = fmt.Sprintf(format, params[2:]...)

	} else {
		body = fmt.Sprint(params[2:]...)
	}

	all = formatHelper(params[0].(string), params[1].(int), f.flag, f.prefix, body)

	if all[len(all)-1] != '\n' {
		all = append(all, '\n')
	}

	size := int64(len(all))

	var (
		err       error
		index     int64
		totalSize int64
	)
	totalSize = f.index + size
	if totalSize > f.logSize {
		rhSize := f.logSize - f.index

		//right hand
		_, err = f.f.Seek(f.index, os.SEEK_SET)
		if err != nil {
			return err
		}
		//f.f.Write(all[:rhSize])
		f.logger.Print(string(all[:rhSize]))
		f.index = f.logSize

		//left hand 
		_, err = f.f.Seek(INDEX_IN_BYTE, os.SEEK_SET)
		if err != nil {
			return err
		}
		//f.f.Write(all[rhSize:])
		f.logger.Print(string(all[rhSize:]))
		index = totalSize - f.logSize + INDEX_IN_BYTE

	} else {
		//right hand
		_, err = f.f.Seek(f.index, os.SEEK_SET)
		if err != nil {
			return err
		}
		//	f.f.Write(all)
		f.logger.Print(string(all))
		index = totalSize

	}

	_, err = f.f.Seek(0, os.SEEK_SET)
	if err != nil {
		return err
	}

	err = binary.Write(f.f, binary.LittleEndian, index)
	if err == nil {
		f.index = index
	}
	return err

}
func (f *file) logHelper(format string, params ...interface{}) {

	if f.isRolling {
		if err := f.logRollingHelper(format, params...); err != nil {
			fmt.Println(err)
		}

	} else {

		var (
			body string
			all  []byte
		)

		if format != "" {
			body = fmt.Sprintf(format, params[2:]...)

		} else {
			body = fmt.Sprint(params[2:]...)
		}

		if all[len(all)-1] != '\n' {
			all = append(all, '\n')
		}

		all = formatHelper(params[0].(string), params[1].(int), f.flag, f.prefix, body)
		//f.f.Write(all)
		f.logger.Print(string(all))

	}

}

func (f *file) Tracef(format string, params ...interface{}) {
	f.logHelper(format, params...)
}

func (f *file) Debugf(format string, params ...interface{}) {
	f.logHelper(format, params...)
}
func (f *file) Infof(format string, params ...interface{}) {
	f.logHelper(format, params...)
}
func (f *file) Warnf(format string, params ...interface{}) {
	f.logHelper(format, params...)
}
func (f *file) Errorf(format string, params ...interface{}) {
	f.logHelper(format, params...)
}
func (f *file) Criticalf(format string, params ...interface{}) {
	f.logHelper(format, params...)
}

func (f *file) Release() {
	if f.f != nil {
		f.f.Close()
		f.f = nil
	}
}

func (f *file) Fatal(v ...interface{}) {
	f.logger.Fatal(v...)

}
func (f *file) Fatalf(format string, v ...interface{}) {
	f.logger.Fatalf(format, v...)
}
func (f *file) Fatalln(v ...interface{}) {
	f.logger.Fatalln(v...)
}
func (f *file) Flags() int {
	f.locker.Lock()
	defer f.locker.Unlock()
	return f.flag

}
func (f *file) Output(calldepth int, s string) error {
	return f.logger.Output(calldepth+1, s)

}
func (f *file) Panic(v ...interface{}) {
	f.logger.Panic(v...)

}
func (f *file) Panicf(format string, v ...interface{}) {
	f.logger.Panicf(format, v...)
}
func (f *file) Panicln(v ...interface{}) {
	f.logger.Panicln(v...)
}
func (f *file) Prefix() string {
	f.locker.Lock()
	defer f.locker.Unlock()
	return f.prefix
}
func (f *file) Print(v ...interface{}) {
	f.logger.Print(v...)
}
func (f *file) Printf(format string, v ...interface{}) {
	f.logger.Printf(format, v...)
}
func (f *file) Println(v ...interface{}) {
	f.logger.Println(v...)
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
