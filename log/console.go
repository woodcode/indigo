package log
import(
	"os"
	"time"
)
//ConsoleProvider write log to console
type consoleProvider struct{
    lg       *logWriter
	Level    int  `json:"level"`
}
//NewConsoleProvider create consoleProvider returning as ProviderInterface.
func NewConsoleProvider() Logger{
    provider := &consoleProvider{
		lg:       newLogWriter(os.Stdout),
		Level:    LevelDebug,
	}
	return provider
}
// Init init console logger.
// jsonConfig like '{"level":LevelTrace}'.
func (c *consoleProvider) Init(jsonConfig string) error {
	return nil
}

// WriteMsg write message in console.
func (c *consoleProvider) Write(when time.Time, msg string, level int) error {
	c.lg.println(when, msg)
	return nil
}
// Destroy implementing method. empty.
func (c *consoleProvider) Destroy() {
}
func init(){
    Register("console", NewConsoleProvider)
}