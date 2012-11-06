package manager

//日志级别
const (
	LEVEL_TRACE = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_CRITICAL

	//默认日志级别为
	LEVEL_DEFAULT = LEVEL_WARN
)

//简单配置项
const (
	DEFAULT_CHAN_SIZE = 10
)

var (
	logPrefixs = []string{
		"[TRACE]",
		"[DEBUG]",
		"[INFO]",
		"[WARN]",
		"[ERROR]",
		"[CRITICAL]",
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
