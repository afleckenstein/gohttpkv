package main

import "net/http"
import "fmt"
import "log"

func SetHandler(db map[string]string, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if len(q) != 1 {
		http.Error(w, fmt.Sprintf("must supply exactly 1 key"), 400)
		return
	}
	var k string
	var v []string
	for k, v = range q {
		break
	}
	if (v[0] == "") || (len(v) != 1) {
		http.Error(w, fmt.Sprintf("must supply exactly 1 value for key"), 400)
		return
	} else {

		db[k] = v[0]
		return
	}
}

func GetHandler(db map[string]string, w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	if len(q) != 1 {
		errmsg := fmt.Sprintf("must supply exactly 1 key")
		http.Error(w, errmsg, 400)
		return
	} else {
		var qk string
		var qv []string
		for qk, qv = range q {
			break
		}
		if qv[0] != "" {
 			http.Error(w, fmt.Sprintf("key cannot have value"), 400)
			return
		}
		v := db[qk]
		if v != "" {
			fmt.Fprint(w, v)
		} else {
			http.Error(w, "", 404)
		}
	}
}

func main() {
	db := make(map[string]string)

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		SetHandler(db, w, r)
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		GetHandler(db, w, r)
	})

	log.Fatal(http.ListenAndServe(":4000", nil))
}
