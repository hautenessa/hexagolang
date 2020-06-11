package hexagolang

// Hexagons implementation interpreted from
// https://www.redblobgames.com/grids/hexagons/implementation.html
// and
// https://www.redblobgames.com/grids/hexagons/

import (
	"image"
	"math"
)

// H is a single hexagon in the grid.
type H struct {
	Q, R int
}

// Delta converts the hex to a delta.
func (h H) Delta() D {
	return D{h.Q, h.R, -h.Q - h.R}
}

// Neighbor one step in a specific direction.
func (h H) Neighbor(d DirectionEnum) H {
	return Add(h, NeighborDelta(d))
}

// Float returns the cube coordinates as float values.
func (h H) Float() (float64, float64, float64) {
	return float64(h.Q), float64(-h.Q - h.R), float64(h.R)
}

// D is the amount of change between two hexagons.
type D struct {
	Q, R, S int
}

// Hex converts the delta to a hex.
func (d D) Hex() H {
	return H{d.Q, d.R}
}

// Abs returns the delta as absolute values. Cmath.Abs(delta)
func (d D) Abs() D {
	return D{
		int(math.Abs(float64(d.Q))),
		int(math.Abs(float64(d.R))),
		int(math.Abs(float64(d.S))),
	}
}

// Add is (a + b)
func Add(a H, b D) H {
	return H{
		Q: a.Q + b.Q,
		R: a.R + b.R,
	}
}

// Subtract the coordinates of the second hexagon from the first hexagon. (a - b)
func Subtract(a, b H) D {
	return D{
		Q: a.Q - b.Q,
		R: a.R - b.R,
		S: -(a.Q - b.Q) - (a.R - b.R),
	}
}

// Multiply a delta by a fixed amount (x(a))
func Multiply(d D, k int) D {
	return D{d.Q * k, d.R * k, d.S * k}
}

// RotateClockwise rotates one point around another point clockwise
func RotateClockwise(origin, moving H) H {
	before := Subtract(moving, origin)
	after := D{-before.R, -before.S, -before.Q}
	return Add(origin, after)
}

// RotateCounterClockwise rotates one point around another point counter clockwise
func RotateCounterClockwise(origin, moving H) H {
	before := Subtract(moving, origin)
	after := D{-before.S, -before.Q, -before.R}
	return Add(origin, after)
}

// Length returns the manhattan distance for a delta
func Length(d D) int {
	abs := d.Abs()
	return (abs.Q + abs.R + abs.S) >> 1
}

// Direction returns the Direction one point is in comparison to another point.
func Direction(d D) DirectionEnum {
	abs := d.Abs()
	if abs.Q >= abs.R && abs.Q >= abs.S {
		if d.Q < 0 {
			return DirectionNegQ
		}
		return DirectionPosQ
	}
	if abs.R >= abs.S {
		if d.R < 0 {
			return DirectionNegR
		}
		return DirectionPosR
	}
	if d.S < 0 {
		return DirectionNegS
	}
	return DirectionPosS
}

// DirectionEnum represents the directions of each of the sides of a hex.
type DirectionEnum int

// String returns the string name of the direction.
func (d DirectionEnum) String() string {
	ret := "DirectionUndefined"
	switch d {
	case DirectionPosQ:
		ret = "DirectionPosQ"
	case DirectionPosR:
		ret = "DirectionPosR"
	case DirectionPosS:
		ret = "DirectionPosS"
	case DirectionNegQ:
		ret = "DirectionNegQ"
	case DirectionNegR:
		ret = "DirectionNegR"
	case DirectionNegS:
		ret = "DirectionNegS"
	}
	return ret
}

// Constants for the directions from a Hex.
const (
	DirectionPosQ DirectionEnum = iota
	DirectionNegR
	DirectionPosS
	DirectionNegQ
	DirectionPosR
	DirectionNegS
	DirectionUndefined
)

var neighbors = []D{
	D{1, 0, -1}, D{1, -1, 0}, D{0, -1, 1}, // positive
	D{-1, 0, 1}, D{-1, 1, 0}, D{0, 1, -1}, // negative
	D{}, // undefined
}

// NeighborDelta returns the delta required to move a single hex in a direction.
func NeighborDelta(d DirectionEnum) D {
	return neighbors[d]
}

// Diagonal represents the direction of each point on a hex.
type Diagonal int

// String returns the string name of the direction.
func (d Diagonal) String() string {
	ret := "DiagonalUndefined"
	switch d {
	case DiagonalPosQ:
		ret = "DiagonalPosQ"
	case DiagonalPosR:
		ret = "DiagonalPosR"
	case DiagonalPosS:
		ret = "DiagonalPosS"
	case DiagonalNegQ:
		ret = "DiagonalNegQ"
	case DiagonalNegR:
		ret = "DiagonalNegR"
	case DiagonalNegS:
		ret = "DiagonalNegS"
	}
	return ret
}

