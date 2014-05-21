package buckets

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"fmt"
	"github.com/iron-io/golog"
)

type Reporter interface {
	Report(*Set)
}

func NewStdoutReporter() *StdoutReporter {
	r := &StdoutReporter{}
	return r
}

type StdoutReporter struct {
}

func (r *StdoutReporter) Report(set *Set) {
	log.Println(set.Name, "total:", set.Total())
}

// key is typically your stathat email
// prefix is a prefix for the stat names for easier viewing in stathat
func NewStathatReporter(key, prefix string) *StathatReporter {
	r := &StathatReporter{}
	r.Key = key
	r.Prefix = prefix
	return r
}

type StathatReporter struct {
	Key    string
	Prefix string
}

func (r *StathatReporter) Report(set *Set) {
	log.Println("Posting to stathat")
	values := url.Values{}
	values.Set("count", strconv.FormatInt(set.Total(), 10))
	values.Set("stat", fmt.Sprintf("%v %v", r.Prefix, set.Name))
	values.Set("ezkey", r.Key)
	resp, err := http.PostForm("http://api.stathat.com/ez", values)
	if err != nil {
		golog.Errorln("error posting to StatHat", err)
		return
	}
	if resp.StatusCode != 200 {
		golog.Errorln("bad status posting to StatHat", resp.StatusCode)
	}
	io.Copy(ioutil.Discard, resp.Body)
	// todo: check error
	resp.Body.Close()
}
