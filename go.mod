module github.com/reyhardy/go-blog

go 1.23.3

require maragu.dev/gomponents v1.1.0 // direct

require (
	github.com/gocql/gocql v1.7.0
	github.com/labstack/echo/v4 v4.13.3
	github.com/scylladb/gocqlx/v3 v3.0.1
	github.com/segmentio/ksuid v1.0.4
	github.com/starfederation/datastar v1.0.0-beta.11
)

replace github.com/gocql/gocql => github.com/scylladb/gocql v1.7.3

require (
	github.com/CAFxX/httpcompression v0.0.9 // indirect
	github.com/a-h/templ v0.3.857 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/delaneyj/gostar v0.8.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	github.com/igrmk/treemap/v2 v2.0.1 // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/labstack/gommon v0.4.2 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/samber/lo v1.49.1 // indirect
	github.com/scylladb/go-reflectx v1.0.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasttemplate v1.2.2 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/exp v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/time v0.11.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
