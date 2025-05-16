package sflogger

// A Level is a logging priority. Higher levels are more important.
type Level uint

const (
	DebugLevel Level = iota
	InfoLevel
	WarningLevel
	ErrorLevel
	FatalLevel
)

var levelNames = map[Level]string{
	DebugLevel:   "DEBUG",
	InfoLevel:    "INFO",
	WarningLevel: "WARNING",
	ErrorLevel:   "ERROR",
	FatalLevel:   "FATAL",
}

var stringToLevel = map[string]Level{
	"DEBUG":   DebugLevel,
	"INFO":    InfoLevel,
	"WARNING": WarningLevel,
	"ERROR":   ErrorLevel,
	"FATAL":   FatalLevel,
}

// Convert the Level to a string.
func (level Level) String() string {
	if name, exists := levelNames[level]; exists {
		return name
	} else {
		return "UNKNOWN"
	}
}

// FromString converts a string to a Level. If the string doesn't match
// any known level, it returns InfoLevel and false.
func FromString(levelStr string) (Level, bool) {
	level, exists := stringToLevel[levelStr]
	return level, exists
}

// UnmarshalText implements encoding.TextUnmarshaler to allow Levels
// to be read from configuration files.
func (l *Level) UnmarshalText(text []byte) error {
	level, ok := stringToLevel[string(text)]
	if !ok {
		*l = InfoLevel // Default to InfoLevel if not found
		return nil
	}
	*l = level
	return nil
}
