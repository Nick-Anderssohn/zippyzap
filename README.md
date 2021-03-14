# Zippy Zap LRU Cache
An LRU cache for go. It is safe for use by multiple goroutines concurrently.
It accomplishes this without any locks. Instead it uses a background goroutine
with chans.

## Benchmarks
coming soon.

## Should I use this in a production service?
Let me run benchmarks first then I'll tell ya. (coming soon).

## Is the name a reference to Pokemon?
Yup.