package hexagolang

import (
	"image"
	"testing"
)

// I need to translate between screen coordinates and hex coordinates.
// I need to translate between hex coordinates and screen coordinates.
// Rational, the grid will be displayed on the screen and the user interacts with the screen.
func TestScreenConversion(t *testing.T) {
	layout := MakeLayout(10, image.Point{0, 0}, OrientationFlat)
	plan := []struct {
		hp  H
		ip  image.Point
		tlp image.Point
	}{
		{H{0, 0}, image.Point{0, 0}, image.Point{-10, -10}},
		{H{2, -1}, image.Point{30, 0}, image.Point{20, -10}},
		{H{-2, 4}, image.Point{-30, 51}, image.Point{-40, 41}},
	}

	for _, expected := range plan {
		if result := layout.CenterFor(expected.hp); expected.ip != result {
			t.Errorf("hex %+v pixel center expected %+v, got %+v.",
				expected.hp, expected.ip, result)
		}
		if result := layout.HexFor(expected.ip); expected.hp != result {
			t.Errorf("image %+v hex for expected %+v, got %+v.",
				expected.ip, expected.hp, result)
		}
		if result := layout.TopLeftFor(expected.hp); expected.tlp != result {
			t.Errorf("hex %+v pixel topleft expected %+v, got %+v.",
				expected.hp, expected.tlp, result)
		}
	}
}

// I need to know the set of hex within a screen distance from a hex.
// Rational, I'd like to use circules instead of Hexagons for the larger areas.
func TestRadiusFor(t *testing.T) {
	layout1, layout2 := MakeLayout(10, image.Point{0, 0}, OrientationFlat),
		MakeLayout(100, image.Point{0, 0}, OrientationPointy)
	plan := []struct {
		lay    Layout
		center H
		rad    int
		pos    []H
		neg    []H
	}{
		{layout1, H{0, 0}, -1,
			[]H{{0, 0}},
			[]H{{1, 0}},
		},
		{layout1, H{0, 0}, 0,
			[]H{{0, 0}},
			[]H{{1, 0}},
		},
		{layout1, H{0, 0}, 1,
			[]H{{0, 0}},
			[]H{{1, 0}},
		},
		{layout1, H{0, 0}, 11,
			[]H{{-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}},
			[]H{{0, 0}, {1, 1}, {1, 2}},
		},
		{layout2, H{200, 200}, 500,
			[]H{
				{196, 202}, {197, 200}, {197, 201}, {197, 202}, {197, 203}, {198, 198},
				{198, 199}, {198, 203}, {199, 198}, {199, 203}, {200, 197}, {200, 203},
				{201, 197}, {201, 202}, {202, 197}, {202, 201}, {203, 197}, {203, 198},
				{203, 199}, {203, 200}},
			[]H{{1, 0}, {200, 200}},
		},
	}

	for tc, params := range plan {
		result := params.lay.RingFor(params.center, params.rad)
		if len(result) != len(params.pos) {
			t.Errorf("index %d: Expected %d results, got %d.", tc, len(params.pos), len(result))
			t.Logf("result was %+v", result)
		}
		for k, v := range params.pos {
			if !result[v] {
				t.Errorf("index %d-%d: positive result expected for %+v. got false", tc, k, v)
				t.Logf("result was %+v", result)
			}
		}
		for k, v := range params.neg {
			if result[v] {
				t.Errorf("index %d-%d: negative result expected for %+v. Got true", tc, k, v)
				t.Logf("result was %+v", result)
			}
		}
	}
}

// I need to be able to get the screen size of a hex in a grid.
// Rational, this will be necessary for consumers to understand the grid
func TestRadius(t *testing.T) {
	plan := []struct {
		l Layout
		r int
	}{
		{MakeLayout(10, image.Point{0, 0}, OrientationPointy), 10},
		{MakeLayout(100, image.Point{0, 0}, OrientationFlat), 100},
		{MakeLayout(1000, image.Point{10, 10}, OrientationPointy), 1000},
	}
	// This will get expanded when I support skewed grids.
	for k, expected := range plan {
		if result := expected.l.Radius; expected.r != result.X {
			t.Errorf("index %d: expected X %d, got %d", k, expected.r, result)
		}
		if result := expected.l.Radius; expected.r != result.Y {
			t.Errorf("index %d: expected Y %d, got %d", k, expected.r, result)
		}
	}
}

// I need to know the distance between two hex.
// I need to know the direction between two hex.
// Rational, items must be able to "face" a direction if necessary.
func TestDistance(t *testing.T) {
	plan := []struct {
		a, b H
		dir  DirectionEnum
		dist int
	}{
		{H{2, 1}, H{10, 0}, DirectionNegQ, 8},
		{H{1, 2}, H{10, 10}, DirectionPosS, 17},
		{H{0, 0}, H{0, 10}, DirectionNegR, 10},
		{H{0, 0}, H{-10, 0}, DirectionPosQ, 10},
		{H{0, 0}, H{-10, -10}, DirectionNegS, 20},
		{H{0, 0}, H{2, -10}, DirectionPosR, 10},
		{H{0, 0}, H{20, -20}, DirectionNegQ, 20},
		{H{0, 0}, H{-10, 10}, DirectionPosQ, 10},
		{H{-10, 10}, H{0, 0}, DirectionNegQ, 10},
		{H{100, 100}, H{-10, -10}, DirectionNegS, 220},
	}

	for k, expected := range plan {
		if result := Length(Subtract(expected.a, expected.b)); expected.dist != result {
			t.Errorf("index %d: expected distance %d, got %d", k, expected.dist, result)
		}
		if result := Direction(Subtract(expected.a, expected.b)); expected.dir != result {
			t.Errorf("index %d: expected direction %s, got %s", k, expected.dir, result)
		}
	}
}

