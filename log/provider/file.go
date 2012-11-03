/*
This driver is about use plain text to save log. 

*/
package provider

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"
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

	//the output destination, required
	path string

	//the output prefix, option
	prefix string

	//rolling or not, option
	isRolling bool

	//the max value in MB, option
	logSize int64

	f      *os.File
	index  int64
	logger *log.Logger
}

func (f *file) initRolling() error {

	var err error
	//open the logger 
	f.f, err = os.OpenFile(f.path, os.O_RDWR, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			f.f, err = os.Create(f.path)

			//skip the index
			if err == nil {
				_, err = f.f.Seek(INDEX_IN_BYTE, os.SEEK_SET)
			}
		}
		return err
	}

	// seek  to the proper index.
	_, err = f.f.Seek(0, os.SEEK_SET)
	if err != nil {
		f.f.Close()
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
		f.f.Close()
		return err
	}

	// if the index is beyond the end. seek to 0
	if f.index >= stat.Size() {
		f.index = 0
	}

	//seek to the end for logging
	_, err = f.f.Seek(f.index, os.SEEK_SET)
	if err != nil {
		f.f.Close()
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

		case LOG_SIZE:
			f.logSize = v.(int64) << 20 // MB ->byte

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

	f.logger = log.New(f.f, f.prefix, DEFAULT_FLAG)
	if f.logger == nil {
		return fmt.Errorf("New logger is failed")
	}

	return nil
}

func (f *file) logRollingHelper(format string, params ...interface{}) error {

	var output string
	if format != "" {
		if format[len(format)-1:] != "\n" {
			format += "\n"
		}
		output = fmt.Sprintf(format, params...)

	} else {
		output = fmt.Sprintln(params...)
	}

	size := int64(len(output))

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

		f.logger.Print(output[:rhSize])

		f.index = f.logSize

		//left hand 
		_, err = f.f.Seek(INDEX_IN_BYTE, os.SEEK_SET)
		if err != nil {
			return err
		}

		f.logger.Print(output[rhSize:])

		index = totalSize - f.logSize + INDEX_IN_BYTE
	} else {

		//right hand
		_, err = f.f.Seek(f.index, os.SEEK_SET)
		if err != nil {
			return err
		}

		f.logger.Print(output)
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
		if format != "" {
			if format[len(format)-1:] != "\n" {
				format += "\n"
			}
			f.logger.Printf(format, params...)

		} else {
			f.logger.Println(params...)
		}

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
