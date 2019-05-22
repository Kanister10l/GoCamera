package main

import (
	"fmt"

	"github.com/ByteArena/poly2tri-go"
	"github.com/akavel/polyclip-go"
)

func main() {
	subject := polyclip.Polygon{{{-2, 1}, {1, 1}, {1, -2}, {-2, -2}}}                      // small square
	clipping := polyclip.Polygon{{{-1 - 10, -1}, {-1 - 10, 2}, {2 - 10, 2}, {2 - 10, -1}}} // overlapping triangle
	result := subject.Construct(polyclip.DIFFERENCE, clipping)

	for _, v := range result {
		contour := []*poly2tri.Point{}
		for _, v2 := range v {
			contour = append(contour, poly2tri.NewPoint(v2.X, v2.Y))
		}

		swctx := poly2tri.NewSweepContext(contour, false)
		swctx.Triangulate()

		fmt.Println(swctx.GetTriangles()[0].Points[1])
	}
}
