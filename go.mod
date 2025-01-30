module github.com/reyhardy/go-blog

go 1.23.3

require maragu.dev/gomponents v1.0.0 // direct

require (
	github.com/gocql/gocql v1.7.0
	github.com/scylladb/gocqlx/v3 v3.0.1
	github.com/segmentio/ksuid v1.0.4
	github.com/starfederation/datastar v1.0.0-beta.1
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.7.3

require (
	github.com/a-h/templ v0.3.819 // indirect
	github.com/delaneyj/gostar v0.8.0 // indirect
	github.com/goccy/go-json v0.10.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/igrmk/treemap/v2 v2.0.1 // indirect
	github.com/samber/lo v1.48.0 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/exp v0.0.0-20250106191152-7588d65b2ba8 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
