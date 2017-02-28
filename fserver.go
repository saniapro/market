package main

import (
	"fmt"
	"net/http"
	"log"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"runtime"
	"sync"
)

type intc uint32

func jsonCMD(res http.ResponseWriter, req *http.Request) {
	//log.Println(req.Body)
	dec := json.NewDecoder(req.Body)
	type Message struct {
		TargetID, ID intc
		Name string
	}
	var m Message
	if err := dec.Decode(&m); err != nil {
		log.Println(err)
		res.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println(m)
	res.WriteHeader(http.StatusOK)
	if m.Name == "update" && m.ID != 0 {
		db, err := getDB()
		if err != nil {
			log.Fatal("DB:", err)
		}
		if ad:=ads.Load(m.ID, db); ad!= nil {
			fmt.Fprintln(res, ad.getJson())
		}
	} else if ad := ads.list[m.ID]; ad != nil {
		fmt.Fprintln(res, ad.getJson())
	} else {
		fmt.Fprintln(res, "{\"Not found ad\":", m.ID, "}")
	}
}

var ads Ads
var users Users

func getDB() (db *sql.DB, err error) {
	//db, err = sql.Open("mysql", "banner:Y9Y50vdm@tcp(statmaster.mi6.kiev.ua:3306)/uho")
	db, err = sql.Open("mysql", "banner:Y9Y50vdm@tcp(192.168.0.31:3306)/uho")
	if err != nil {
		log.Fatal("DB:", err)
	}
	return db, err
}

func main() {
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func() {
		defer wg.Done()
			db, err := getDB()
		if err == nil {
			users.LoadAll(db)
		}
	}()

	wg.Add(1)
	go func() {
		db1,err := getDB()
		defer wg.Done()
		if err == nil {
			ads.LoadAll(db1)
		}
		db1.Close()
	}()
	
	wg.Wait()
	log.Print("All load")
	var m runtime.MemStats
  runtime.GC()
  runtime.ReadMemStats(&m)
  log.Printf("Alloc %d kb; TotalAlloc %d kb\n", m.Alloc / 1000, m.TotalAlloc / 1000)
	http.HandleFunc("/cmd/", jsonCMD)
	port:=":8077"
	log.Println("vserver started", port)
	err = http.ListenAndServe(port, nil)
  if err != nil {
    log.Fatal("ListenAndServe:", err)
  }
}
