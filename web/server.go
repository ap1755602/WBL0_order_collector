package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"wildberries_L0/model"
)

func Server(content *map[string]*model.Content, mutex *sync.Mutex) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		file, err := template.ParseFiles("../html/index.html")
		if err != nil {
			log.Fatal(err)
		}
		file.Execute(w, nil)
	})
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		if id := r.FormValue("id"); id != "" && (*content)[id] != nil {
			tmpl, ok := template.ParseFiles("../html/result_t.html")
			if ok != nil {
				log.Fatal(ok)
			}
			var buf bytes.Buffer
			marshal, err := json.Marshal((*content)[id])
			if err != nil {
				return
			}
			err = json.Indent(&buf, marshal, "", "\t")
			if err != nil {
				return
			}
			err = tmpl.Execute(w, buf.String())
			if err != nil {
				fmt.Println(err)
			}

		} else {
			http.ServeFile(w, r, "../html/result_f.html")
		}
		mutex.Unlock()
	})
	fmt.Println("Server is listening...")
	err := http.ListenAndServe(":8181", nil)
	if err != nil {
		return
	}
}
