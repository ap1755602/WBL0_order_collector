package main

import (
	"flag"
	"github.com/nats-io/stan.go"
	"io/ioutil"
	"log"
)

var usageStr = `
Usage: stan-subAndServ [options] <subject> <path json order file>
Options:
	-cid,  <cluster name>   NATS Streaming cluster name
	-id, <client ID>      NATS Streaming client ID
`

func usage() {
	log.Fatalf(usageStr)
}

func main() {
	var (
		clusterID string
		clientID  string
		subj      string
		msg       string // single json file only
	)
	flag.StringVar(&clusterID, "cid", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-subAndServ", "The NATS Streaming client ID to connect with")
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Printf("Error: A subject must be specified.")
		usage()
	}
	subj, msg = args[0], args[1]
	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect: %v.\n", err)
	}
	defer sc.Close()
	rawJson, err := ioutil.ReadFile(msg)
	if err != nil {
		log.Fatalf("Error during open file: %v\n", err)
	}
	err = sc.Publish(subj, rawJson)
}
