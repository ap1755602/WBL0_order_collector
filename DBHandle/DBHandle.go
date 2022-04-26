package DBHandle

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"wildberries_L0/model"
)

type ConfigDB struct {
	User, Pass, Addr, Name string
}

func (c ConfigDB) Connect() *sql.DB {

	url := "postgresql://" + c.User + ":" + c.Pass + "@" + c.Addr + "/" + c.Name + "?sslmode=disable"
	fmt.Println(url)
	open, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	return open
}

func (c ConfigDB) LoadCache(cache *map[string]*model.Content, open *sql.DB) {
	query, err := open.Query("select content from order_wb;")
	if err != nil {
		log.Fatal(err)
	}
	defer func(query *sql.Rows) {
		err := query.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(query)
	for query.Next() {
		var tmpSLB []byte
		err := query.Scan(&tmpSLB)
		if err != nil {
			log.Fatal(err)
		}
		tmpJSON := new(model.Content)
		err = json.Unmarshal(tmpSLB, tmpJSON)
		if err != nil {
			log.Fatal(err)
		}
		(*cache)[tmpJSON.OrderUid] = tmpJSON
	}
}

func (c ConfigDB) AddNewOrder(cache *map[string]*model.Content, open *sql.DB, rawOrder []byte) {
	cont := new(model.Content)
	err := json.Unmarshal(rawOrder, cont)
	if err != nil || cont.OrderUid == "" {
		fmt.Println("Invalid file")
		return
	}
	if _, ok := (*cache)[cont.OrderUid]; ok == true {
		fmt.Println("Order id is already bound")
		return
	}
	(*cache)[cont.OrderUid] = cont
	_, err = open.Exec("INSERT INTO order_wb (order_uid, content) VALUES ($1, $2);",
		cont.OrderUid, rawOrder)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Order is wrote")
}
