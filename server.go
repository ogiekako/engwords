package hello

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
)

const mappingKind = "words"

type Entry struct {
	Word string
}

func init() {
	http.HandleFunc("/submit", submit)
	http.Handle("/", http.RedirectHandler("/submit", http.StatusFound))
}

func submit(w http.ResponseWriter, r *http.Request) {
	word := r.FormValue("word")
	if r.Method != "POST" || word == "" {
		fmt.Fprint(w, `<html><form action="submit" method="post"><input type="text" name="word" placeholder="word"><input type="submit" value="submit"/></form></html>`)
		return
	}
	c := appengine.NewContext(r)
	key := datastore.NewKey(c, mappingKind, word, 0, nil)
	e := Entry{
		Word: word,
	}
	key, err := datastore.Put(c, key, &e)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = datastore.Get(c, key, &e); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Submit %v value %v\n", r.Method, e.Word)
}
