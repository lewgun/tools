package manager

//日志级别
const (
	LEVEL_TRACE = 1 << iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_CRITICAL
	LEVEL_ALL = LEVEL_TRACE | LEVEL_DEBUG | LEVEL_INFO | LEVEL_WARN | LEVEL_ERROR | LEVEL_CRITICAL

	//默认日志级别为
	LEVEL_DEFAULT = LEVEL_WARN
)

//简单配置项
const (
	DEFAULT_CHAN_SIZE = 10
)

var (
	logPrefixs = map[int]string{
		LEVEL_TRACE:    "[TRACE]",
		LEVEL_DEBUG:    "[DEBUG]",
		LEVEL_INFO:     "[INFO]",
		LEVEL_WARN:     "[WARN]",
		LEVEL_ERROR:    "[ERROR]",
		LEVEL_CRITICAL: "[CRITICAL]",
	}
)

//输出参数
type LogArgs struct {
	Level  int
	File   string
	Line   int
	Format string
	Params []interface{}
}
