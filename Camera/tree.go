package Camera

import "math"

type Node struct {
	Polygons  []Polygon
	PolyList  []Polygon
	UsedPoly  []bool
	LeftNode  *Node
	RightNode *Node
	Leaf      bool
}

var EPSILON float32 = 0.000000001

func BuildTree(polygons []Polygon) *Node {
	node := &Node{}
	leftPolygons := []Polygon{}
	rightPolygons := []Polygon{}
	isLeaf := true

	for i := 1; i < len(polygons); i++ {
		if IsIntersected(polygons[0], polygons[i]) {
			lp, rp := DividePolygon(polygons[0], polygons[i])
			leftPolygons = append(leftPolygons, lp...)
			rightPolygons = append(rightPolygons, rp...)
		} else if IsBehind(polygons[0], polygons[i]) {
			leftPolygons = append(leftPolygons, polygons[i])
		} else if !IsBehind(polygons[0], polygons[i]) {
			rightPolygons = append(rightPolygons, polygons[i])
		}
	}

	if len(rightPolygons) > 0 {
		node.RightNode = BuildTree(rightPolygons)
		isLeaf = false
	}

	if len(leftPolygons) > 0 {
		node.LeftNode = BuildTree(leftPolygons)
		isLeaf = false
	}

	node.Leaf = isLeaf

	return node
}

func IsIntersected(poli1, poli2 Polygon) bool {
	return true
}

func IsBehind(poli1, poli2 Polygon) bool {
	return true
}

func DividePolygon(divisor Polygon, base Polygon) ([]Polygon, []Polygon) {
	return []Polygon{}, []Polygon{}
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
