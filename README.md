# Structured Logs Search (slearch)
This is a simple  utility to search through structured JSON logs that come out of your logging systems; eg: `kubectl logs` or `docker logs`.

Slearch will read structured JSON logs via `stdin` that are separated newlines and allow you to filter out exact or regex matches per log line.

## Usage
```
$ go-slearch --help
Read stuctured logs from STDIN and filter out lines based on exact and regex matches. Currently only supports JSON logs.

Usage:
  go-slearch [flags]

Flags:
  -h, --help                   help for go-slearch
  -d, --key_delimiter string   the string to split the key on for complex key queries (default ".")
  -m, --match strings          key and value to match on. eg: label1=value1
  -p, --print_keys strings     keys to print if a match is found
  -r, --regexp strings         key and value to regex match on. eg: label1=value*
  -s, --search_type string     the search type to use: 'and' or 'or' (default "and")

```

## Example
```
$ make build
$ cat example.txt
{"severity": "info", "key1": "value1", "key2": "value2"}
{"severity": "debug", "key1": "value12", "key2": "value2", "key3": "value3"}
{"severity": "info", "key1": "value13", "key3": "value3"}
{"severity": "info", "key1": "value14", "key4": "value4"}
{"severity": "info", "key2": "value2", "key5": 5}
{"key2": "value2", "complexKey": {"key1": "value1", "complexKey1": {"key3": "value3"}}}

# Exact match search with -m
$ cat example.txt | go-slearch -m key1="value1"
{"severity": "info", "key1": "value1", "key2": "value2"}

# Regex match search with -r
$ cat example.txt | go-slearch -r key1="value1"
{"severity": "debug", "key1": "value12", "key2": "value2", "key3": "value3"}
{"severity": "info", "key1": "value13", "key3": "value3"}
{"severity": "info", "key1": "value14", "key4": "value4"}
{"severity": "info", "key1": "value1", "key2": "value2"}


# Exact match on multiple keys, but because by default it queries with an AND we get no results
$ cat example.txt | go-slearch -m key1=value1,key3=value3

# Exact match on multiple key WITH OR
$ cat example.txt | go-slearch -m key1=value1,key3=value3 -s or
{"severity": "debug", "key1": "value12", "key2": "value2", "key3": "value3"}
{"severity": "info", "key1": "value1", "key2": "value2"}
{"severity": "info", "key1": "value13", "key3": "value3"}

# Search on a complex ket
$ cat example.txt | go-slearch -m complexKey.key1=value1
{"key2": "value2", "complexKey": {"key1": "value1", "complexKey1": {"key3": "value3"}}}

# Only print certain keys
$ cat example.txt | go-slearch -m key1=value1,key3=value3 -s or -p severity,key1,key2
{"severity":"debug", "key1":"value12", "key2":"value2"}
{"severity":"info", "key1":"value1", "key2":"value2"}
{"severity":"info", "key1":"value13"}
```

## TODO
```
- Make it work for structured text format
- Tests
- Better controlled concurrency
```
