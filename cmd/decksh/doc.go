/*
decksh: a little language for deck markup

```decksh``` is a domain-specific language (DSL) for generating [```deck```](https://github.com/ajstarks/deck/blob/master/README.md) markup.

Running

```decksh``` reads from the specified input, and writes deck markup to the specified output destination:

    $ decksh                   # input from stdin, output to stdout
    $ decksh -o foo.xml        # input from stdin, output to foo.xml
    $ decksh foo.sh            # input from foo.sh output to stdout
    $ decksh -o foo.xml foo.sh # input from foo.sh output to foo.xml

Typically, ```decksh``` acts as the head of a rendering pipeline:

    $ decksh text.dsh | pdf -pagesize 1200,900

Example input

This deck script:

    // Example deck
    midx=50
    midy=50
    iw=640
    ih=480

    imfile="follow.jpg"
    imlink="https://budnitzbicycles.com"
    imscale=58
    dtop=87

    opts="-fulldeck=f -textsize 1  -xlabel=2  -barwidth 1.5"
    deck
        slide "white" "black"
            ctext "Deck elements" midx dtop 5
            cimage "follow.jpg" "Dreams" 72 midy iw ih imscale imlink
            textblock "Budnitz #1, Plainfield, NJ, May 10, 2015" 55 35 10 1 "serif" "white"

            // List
            blist 10 75 3
                li "text, image, list"
                li "rect, ellipse, polygon"
                li "line, arc, curve"
            elist

            // Graphics
            gy=10
            c1="red"
            c2="blue"
            c3="green"
            rect    15 gy 8 6              c1
            ellipse 27.5 gy 8 6            c2
            polygon "37 37 45" "7 13 10"   c3
            line    50 gy 60 gy 0.25       c1
            arc     70 gy 10 8 0 180 0.25  c2
            curve   80 gy 95 25 90 gy 0.25 c3


            // Chart
            chartleft=10
            chartright=45
            charttop=42
            chartbottom=28
            dchart -left chartleft -right chartright -top charttop -bottom chartbottom opts AAPL.d
        eslide
    edeck


Produces:

![exampledeck](exampledeck.png)

Text, font, color, caption and link arguments follow Go convetions (surrounded by double quotes).
Colors are in rgb format ("rgb(n,n,n)"), or [SVG color names](https://www.w3.org/TR/SVG11/types.html#ColorKeywords).

Coordinates, dimensions, scales and opacities are floating point numbers ranging from from 0-100
(they represent percentages on the canvas and percent opacity).  Some arguments are optional, and
if omitted defaults are applied (black for text, gray for graphics, 100% opacity).

Canvas size and image dimensions are in pixels.

Simple assignments

```id=<number>``` defines a constant, which may be then subtitited. For example:

    x=10
    y=20
    text "hello, world" x y 5

Assignment operations

```id+=<number>``` increment the value of ```id``` by ```<number>```

    x+=5

```id-=<number>``` decrement the value of ```id``` by ```<number>```

    x-=10

```id*=<number>``` multiply the value of ```id``` by ```<number>```

    x*=50

```id*=<number>``` divide the value of ```id``` by ```<number>```

    x/=100

Binary operations

Addition ```id=<id> + number or <id>```

    tx=10
    spacing=1.2

    sx=tx-10
    vx=tx+spacing

Subtraction ```id=<id> - number or <id>```

    a=x-10

Muliplication ```id=<id> * number or <id>```

    a=x*10

Division ```id=<id> / number or <id>```

    a=x/10

Begin or end a deck.

    deck
    edeck

Begin, end a slide with optional background and text colors.

    slide [bgcolor] [fgcolor]
    eslide

Specify the size of the canvas.

    canvas w h


Random Number

	x=random min max

assign a random number in the specified range

Mapping

    x=vmap v vmin vmax min max

For value ```v```, map the range ```vmin-vmax``` to ```min-max```.

Polar Coordinates

    x=polarx cx cy r theta
    y=polary cx cy r theta

Return the polar coordinate given a center at ```(cx, cy)```, radius ```r```, and angle ```theta``` (in degrees)



Loops

Loop over ```statements```, with ```x``` starting at ```begin```, ending at ```end``` with an optional ```increment``` (if omitted the increment is 1).
Substitution of ```x``` will occur in statements.

    for x=begin end [increment]
        statements
    efor

Loop over ```statements```, with ```x``` ranging over the contents of items within ```[]```.
Substitution of ```x``` will occur in statements.

    for x=["abc" "def" "ghi"]
        statements
    efor

Loop over ```statements```, with ```x``` ranging over the contents ```"file"```.
Substitution of ```x``` will occur in statements.

    for x="file"
        statements
    efor


Text

Left, centered, end, rotated, block-aligned text or a file's contents with
optional font ("sans", "serif", "mono", or "symbol"), color and opacity.

Also, show blocks of code on a gray background.

    text       "text"     x y       size [font] [color] [opacity] [link]
    ctext      "text"     x y       size [font] [color] [opacity] [link]
    etext      "text"     x y       size [font] [color] [opacity] [link]
    rtext      "text"     x y angle size [font] [color] [opacity] [link]
    textblock  "text"     x y width size [font] [color] [opacity] [link]
    textfile   "filename" x y       size [font] [color] [opacity] [linespacing]
    textcode   "filename" x y width size [color]

Images

Plain and captioned, with optional scales and links

    image  "file"           x y width height [scale] [link]
    cimage "file" "caption" x y width height [scale] [link]

Lists

(plain, bulleted, numbered, centered). Optional arguments specify the color, opacity, line spacing, link and rotation (degrees)

    list   x y size [font] [color] [opacity] [linespacing] [link] [rotation]
    blist  x y size [font] [color] [opacity] [linespacing] [link] [rotation]
    nlist  x y size [font] [color] [opacity] [linespacing] [link] [rotation]
    clist  x y size [font] [color] [opacity] [linespacing] [link] [rotation]

### list items, and ending the list

    li "text"
    elist

Graphics

Rectangles, ellipses, squares and circles: specify the location and dimensions with optional color and opacity.
The default color and opacity is gray, 100%

    rect    x y w h [color] [opacity]
    ellipse x y w h [color] [opacity]

    square  x y w   [color] [opacity]
    circle  x y w   [color] [opacity]


Rounded rectangles are similar, with the added radius for the corners: (solid colors only)

    rrect   x y w h r [color]


For polygons, specify the x and y coordinates as a series of numbers, with optional color and opacity.

    polygon "xcoords" "ycoords" [color] [opacity]

For lines, specify the coordinates for the beginning and end points.
For horizonal and vertical lines specify the initial point and the length.
Line thickness, color and opacity are optional, with defaults

    line    x1 y1 x2 y2 [size] [color] [opacity]
    hline   x y length  [size] [color] [opacity]
    vline   x y length  [size] [color] [opacity]

Curve is a quadratic bezier: specify the beginning location, the control point, and ending location.
For arcs, specify the location of the center point, the width and height, and the beginning and ending angles.
Line thickness, color and opacity are optional, with defaults (0.2, gray, 100%).

    curve   x1 y1 x2 y2 x3 y3 [size] [color] [opacity]
    arc     x y w h a1 a2     [size] [color] [opacity]

Arrows

Arrows with optional linewidth, width, height, color, and opacity.
Default linewidth is 0.2, default arrow width and height is 3, default color and opacity is gray, 100%.
The curve variants use the same syntax for specifying curves.

    arrow   x1 y1 x2 y2       [linewidth] [arrowidth] [arrowheight] [color] [opacity]
    lcarrow x1 y1 x2 y2 x3 y3 [linewidth] [arrowidth] [arrowheight] [color] [opacity]
    rcarrow x1 y1 x2 y2 x3 y3 [linewidth] [arrowidth] [arrowheight] [color] [opacity]
    ucarrow x1 y1 x2 y2 x3 y3 [linewidth] [arrowidth] [arrowheight] [color] [opacity]
    dcarrow x1 y1 x2 y2 x3 y3 [linewidth] [arrowidth] [arrowheight] [color] [opacity]

Braces

Left, right, up and down-facing braces.
(x, y) is the location of the point of the brace, and linewidth, color and opacity are optional
(defaults are gray, 100%)

    lbrace x y height aw ah [linewidth] [color] [opacity]
    rbrace x y height aw ah [linewidth] [color] [opacity]
    ubrace x y width  aw ah [linewidth] [color] [opacity]
    dbrace x y width  aw ah [linewidth] [color] [opacity]

Charts

Run the [dchart](https://github.com/ajstarks/deck/blob/master/cmd/dchart/README.md) command with the specified arguments.

    dchart [args]

Legend

Show a colored legend

    legend "text" x y size [font] [color]


Include decksh markup from a file

    include "file"

places the contents of ```"file"``` inline.

Data: Make a file

    data "foo.d"
    uno    100
    dos    200
    tres   300
    edata

makes a file named ```foo.d``` with the lines between ```data``` and ```edata```.

Grid: Place objects on a grid

    grid "file.dsh" x y xskip yskip limit

The first file argument (```"file.dsh"``` above) specifies a file with decksh commands; each item in the file must include the arguments "x" and "y". Normal variable substitution occurs for other arguments. For example if the contents of ```file.dsh``` has six items:

    circle x y 5
    circle x y 10
    circle x y 15
    square x y 5
    square x y 10
    square x y 15

The line:

    grid "file.dsh" 10 80 20 30 50

creates two rows: three circles and then three squares

```x, y``` specify the beginning location of the items, ```xskip``` is the horizontal spacing between items.
```yinternal``` is the vertical spacing between items and ```limit``` the the horizontal limit. When the ```limit``` is reached,
a new row is created.
*/
package main
