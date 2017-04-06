package main

import (
	"encoding/json"
	"fmt"
	"github.com/httpreserve/httpreserve"
	"os"
	"sync"
	"time"
)

// This structure is used to communicate with the server
// we may also use some static storage in the form of Bolt DB
// the final signal to the webapp will be a empty payload
// but with the complete flag set to true so that we know there
// is no more work to be processed. ls contains a link stat
// data structure if we can recreate one from the JSON we receive
// else the js variable will contain a single JSON document to
// be processed.
type processLog struct {
	complete bool
	ls       httpreserve.LinkStats
	js       string
	lmap     map[string]interface{}
}

var serverWG sync.WaitGroup

// webapprun lets us start the server for the user to access
func webappRun() {
	defer serverWG.Done()

	// pause to let our buffers begin to fill...
	// TODO: look for safer, more idiomatic ways to solve...
	time.Sleep(3 * time.Second)
	fmt.Fprintln(os.Stderr, "server starting on localhost: http://127.0.0.1:2041")

	if port == "" {
		port = "2041"
	}

	err := DefaultServer(port)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

var processedSlices []processLog

func processStats(ce string) {
	defer serverWG.Done()

	// Pause to initialize the programme first and then let we
	// can begin polling websites...
	time.Sleep(2 * time.Second)

	var ls httpreserve.LinkStats

	//for range linkmap {
	//ce := <-ch // json from channel
	err := json.Unmarshal([]byte(ce), &ls)
	if err != nil {
		fmt.Fprintln(os.Stderr, "problem unmarshalling data.", err)
	}

	// retrieve a map from the structure and write it out to the
	// http server...
	lmap := storeStruct(ls, ce)
	if len(lmap) > 0 {
		var ps processLog
		ps.js = ce
		ps.ls = ls
		ps.lmap = lmap
		processedSlices = append(processedSlices, ps)
	}
	//}

	var ps processLog
	ps.complete = true
	processedSlices = append(processedSlices, ps)
}

var clock string

func clockOut() {
	t := time.Now()
	clock = t.Format("Mon Jan _2 15:04:05 2006")
}

// webappHanlder enables us to establish the web server and create
// the structures we need to present our data to the user...
func webappHandler(ch string) {

	serverWG.Add(1)
	go webappRun()

	serverWG.Add(1)
	go processStats(ch)

	serverWG.Wait()
	return
}
