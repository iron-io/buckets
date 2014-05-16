package buckets

import (
	"log"
	"time"
)

func New(numBuckets int) *BucketMaster {
	b := &BucketMaster{}
	b.numBuckets = numBuckets
	b.sets = make(map[string]*Set)
	b.requests = make(chan Occurrence, 1000)
	return b
}

type BucketMaster struct {
	requests    chan Occurrence
	firstSecond int64
	numBuckets  int
	sets        map[string]*Set
}

// todo: need better names
type Set struct {
	Buckets []int64
}

func (b *BucketMaster) AddSet(name string) {
	r := &Set{}
	r.Buckets = make([]int64, b.numBuckets)
	b.sets[name] = r
}

func (b *BucketMaster) Start() {
	b.firstSecond = time.Now().Unix()
	go func() {
		for t := range b.requests {
			rd := b.sets[t.Name()]
			if rd == nil {
				log.Panicln("Runner", t.Name(), "not found, be sure to call AddRunner for each bucket list.")
			}
			//			log.Println("Got request", rd)
			rd.Buckets[t.Unix()-b.firstSecond] += 1
		}
	}()
}

func (b *BucketMaster) Stop() {
	close(b.requests)
}

func (b *BucketMaster) Inc(o Occurrence) {
	b.requests <- o
}

func (b *BucketMaster) Get(name string) *Set {
	return b.sets[name]
}

type Occurrence interface {
	Name() string
	Unix() int64
}
