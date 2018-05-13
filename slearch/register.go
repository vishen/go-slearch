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

func GetAllFormatters() []StructuredLogFormatter {
	formattersList := make([]StructuredLogFormatter, len(formatters))
	i := 0
	for _, f := range formatters {
		formattersList[i] = f
		i++
	}
	return formattersList
}

func getFormatter(key string) (StructuredLogFormatter, bool) {
	formattersMu.Lock()
	defer formattersMu.Unlock()

	slf, ok := formatters[key]
	return slf, ok
}
