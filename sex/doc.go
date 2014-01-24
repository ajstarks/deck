/*

sex is a server program that provides an API for slide decks
To install:
	
	go get github.com/ajstarks/deck/sex

Command line options control the working directory and address:port

-port Address:port (default: localhost:1958) 

-dir working directory (default: ".")

GET /deck lists information on slide decks, (filename, file size, modification time) in JSON

POST /deck/file.xml?cmd=[duration]  starts up a deck; the deck, duration, and process id are returned in JSON

POST /deck?cmd=stop stops the running  deck

DELETE /deck/file.xml  removes a deck

PUT or POST to /upload with the Deck: header sets uploads the contents of the Deck: header to the server

*/ 
package main
