package zippyzap

// A node with no timestamp tied to it
type node struct {
	key interface{}
	val interface{}

	older   *node
	younger *node
}

func insert(younger, older *node, nNode *node) {
	// insertTimestampedNode the new timestampedNode
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
