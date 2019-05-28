package Camera

import (
	"log"
	"math"

	"github.com/kanister10l/GoCamera/Helpers"
)

type Camera struct {
	Rotation
	Position
	Axis
	Fov
	LookAt   Position
	DrawType int
}

type Rotation struct {
	XDeg float32
	YDeg float32
	ZDeg float32
}

type Position struct {
	X float32
	Y float32
	Z float32
}

type Axis struct {
	XVector []float32
	YVector []float32
	ZVector []float32
}

type Fov struct {
	HorizontalFov float32
	VerticalFov   float32
	FovRatio      float32
}

func NewCameraAt(x, y, z, fov, screenRatio float32) *Camera {
	log.Println("Preparing new camera")
	camera := Camera{}

	camera.DrawType = 0

	camera.X = x
	camera.Y = y
	camera.Z = z

	camera.XDeg = 0.0
	camera.YDeg = 0.0
	camera.ZDeg = 0.0

	camera.FovRatio = screenRatio
	camera.HorizontalFov = fov
	camera.VerticalFov = fov / camera.FovRatio

	camera.UpdateCamera()

	log.Println("New camera data:", camera)

	return &camera
}

func (camera *Camera) UpdateCamera() {
	camera.NewBaseAxis()
	camera.RotateAxis(camera.Rotation)
	camera.LookAt.X = camera.X
	camera.LookAt.Y = camera.Y
	camera.LookAt.Z = camera.Z
	camera.LookAt.Translate(camera.ZVector[0], camera.ZVector[1], camera.ZVector[2])
}

func (camera *Camera) Reset() {
	camera.XDeg = 0.0
	camera.YDeg = 0.0
	camera.ZDeg = 0.0

	camera.X = 0
	camera.Y = 0
	camera.Z = 0
}

func (camera *Camera) ChangeDrawType() {
	if camera.DrawType == 0 {
		camera.DrawType = 1
	} else {
		camera.DrawType = 0
	}
}

func (camera *Camera) SphereDrawType() {
	camera.DrawType = 2
}

func (r *Rotation) Rotate(xPlane, yPlane, zPlane float32) {
	if r.XDeg+xPlane == 90 || r.XDeg+xPlane == 270 {
		xPlane -= 0.01
	}
	if r.YDeg+yPlane == 90 || r.YDeg+yPlane == 270 {
		yPlane -= 0.01
	}
	if r.ZDeg+zPlane == 90 || r.ZDeg+zPlane == 270 {
		zPlane -= 0.01
	}

	if r.XDeg+xPlane >= 360.0 {
		r.XDeg = xPlane - (360.0 - r.XDeg)
	} else if r.XDeg+xPlane < 0 {
		r.XDeg = 360.0 + (xPlane + r.XDeg)
	} else {
		r.XDeg = r.XDeg + xPlane
	}

	if r.YDeg+yPlane >= 360.0 {
		r.YDeg = yPlane - (360.0 - r.YDeg)
	} else if r.YDeg+yPlane < 0 {
		r.YDeg = 360.0 + (yPlane + r.YDeg)
	} else {
		r.YDeg = r.YDeg + yPlane
	}

	if r.ZDeg+zPlane >= 360.0 {
		r.ZDeg = zPlane - (360.0 - r.ZDeg)
	} else if r.ZDeg+zPlane < 0 {
		r.ZDeg = 360.0 + (zPlane + r.ZDeg)
	} else {
		r.ZDeg = r.ZDeg + zPlane
	}
}

func (p *Position) Translate(x, y, z float32) {
	p.X = p.X + x
	p.Y = p.Y + y
	p.Z = p.Z + z
}

func (f *Fov) AdjustFov(value float32) {
	if f.HorizontalFov+value > 180 || f.HorizontalFov+value < 0 {
		return
	}
	f.HorizontalFov += value
	f.VerticalFov += value / f.FovRatio
}

func (a *Axis) NewBaseAxis() {
	a.XVector = []float32{1.0, 0.0, 0.0}
	a.YVector = []float32{0.0, 1.0, 0.0}
	a.ZVector = []float32{0.0, 0.0, 1.0}
}

func (a *Axis) RotateAxis(rotation Rotation) {
	a.XVector = RotateVector3D(a.XVector, rotation)
	a.YVector = RotateVector3D(a.YVector, rotation)
	a.ZVector = RotateVector3D(a.ZVector, rotation)
}

func RotateVector3D(vector []float32, rotation Rotation) []float32 {
	tmpVector := []float32{vector[0], vector[1], vector[2]}

	//Rotate around x Axis
	xPrim := tmpVector[0]
	yPrim := tmpVector[1]*float32(math.Cos(float64(Helpers.DegToRad(rotation.XDeg)))) - tmpVector[2]*float32(math.Sin(float64(Helpers.DegToRad(rotation.XDeg))))
	zPrim := tmpVector[1]*float32(math.Sin(float64(Helpers.DegToRad(rotation.XDeg)))) + tmpVector[2]*float32(math.Cos(float64(Helpers.DegToRad(rotation.XDeg))))
	tmpVector[0] = xPrim
	tmpVector[1] = yPrim
	tmpVector[2] = zPrim

	//Rotate around y Axis
	xPrim = tmpVector[0]*float32(math.Cos(float64(Helpers.DegToRad(rotation.YDeg)))) + tmpVector[2]*float32(math.Sin(float64(Helpers.DegToRad(rotation.YDeg))))
	yPrim = tmpVector[1]
	zPrim = (-tmpVector[0])*float32(math.Sin(float64(Helpers.DegToRad(rotation.YDeg)))) + tmpVector[2]*float32(math.Cos(float64(Helpers.DegToRad(rotation.YDeg))))
	tmpVector[0] = xPrim
	tmpVector[1] = yPrim
	tmpVector[2] = zPrim

	//Rotate around z Axis
	xPrim = tmpVector[0]*float32(math.Cos(float64(Helpers.DegToRad(rotation.ZDeg)))) - tmpVector[1]*float32(math.Sin(float64(Helpers.DegToRad(rotation.ZDeg))))
	yPrim = tmpVector[0]*float32(math.Sin(float64(Helpers.DegToRad(rotation.ZDeg)))) + tmpVector[1]*float32(math.Cos(float64(Helpers.DegToRad(rotation.ZDeg))))
	zPrim = tmpVector[2]
	tmpVector[0] = xPrim
	tmpVector[1] = yPrim
	tmpVector[2] = zPrim

	return tmpVector
}
