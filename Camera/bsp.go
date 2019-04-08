package Camera

import (
	"github.com/kanister10l/GoCamera/World"
	"github.com/pkg/errors"
	"log"
)

var SideMarker = [6][4]int{
	{0, 1, 4, 5},
	{1, 2, 5, 6},
	{2, 3, 6, 7},
	{0, 3, 4, 7},
	{0, 1, 2, 3},
	{4, 5, 6, 7},
}

type BSPFigure struct {
	Frames []BSPFrame
}

type BSPFrame struct {
	Lines    []BSPLine
	Enclosed bool
	Fixed    bool
	Side     int
	Visible  bool
}

type BSPLine struct {
	Line     World.Line
	P1Dist   float32
	P2Dist   float32
	P1AngleX float32
	P1AngleY float32
	P2AngleX float32
	P2AngleY float32
}

func (f *BSPFigure) Init() {
	f.Frames = make([]BSPFrame, 6)

	for k := range f.Frames {
		f.Frames[k] = BSPFrame{}
		f.Frames[k].Init(k)
	}
}

func (f *BSPFigure) AddLine(line World.Line, d1, d2, p1AngleX, p1AngleY, p2AngleX, p2AngleY float32, visible bool) {
	s1, s2 := f.PossibleSides(line)
	f.Frames[s1].AddLine(line, d1, d2, p1AngleX, p1AngleY, p2AngleX, p2AngleY, visible)
	f.Frames[s2].AddLine(line, d1, d2, p1AngleX, p1AngleY, p2AngleX, p2AngleY, visible)
}

func (f *BSPFigure) PossibleSides(line World.Line) (int, int) {
	s1 := -1
	s2 := -1

	p1 := line.P1
	p2 := line.P2

	for k, v := range SideMarker {
		matchCounter := 0
		for _, v2 := range v {
			if p1 == v2 || p2 == v2 {
				matchCounter++
			}
		}

		if matchCounter == 2 {
			if s1 == -1 {
				s1 = k
			} else {
				s2 = k
			}
		}
	}

	if s1 == -1 || s2 == -1 {
		log.Println("Error Finding Possible Sides")
	}

	return s1, s2
}

func (f *BSPFrame) Init(side int) {
	f.Lines = []BSPLine{}
	f.Enclosed = false
	f.Fixed = true
	f.Side = side
	f.Visible = false
}

func (f *BSPFrame) AddLine(line World.Line, d1, d2, p1AngleX, p1AngleY, p2AngleX, p2AngleY float32, visible bool) {
	f.Lines = append(f.Lines, BSPLine{Line: line, P1Dist: d1, P2Dist: d2, P1AngleX: p1AngleX, P1AngleY: p1AngleY, P2AngleX: p2AngleX, P2AngleY: p2AngleY})
	if !f.Visible {
		f.Visible = visible
	}
}

