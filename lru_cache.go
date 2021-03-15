package zippyzap

// LRUCache is a least-recently-used cache.
// It is safe for use by multiple goroutines
// concurrently. It accomplishes this without
// any locks. Instead, it has a background goroutine
// consuming chans.
type LRUCache struct {
	maxSize int

	oldest   *node
	youngest *node

	lookupMap map[interface{}]*node

	operationChan chan operation
	shutdownChan  chan chan bool
}

type operation struct {
	opType opType

	key interface{}
	val interface{}

	resultChan chan result
}

type opType int

const (
	opTypePut opType = iota
	opTypeGet
	opTypeRemove
)

type result struct {
	found bool

	target *node
}

// CreateAndStartLRUCache creates a new cache and spins up
// a goroutine that is used as part of the internal implementation
// of the cache. You can stop the goroutine via the Shutdown func.
func CreateAndStartLRUCache(maxSize int) *LRUCache {
	cache := LRUCache{
		maxSize:       maxSize,
		lookupMap:     map[interface{}]*node{},
		operationChan: make(chan operation, maxSize),
		shutdownChan:  make(chan chan bool, 1),
	}

	go cache.run()

	return &cache
}

// Shutdown stops the internal goroutine of the cache.
// Do not re-use a cache after calling Shutdown. Be sure
// to nil out any references you have to the cache so that
// it gets garbage collected.
func (l *LRUCache) Shutdown() {
	waitForShutdown := make(chan bool)

	l.shutdownChan <- waitForShutdown

	<-waitForShutdown
}

// ContainsKey returns true if the key exists in the cache.
// Otherwise, it returns false.
func (l *LRUCache) ContainsKey(key interface{}) bool {
	_, containsKey := l.lookupMap[key]
	return containsKey
}

// Get retrieves the value found for the provided key. It will
// return false as the second return value if it could not find
// the key in the cache.
func (l *LRUCache) Get(key interface{}) (val interface{}, found bool) {
	result := l.executeOperation(opTypeGet, key, nil)

	if result.found {
		val = result.target.val
	}

	return val, result.found
}

// Put is an stores key and value in the cache. If the key
// already exists in the cache, then its value and timestamp
// will be updated. Otherwise, the new value will be inserted.
func (l *LRUCache) Put(key, val interface{}) {
	l.executeOperation(opTypePut, key, val)
}

// Remove removes the key/val from the cache. It will return
// the value found under the key and a bool indicating whether
// or not it could find the key in the cache.
func (l *LRUCache) Remove(key interface{}) (val interface{}, found bool) {
	result := l.executeOperation(opTypeRemove, key, nil)

	if result.found {
		val = result.target.val
	}

	return val, result.found
}

// Size returns the current size of the cache
func (l *LRUCache) Size() int {
	return len(l.lookupMap)
}

func (l *LRUCache) executeOperation(opT opType, key, val interface{}) result {
	op := operation{
		opType:     opT,
		key:        key,
		val:        val,
		resultChan: make(chan result),
	}

	l.operationChan <- op

	return <-op.resultChan
}

func (l *LRUCache) run() {
	for {
		select {
		case op := <-l.operationChan:
			var runResult result

			switch op.opType {
			case opTypePut:
				runResult = l.put(op)
			case opTypeGet:
				runResult = l.get(op)
			case opTypeRemove:
				runResult = l.remove(op)
			}

			op.resultChan <- runResult

		case resultChan := <-l.shutdownChan:
			resultChan <- true
			return
		}
	}
}

func (l *LRUCache) put(op operation) result {
	// Check if exists
	if l.ContainsKey(op.key) {
		return l.update(op)
	} else {
		nNode := &node{
			key: op.key,
			val: op.val,
		}

		l.insert(nNode)

		return result{
			found:  true,
			target: nNode,
		}
	}
}

func (l *LRUCache) update(op operation) result {
	// pull it out of the cache
	existing := l.remove(op)

	// update its value
	existing.target.val = op.val

	l.insert(existing.target)

	return existing
}

func (l *LRUCache) get(op operation) result {
	if !l.ContainsKey(op.key) {
		return result{}
	}

	// Update the timestamp and move to correct
	// spot in linked list
	existing := l.lookupMap[op.key]
	l.removeNode(existing)
	l.insert(existing)

	return result{
		found:  true,
		target: existing,
	}
}

func (l *LRUCache) insert(target *node) {
	// insert into doubly linked list
	insert(nil, l.youngest, target)

	// insert into lookup map
	l.lookupMap[target.key] = target

	// adjust size and remove oldest if necessary
	if l.Size() > l.maxSize {
		newOldest := l.oldest.younger
		l.removeNode(l.oldest)
		l.oldest = newOldest
	}

	l.youngest = target

	if l.oldest == nil {
		l.oldest = target
	}
}

func (l *LRUCache) remove(op operation) result {
	// Try to find the simpleNode based off of key
	getResult := l.get(op)

	if !getResult.found {
		return getResult
	}

	l.removeNode(getResult.target)

	return getResult
}

func (l *LRUCache) removeNode(target *node) {
	// adjust youngest if necessary
	if l.youngest == target {
		l.youngest = target.older
	}

	// adjust oldest if necessary
	if l.oldest == target {
		l.oldest = target.younger
	}

	// remove it from the cache
	remove(target)
	delete(l.lookupMap, target.key)
}
