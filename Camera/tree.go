package Camera

import (
	"fmt"
	"log"
	"math"
	"math/rand"

	"github.com/akavel/polyclip-go"
	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/kanister10l/GoCamera/Helpers"
	"github.com/tchayen/triangolatte"
)

type Node struct {
	RightNode   *Node
	Leaf        bool
	ToRender    []Polygon
	RenderReady bool
}

var EPSILON float32 = 0.000000001

func recoverName() {
	if r := recover(); r != nil {
		fmt.Println("recovered from ", r)
	}
}

func BuildTree(polygons []Polygon) *Node {
	node := &Node{}
	node.ToRender = []Polygon{}
	rightPolygons := []Polygon{}
	isLeaf := true
	reversed := false

	for i := 1; i < len(polygons); i++ {
		if n := IsIntersected(polygons[0], polygons[i]); n > 0 && n < 3 {
			rp, rev := DividePolygon(polygons[0], polygons[i])
			if !reversed && rev {
				reversed = true
			}
			rightPolygons = append(rightPolygons, rp...)
		} else if !IsBehind(polygons[0], polygons[i]) && n == 0 {
			rightPolygons = append(rightPolygons, polygons[i])
		} else if n1 := IsIntersected(polygons[i], polygons[0]); n == 3 && n1 == 3 && !IsBehind(polygons[0], polygons[i]) {
			node.ToRender = append(node.ToRender, polygons[i])
			node.RenderReady = true
		}
	}

	if reversed {
		rightPolygons = append(rightPolygons, polygons[0])
	} else {
		node.ToRender = append(node.ToRender, polygons[0])
		node.RenderReady = true
	}

	if len(rightPolygons) > 0 {
		node.RightNode = BuildTree(rightPolygons)
		isLeaf = false
	}

	node.Leaf = isLeaf

	return node
}

func Traverse(node *Node) {
	if !node.Leaf {
		Traverse(node.RightNode)
	}

	if node.RenderReady {
		for _, v := range node.ToRender {
			gl.BindVertexArray(Helpers.MakeVao(v.Drawer, []float32{}, true))
			gl.DrawArrays(gl.TRIANGLES, 0, int32(len(v.Drawer)/3))
		}
	}
}

func IsIntersected(poli1, poli2 Polygon) int {
	n := 0
	for i := 0; i < 8; i += 3 {
		if !IsInside(poli1, poli2.Drawer[i], poli2.Drawer[i+1]) {
			n++
		}
	}
	return n
}

func IsBehind(poli1, poli2 Polygon) bool {
	if poli1.Dist < poli2.Dist {
		return true
	}

	return false
}

func DividePolygon(divisor Polygon, base Polygon) ([]Polygon, bool) {
	defer recoverName()
	var subject polyclip.Polygon
	var clipper polyclip.Polygon
	reversed := false
	d := base.Dist

	if IsBehind(divisor, base) {
		subject = polyclip.Polygon{{{float64(base.Drawer[0]), float64(base.Drawer[1])},
			{float64(base.Drawer[3]), float64(base.Drawer[4])},
			{float64(base.Drawer[6]), float64(base.Drawer[7])}}}
		clipper = polyclip.Polygon{{{float64(divisor.Drawer[0]), float64(divisor.Drawer[1])},
			{float64(divisor.Drawer[3]), float64(divisor.Drawer[4])},
			{float64(divisor.Drawer[6]), float64(divisor.Drawer[7])}}}
	} else {
		reversed = true
		d = divisor.Dist
		subject = polyclip.Polygon{{{float64(divisor.Drawer[0]), float64(divisor.Drawer[1])},
			{float64(divisor.Drawer[3]), float64(divisor.Drawer[4])},
			{float64(divisor.Drawer[6]), float64(divisor.Drawer[7])}}}

		clipper = polyclip.Polygon{{{float64(base.Drawer[0]), float64(base.Drawer[1])},
			{float64(base.Drawer[3]), float64(base.Drawer[4])},
			{float64(base.Drawer[6]), float64(base.Drawer[7])}}}
	}

	result := subject.Construct(polyclip.DIFFERENCE, clipper)

	newTr := []Polygon{}

	/*for _, v := range result {
		contour := []*poly2tri.Point{}
		for _, v2 := range v {
			contour = append(contour, poly2tri.NewPoint(v2.X, v2.Y))
		}

		swctx := poly2tri.NewSweepContext(contour, false)
		swctx.Triangulate()

		tr := swctx.GetTriangles()
		for _, v2 := range tr {
			newTr = append(newTr, Polygon{Drawer: []float32{
				float32(v2.Points[0].X), float32(v2.Points[0].Y), 0,
				float32(v2.Points[1].X), float32(v2.Points[1].Y), 0,
				float32(v2.Points[2].X), float32(v2.Points[2].Y), 0,
			}, Dist: d})
		}
	}*/

	for _, v := range result {
		contour := []triangolatte.Point{}
		for _, v2 := range v {
			contour = append(contour, triangolatte.Point{X: v2.X + (rand.Float64()*2-1)*0.01, Y: v2.Y + (rand.Float64()*2-1)*0.01})
		}

		tr, err := triangolatte.Polygon(contour)
		if err != nil {
			log.Println(err)
		}

		for i := 0; i < len(tr); i += 6 {
			newTr = append(newTr, Polygon{Drawer: []float32{
				float32(tr[i]), float32(tr[i+1]), 0,
				float32(tr[i+2]), float32(tr[i+3]), 0,
				float32(tr[i+4]), float32(tr[i+5]), 0,
			}, Dist: d})
		}
	}

	return newTr, reversed
}

func Orientation(p1x, p1y, p2x, p2y, p3x, p3y float32) int {
	val := (p2y-p1y)*(p3x-p2x) - (p2x-p1x)*(p3y-p2y)
	if FloatEqual(val, 0) {
		return 0
	} else if val > 0 {
		return 1
	}

	return 2
}

func Area(x1, y1, x2, y2, x3, y3 float32) float32 {
	return float32(math.Abs(float64((x1*(y2-y3) + x2*(y3-y1) + x3*(y1-y2)) / 2.0)))
}

func IsInside(poly Polygon, x, y float32) bool {
	a := Area(poly.Drawer[0], poly.Drawer[1], poly.Drawer[3], poly.Drawer[4], poly.Drawer[6], poly.Drawer[7])
	a1 := Area(x, y, poly.Drawer[3], poly.Drawer[4], poly.Drawer[6], poly.Drawer[7])
	a2 := Area(poly.Drawer[0], poly.Drawer[1], x, y, poly.Drawer[6], poly.Drawer[7])
	a3 := Area(poly.Drawer[0], poly.Drawer[1], poly.Drawer[3], poly.Drawer[4], x, y)

	return FloatEqual(a1+a2+a3, a)
}

func FloatEqual(a, b float32) bool {
	return float32(math.Abs(float64(a-b))) < EPSILON
}
