package main

import (
	"fmt"
	"os"
	"github.com/boltdb/bolt"
	"log"
	"strconv"
)

func ReadFlatPage(fid string) string {
	post := ""
	db, err := bolt.Open(Database, 0644, nil)
	if err != nil {
		log.Fatal(err, "990")
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("flats"))
		if bucket == nil {
			fmt.Println("Bucket <flats> not found in Db !") //, err)
		}
		fmt.Println(fid)

		post = string(bucket.Get([]byte(fid)))
		return nil
	})
	return post
}

func ReadAllFlats() string {
	post := ""
	db, err := bolt.Open(Database, 0644, nil)
	if err != nil {
		log.Fatal(err, "990")
	}
	defer db.Close()
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("flats"))
		if bucket == nil {
			fmt.Println("Bucket <flats> not found in Db !") //, err)
		}
		//	fmt.Println(fid)
		//	lis := list.New()
		count := 1

		c := bucket.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			//		fmt.Println("key=", k)
			post += "<div class=\"r\"><a href=\"/page/" + string(k) + "\">" + strconv.Itoa(count) + " : " + string(k) + "</a></div><br>"
			//		lis.InsertAfter(count, lis.PushFront(string(k)))
			//		for el := lis.Front(); el != nil; el = el.Next() {
			//			fmt.Println(count, "= ", el.Value)
			//		}
			count++

		}
		return nil
	})
	return post
}

func getFlat(fid string) string {
	page := ""
	db, err := bolt.Open(Database, 0644, nil)
	if err != nil {
		log.Fatal(err, "777")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("flats"))
		if bucket == nil {
			fmt.Println("Bucket <flats> not found in Db !", err)
		}

		chunk := string(bucket.Get([]byte(fid)))
		//		fmt.Println("page in getFlat() = ", chunk)
		if len(chunk) > 0 {
			page = chunk
		}
		fmt.Println("getFlat(" + fid + "): OK")

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
	return page
}

func Format(p string) string {
	return p
}

func savePrint(p string, t string) {
	f, err := os.Create("out/" + p + ".x")
	if err != nil {
		fmt.Println(err)
		return
	}
	l, err := f.WriteString(t)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
