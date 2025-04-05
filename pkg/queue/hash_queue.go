package queue

type HashQueue struct {
	hashes chan string
}

func NewHashQueue(capacity int) *HashQueue {
	return &HashQueue{
		hashes: make(chan string, capacity),
	}
}

func (q *HashQueue) Push(hash string) {
	q.hashes <- hash
}

func (q *HashQueue) Pop() (string, bool) {
	hash, ok := <-q.hashes
	return hash, ok
}

func (q *HashQueue) Size() int {
	return len(q.hashes)
}

func (q *HashQueue) PushAll(hashes []string) {
	for _, hash := range hashes {
		q.hashes <- hash
	}
}
