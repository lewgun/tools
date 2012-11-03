package provider

import (
	"log"
)

//默认配置项
const (
	DEFAULT       = CONSOLE
	DEFAULT_FLAG  = log.LstdFlags
	INDEX_IN_BYTE = 8
)

const (
	PATH       = "path"
	PREFIX     = "prefix"
	IS_ROLLING = "is_rolling"
	LOG_SIZE   = "log_size"
)

const (
	DEFAULT_FORMAT = 0
	WITH_FORMAT    = 1
)

//记录器配置参数
type Arg struct {

	//将定义的driver的名称，全局唯一
	Driver string

	//此参数随driver不同，而内容不同，由driver自身处理
	Extras map[string]interface{}
}
