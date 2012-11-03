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

const (
	_ = 1 << iota
	Lnofile
	Llongfile  // full file name and line number: /a/b/c/d.go:23
	Lshortfile // final file name element and line number: d.go:23. overrides Llongfile
	Lwithfile  = Lshortfile | Llongfile
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
	Level    int
	FileLine string
	Format   string
	Params   []interface{}
}