// Constants for the ddiagonal from a Hex
const (
	DiagonalPosQ Diagonal = iota
	DiagonalNegR
	DiagonalPosS
	DiagonalNegQ
	DiagonalPosR
	DiagonalNegS
	DiagonalUndefined
)

var diagonals = []D{
	D{2, -1, -1}, D{1, -2, 1}, D{-1, -1, 2}, // positive
	D{-2, 1, 1}, D{-1, 2, -1}, D{1, 1, -2}, // negative
	D{}, // undefined
}

// DiagonalDelta returns the delta required to move a single hex in a direction.
func DiagonalDelta(d DirectionEnum) D {
	return diagonals[d]
}

// Line gets the hexagons on a line between two hex.
func Line(a, b H) []H {
	n := Length(Subtract(a, b))
	results := make([]H, 0, n)
	visited := make(map[H]bool, n)
	ax, ay, az := a.Float()
	bx, by, bz := b.Float()
	x, y, z := bx-ax, by-ay, bz-az

	step := 1. / float64(n+1)
	for h := 0; h <= n; h++ {
		t := step * float64(h)
		pnt := unfloat(ax+x*t, ay+y*t, az+z*t)
		if !visited[pnt] {
			results = append(results, pnt)
			visited[pnt] = true
		}
	}
	if !visited[b] {
		results = append(results, b)
	}

	return results
}

// Range returns the slice of all points in a distance from a point.
func Range(h H, rad int) map[H]bool {
	results := make(map[H]bool, rad*rad)
	if rad < 1 {
		return results
	}
	for x := -rad; x <= rad; x++ {
		for y := intMax(-rad, -x-rad); y <= intMin(rad, -x+rad); y++ {
			z := -x - y
			delta := D{
				Q: int(x),
				R: int(z),
				S: int(y),
			}
			results[Add(h, delta)] = true
		}
	}
	return results
}

// Ring returns the ring of hex points specific manhattan distance from h.
func Ring(h H, rad int) map[H]bool {
	results := make(map[H]bool)
	if rad < 1 {
		return results
	}

	h = Add(h, Multiply(NeighborDelta(DirectionPosS), rad))
	results[h] = true
	if rad > 1 {
		for i := 0; i < 6; i++ {
			for j := 0; j < rad; j++ {
				h = Add(h, NeighborDelta(DirectionEnum(i)))
				results[h] = true
			}
		}
	}
	return results
}

// unfloat returns a tuple as a Point, Rounded.
func unfloat(x, y, z float64) H {
	rx, ry, rz := math.Round(x), math.Round(y), math.Round(z)
	dx, dy, dz := math.Abs(rx-x), math.Abs(ry-y), math.Abs(rz-z)

	if dx > dz && dx > dy {
		rx = -rz - ry
	} else if dz > dy {
		rz = -rx - ry
	} else {
		ry = -rx - rz
	}
	return H{
		Q: int(math.Round(rx)),
		R: int(math.Round(rz)),
	}
}