// I need to draw a line between two hex.
// Rational, needed in path planning, and drawing screen circles.
func TestLine(t *testing.T) {
	plan := []struct {
		a, b H
		line []H
	}{
		{H{2, 1}, H{10, 10},
			[]H{
				{2, 1}, {2, 2}, {3, 2}, {3, 3}, {4, 3}, {4, 4}, {5, 4}, {5, 5},
				{6, 5}, {6, 6}, {7, 6}, {7, 7}, {8, 7}, {8, 8}, {9, 8}, {9, 9},
				{10, 9}, {10, 10}}},
		{H{10, 10}, H{3, 1},
			[]H{
				{10, 10}, {10, 9}, {9, 9}, {9, 8}, {8, 8}, {8, 7}, {7, 7}, {7, 6},
				{7, 5}, {6, 5}, {6, 4}, {5, 4}, {5, 3}, {4, 3}, {4, 2}, {3, 2},
				{3, 1}}},
		{H{-4, 3}, H{0, 0},
			[]H{
				{-4, 3}, {-3, 2}, {-2, 2}, {-1, 1}, {0, 0}}},
		{H{4, 1}, H{16, 4},
			[]H{
				{4, 1}, {5, 1}, {6, 1}, {6, 2}, {7, 2}, {8, 2}, {9, 2}, {10, 2},
				{10, 3}, {11, 3}, {12, 3}, {13, 3}, {14, 3}, {14, 4}, {15, 4},
				{16, 4}}},
		{H{3, 3}, H{4, 12},
			[]H{
				{3, 3}, {3, 4}, {3, 5}, {3, 6}, {3, 7}, {4, 7}, {4, 8}, {4, 9},
				{4, 10}, {4, 11}, {4, 12}}},
		{H{9, 5}, H{15, 11},
			[]H{
				{9, 5}, {10, 5}, {10, 6}, {11, 6}, {11, 7}, {12, 7}, {12, 8},
				{13, 8}, {13, 9}, {14, 9}, {14, 10}, {15, 10}, {15, 11}}},
	}

	for tc, expected := range plan {
		result := Line(expected.a, expected.b)
		if len(result) != len(expected.line) {
			t.Errorf("Index %d: Expected %d steps, got %d", tc, len(expected.line), len(result))
			t.Logf("\tResults were %#v", result)
		}
		for k := 0; k < intMin(len(expected.line), len(result)); k++ {
			if expected.line[k] != result[k] {
				t.Errorf("Index %d-%d: Expected %+v, got %+v", tc, k, expected.line[k], result[k])
				t.Logf("\tResults were %#v", result)
			}
		}
	}
}

// I need to know the set of hex within a hex distance from a given hex.
// Rational, needed for calculating move distance and weapon range.
func TestArea(t *testing.T) {
	plan := []struct {
		a   H
		rad int
		pos []H
		neg []H
	}{
		{H{0, 0}, 2,
			[]H{
				{1, 0}, {2, -2}, {-1, 2}, {0, 2}, {0, 1}, {0, -2}, {-2, 1}, {-1, 0},
				{0, 0}, {1, -2}, {-2, 0}, {-1, 1}, {-1, -1}, {0, -1}, {2, -1}, {-2, 2},
				{1, 1}, {1, -1}, {2, 0},
			},
			[]H{
				{2, 1}, {1, 2}, {-1, -2}, {-2, -1}, {2, 2}, {-2, -2},
			},
		},
	}

	for tc, params := range plan {
		result := Range(params.a, params.rad)
		if len(result) != len(params.pos) {
			t.Errorf("index %d: Expected %d results, got %d.", tc, len(params.pos), len(result))
			t.Logf("result was %+v", result)
		}
		for k, v := range params.pos {
			if !result[v] {
				t.Errorf("index %d-%d: positive result expected for %+v. got false", tc, k, v)
				t.Logf("result was %+v", result)
			}
		}
		for k, v := range params.neg {
			if result[v] {
				t.Errorf("index %d-%d: negative result expected for %+v. Got true", tc, k, v)
				t.Logf("result was %+v", result)
			}
		}
	}
}

