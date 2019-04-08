package Camera

type Polygon struct {
	Drawer []float32
	Dist   float32
}

func MakePolygon(drawer []float32) Polygon {
	poli := Polygon{Drawer: drawer}
	return poli
}
