# go-redis-pool

go-redis-pool was designed to implement the read/write split in Redis master-slave mode, and easy way to sharding the data.

## Installation

go-redis-pool requires a Go version with [Modules](https://github.com/golang/go/wiki/Modules) support and uses import versioning. So please make sure to initialize a Go module before installing go-redis-pool:

```shell
go get github.com/gismeteo/redis-pool
```

## Quick Start

API documentation and examples are available via [godoc](https://godoc.org/github.com/bitleak/go-redis-pool/v2)

### Setup The Master-Slave Pool

```go
pool, err := pool.NewHA(&pool.HAConfig{
    Master: "127.0.0.1:6379",
    Slaves: []string{
        "127.0.0.1:6380",
        "127.0.0.1:6381",
    },
    Password: "", // set master password
    ReadonlyPassword: "", // use password if no set
})

pool.Set(ctx, "foo", "bar", 0)
```

The read-only commands would go throught slaves and write commands would into the master instance. We use the Round-Robin as default when determing which slave to serve the readonly request, and currently supports:

* RoundRobin (default)
* Random
* Weight

For example, we change the distribution type to `Weight`:

```go
pool, err := pool.NewHA(&pool.HAConfig{
    Master: "127.0.0.1:6379",
    Slaves: []string{
        "127.0.0.1:6380",  // default weight is 100 if missing
        "127.0.0.1:6381:200", // weight is 200
        "127.0.0.1:6382:300", // weigght is 300
    },
    PollType: pool.PollByWeight,
})
```

The first slave would serve 1/6 reqeusts, and second slave would serve 2/6, last one would serve 3/6. 

##### Auto Eject The Failure Host 

```go
pool, err := pool.NewHA(&pool.HAConfig{
    Master: "127.0.0.1:6379",
    Slaves: []string{
        "127.0.0.1:6380",  // default weight is 100 if missing
        "127.0.0.1:6381:200", // weight is 200
        "127.0.0.1:6382:300", // weigght is 300
    },
    AutoEjectHost: true,
    ServerFailureLimit: 3,
    ServerRetryTimeout: 5 * time.Second,
    MinServerNum: 2,
})
```

The pool would evict the host if reached `ServerFailureLimit` times of failure and retry the host after `ServerRetryTimeout`. The
`MinServerNum` was used to avoid evicting too many and would overrun other alive servers. 

