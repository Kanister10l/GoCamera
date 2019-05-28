package Camera

import (
	"math"

	"github.com/gerow/go-color"
	"github.com/go-gl/mathgl/mgl32"
)

type SpherePoint struct {
	X         float32
	Y         float32
	Z         float32
	Nvector   mgl32.Vec3
	Intensity float32
	Layer     int
}

type SphereWorld struct {
	X       float32
	Y       float32
	Z       float32
	XOrigin float32
	YOrigin float32
	ZOrigin float32
	AngleH  float32
	AngleV  float32
	Radius  float32
	Ka      float32
	Kd      float32
	Ks      float32
	Hue     float32
	N       int
}

type SpherePolygon struct {
	Drawer []float32
	Color  []float32
}

func CreateSphereWorld(xOrigin, yOrigin, zOrigin, dist float32) *SphereWorld {
	sp := &SphereWorld{}
	sp.XOrigin = xOrigin
	sp.YOrigin = yOrigin
	sp.ZOrigin = zOrigin
	sp.Radius = dist
	sp.AngleH = float32(math.Pi)
	sp.AngleV = 0
	sp.Hue = 0.5
	sp.Ka = 0
	sp.Kd = 0.5
	sp.Ks = 0.5
	sp.N = 50
	sp.Update()
	return sp
}

func (s *SphereWorld) Rotate(horizontalDelta float32, verticalDelta float32) {
	s.AngleH += horizontalDelta
	s.AngleV += verticalDelta

	if s.AngleH < 0 {
		s.AngleH = float32(2*math.Pi) + s.AngleH
	} else if s.AngleH > float32(2*math.Pi) {
		s.AngleH = s.AngleH - float32(2*math.Pi)
	}

	if s.AngleV < 0 {
		s.AngleV = float32(2*math.Pi) + s.AngleV
	} else if s.AngleV > float32(2*math.Pi) {
		s.AngleV = s.AngleV - float32(2*math.Pi)
	}

	s.Update()
}

func (sp *SphereWorld) ModifyConstant(a, d, s, h float32, n int) {
	sp.Ka += a
	sp.Kd += d
	sp.Ks += s
	sp.N += n
	sp.Hue += h

	if sp.Ka < 0 {
		sp.Ka = 0
	} else if sp.Ka > 1 {
		sp.Ka = 1
	}

	if sp.Kd < 0 {
		sp.Kd = 0
	} else if sp.Kd > 1 {
		sp.Kd = 1
	}

	if sp.Ks < 0 {
		sp.Ks = 0
	} else if sp.Ks > 1 {
		sp.Ks = 1
	}

	if sp.Hue < 0 {
		sp.Hue = 0
	} else if sp.Hue > 1 {
		sp.Hue = 1
	}

	if sp.N < 1 {
		sp.N = 1
	}
}

func (s *SphereWorld) Update() {
	s.X = s.Radius * float32(math.Sin(float64(s.AngleH))) * float32(math.Cos(float64(s.AngleV)))
	s.Y = s.Radius * float32(math.Sin(float64(s.AngleH))) * float32(math.Sin(float64(s.AngleV)))
	s.Z = s.Radius * float32(math.Cos(float64(s.AngleH)))
}

func GenerateSphere(r, xOrigin, yOrigin, zOrigin float32, resolution int, angleResolution int) []SpherePoint {
	points := []SpherePoint{}
	gap := r / float32(resolution-1)
	angleDelta := math.Pi * 2 / float32(angleResolution)
	position := float32(0)

	newPoint := SpherePoint{}
	newPoint.X = xOrigin
	newPoint.Y = yOrigin
	newPoint.Z = zOrigin - r
	newPoint.Nvector = mgl32.NewVecNFromData([]float32{newPoint.X - xOrigin, newPoint.Y - yOrigin, newPoint.Z - zOrigin}).Vec3().Normalize()
	newPoint.Layer = 0
	newPoint.Intensity = 0
	points = append(points, newPoint)
	position += gap

	for k := 1; k < resolution; k++ {
		angle := float32(0)
		for o := 0; o < angleResolution; o++ {
			newPoint := SpherePoint{}
			newPoint.X = xOrigin + position*float32(math.Cos(float64(angle)))
			newPoint.Y = yOrigin + position*float32(math.Sin(float64(angle)))
			zOffset := float32(math.Sqrt(math.Pow(float64(r), 2) - math.Pow(float64(position), 2)))
			if math.IsNaN(float64(zOffset)) {
				zOffset = 0
			}
			newPoint.Z = zOrigin - zOffset
			newPoint.Nvector = mgl32.NewVecNFromData([]float32{newPoint.X - xOrigin, newPoint.Y - yOrigin, newPoint.Z - zOrigin}).Vec3().Normalize()
			newPoint.Layer = k
			newPoint.Intensity = 0
			points = append(points, newPoint)
			angle += angleDelta
		}

		position += gap
	}

	return points
}

