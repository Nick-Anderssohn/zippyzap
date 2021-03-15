# Zippy Zap LRU Cache
An LRU cache for go. It is safe for use by multiple goroutines concurrently.
It accomplishes this without any locks. Instead it uses a background goroutine
with chans.

## Benchmarks
See lru_cache_test.go. These were the results on my 2017 macbook pro
```
goos: darwin
goarch: amd64
pkg: github.com/Nick-Anderssohn/zippyzap
BenchmarkLRUCache_Put_SameInput-8        	  908406	      1341 ns/op
BenchmarkLRUCache_Put_RandomInput-8      	 1000000	      1076 ns/op
BenchmarkLRUCache_Put_500_Concurrent-8   	    1410	    820978 ns/op
PASS
ok  	github.com/Nick-Anderssohn/zippyzap	3.797s
```
Planning on writing these same benchmarks for
[hashicorp's LRU lib](https://github.com/hashicorp/golang-lru). I'll
add the results here once I have done that.

## Should I use this in a production service?
Initial results seen above are promising, but the benchmarks are pretty bare-bones.
It needs to be compared to another lib to see how it stacks up.

## Is the name a reference to Pokemon?
Yup.