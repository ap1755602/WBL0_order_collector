package main

import (
	"database/sql"
	"flag"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"sync"
	"wildberries_L0/DBHandle"
	"wildberries_L0/model"
	"wildberries_L0/web"
)

var usageStr = `
Usage: stan-subAndServ [options] <subject>
Options:
	-cid, 	<cluster name>		NATS Streaming cluster name
	-id,	<client ID>      	NATS Streaming client ID
	-n,		<name DB>			Name of the data-base 
	-u,		<DB username>		Username of the connecting database
	-p,		<DB user password>	Password of the connecting database user
	-a,		<DB address>		Database address
`

func usage() {
	log.Fatalf(usageStr)
}

func printMsg(m *stan.Msg, i int) {
	log.Printf("[#%d] Received: %s\n", i, m)
}

func main() {
	var (
		clusterID, clientID, subj string
		db                        *DBHandle.ConfigDB
		mutex                     sync.Mutex
	)

	db = new(DBHandle.ConfigDB)
	cache := make(map[string]*model.Content, 0)

	flag.StringVar(&clusterID, "cid", "test-cluster", "The NATS Streaming cluster ID")
	flag.StringVar(&clientID, "id", "stan-subAndServ", "The NATS Streaming client ID to connect with")
	flag.StringVar(&db.Pass, "p", "frodney", "Password of the connecting database user")
	flag.StringVar(&db.User, "u", "frodney", "Username of the connecting database")
	flag.StringVar(&db.Addr, "a", "localhost", "Database address")
	flag.StringVar(&db.Name, "n", "wb_l0_db", "Name of the data-base")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Printf("Error: A subject must be specified.")
		usage()
	}
	subj = args[0]

	// Connect to the DB and download its content

	open := db.Connect()
	db.LoadCache(&cache, open)
	defer func(open *sql.DB) {
		err := open.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(open)

	// Connect and subscribe to the nats-streaming-server

	sc, err := stan.Connect(clusterID, clientID)
	if err != nil {
		log.Fatalf("Can't connect: %v.\nMake sure a NATS Streaming Server is running at: %s", err, " 0.0.0.0:4222")
	}

	_, err = sc.Subscribe(subj, func(msg *stan.Msg) {
		mutex.Lock()
		db.AddNewOrder(&cache, open, msg.Data)
		mutex.Unlock()

	})
	if err != nil {
		log.Fatal(err)
	}

	// Start http server listening

	web.Server(&cache, &mutex)
}