// I need to perform Vertex operations on a hex.
// rational, needed to draw the grid and this allows me to compute triangles.
func TestVertices(t *testing.T) {
	plan := []struct {
		lay Layout
		a   H
		v   []F
	}{
		{MakeLayout(10, image.Point{0, 0}, OrientationPointy), H{0, 0},
			[]F{
				{8.6603, 5}, {0, 10}, {-8.6603, 5},
				{-8.6603, -5}, {0, -10}, {8.6603, -5},
				{0, 0},
			},
		},
		{MakeLayout(20, image.Point{0, 0}, OrientationPointy), H{0, 0},
			[]F{
				{17.3205, 10}, {0, 20}, {-17.3205, 10},
				{-17.3205, -10}, {0, -20}, {17.3205, -10},
				{0, 0},
			},
		},
		{MakeLayout(10, image.Point{0, 0}, OrientationFlat), H{3, -2},
			[]F{
				{55, -8.6603}, {50, 0}, {40, 0},
				{35, -8.6603}, {40, -17.3205}, {50, -17.3205},
				{45, -8.6603},
			},
		},
		{MakeLayout(10, image.Point{0, 0}, OrientationFlat), H{3, -1},
			[]F{
				{55, 8.6603}, {50, 17.3205}, {40, 17.3205},
				{35, 8.6603}, {40, 0}, {50, 0},
				{45, 8.6603},
			},
		},
		{MakeLayout(20, image.Point{40, 40}, OrientationPointy), H{4, 6},
			[]F{
				{299.8076, 230}, {282.4871, 240}, {265.1666, 230},
				{265.1666, 210}, {282.4871, 200}, {299.8076, 210},
				{282.4871, 220},
			},
		},
	}

	for tc, params := range plan {
		result := params.lay.Vertices(params.a)
		for k := range result {
			deltaX := result[k].X - params.v[k].X
			deltaY := result[k].Y - params.v[k].Y
			if -0.0001 > deltaX || deltaX > 0.0001 || -0.0001 > deltaY || deltaY > 0.0001 {
				t.Errorf("index %d: vertex %d expected %+v, got %+v", tc, k, params.v[k], result[k])
				t.Logf("\t%#v", result)
			}
		}
	}

	plan2 := []struct {
		lay    Layout
		a, b   H
		va, vb int
	}{
		{MakeLayout(32, image.Point{}, OrientationPointy), H{0, 0}, H{1, 0}, 0, 2},
		{MakeLayout(32, image.Point{}, OrientationPointy), H{1, 0}, H{2, 0}, 0, 2},
		{MakeLayout(32, image.Point{}, OrientationPointy), H{2, 0}, H{3, 0}, 0, 2},
		{MakeLayout(32, image.Point{}, OrientationFlat), H{0, 0}, H{1, 0}, 0, 4},
		{MakeLayout(32, image.Point{}, OrientationFlat), H{1, 0}, H{2, 0}, 0, 4},
		{MakeLayout(32, image.Point{}, OrientationFlat), H{2, 0}, H{3, 0}, 0, 4},
	}

	for tc, params := range plan2 {
		r1, r2 := params.lay.Vertices(params.a), params.lay.Vertices(params.b)
		deltaX := r1[params.va].X - r2[params.vb].X
		deltaY := r1[params.va].Y - r2[params.vb].Y
		if -0.0001 > deltaX || deltaX > 0.0001 || -0.0001 > deltaY || deltaY > 0.0001 {
			t.Errorf("index %d: vertex %v doesn't equal matching neighbor vertex %v", tc, r1[params.va], r2[params.vb])
			t.Logf("\t%+v vs %+v", r1, r2)
			for k := range r1 {
				t.Logf("%d: %f - %f == %f, %f", k, r2[k], r1[k], r2[k].X-r1[k].X, r2[k].Y-r1[k].Y)
			}
		}
	}
}

// I need all features to be fast for a game engine.
// needs definition of fast.
// rational, these functions will be invoked frequently as part of calculating the game.
func BenchmarkScreenToHex(b *testing.B) {
	layout := MakeLayout(64, image.Point{}, OrientationPointy)
	for h := 0; h < b.N; h++ {
		layout.HexFor(image.Point{h, h})
	}
}

func BenchmarkScreenRing(b *testing.B) {
	layout := MakeLayout(64, image.Point{}, OrientationPointy)
	for h := 0; h < b.N; h++ {
		layout.RingFor(H{h, h}, 512)
	}
}

func BenchmarkScreenArea(b *testing.B) {
	layout := MakeLayout(64, image.Point{}, OrientationPointy)
	for h := 0; h < b.N; h++ {
		layout.AreaFor(H{h, h}, 512)
	}
}

func BenchmarkHexToScreen(b *testing.B) {
	layout := MakeLayout(64, image.Point{}, OrientationPointy)
	for h := 0; h < b.N; h++ {
		layout.CenterFor(H{h, h})
	}
}

func BenchmarkLine(b *testing.B) {
	for h := 0; h < b.N; h++ {
		Line(H{256, 256}, H{-256, 256})
	}
}

func BenchmarkRing(b *testing.B) {
	for h := 0; h < b.N; h++ {
		Ring(H{20, 20}, 40)
	}
}

func BenchmarkArea(b *testing.B) {
	for h := 0; h < b.N; h++ {
		Range(H{20, 20}, 40)
	}
}
