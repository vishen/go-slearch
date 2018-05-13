package slearch

type StructuredLogMatchType int

const (
	StructuredLogMatchTypeAnd = StructuredLogMatchType(iota)
	StructuredLogMatchTypeOr
)

type KV struct {
	Key         string
	Value       string
	RegexString string
}

type Config struct {
	// Defines which 'StructuredLogFormatter' to use
	LogFormatterType string

	// Whether this is an AND or OR matching
	MatchType StructuredLogMatchType

	// Values to match on
	MatchOn []KV

	// Which keys to print for matching records
	PrintKeys []string

	// String to split the key on
	KeySplitString string

	// Will print all debug statements
	Verbose bool
}
