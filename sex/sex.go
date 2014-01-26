// sex: slide execution
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

var (
	port    = flag.String("port", ":1958", "http service address")
	deckdir = flag.String("dir", ".", "directory for decks")
	deckrun = false
	deckpid int
)

const timeformat = "Jan 2, 2006, 3:04pm (MST)"

type layout struct {
	x     float64
	align string
}

func main() {
	flag.Parse()
	err := os.Chdir(*deckdir)
	if err != nil {
		log.Fatal("Set Directory:", err)
	}
	log.Print("Startup...")
	http.Handle("/deck/", http.HandlerFunc(deck))
	http.Handle("/upload/", http.HandlerFunc(upload))
	http.Handle("/table/", http.HandlerFunc(processtable))

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
				w.Write([]byte(",\n"))
			}
			w.Write([]byte(fmt.Sprintf("{\"name\":\"%s\", \"size\":%d, \"date\":\"%s\"}",
				s.Name(), s.Size(), s.ModTime().Format(timeformat))))
		}
	}
	w.Write([]byte("]}\n"))
}

// maketable creates a deck file from a tab separated list
// that includes a specification in the first record
func maketable(w io.Writer, r io.Reader) {
	y := 90.0
	linespacing := 8.0
	textsize := 3.0
	tightness := 3.5
	showrule := true

	l := make([]layout, 10)
	fmt.Fprintf(w, "<deck><slide>\n")
	scanner := bufio.NewScanner(r)
	for nr := 0; scanner.Scan(); nr++ {
		data := scanner.Text()
		fields := strings.Split(data, "\t")
		nf := len(fields)
		if nf > 10 {
			nf = 10
		}

		if nr == 0 {
			for i := 0; i < nf; i++ {
				c := strings.Split(fields[i], ":")
				x, _ := strconv.ParseFloat(c[0], 64)
				l[i].x = x
				l[i].align = c[1]
			}
		} else {
			ty := y - (linespacing / tightness)
			for i := 0; i < nf; i++ {
				fmt.Fprintf(w, "<text xp=\"%g\" yp=\"%g\" sp=\"%g\" align=\"%s\">%s</text>\n",
					l[i].x, y, textsize, l[i].align, fields[i])
			}
			if showrule {
				fmt.Fprintf(w, "<line xp1=\"%g\" yp1=\"%.2f\" xp2=\"%g\" yp2=\"%.2f\" sp=\"0.05\"/>\n",
					l[0].x, ty, l[nf-1].x+5, ty)
			}
		}
		y -= linespacing
	}
	fmt.Fprintf(w, "</slide></deck>\n")
}

// table makes a table from POSTed data
func processtable(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	if req.Method == "POST" {
		deckpath := req.Header.Get("Deck")
		defer req.Body.Close()
		f, err := os.Create(deckpath)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%s %v", requester, err)
			return
		}
		maketable(f, req.Body)
		f.Close()
		w.Write([]byte(fmt.Sprintf("{\"table\":\"%s\"}\n", deckpath)))
		log.Printf("%s table %s", requester, deckpath)
	}
}

// upload uploads decks from POSTed data
func upload(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	if req.Method == "POST" || req.Method == "PUT" {
		deckpath := req.Header.Get("Deck")
		deckdata, err := ioutil.ReadAll(req.Body)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%s %v", requester, err)
			return
		}
		defer req.Body.Close()
		err = ioutil.WriteFile(deckpath, deckdata, 0644)
		if err != nil {
			w.WriteHeader(500)
			log.Printf("%s %v", requester, err)
			return
		}
		w.Write([]byte(fmt.Sprintf("{\"upload\":\"%s\"}\n", deckpath)))
		log.Printf("%s write: %#v, %d bytes", requester, deckpath, len(deckdata))
	}
}

// deck processes slide decks
// GET /deck  -- list information
// POST /deck/file.xml?cmd=[duration] -- starts a deck
// POST /deck?cmd=stop -- stops a deck
// DELETE /deck/file.xml  --  removes a deck
func deck(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	query := req.URL.Query()
	deck := path.Base(req.URL.Path)
	cmd := query["cmd"]
	method := req.Method
	postflag := method == "POST" && len(cmd) == 1 && deck != "deck"
	log.Printf("%s %s %#v %#v", requester, method, deck, cmd)
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
		w.Write([]byte(fmt.Sprintf("{\"deckpid\":\"%d\", \"deck\":\"%s\", \"duration\":\"%s\"}\n", deckpid, deck, cmd[0])))
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
		w.Write([]byte(fmt.Sprintf("{\"deckpid\":\"%d\"}\n", deckpid)))
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
		w.Write([]byte(fmt.Sprintf("{\"remove\":\"%s\"}\n", deck)))
		log.Printf("%s remove %s", requester, deck)
		return
	}
}
