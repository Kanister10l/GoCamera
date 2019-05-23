package Camera

import (
	"math"

	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/kanister10l/GoCamera/Helpers"
	"github.com/kanister10l/GoCamera/World"
)

func (camera *Camera) DrawWorld(world *World.World) {
	for _, entity := range world.Entities {
		for _, line := range entity.Lines {
			p1Visible, p1AngleX, p1AngleY := camera.CheckVisibility(entity.Points[line.P1])
			p2Visible, p2AngleX, p2AngleY := camera.CheckVisibility(entity.Points[line.P2])

			if p1Visible || p2Visible {
				x1, y1 := Helpers.NormalizePosition(p1AngleX, p1AngleY, camera.HorizontalFov/2, camera.VerticalFov/2)
				x2, y2 := Helpers.NormalizePosition(p2AngleX, p2AngleY, camera.HorizontalFov/2, camera.VerticalFov/2)
				drawLine := []float32{
					x1, y1, 0,
					x2, y2, 0,
					x2, y2, 0,
				}

				gl.BindVertexArray(Helpers.MakeVao(drawLine))
				gl.DrawArrays(gl.TRIANGLES, 0, int32(len(drawLine)/3))
			}
		}
	}
}

func (camera *Camera) DrawFullWorld(world *World.World) {
	figures := []BSPFigure{}

	polygons := []Polygon{}
	for k1, entity := range world.Entities {
		newFigure := BSPFigure{}
		newFigure.Init()
		figures = append(figures, newFigure)
		for _, line := range entity.Lines {
			visible := false
			p1Visible, p1AngleX, p1AngleY := camera.CheckVisibility(entity.Points[line.P1])
			p2Visible, p2AngleX, p2AngleY := camera.CheckVisibility(entity.Points[line.P2])

			if p1Visible || p2Visible {
				visible = true
			}

			x1, y1 := Helpers.NormalizePosition(p1AngleX, p1AngleY, camera.HorizontalFov/2, camera.VerticalFov/2)
			x2, y2 := Helpers.NormalizePosition(p2AngleX, p2AngleY, camera.HorizontalFov/2, camera.VerticalFov/2)

			d1 := mgl32.NewVecNFromData([]float32{entity.Points[line.P1].X - camera.X, entity.Points[line.P1].Y - camera.Y, entity.Points[line.P1].Z - camera.Z}).Vec3().Len()
			d2 := mgl32.NewVecNFromData([]float32{entity.Points[line.P2].X - camera.X, entity.Points[line.P2].Y - camera.Y, entity.Points[line.P2].Z - camera.Z}).Vec3().Len()
			figures[k1].AddLine(line, d1, d2, x1, y1, x2, y2, visible)
		}

		for k2 := range figures[k1].Frames {
			if figures[k1].Frames[k2].Visible {
				poly1, poly2, err := figures[k1].Frames[k2].ConvertToPolygons()
				if err == nil {
					polygons = append(polygons, poly1, poly2)
				}
			}
		}
	}

	polyLen := len(polygons)

	for i := 0; i < polyLen; i++ {
		min := polygons[0].Dist
		toRender := 0
		for k, v := range polygons {
			if v.Dist > min {
				min = v.Dist
				toRender = k
			}
		}

		gl.BindVertexArray(Helpers.MakeVao(polygons[toRender].Drawer))
		gl.DrawArrays(gl.TRIANGLES, 0, int32(len(polygons[toRender].Drawer)/3))
		polygons = append(polygons[:toRender], polygons[toRender+1:]...)
	}

	/*tree := BuildTree(polygons)

	Traverse(tree)*/
}

func (camera *Camera) CheckVisibility(point World.Point) (bool, float32, float32) {
	poi := mgl32.NewVecNFromData([]float32{point.X - camera.X, point.Y - camera.Y, point.Z - camera.Z}).Vec3()
	vNorm := mgl32.NewVecNFromData([]float32{camera.XVector[0], camera.XVector[1], camera.XVector[2]}).Vec3().Normalize()
	vDist := poi.Dot(vNorm)
	vNormExt := vNorm.Mul(vDist)
	vX := point.X - vNormExt.X()
	vY := point.Y - vNormExt.Y()
	vZ := point.Z - vNormExt.Z()

	hNorm := mgl32.NewVecNFromData([]float32{camera.YVector[0], camera.YVector[1], camera.YVector[2]}).Vec3().Normalize()
	hDist := poi.Dot(hNorm)
	hNormExt := hNorm.Mul(hDist)
	hX := point.X - hNormExt.X()
	hY := point.Y - hNormExt.Y()
	hZ := point.Z - hNormExt.Z()

	vCameraVector := mgl32.NewVecNFromData([]float32{camera.ZVector[0], camera.ZVector[1], camera.ZVector[2]}).Vec3()
	vPointVector := mgl32.NewVecNFromData([]float32{vX - camera.X, vY - camera.Y, vZ - camera.Z}).Vec3()
	vPlaneAngle := mgl32.RadToDeg(float32(math.Acos(float64(vCameraVector.Dot(vPointVector) / (vCameraVector.Len() * vPointVector.Len())))))

	hCameraVector := mgl32.NewVecNFromData([]float32{camera.ZVector[0], camera.ZVector[1], camera.ZVector[2]}).Vec3()
	hPointVector := mgl32.NewVecNFromData([]float32{hX - camera.X, hY - camera.Y, hZ - camera.Z}).Vec3()
	hPlaneAngle := mgl32.RadToDeg(float32(math.Acos(float64(hCameraVector.Dot(hPointVector) / (hCameraVector.Len() * hPointVector.Len())))))

	if !camera.IsRightSide(hX, hZ, camera.X, camera.Z, camera.LookAt.X, camera.LookAt.Z, camera.X+1, camera.Z) && (camera.Rotation.YDeg <= 90 || camera.Rotation.YDeg >= 270) {
		hPlaneAngle = -hPlaneAngle
	} else if camera.IsRightSide(hX, hZ, camera.X, camera.Z, camera.LookAt.X, camera.LookAt.Z, camera.X+1, camera.Z) && (camera.Rotation.YDeg > 90 && camera.Rotation.YDeg < 270) {
		hPlaneAngle = -hPlaneAngle
	}

	if camera.IsRightSide(vZ, vY, camera.Z, camera.Y, camera.LookAt.Z, camera.LookAt.Y, camera.Z+1, camera.Y) && camera.Rotation.XDeg <= 180 {
		vPlaneAngle = -vPlaneAngle
	} else if !camera.IsRightSide(vZ, vY, camera.Z, camera.Y, camera.LookAt.Z, camera.LookAt.Y, camera.Z+1, camera.Y) && camera.Rotation.XDeg > 180 {
		vPlaneAngle = -vPlaneAngle
	}

	if hPlaneAngle <= camera.HorizontalFov/2 && vPlaneAngle <= camera.VerticalFov/2 {
		return true, hPlaneAngle, vPlaneAngle
	}

	return false, hPlaneAngle, vPlaneAngle
}

func (camera *Camera) IsRightSide(x, y, x1, y1, x2, y2, baseX, baseY float32) bool {
	d := (x-x1)*(y2-y1) - (y-y1)*(x2-x1)
	baseD := (baseX-x1)*(y2-y1) - (baseY-y1)*(x2-x1)

	if baseD > 0 {
		if d > 0 {
			return true
		} else {
			return false
		}
	} else {
		if d < 0 {
			return true
		} else {
			return false
		}
	}
}
