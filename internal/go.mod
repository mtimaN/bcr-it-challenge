module internal

go 1.23.9

replace db => ./database

require db v0.0.0-00010101000000-000000000000

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/gocql/gocql v1.7.0 // indirect
	github.com/golang/snappy v0.0.3 // indirect
	github.com/hailocab/go-hostpool v0.0.0-20160125115350-e80d13ce29ed // indirect
	golang.org/x/crypto v0.38.0 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
)
