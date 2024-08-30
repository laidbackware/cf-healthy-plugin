package sig_check

// Logger is used for outputting log-cache results and errors
type Logger interface {
	Fatalf(format string, args ...interface{})
	Printf(format string, args ...interface{})
}
