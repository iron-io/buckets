package buckets

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/iron-io/golog"
)

var config *Config

type Config struct {
	Reporters []*ReporterConfig
}

type ReporterConfig struct {
	Service string
	Key     string
}

func init() {
	config = &Config{}
	configFile := "test_config.json"
	config_s, err := ioutil.ReadFile(configFile)
	if err != nil {
		golog.Fatalln("Couldn't find config at:", configFile)
	}

	err = json.Unmarshal(config_s, &config)
	if err != nil {
		golog.Fatalln("Couldn't unmarshal config!", err)
	}
	golog.Infoln("config:", config)
}

func TestBuckets(t *testing.T) {
	log.Println("Starting test...")
	seconds := 90
	bm := New()
	// todo: move this into the lib, others probably would do the same thing
	for _, r := range config.Reporters {
		log.Println("Adding reporter", r.Service)
		switch r.Service {
		case "stdout":
			//			log.Println("Adding stdout")
			bm.AddReporter(NewStdoutReporter())
		case "stathat":
			//			log.Println("Adding stathat")
			bm.AddReporter(NewStathatReporter(r.Key, "test"))
		}
	}
	bm.ReportEvery(5 * time.Second)
	bm.AddSet("set1")
	bm.AddSet("set2")
	stopAt := time.Now().Add(time.Duration(seconds) * time.Second)
	j := 0
	for {
		now := time.Now()
		if now.After(stopAt) {
			break
		}
		bm.Inc(&DefaultOccurrence{"set1", time.Now()})
		if j%2 == 0 {
			bm.Inc(&DefaultOccurrence{"set2", time.Now()})
		}
		j++
	}
	time.Sleep(500 * time.Millisecond)
	bm.Report()

}
