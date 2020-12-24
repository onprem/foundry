package furnace

import "sync"

type queue struct {
	sync.RWMutex
	items []Package
}

func (q *queue) enqueue(packages []Package) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, packages...)
}

// dequeue removes a package from the queue and returns it.
func (q *queue) dequeue() (Package, bool) {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return Package{}, false
	}
	pkg := q.items[0]
	q.items = q.items[1:]
	return pkg, true
}

// GetItems returns all elements in a thread-safe way.
func (q *queue) getItems() []Package {
	q.RLock()
	defer q.RUnlock()
	return q.items
}

func (q *queue) isQueued(pkg Package) bool {
	for _, v := range q.getItems() {
		if v.Name == pkg.Name {
			return true
		}
	}
	return false
}
