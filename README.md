# Structured Logs Search (slearch)
This is a simple  utility to search through structured JSON or text logs that come out of your logging systems; eg: `kubectl logs` or `docker logs`.

Slearch will read structured JSON or text logs via `stdin` that are separated by newlines and allow you to filter out exact or regex matches per log line. It will print the matching results in the same order they are read.

By default it will match queries as `AND`, meaning it will only show results where all queries for a line match. This can be changed to an `OR` query with the `-s` flag.

Matched queries can either be exact with `-m` or by regex (using golang stdlib regexp package) with `-r`.

It will attempt to autodetect what format the log line is in, and match based on that format. If you want to specify a particular format you can with the `-t` argument.

## Installing
```
$ go get -u github.com/vishen/go-slearch
```

## Installing from source
```
$ make build
```

## Usage
```
Read stuctured logs from STDIN and filter out lines based on exact and regex matches. Currently only supports JSON and text logs.

Usage:
  go-slearch [flags]

Flags:
  -h, --help                   help for go-slearch
  -d, --key_delimiter string   the string to split the key on for complex key queries
  -m, --match strings          key and value to match on. eg: label1=value1
  -p, --print_keys strings     keys to print if a match is found
  -r, --regexp strings         key and value to regex match on. eg: label1=value*
  -s, --search_type string     the search type to use: 'and' or 'or' (default "and")
  -t, --type string            the log type to use: 'json' or 'text'. If unspecified it will attempt to use all log types
  -v, --verbose                verbose output
```

## Example
```
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

# Search on a complex key, the -d is the delimiter to use for a complex key, which is '.' in this case
$ cat example.txt | go-slearch -m complexKey.key1=value1 -d .
{"key2": "value2", "complexKey": {"key1": "value1", "complexKey1": {"key3": "value3"}}}

# Only print certain keys
$ cat example.txt | go-slearch -m key1=value1,key3=value3 -s or -p severity,key1,key2
{"severity":"debug", "key1":"value12", "key2":"value2"}
{"severity":"info", "key1":"value1", "key2":"value2"}
{"severity":"info", "key1":"value13"}
```

## TODO
```
- Bold what was matched
- Sort based on keys
- Tests and documentation and examples
    - Document regex case insensitive match ("(?i)" a the start of the regex)
- Add new formatters:
    - Add log formatter that splits stuff by delimeter, and you can the select which element(s) via $n
    - Add googles glog format
- Better controlled concurrency
- Ignore parts of the line that don't match a specified format?
- Add integration support (might not be worth it if we have to pull in the client libraries for these)?
    - kubernetes
    - docker
    - stackdriver logs
```
