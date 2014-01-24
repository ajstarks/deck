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
	"path"
	"strings"
)

var port = flag.String("port", ":1958", "http service address")
var deckdir = flag.String("dir", ".", "directory for decks")
var deckrun = false
var deckpid int

const timeformat = "Jan 2, 2006, 3:04pm (MST)"

func main() {
	flag.Parse()
	err := os.Chdir(*deckdir)
	if err != nil {
		log.Fatal("Set Directory:", err)
	}
	log.Print("Startup...")
	http.Handle("/deck/", http.HandlerFunc(deck))
	http.Handle("/upload/", http.HandlerFunc(upload))

	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// writedeckinfo returns information (file, size, date) for a .xml files in the deck directory
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

// upload uploads decks
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

// deck processes slide decks
// GET /slide  -- list information
// POST /slide/file.xml?cmd=[duration] -- starts a deck
// POST /slide?cmd=stop -- stops a deck
// DELETE /slide/file.xml  --  removes a deck
func deck(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	query := req.URL.Query()
	deck := path.Base(req.URL.Path)
	cmd := query["cmd"]
	method := req.Method
	postflag := method == "POST" && len(cmd) == 1
	log.Printf("%s %s %#v %#v", requester, method, deck, cmd)
	if deck == "deck" {
		return
	}
	switch {
	case postflag && !deckrun && cmd[0] != "stop":
		if deck == "" {
			w.WriteHeader(406)
			return
		}
		command := exec.Command("vgdeck", "-loop", cmd[0], deck)
		err := command.Start()
		if err != nil {
			log.Printf("%s %v", requester, err)
			w.WriteHeader(500)
			return
		}
		deckpid = command.Process.Pid
		deckrun = true
		log.Printf("%s deck: %#v, duration: %#v, pid: %d", requester, deck, cmd[0], deckpid)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"DeckPid":"%d", "Deck":"%s", "Duration":"%s"}`, deckpid, deck, cmd[0])))
		return
	case postflag && deckrun && cmd[0] == "stop":
		kp, err := os.FindProcess(deckpid)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%s %v", requester, err)
			return
		}
		err = kp.Kill()
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%s %v", requester, err)
			return
		}
		log.Printf("%s kill %d", requester, deckpid)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(fmt.Sprintf(`{"DeckPid":"%d"}`, deckpid)))
		deckrun = false
		return
	case method == "GET":
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
		return
	case method == "DELETE":
		if deck == "" {
			log.Printf("%s need the name to remove", requester)
                        w.WriteHeader(406)
			return
		}
		err := os.Remove(deck)
		if err != nil {
			log.Printf("%s %v", requester, err)
			w.WriteHeader(500)
			return
		}
		log.Printf("%s remove %s", requester, deck)
		return
	}
}
