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

```

## License

Copyright (c) 2020, hautenessa
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:
1. Redistributions of source code must retain the above copyright
   notice, this list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright
   notice, this list of conditions and the following disclaimer in the
   documentation and/or other materials provided with the distribution.
3. All advertising materials mentioning features or use of this software
   must display the following acknowledgement:
   This product includes software developed by the hautenessa.
4. Neither the name of the hautenessa nor the
   names of its contributors may be used to endorse or promote products
   derived from this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY hautenessa ''AS IS'' AND ANY
EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL hautenessa BE LIABLE FOR ANY
DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.