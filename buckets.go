package buckets

import (
	"log"
	"sync"
	"time"
)

// todo: should allow user to choose bucket size, like time.Second or time.Minute
func New() *BucketMaster {
	b := &BucketMaster{}
	b.numBuckets = 60 * 5 // store 5 minutes worth
	b.sets = make(map[string]*Set)
	b.requests = make(chan Occurrence, 1000)
	b.start()
	return b
}

type BucketMaster struct {
	Reporters []Reporter

	requests    chan Occurrence
	firstSecond int64
	numBuckets  int
	sets        map[string]*Set
	reporting   bool
	ticker      *time.Ticker
	doneChan    chan bool
	mutex       sync.Mutex
}

func (b *BucketMaster) AddReporter(reporter Reporter) {
	b.Reporters = append(b.Reporters, reporter)
}

// Will run Report() ever duration
func (b *BucketMaster) ReportEvery(duration time.Duration) {
	// might want to add a mutex, but doubt this will be called more than once anyways
	if b.reporting {
		b.ticker.Stop()
		b.doneChan <- true
	}
	b.doneChan = make(chan bool)
	b.ticker = time.NewTicker(duration)
	go func() {
		defer log.Println("ticker exited")
		for {
			select {
			case <-b.ticker.C:
				b.Report()
			case <-b.doneChan:
				log.Println("Done")
				return
			}
		}
	}()
	b.reporting = true
}

func (b *BucketMaster) Report() {
	b.mutex.Lock()
	// lock so we reset everything properly
	sets := b.cloneSets()
	b.reset() // will reuse the slices in the sets
	b.mutex.Unlock()
	// Not sure if we should pass in all the sets to the reporter and let it deal with it?
	//	log.Println("sets", b.sets)
	for _, s := range sets {
		//		log.Println("set", s)
		for _, r := range b.Reporters {
			//			log.Println("reporter", r)
			r.Report(s)
		}
	}
}

func (b *BucketMaster) cloneSets() []*Set {
	sets := make([]*Set, len(b.sets))
	i := 0
	for _, s := range b.sets {
		s2 := s.Clone()
		sets[i] = s2
		i++
	}
	return sets
}

type Set struct {
	Name    string
	Buckets []int64

	mutex sync.Mutex
}

func (s *Set) Total() int64 {
	x := int64(0)
	for _, v := range s.Buckets {
		x += v
	}
	return x
}

func (s *Set) Clone() *Set {
//	log.Println("Cloning", s.Name)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s2 := &Set{Name: s.Name}
	b := make([]int64, len(s.Buckets))
	copy(b, s.Buckets)
	s2.Buckets = b
	//	s.Buckets = s.Buckets[:0]
	return s2
}

func (b *BucketMaster) AddSet(name string) {
	r := &Set{}
	r.Name = name
	r.Buckets = make([]int64, b.numBuckets)
	b.sets[name] = r
}

func (b *BucketMaster) reset() {
	b.firstSecond = time.Now().Unix()
}

func (b *BucketMaster) start() {
	b.reset()
	go func() {
		for t := range b.requests {
			rd := b.sets[t.Name()]
			if rd == nil {
				log.Panicln("Set", t.Name(), "not found, be sure to call AddSet for each bucket list.")
			}
			rd.Buckets[t.Unix()-b.firstSecond] += 1
		}
	}()
}

// explicit shutdown
func (b *BucketMaster) Stop() {
	close(b.requests)
}

// Adds an Occurence to a set
// todo: should think of a better name for this function
func (b *BucketMaster) Inc(o Occurrence) {
	b.requests <- o
}

// Get a set
func (b *BucketMaster) Get(name string) *Set {
	return b.sets[name]
}

type Occurrence interface {
	Name() string
	Unix() int64
}

type DefaultOccurrence struct {
	Nam  string
	Time time.Time
}

func (ri *DefaultOccurrence) Name() string {
	return ri.Nam
}

func (ri *DefaultOccurrence) Unix() int64 {
	return ri.Time.Unix()
}
