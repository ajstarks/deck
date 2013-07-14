/*
vgdeck is a program for showing presentations on the Raspberry Pi, using the deck and openvg libraries.
To install:

	go install github.com/ajstarks/deck/vgdeck

To run vgdeck, specify one or more files (marked up in deck xml) on the command line, and each will be shown in turn.

	$ vgdeck sales.xml program.xml architecture.xml

Here are the vgdeck commands:

      Next slide: +, Ctrl-N, [Return]
      Previous slide, -, Ctrl-P, [Backspace]
      First slide: ^, Ctrl-A
      Last slide: $, Ctrl-E
      Reload: r, Ctrl-R
      X-Ray: x, Ctrl-X
      Search: /, Ctrl-F
      Save: s, Ctrl-S
      Quit: q

All commands are a single keystroke, acted on immediately
(only the search command waits until you hit [Return] after entering your search text)
To cycle through the deck, repeatedly tap [Return] key
*/
package main
