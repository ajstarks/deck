/*

sex is a server program that provides an API for slide decks. Responses are encoded in JSON.
To install:

	go get github.com/ajstarks/deck/sex

Command line options control the working directory and address:port

-listen Address:port (default: localhost:1958)

-dir working directory (default: ".")

-maxupload maximum upload size (bytes)

GET / lists the API

GET /deck lists information on content, (filename, file size, modification time) in JSON

GET /deck?filter=[type] filter content list by type (std, deck, image, video)

POST /deck/file.xml?cmd=[duration]  starts up a deck; the deck, duration, and process id are returned in JSON

POST /deck?cmd=stop stops the running  deck

DELETE /deck/file.xml  removes a deck

PUT or POST to /upload  uploads the contents of the Deck: header to the server

POST /table with the content of a tab-separated list, creates a slide with a formatted table, the Deck: header specifies the resulting deck file

POST /media plays the media file specified in the Media: header
*/
package main
