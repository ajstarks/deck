/*

vgdeck is a program for showing presentations on the Raspberry Pi, using the deck and openvg libraries.
To install:

	go get github.com/ajstarks/deck/vgdeck

To run vgdeck, specify one or more files (marked up in deck xml) on the command line, and each will be shown in turn.

	$ vgdeck -loop <pause> -g <percent> sales.xml program.xml architecture.xml

The -g (grid) option specified te scale of the x-ray grid.
The loop option pauses the specified duration between slides. If loop is not specified, then vgdeck enters
an interactive mode using these commands:

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
