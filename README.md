# Zippy Zap LRU Cache
An LRU cache for go. It is safe for use by multiple goroutines concurrently.
It accomplishes this without any locks. Instead, it uses a background goroutine
with chans.

## Benchmarks
See lru_cache_test.go. These were the results on my 2017 macbook pro:
```
~/go/src/zippyzap$ go test -bench=.
goos: darwin
goarch: amd64
pkg: github.com/Nick-Anderssohn/zippyzap
BenchmarkLRU_Rand-8   	  536053	      2103 ns/op
--- BENCH: BenchmarkLRU_Rand-8
    lru_cache_test.go:110: hit: 0 miss: 1 ratio: 0.000000
    lru_cache_test.go:110: hit: 0 miss: 100 ratio: 0.000000
    lru_cache_test.go:110: hit: 1420 miss: 8580 ratio: 0.165501
    lru_cache_test.go:110: hit: 132399 miss: 403654 ratio: 0.328001
BenchmarkLRU_Freq-8   	  566115	      2040 ns/op
--- BENCH: BenchmarkLRU_Freq-8
    lru_cache_test.go:143: hit: 1 miss: 0 ratio: +Inf
    lru_cache_test.go:143: hit: 100 miss: 0 ratio: +Inf
    lru_cache_test.go:143: hit: 9808 miss: 192 ratio: 51.083333
    lru_cache_test.go:143: hit: 178565 miss: 387550 ratio: 0.460753
PASS
ok  	github.com/Nick-Anderssohn/zippyzap	2.460s
```
These are the results using [hashicorp's LRU lib](https://github.com/hashicorp/golang-lru)
instead of this lib. They use a normal impl with locks instead of chans:
```
~/go/src/golang-lru$ go test -bench=BenchmarkLRU
goos: darwin
goarch: amd64
pkg: github.com/hashicorp/golang-lru
BenchmarkLRU_Rand-8   	 3151484	       374 ns/op
--- BENCH: BenchmarkLRU_Rand-8
    lru_test.go:34: hit: 0 miss: 1 ratio: 0.000000
    lru_test.go:34: hit: 0 miss: 100 ratio: 0.000000
    lru_test.go:34: hit: 1426 miss: 8574 ratio: 0.166317
    lru_test.go:34: hit: 248966 miss: 751034 ratio: 0.331498
    lru_test.go:34: hit: 787296 miss: 2364188 ratio: 0.333009
BenchmarkLRU_Freq-8   	 3306544	       345 ns/op
--- BENCH: BenchmarkLRU_Freq-8
    lru_test.go:66: hit: 1 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 100 miss: 0 ratio: +Inf
    lru_test.go:66: hit: 9889 miss: 111 ratio: 89.090090
    lru_test.go:66: hit: 313616 miss: 686384 ratio: 0.456910
    lru_test.go:66: hit: 1016361 miss: 2290183 ratio: 0.443790
PASS
ok  	github.com/hashicorp/golang-lru	3.587s
```
Looks like a normal impl of LRU using locks instead of chans is faster; at least it
is when you test by hitting the cache sequentially in a single goroutine. It's probably
faster in a concurrent situation too, but I want to write a benchmark for that situation
as well.

## Should I use this in a production service?
No. Given the above results, a traditional implementation using
locks gives you better perf. I would like to test hitting the cache concurrently
though. I doubt the results would be any different, but I'm still curious.

## Is the name a reference to Pokemon?
Yup.
