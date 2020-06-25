# hexagolang

hexagolang is a hexagon library for golang based on https://www.redblobgames.com/grids/hexagons/.

My typical usage looks like this:

```golang
package main

import (
	"fmt"
	"image"

	hex "github.com/hautenessa/hexagolang"
)

func main() {
	hexagonRadius := 32           // The radius of the hexagon, also the length of each side.
	screenOrigin := image.Point{} // origin to use when converting between screen and hexagon coordinates.

	// The layout object is responsible for converting between screen and hexagon coordinates.
	layout := hex.MakeLayout(hexagonRadius, screenOrigin, hex.OrientationFlat)

	myFirstHex := hex.H{1, 0}                   // Uses axial coordinates.
	screenPoint := layout.CenterFor(myFirstHex) // convert the hexagon center into screen coordinates.
	fmt.Printf("screen point %+v\n", screenPoint)

	// Create a line between two points and list the hexagon steps.
	mySecondHex := hex.H{4, 2}
	for h, gon := range hex.Line(myFirstHex, mySecondHex) {
		fmt.Printf("%d: %v", h, gon)
	}
	fmt.Printf("\n")

    // Convert a random screen position into a hexagon coordinate. 
	mousePoint := image.Point{800, 600}
	myThirdHex := layout.HexFor(mousePoint)
	fmt.Printf("hex for mousepoint %v\n", myThirdHex)
}