func intMax(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func intMin(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Orientation is the orientation of the hexagon map
type Orientation struct {
	f, b [4]float64
	a    float64
	c    [6]float64
	s    [6]float64
}

// Define the default set of orientations.
var (
	OrientationPointy Orientation = Orientation{
		f: [4]float64{math.Sqrt(3.), math.Sqrt(3.) / 2., 0.0, 3. / 2.},
		b: [4]float64{math.Sqrt(3.) / 3., -1. / 3., 0.0, 2. / 3.},
		a: 0.5,
		c: [6]float64{
			math.Cos(2. * math.Pi * 0.5 / 6),
			math.Cos(2. * math.Pi * 1.5 / 6),
			math.Cos(2. * math.Pi * 2.5 / 6),
			math.Cos(2. * math.Pi * 3.5 / 6),
			math.Cos(2. * math.Pi * 4.5 / 6),
			math.Cos(2. * math.Pi * 5.5 / 6),
		},
		s: [6]float64{
			math.Sin(2. * math.Pi * 0.5 / 6),
			math.Sin(2. * math.Pi * 1.5 / 6),
			math.Sin(2. * math.Pi * 2.5 / 6),
			math.Sin(2. * math.Pi * 3.5 / 6),
			math.Sin(2. * math.Pi * 4.5 / 6),
			math.Sin(2. * math.Pi * 5.5 / 6),
		},
	}
	OrientationFlat Orientation = Orientation{
		f: [4]float64{3. / 2., 0.0, math.Sqrt(3.) / 2., math.Sqrt(3.)},
		b: [4]float64{2. / 3., 0.0, -1. / 3., math.Sqrt(3.) / 3.},
		a: 0.0,
		c: [6]float64{
			math.Cos(2. * math.Pi * 0. / 6),
			math.Cos(2. * math.Pi * 1. / 6),
			math.Cos(2. * math.Pi * 2. / 6),
			math.Cos(2. * math.Pi * 3. / 6),
			math.Cos(2. * math.Pi * 4. / 6),
			math.Cos(2. * math.Pi * 5. / 6),
		},
		s: [6]float64{
			math.Sin(2. * math.Pi * 0. / 6),
			math.Sin(2. * math.Pi * 1. / 6),
			math.Sin(2. * math.Pi * 2. / 6),
			math.Sin(2. * math.Pi * 3. / 6),
			math.Sin(2. * math.Pi * 4. / 6),
			math.Sin(2. * math.Pi * 5. / 6),
		},
	}
)

// F represents a floating point point, used for polygon drawing functions.
type F struct {
	X, Y float64
}

// Layout is the layout of the hex grid.
type Layout struct {
	size   image.Point
	origin image.Point
	m      Orientation
}

// MakeLayout for rendering on the screen.
func MakeLayout(hexSize int, originCenter image.Point, orientation Orientation) Layout {
	return Layout{
		size:   image.Point{hexSize, hexSize},
		origin: originCenter,
		m:      orientation,
	}
}

// CenterFor returns the point at the center (as an int) of the hex based on the layout.
func (l Layout) CenterFor(h H) image.Point {
	q, r :=
		float64(h.Q),
		float64(h.R)
	x := (l.m.f[0]*q + l.m.f[1]*r) * float64(l.size.X)
	y := (l.m.f[2]*q + l.m.f[3]*r) * float64(l.size.Y)
	return image.Point{int(x) + l.origin.X, int(y) + l.origin.Y}
}

// CntrFor returns the point at the center (as a float) of the hex based on the layout.
func (l Layout) CntrFor(h H) F {
	q, r :=
		float64(h.Q),
		float64(h.R)
	x := (l.m.f[0]*q + l.m.f[1]*r) * float64(l.size.X)
	y := (l.m.f[2]*q + l.m.f[3]*r) * float64(l.size.Y)
	return F{x + float64(l.origin.X), y + float64(l.origin.Y)}
}

// TopLeftFor returns the point at the top left of the hex based on the layout.
func (l Layout) TopLeftFor(h H) image.Point {
	p := l.CenterFor(h)
	p.X -= l.size.X
	p.Y -= l.size.Y
	return p
}

// HexFor an image.Point returns the Hex around a screen point.
func (l Layout) HexFor(p image.Point) H {
	x, y :=
		float64(p.X-l.origin.X),
		float64(p.Y-l.origin.Y)
	q := (l.m.b[0]*x + l.m.b[1]*y) / float64(l.size.X)
	r := (l.m.b[2]*x + l.m.b[3]*y) / float64(l.size.Y)
	return unfloat(q, -q-r, r)
}

// RingFor returns a set of hex within rad pixel distance of center.
func (l Layout) RingFor(center H, rad int) map[H]bool {
	result := make(map[H]bool, 1)
	if rad < l.size.X && rad < l.size.Y {
		result[center] = true
		return result
	}
	cp := l.CenterFor(center)
	P := 1 - rad
	pxl := image.Point{rad, 0}
	points := []image.Point{}
	for ; pxl.X > pxl.Y; pxl.Y++ {
		if P <= 0 {
			P = P + 2*pxl.Y + 1
		} else {
			pxl.X--
			P = P + 2*pxl.Y - 2*pxl.X + 1
		}

		if pxl.X < pxl.Y {
			break
		}

		points = []image.Point{
			image.Point{pxl.X + cp.X, pxl.Y + cp.Y},
			image.Point{-pxl.X + cp.X, pxl.Y + cp.Y},
			image.Point{pxl.X + cp.X, -pxl.Y + cp.Y},
			image.Point{-pxl.X + cp.X, -pxl.Y + cp.Y},
			image.Point{pxl.Y + cp.X, pxl.X + cp.Y},
			image.Point{-pxl.Y + cp.X, pxl.X + cp.Y},
			image.Point{pxl.Y + cp.X, -pxl.X + cp.Y},
			image.Point{-pxl.Y + cp.X, -pxl.X + cp.Y},
		}
		for _, v := range points {
			result[l.HexFor(v)] = true
		}
	}
	return result
}

// AreaFor returns all hex in the area of a screen circle.
func (l Layout) AreaFor(center H, rad int) map[H]bool {
	loop := l.RingFor(center, rad)
	result := make(map[H]bool)
	for k, v := range loop {
		if v == true {
			result[k] = true
			for _, inside := range Line(k, center) {
				result[inside] = true
			}
		}
	}
	return result
}

// Vertices returns the location of all verticies for a given hexagon.
func (l Layout) Vertices(h H) []F {
	result := make([]F, 6, 7)
	center := l.CntrFor(h)
	for k := range result {
		result[k] = F{
			X: center.X + float64(l.size.X)*l.m.c[k],
			Y: center.Y + float64(l.size.Y)*l.m.s[k],
		}
	}
	result = append(result, center)
	return result
}

// Radius is the size of the hexagons.
func (l Layout) Radius() image.Point {
	return l.size
}
