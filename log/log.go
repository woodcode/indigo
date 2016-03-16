package log
import(
    "time"
    "fmt"
    "os"
    "runtime"
    "path"
    "strconv"
)
// RFC5424 log message levels.
const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)
type providerType func() Logger
// Logger defines the behavior of a log provider.
type Logger interface{
    Init(config string) error
    Write(t time.Time, msg string, level int) error
    Destroy()
}

var providers = make(map[string]providerType)
// Register makes a log provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, provider providerType) {
	if provider == nil {
		panic("log: Register provider is nil")
	}
	if _, dup := providers[name]; dup {
		panic("log: Register called twice for provider " + name)
	}
	providers[name] = provider
}

type loggerProvider struct {
	Logger
	name string
}
// IndigoLogger is default logger in indigo application.
// it can contain several providers and log message into all providers.
type IndigoLogger struct{
    level               int
	enableFuncCallDepth bool
	loggerFuncCallDepth int
    provider            *loggerProvider
}
// NewLogger returns a new IndigoLogger.
// provider a given logger provider name
// config need to be correct JSON as string: {"interval":360}.
func NewLogger(provider string, config string) (*IndigoLogger, error){
    logger:=new(IndigoLogger)
    logger.level = LevelError
    logger.enableFuncCallDepth = true
    logger.loggerFuncCallDepth = 2
    log, ok := providers[provider]
	if !ok {
		return nil, fmt.Errorf("logs: unknown provider %q (forgotten Register?)", provider)
	}
    lg := log()
    err := lg.Init(config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "log.NewLogger: "+err.Error())
		return nil, err
	}
    logger.provider = &loggerProvider{Logger:lg, name:provider}
    
    return logger, nil
}

func (logger *IndigoLogger) writeToLoggers(when time.Time, msg string, level int) {
		err := logger.provider.Logger.Write(when, msg, level)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to WriteMsg to provider:%v,error:%v\n", logger.provider.name, err)
		}
	
}
func (logger *IndigoLogger) writeMsg(logLevel int, msg string) error {
	when := time.Now()
	if logger.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(logger.loggerFuncCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		msg = "[" + filename + ":" + strconv.FormatInt(int64(line), 10) + "]" + msg
	}
	logger.writeToLoggers(when, msg, logLevel)
	return nil
}
// Close close logger, flush all chan data and destroy all provider in IndigoLogger.
func (logger *IndigoLogger) Close() {
	 logger.provider.Destroy()
     logger.provider = nil
}

// Error Log ERROR level message.
func (logger *IndigoLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf("[E] "+format, v...)
	logger.writeMsg(LevelError, msg)
}