func CalculateLightIntensity(points []SpherePoint, xSource, ySource, zSource, xWatch, yWatch, zWatch, ka, kd, ks float32, n int) []SpherePoint {
	intensity := float32(0.0)
	max := float32(0.0)
	for k := range points {
		lightVec := mgl32.NewVecNFromData([]float32{points[k].X - xSource, points[k].Y - ySource, points[k].Z - zSource}).Vec3().Normalize()
		watchVec := mgl32.NewVecNFromData([]float32{xWatch - points[k].X, yWatch - points[k].Y, zWatch - points[k].Z}).Vec3().Normalize()

		if points[k].Nvector.Dot(watchVec) >= 0 && points[k].Nvector.Dot(lightVec.Mul(-1)) >= 0 {
			reflectionVec := lightVec.Sub(points[k].Nvector.Mul(lightVec.Dot(points[k].Nvector) * 2))
			intensity = ka + kd*points[k].Nvector.Dot(lightVec.Mul(-1)) + ks*float32(math.Pow(float64(reflectionVec.Dot(watchVec)), float64(n)))
			points[k].Intensity = intensity
			if intensity > max {
				max = intensity
			}
		} else {
			points[k].Intensity = ka
		}
	}

	if max > 1 {
		for k := range points {
			points[k].Intensity /= max
		}
	}

	return points
}

func Polygonyfy(points []SpherePoint, xCanvasSize, yCanvasSize, hue float32) []SpherePolygon {
	poly := []SpherePolygon{}
	hsl := color.HSL{H: float64(hue), S: 1, L: 0}

	//calc layer capacity
	lc := 0
	for _, point := range points {
		if point.Layer == 1 {
			lc++
		} else if point.Layer > 1 {
			break
		}
	}

	for i := 0; i < lc; i++ {
		xp := i + 2
		if i == lc-1 {
			xp = 1
		}

		hsl.L = float64(points[0].Intensity)
		p1L := hsl.ToRGB()
		hsl.L = float64(points[i+1].Intensity)
		p2L := hsl.ToRGB()
		hsl.L = float64(points[xp].Intensity)
		p3L := hsl.ToRGB()
		poly = append(poly, SpherePolygon{
			Drawer: []float32{
				points[0].X / xCanvasSize, points[0].Y / yCanvasSize, 0,
				points[i+1].X / xCanvasSize, points[i+1].Y / yCanvasSize, 0,
				points[xp].X / xCanvasSize, points[xp].Y / yCanvasSize, 0,
			},
			Color: []float32{
				float32(p1L.R), float32(p1L.G), float32(p1L.B),
				float32(p2L.R), float32(p2L.G), float32(p2L.B),
				float32(p3L.R), float32(p3L.G), float32(p3L.B),
			},
		})
	}

	for i := 1; i < len(points)-lc; i += lc {
		for j := 0; j < lc; j++ {
			xp := j + 1
			xm := j - 1
			if j == 0 {
				xm = lc - 1
			}
			if j == lc-1 {
				xp = 0
			}

			//Polygon 1
			hsl.L = float64(points[i+j].Intensity)
			p1L := hsl.ToRGB()
			hsl.L = float64(points[i+j+lc].Intensity)
			p2L := hsl.ToRGB()
			hsl.L = float64(points[i+xp+lc].Intensity)
			p3L := hsl.ToRGB()
			poly = append(poly, SpherePolygon{
				Drawer: []float32{
					points[i+j].X / xCanvasSize, points[i+j].Y / yCanvasSize, 0,
					points[i+j+lc].X / xCanvasSize, points[i+j+lc].Y / yCanvasSize, 0,
					points[i+xp+lc].X / xCanvasSize, points[i+xp+lc].Y / yCanvasSize, 0,
				},
				Color: []float32{
					float32(p1L.R), float32(p1L.G), float32(p1L.B),
					float32(p2L.R), float32(p2L.G), float32(p2L.B),
					float32(p3L.R), float32(p3L.G), float32(p3L.B),
				},
			})

			//Polygon 2
			hsl.L = float64(points[i+j].Intensity)
			p1L = hsl.ToRGB()
			hsl.L = float64(points[i+j+lc].Intensity)
			p2L = hsl.ToRGB()
			hsl.L = float64(points[i+xm].Intensity)
			p3L = hsl.ToRGB()
			poly = append(poly, SpherePolygon{
				Drawer: []float32{
					points[i+j].X / xCanvasSize, points[i+j].Y / yCanvasSize, 0,
					points[i+j+lc].X / xCanvasSize, points[i+j+lc].Y / yCanvasSize, 0,
					points[i+xm].X / xCanvasSize, points[i+xm].Y / yCanvasSize, 0,
				},
				Color: []float32{
					float32(p1L.R), float32(p1L.G), float32(p1L.B),
					float32(p2L.R), float32(p2L.G), float32(p2L.B),
					float32(p3L.R), float32(p3L.G), float32(p3L.B),
				},
			})
		}
	}

	return poly
}
