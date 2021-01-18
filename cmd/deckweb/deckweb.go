// deck -- command line access to the deck web API
package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	defurl := os.Getenv("DECKS")
	if defurl == "" {
		defurl = "http://localhost:1958"
	}
	var url = flag.String("url", defurl, "deck service url")
	flag.Parse()
	dispatch(*url, flag.Args())
}

// dispatch parses commands and calls the correct function
func dispatch(url string, args []string) {
	if len(args) < 1 {
		usage()
		return
	}
	switch args[0] {
	case "list":
		list(url, args)
	case "upload", "up", "load":
		upload(url, args)
	case "play", "start":
		play(url, args)
	case "remove", "delete", "del":
		remove(url, args)
	case "table", "tab":
		table(url, args)
	case "video", "media":
		video(url, args)
	case "stop", "kill":
		stop(url)
	default:
		usage()
	}
}

// list performs content listings
func list(url string, files []string) {
	ltype := "std"
	if len(files) > 1 {
		ltype = files[1]
	}
	resp, err := http.Get(url + "/deck/?filter=" + ltype)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	showrb(resp)
}

// play starts up a deck
func play(url string, files []string) {
	if len(files) < 2 {
		fmt.Println("specify a file to play")
		return
	}
	duration := "1s"
	slide := "0"
	if len(files) > 2 {
		duration = files[2]
	}
	if len(files) > 3 {
		slide = files[3]
	}
	resp, err := http.Post(url+"/deck/"+files[1]+"?cmd="+duration+"&slide="+slide, "application/octet-stream", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	showrb(resp)
}

// video plays video
func video(url string, files []string) {
	if len(files) < 2 {
		fmt.Println("specify a file to play")
		return
	}
	client := &http.Client{}
	reqhead(client, "POST", url+"/media/", "Media", filepath.Base(files[1]), nil)
}

// upload copies content to the server
func upload(url string, files []string) {
	if len(files) < 2 {
		fmt.Println("specify files to upload")
		return
	}
	client := &http.Client{}
	for _, filename := range files[1:] {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}
		reqhead(client, "POST", url+"/upload/", "Deck", filepath.Base(filename), f)
		f.Close()
	}
}

// remove content from the server
func remove(url string, files []string) {
	if len(files) < 2 {
		fmt.Fprintln(os.Stderr, "specify files to be removed")
		return
	}
	client := &http.Client{}
	for _, filename := range files[1:] {
		reqhead(client, "DELETE", url+"/deck/"+filename, "Deck", filepath.Base(filename), nil)

	}
}

// table makes tabular slides
func table(url string, files []string) {
	if len(files) < 2 {
		fmt.Println("specify a table file")
		return
	}
	textsize := "1.4"
	if len(files) > 2 {
		textsize = files[2]
	}
	f, err := os.Open(files[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	defer f.Close()
	client := &http.Client{}
	reqhead(client, "POST", url+"/table/?textsize="+textsize, "Deck", filepath.Base(files[1]), f)
}

// stop a deck
func stop(url string) {
	resp, err := http.Post(url+"/deck/?cmd=stop", "", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	showrb(resp)
}

// showrb prints the response body
func showrb(r *http.Response) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	fmt.Printf("%s", string(data))
	r.Body.Close()
}

// reqhead makes a HTTP request, with the specified header, printing the response
func reqhead(client *http.Client, method, url, header, hval string, r io.Reader) {
	var req *http.Request
	var resp *http.Response
	var err error

	req, err = http.NewRequest(method, url, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	req.Header.Add(header, hval)
	resp, err = client.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}
	showrb(resp)
}

// usage prints the usage message
func usage() {
	fmt.Fprintf(os.Stderr, "%s",
		`Usage:
	List:    deck list [image|deck|video]
	Play:    deck play file [duration] [slide number]
	Stop:    deck stop
	Upload:  deck upload files...
	Remove:  deck remove files...
	Video:   deck video file
	Table:   deck table file [textsize]
`)
}
