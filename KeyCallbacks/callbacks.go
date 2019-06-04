package KeyCallbacks

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/kanister10l/GoCamera/Camera"
	"github.com/kanister10l/GoCamera/World"
)

func SetCallbacks(window *glfw.Window, camera *Camera.Camera, world *World.World, sp *Camera.SphereWorld) {
	window.SetKeyCallback(func(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
		if action == 1 || action == 2 {
			if key == glfw.KeyD {
				transValue := Camera.RotateVector3D([]float32{0.03, 0.0, 0.0}, camera.Rotation)
				camera.Translate(transValue[0], transValue[1], transValue[2])
				camera.UpdateCamera()
			} else if key == glfw.KeyA {
				transValue := Camera.RotateVector3D([]float32{-0.03, 0.0, 0.0}, camera.Rotation)
				camera.Translate(transValue[0], transValue[1], transValue[2])
				camera.UpdateCamera()
			} else if key == glfw.KeyW {
				transValue := Camera.RotateVector3D([]float32{0.0, 0.0, 0.03}, camera.Rotation)
				camera.Translate(transValue[0], transValue[1], transValue[2])
				camera.UpdateCamera()
			} else if key == glfw.KeyS {
				transValue := Camera.RotateVector3D([]float32{0.0, 0.0, -0.03}, camera.Rotation)
				camera.Translate(transValue[0], transValue[1], transValue[2])
				camera.UpdateCamera()
			} else if key == glfw.KeyU {
				transValue := Camera.RotateVector3D([]float32{0.0, -0.03, 0.0}, camera.Rotation)
				camera.Translate(transValue[0], transValue[1], transValue[2])
				camera.UpdateCamera()
			} else if key == glfw.KeyJ {
				transValue := Camera.RotateVector3D([]float32{0.0, 0.03, 0.0}, camera.Rotation)
				camera.Translate(transValue[0], transValue[1], transValue[2])
				camera.UpdateCamera()
			} else if key == glfw.KeyRight {
				camera.Rotate(0, 1, 0)
				camera.UpdateCamera()
			} else if key == glfw.KeyLeft {
				camera.Rotate(0, -1, 0)
				camera.UpdateCamera()
			} else if key == glfw.KeyUp {
				camera.Rotate(1, 0, 0)
				camera.UpdateCamera()
			} else if key == glfw.KeyDown {
				camera.Rotate(-1, 0, 0)
				camera.UpdateCamera()
			} else if key == glfw.KeyH {
				camera.AdjustFov(1)
				camera.UpdateCamera()
			} else if key == glfw.KeyY {
				camera.AdjustFov(-1)
				camera.UpdateCamera()
			} else if key == glfw.KeyR {
				camera.Reset()
				camera.UpdateCamera()
			} else if key == glfw.KeyEscape {
				window.SetShouldClose(true)
			} else if key == glfw.KeyPageDown {
				camera.ChangeDrawType()
			} else if key == glfw.KeyPageUp {
				camera.SphereDrawType()
			} else if key == glfw.KeyKP8 {
				sp.Rotate(0, 0.0175)
			} else if key == glfw.KeyKP2 {
				sp.Rotate(0, -0.0175)
			} else if key == glfw.KeyKP6 {
				sp.Rotate(0.0175, 0)
			} else if key == glfw.KeyKP4 {
				sp.Rotate(-0.0175, 0)
			} else if key == glfw.Key1 {
				sp.ModifyConstant(0, 0, 0, -0.02, 0)
			} else if key == glfw.Key2 {
				sp.ModifyConstant(0, 0, 0, 0.02, 0)
			} else if key == glfw.Key3 {
				sp.ModifyConstant(-0.02, 0, 0, 0, 0)
			} else if key == glfw.Key4 {
				sp.ModifyConstant(0.02, 0, 0, 0, 0)
			} else if key == glfw.Key5 {
				sp.ModifyConstant(0, -0.02, 0, 0, 0)
			} else if key == glfw.Key6 {
				sp.ModifyConstant(0, 0.02, 0, 0, 0)
			} else if key == glfw.Key7 {
				sp.ModifyConstant(0, 0, -0.02, 0, 0)
			} else if key == glfw.Key8 {
				sp.ModifyConstant(0, 0, 0.02, 0, 0)
			} else if key == glfw.Key9 {
				sp.ModifyConstant(0, 0, 0, 0, -1)
			} else if key == glfw.Key0 {
				sp.ModifyConstant(0, 0, 0, 0, 1)
			} else if key == glfw.KeyM {
				sp.SelectNextMaterial()
			}
		}
	})
}
