package gonyan

// LogLevel is used to define supported logging levels.
type LogLevel int

// Log level defitions:
//
//  * Debug
//  * Verbose
//  * Info
//  * Warning
//  * Error
//  * Fatal
//  * Panic
const (
	Debug   LogLevel = iota
	Verbose LogLevel = iota
	Info    LogLevel = iota
	Warning LogLevel = iota
	Error   LogLevel = iota
	Fatal   LogLevel = iota
	Panic   LogLevel = iota
)

// GetLevelLabel returns a string label for provided level.
func GetLevelLabel(level LogLevel) string {
	switch level {
	case Debug:
		return "Debug"
	case Verbose:
		return "Verbose"
	case Info:
		return "Info"
	case Warning:
		return "Warning"
	case Error:
		return "Error"
	case Fatal:
		return "Fatal"
	case Panic:
		return "Panic"
	default:
		return ""
	}
}
