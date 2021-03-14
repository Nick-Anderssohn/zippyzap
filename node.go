package zippyzap

import "time"

type node struct {
	timestamp time.Time

	key interface{}
	val interface{}

	older *node
	younger *node
}

func newNode(key, val interface{}) *node {
	return &node{
		timestamp: time.Now(),
		key:       key,
		val:       val,
	}
}

// insert the node in order based off of timestamp.
// This walks from the youngest end towards the oldest
func insert(younger, older *node, nNode *node) {
	// walk to the correct place in the linked list
	// if new's timestamp is before older, than older isn't actually
	// older than new
	if older != nil && nNode.timestamp.Before(older.timestamp) {
		insert(older, older.older, nNode)
		return
	}

	// insert the new node
	nNode.younger = younger
	nNode.older = older

	if older != nil {
		older.younger = nNode
	}

	if younger != nil {
		younger.older = nNode
	}
}

func remove(target *node) {
	if target.older != nil {
		target.older.younger = target.younger
	}

	if target.younger != nil {
		target.younger.older = target.older
	}

	// nice to completely detach. then, if for some reason
	// we want to check if target is still in the linked list,
	// we can inspect these two pointers.
	target.older = nil
	target.younger = nil
}