func (f *BSPFrame) ConvertToPolygons() (Polygon, Polygon, error) {
	if len(f.Lines) < 4 {
		return Polygon{}, Polygon{}, errors.New("Not enough lines")
	}
	var poli1 Polygon
	var poli2 Polygon
mainFor:
	for i := 1; i < 4; i++ {
		if f.Lines[i].Line.P1 == f.Lines[0].Line.P1 || f.Lines[i].Line.P2 == f.Lines[0].Line.P1 {
			poli1 = MakePolygon([]float32{
				f.Lines[i].P1AngleX, f.Lines[i].P1AngleY, 0,
				f.Lines[i].P2AngleX, f.Lines[i].P2AngleY, 0,
				f.Lines[0].P2AngleX, f.Lines[0].P2AngleY, 0,
			})

			pMaker := i
			sp1 := 0
			sp2 := 0
			for k := 1; k < 4; k++ {
				if k != pMaker && sp1 == 0 {
					sp1 = k
				} else if k != pMaker {
					sp2 = k
				}
			}

			if f.Lines[sp1].Line.P1 == f.Lines[sp2].Line.P1 || f.Lines[sp1].Line.P2 == f.Lines[sp2].Line.P1 {
				poli2 = MakePolygon([]float32{
					f.Lines[sp1].P1AngleX, f.Lines[sp1].P1AngleY, 0,
					f.Lines[sp1].P2AngleX, f.Lines[sp1].P2AngleY, 0,
					f.Lines[sp2].P2AngleX, f.Lines[sp2].P2AngleY, 0,
				})
			} else if f.Lines[sp1].Line.P1 == f.Lines[sp2].Line.P2 || f.Lines[sp1].Line.P2 == f.Lines[sp2].Line.P2 {
				poli1 = MakePolygon([]float32{
					f.Lines[sp1].P1AngleX, f.Lines[sp1].P1AngleY, 0,
					f.Lines[sp1].P2AngleX, f.Lines[sp1].P2AngleY, 0,
					f.Lines[sp2].P1AngleX, f.Lines[sp2].P1AngleY, 0,
				})
			}
		} else if f.Lines[i].Line.P1 == f.Lines[0].Line.P2 || f.Lines[i].Line.P2 == f.Lines[0].Line.P2 {
			poli1 = MakePolygon([]float32{
				f.Lines[i].P1AngleX, f.Lines[i].P1AngleY, 0,
				f.Lines[i].P2AngleX, f.Lines[i].P2AngleY, 0,
				f.Lines[0].P1AngleX, f.Lines[0].P1AngleY, 0,
			})

			pMaker := i
			sp1 := 0
			sp2 := 0
			for k := 1; k < 4; k++ {
				if k != pMaker && sp1 == 0 {
					sp1 = k
				} else if k != pMaker {
					sp2 = k
				}
			}

			if f.Lines[sp1].Line.P1 == f.Lines[sp2].Line.P1 || f.Lines[sp1].Line.P2 == f.Lines[sp2].Line.P1 {
				poli2 = MakePolygon([]float32{
					f.Lines[sp1].P1AngleX, f.Lines[sp1].P1AngleY, 0,
					f.Lines[sp1].P2AngleX, f.Lines[sp1].P2AngleY, 0,
					f.Lines[sp2].P2AngleX, f.Lines[sp2].P2AngleY, 0,
				})
			} else if f.Lines[sp1].Line.P1 == f.Lines[sp2].Line.P2 || f.Lines[sp1].Line.P2 == f.Lines[sp2].Line.P2 {
				poli1 = MakePolygon([]float32{
					f.Lines[sp1].P1AngleX, f.Lines[sp1].P1AngleY, 0,
					f.Lines[sp1].P2AngleX, f.Lines[sp1].P2AngleY, 0,
					f.Lines[sp2].P1AngleX, f.Lines[sp2].P1AngleY, 0,
				})
			}
		}
	}

	return poli1, poli2, nil
}

func (f *BSPFrame) FindHangingPoints() (int, int, float32, float32, float32, float32, float32, float32) {
	res1 := -2
	res2 := -2
	var dist1 float32 = 0.0
	var dist2 float32 = 0.0
	var p1x float32 = 0.0
	var p1y float32 = 0.0
	var p2x float32 = 0.0
	var p2y float32 = 0.0
	for k, v := range f.Lines {
		hanger1 := v.Line.P1
		hanger2 := v.Line.P2
		for k2, v2 := range f.Lines {
			if k != k2 {
				if (hanger1 == v2.Line.P1 || hanger1 == v2.Line.P2) && (hanger1 != res1 || hanger1 != res2) {
					hanger1 = -1
				}

				if (hanger2 == v2.Line.P1 || hanger2 == v2.Line.P2) && (hanger2 != res1 || hanger2 != res2) {
					hanger2 = -1
				}
			}
		}
		if hanger1 != -1 {
			res1 = hanger1
			dist1 = v.P1Dist
			p1x = v.P1AngleX
			p1y = v.P1AngleY
		}
		if hanger2 != -1 {
			res2 = hanger2
			dist2 = v.P2Dist
			p2x = v.P2AngleX
			p2y = v.P2AngleY
		}
		if res1 != res2 {
			break
		}
	}

	return res1, res2, dist1, dist2, p1x, p1y, p2x, p2y
}

func (f *BSPFrame) EncloseFrame() {
	if len(f.Lines) > 1 {
		h1, h2, d1, d2, p1x, p1y, p2x, p2y := f.FindHangingPoints()
		f.Lines = append(f.Lines, BSPLine{Line: World.Line{P1: h1, P2: h2}, P1Dist: d1, P2Dist: d2, P1AngleX: p1x, P1AngleY: p1y, P2AngleX: p2x, P2AngleY: p2y})
		f.Enclosed = true
	} else {
		f.Lines = []BSPLine{}
		f.Enclosed = true
		f.Visible = false
	}
}

func (l *BSPLine) FindClosePoint() int {
	if l.P1Dist < l.P2Dist {
		return 1
	}

	return 2
}

func (l *BSPLine) FindDistantPoint() int {
	if l.P1Dist > l.P2Dist {
		return 1
	}

	return 2
}
