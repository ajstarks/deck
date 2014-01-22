// sex: slide execution
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

var port = flag.String("port", ":1958", "http service address")
var deckdir = flag.String("dir", ".", "directory for decks")
var deckpid = -1
const timeformat = "Jan 2, 2006, 3:04pm (MST)"

func main() {
	flag.Parse()
	err := os.Chdir(*deckdir)
	if  err != nil {
		log.Fatal("Set Directory", err)
	}
	log.Print("Startup...")
	http.Handle("/deck/", http.HandlerFunc(deck))
	http.Handle("/upload/", http.HandlerFunc(upload))

	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func writedeckinfo(w http.ResponseWriter, data []os.FileInfo, suffix string) {
	w.Write([]byte(`{"decks":[`))
	nf := 0
	for _, s := range data {
		if strings.HasSuffix(s.Name(), suffix) {
			nf++
			if nf > 1 {
				w.Write([]byte(`,`))
			}
			w.Write([]byte(fmt.Sprintf(`{"name":"%s", "size":%d, "date":"%s"}`,
				s.Name(), s.Size(), s.ModTime().Format(timeformat))))
		}
	}
	w.Write([]byte(`]}`))
}
func upload(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	if req.Method == "POST" || req.Method == "PUT" {
		path := req.Header.Get("Deck")
		deckdata, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%s %v", requester, err)
			return
		}
		defer req.Body.Close()
		err = ioutil.WriteFile(path, deckdata, 0644)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%s %v", requester, err)
			return
		}
		log.Printf("%s Write: %#v, %d bytes", requester, path, len(deckdata))
	}
}
func deck(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	switch req.Method {
	case "POST":
		deck := req.Header.Get("Deck")
		duration := req.Header.Get("Duration")
		if deckpid == -1 {
			if duration == "" || deck == "" {
				w.WriteHeader(406)
				return
			}
			cmd := exec.Command("vgdeck", "-loop", duration, deck)
			err := cmd.Start()
			if err != nil {
				log.Printf("%s %v", requester, err)
				w.WriteHeader(500)
			}
			deckpid = cmd.Process.Pid
			log.Printf("%s deck: %#v, duration: %#v, pid: %d", requester, deck, duration, deckpid)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf(`{"DeckPid":"%d", "Deck":"%s", "Duration":"%s"}`, deckpid, deck, duration)))
		}

		if req.Header.Get("Kill") != "" {
			kp, err := os.FindProcess(deckpid)
			if err == nil {
				kp.Kill()
				log.Printf("%s kill %d", requester, deckpid)
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(fmt.Sprintf(`{"DeckPid":"%d"}`, deckpid)))
				deckpid = -1
			} else {
				w.WriteHeader(500)
				log.Printf("%s %v", requester, err)
			}
		}
	case "GET":
		f, err := os.Open(*deckdir)
		if err != nil {
			log.Printf("%s %v", requester, err)
			w.WriteHeader(500)
			return
		}
		names, err := f.Readdir(-1)
		if err != nil {
			log.Printf("%s %v", requester, err)
			w.WriteHeader(500)
			return
		}
		log.Printf("%s list decks", requester)
		w.Header().Set("Content-Type", "application/json")
		writedeckinfo(w, names, ".xml")
	case "DELETE":
		remdeck := req.Header.Get("Deck")
		if remdeck != "" {
			err := os.Remove(remdeck)
			if err != nil {
				log.Printf("%s %v", requester, err)
				w.WriteHeader(500)
				return
			}
			log.Printf("%s remove %s", requester, remdeck)
		} else {
			log.Printf("%s need the name to remove", requester)
			w.WriteHeader(406)
		}
	}
}
