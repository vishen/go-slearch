package slearch

import "sync"

var (

	// Map of structured log formatter
	formatters   = map[string]StructuredLogFormatter{}
	formattersMu = sync.Mutex{}
)

func Register(key string, structuredLogFormatter StructuredLogFormatter) {
	formattersMu.Lock()
	formatters[key] = structuredLogFormatter
	formattersMu.Unlock()
}

func getFormatter(key string) (StructuredLogFormatter, bool) {
	formattersMu.Lock()
	defer formattersMu.Unlock()

	slf, ok := formatters[key]
	return slf, ok
}
