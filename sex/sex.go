// sex: slide execution
package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ajstarks/deck"
)

const (
	timeformat = "Jan 2, 2006, 3:04pm (MST)"
	deckpat    = `\.xml$`
	imgpat     = `\.png$|\.jpg$|\.jpeg$`
	stdpat     = `\.xml$|` + imgpat
	vidpat     = `\.mov$|\.mp4$|\.m4v$|\.avi$|\.h264$`
	inforesp   = "{\"API\":[{\"deck\":\"/deck/\"},{\"upload\":\"/upload/\"},{\"media\":\"/media/\"},{\"table\":\"/table/\"}]}\n"
)

var (
	listen    = flag.String("listen", ":1958", "http service address")
	sdir      = flag.String("dir", ".", "directory for decks")
	maxupload = flag.Int64("maxupload", 50*1024*1024, "maximum upload size")
	deckrun   = false
	deckpid   int
	filepats  = map[string]string{"std": stdpat, "deck": deckpat, "image": imgpat, "video": vidpat}
)

type layout struct {
	x     float64
	align string
}

func main() {
	flag.Parse()
	deckdir, err := filepath.Abs(*sdir)
	if err != nil {
		log.Fatal("Directory:", err)
	}
	err = os.Chdir(deckdir)
	if err != nil {
		log.Fatal("Set Directory:", err)
	}
	http.Handle("/", http.HandlerFunc(info))
	http.Handle("/deck/", http.HandlerFunc(dodeck))
	http.Handle("/upload/", http.HandlerFunc(upload))
	http.Handle("/table/", http.HandlerFunc(table))
	http.Handle("/media/", http.HandlerFunc(media))

	log.Printf("Serving from %q, upload limit: %d", deckdir, *maxupload)
	err = http.ListenAndServe(*listen, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

// stopProcess kills and waits for a running process to exit
func stopProcess(w http.ResponseWriter, pid int, requester string) {
	kp, err := os.FindProcess(pid)
	if err != nil {
		eresp(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%s %v", requester, err)
		return
	}
	err = kp.Kill()
	if err != nil {
		eresp(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%s %v", requester, err)
		return
	}
	var ps *os.ProcessState
	ps, err = kp.Wait()
	if err != nil {
		eresp(w, err.Error(), http.StatusInternalServerError)
		log.Printf("%s %v", requester, err)
		return
	}

	log.Printf("%s deck: %v %d exited=%v", requester, ps, pid, ps.Exited())
}

// validpath returns the base of path, or the empty string for the current path
func validpath(s string) string {
	b := filepath.Base(s)
	if b == "." {
		return ""
	}
	return b
}

// eresp sends the client a JSON encoded error
func eresp(w http.ResponseWriter, err string, code int) {
	http.Error(w, fmt.Sprintf("{\"error\": \"%s\"}", err), code)
}

func metadata(filename string) string {
	stdmeta := "---"
	if strings.HasSuffix(filename, ".xml") {
		d, err := deck.Read(filename, 0, 0)
		if err != nil {
			return stdmeta
		} else {
			return fmt.Sprintf("%d slides", len(d.Slide))
		}
	}

	if strings.HasSuffix(filename, ".png") || strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") {
		var err error
		f, err := os.Open(filename)
		if err != nil {
			return stdmeta
		}
		defer f.Close()
		im, _, err := image.DecodeConfig(f)
		if err != nil {
			return stdmeta
		} else {
			return fmt.Sprintf("(%d x %d)", im.Width, im.Height)
		}
	}
	return stdmeta
}

// deckinfo returns information (file, size, date) for a deck and movie files in the deck directory
func deckinfo(w http.ResponseWriter, data []os.FileInfo, pattern string) {
	io.WriteString(w, `{"decks":[`)
	nf := 0
	for _, s := range data {
		matched, err := regexp.MatchString(pattern, s.Name())
		if err == nil && matched {
			nf++
			if nf > 1 {
				io.WriteString(w, ",\n")
			}
			io.WriteString(w, fmt.Sprintf(`{"name":"%s", "meta":"%s", "size":%d, "date":"%s"}`,
				s.Name(), metadata(s.Name()), s.Size(), s.ModTime().Format(timeformat)))
		}
	}
	io.WriteString(w, "]}\n")
}

// info show API information
// GET /
func info(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "GET" {
		io.WriteString(w, inforesp)
		log.Printf("%s info", req.RemoteAddr)
	}
}

// maketable creates a deck file from a tab separated list
// that includes a specification in the first record
func maketable(w io.Writer, r io.Reader, textsize float64) {
	y := 95.0
	linespacing := (textsize * 2.0) + 2.0
	// tightness :=  0.0 // linespacing / 2.0

	l := make([]layout, 10)
	fmt.Fprintf(w, "<deck><slide>\n")
	scanner := bufio.NewScanner(r)
	for nr := 0; scanner.Scan(); nr++ {
		data := scanner.Text()
		fields := strings.Split(data, "\t")
		nf := len(fields)
		if nf > 10 || nf < 1 {
			nf = 10
		}
		if nr == 0 {
			for i := 0; i < nf; i++ {
				c := strings.Split(fields[i], ":")
				if len(c) != 2 {
					return
				}
				x, _ := strconv.ParseFloat(c[0], 64)
				l[i].x = x
				l[i].align = c[1]
			}
		} else {
			for i := 0; i < nf; i++ {
				fmt.Fprintf(w, "<text xp=\"%g\" yp=\"%g\" sp=\"%g\" align=\"%s\">%s</text>\n",
					l[i].x, y, textsize, l[i].align, fields[i])
			}
			/*
				ty := y - (linespacing - tightness)
				if showrule {
					fmt.Fprintf(w, "<line xp1=\"%g\" yp1=\"%.2f\" xp2=\"%g\" yp2=\"%.2f\" sp=\"0.05\"/>\n",
						l[0].x, ty, l[nf-1].x+5, ty)
				}
			*/
		}
		y -= linespacing
	}
	fmt.Fprintf(w, "</slide></deck>\n")
}

// table makes a table from POSTed data
// POST /table/, Deck:<input>
// POST /table/?textsize=[size] Deck:<input> -- speciify the text size
func table(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	query := req.URL.Query()
	p, ok := query["textsize"]
	var param string
	if ok {
		param = p[0]
	}
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "POST" {
		defer req.Body.Close()
		deckpath := validpath(req.Header.Get("Deck"))
		if deckpath == "" {
			eresp(w, "table: no deckpath", http.StatusInternalServerError)
			log.Printf("%s table error: no deckpath", requester)
			return
		}
		f, err := os.Create(deckpath)
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		var textsize float64
		textsize, err = strconv.ParseFloat(param, 64)
		if err != nil || textsize > 20 || textsize < 0.5 {
			textsize = 3.0
		}
		maketable(f, req.Body, textsize)
		f.Close()
		io.WriteString(w, fmt.Sprintf("{\"table\":\"%s (%.1f)\"}\n", deckpath, textsize))
		log.Printf("%s table: %s size:%.1f", requester, deckpath, textsize)
	}
}

// upload uploads decks from POSTed data
// POST /upload, Deck:<file>
func upload(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	w.Header().Set("Content-Type", "application/json")
	if req.Method == "POST" || req.Method == "PUT" {
		deckpath := validpath(req.Header.Get("Deck"))
		if deckpath == "" {
			eresp(w, "upload: no deckpath", http.StatusInternalServerError)
			log.Printf("%s upload error: no deckpath", requester)
			return
		}
		deckdata, err := ioutil.ReadAll(req.Body)
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		defer req.Body.Close()
		dl := int64(len(deckdata))
		if dl > *maxupload {
			msg := fmt.Sprintf("upload: %d bytes over the limit of %d", dl-*maxupload, *maxupload)
			eresp(w, msg, http.StatusForbidden)
			log.Printf(requester + " " + msg)
			return
		}
		err = ioutil.WriteFile(deckpath, deckdata, 0644)
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		io.WriteString(w, fmt.Sprintf("{\"upload\":\"%s\", \"size\": %d}\n", deckpath, dl))
		log.Printf("%s upload: %#v, %d bytes", requester, deckpath, dl)
	}
}

// media plays video
// POST /media Media:<file>
func media(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	w.Header().Set("Content-Type", "application/json")
	media := validpath(req.Header.Get("Media"))
	method := req.Method
	query := req.URL.Query()
	p, ok := query["cmd"]
	var param string
	if ok {
		param = p[0]
	}
	if method == "POST" && param == "" && media != "" {
		log.Printf("%s media: running %s", requester, media)
		command := exec.Command("omxplayer", "-o", "both", media)
		err := command.Start()
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		deckpid = command.Process.Pid
		log.Printf("%s media: %#v, pid: %d", requester, media, deckpid)
		io.WriteString(w, fmt.Sprintf("{\"deckpid\":\"%d\", \"media\":\"%s\"}\n", deckpid, media))
		return
	}

	if method == "POST" && param == "stop" {
		stopProcess(w, deckpid, requester)
		io.WriteString(w, fmt.Sprintf("{\"stop\":\"%d\"}\n", deckpid))
		return
	}
}

// deck processes slide decks
// GET /deck  -- list information
// POST /deck/file.xml?cmd=[duration] -- starts a deck
// POST /deck?cmd=stop -- stops a deck
// DELETE /deck/file.xml  --  removes a deck
func dodeck(w http.ResponseWriter, req *http.Request) {
	requester := req.RemoteAddr
	w.Header().Set("Content-Type", "application/json")
	query := req.URL.Query()
	dpath := strings.Split(req.URL.Path, "/")
	if len(dpath) < 3 {
		eresp(w, "malformed URL", http.StatusForbidden)
		log.Printf("%s malformed URL", requester)
		return
	}
	deck := dpath[2]
	p, ok := query["cmd"]
	var param string
	if ok {
		param = p[0]
	}
	method := req.Method
	postflag := method == "POST" && len(param) > 0
	switch {
	case postflag && !deckrun && param != "stop":
		if deck == "" {
			eresp(w, "deck: need a deck", http.StatusForbidden)
			log.Printf("%s deck: need a deck", requester)
			return
		}
		pausetime, err := time.ParseDuration(param)
		if err != nil {
			eresp(w, err.Error(), http.StatusForbidden)
			log.Printf("%s %v", requester, err)
			return
		}
		if pausetime > 24*time.Hour {
			eresp(w, "deck: pause time too long", http.StatusForbidden)
			log.Printf("%s deck: pause time too long", requester)
			return
		}
		var slidenum string
		sn, ok := query["slide"]
		if ok {
			slidenum = sn[0]
		} else {
			slidenum = "0"
		}
		log.Printf("vgdeck -loop %s -slide %s %s", param, slidenum, deck)
		command := exec.Command("vgdeck", "-loop", param, "-slide", slidenum, deck)
		err = command.Start()
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		deckpid = command.Process.Pid
		deckrun = true
		log.Printf("%s deck: %#v, duration: %#v, pid: %d", requester, deck, param, deckpid)
		io.WriteString(w, fmt.Sprintf("{\"deckpid\":\"%d\", \"deck\":\"%s\", \"duration\":\"%s\"}\n", deckpid, deck, param))
		return
	case postflag && deckrun && param == "stop":
		stopProcess(w, deckpid, requester)
		io.WriteString(w, fmt.Sprintf("{\"stop\":\"%d\"}\n", deckpid))
		deckrun = false
		return
	case method == "GET":
		names, err := ioutil.ReadDir(".")
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		filepattern := stdpat
		filter, fok := query["filter"]
		if fok {
			fp, pok := filepats[filter[0]]
			if pok {
				filepattern = fp
			}
		}
		log.Printf("%s deck: list content %s", requester, filepattern)
		deckinfo(w, names, filepattern)
		return
	case method == "DELETE":
		if deck == "" {
			eresp(w, "deck delete: specify a name", http.StatusForbidden)
			log.Printf("%s delete error: specify a name", requester)
			return
		}
		fs, err := os.Stat(deck)
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		if fs.IsDir() {
			eresp(w, "cannot remove directories", http.StatusInternalServerError)
			log.Printf("%s cannot remove directories", requester)
			return
		}
		err = os.Remove(deck)
		if err != nil {
			eresp(w, err.Error(), http.StatusInternalServerError)
			log.Printf("%s %v", requester, err)
			return
		}
		io.WriteString(w, fmt.Sprintf("{\"remove\":\"%s\"}\n", deck))
		log.Printf("%s deck: remove %s", requester, deck)
		return
	}
}
