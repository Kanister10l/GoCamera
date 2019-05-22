package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/kanister10l/GoCamera/Camera"
	"github.com/kanister10l/GoCamera/KeyCallbacks"
	"github.com/kanister10l/GoCamera/World"

	"github.com/go-gl/gl/v4.1-compatibility/gl" // OR: github.com/go-gl/gl/v2.1/gl
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	/*vertexShaderSource = `
		#version 410
		in vec3 vp;
		void main() {
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"*/
	vertexShaderSource = `
		#version 410
		layout(location = 0) in vec3 vp;
		layout(location = 1) in vec3 vertex_colour;

		out vec3 colour;

		void main() {
			colour = vertex_colour;
			gl_Position = vec4(vp, 1.0);
		}
	` + "\x00"

	fragmentShaderSource = `
		#version 410
		in vec3 colour;
		out vec4 frag_colour;
		void main() {
			frag_colour = vec4(colour, 1.0);
		}
	` + "\x00"
)

func main() {
	rand.Seed(time.Now().Unix())
	runtime.LockOSThread()

	widthPtr := flag.Int("width", 1280, "Width of the window in pixels")
	heightPtr := flag.Int("height", 720, "Height of the window in pixels")

	flag.Parse()

	width := *widthPtr
	height := *heightPtr

	camera := Camera.NewCameraAt(0.0, 0.0, 0.0, 75, float32(width)/float32(height))

	world := World.NewWorld()
	err := world.Build("worldDescriptor.json")
	if err != nil {
		os.Exit(127)
	}

	window := initGlfw(width, height)
	defer glfw.Terminate()
	program := initOpenGL()
	KeyCallbacks.SetCallbacks(window, camera, world)

	log.Println(`
	KeyBindings:
	A ---> Move Left
	D ---> Move Right
	W ---> Move Forward
	S ---> Move Backward
	U ---> Move Up
	J ---> Move Down
	Left Arrow ---> Look Left
	Right Arrow ---> Look Right
	Up Arrow ---> Look Up
	Down Arrow ---> Look Down
	Y ---> Increase Field of View (ZOOM)
	H ---> Decrease Field of View (ZOOM)
	R ---> Reset Camera to Original Position
	ESC ---> Quit`)

	for !window.ShouldClose() {
		draw(window, program, camera, world)
	}
}

func draw(window *glfw.Window, program uint32, camera *Camera.Camera, world *World.World) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	if camera.DrawType == 0 {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)
		camera.DrawWorld(world)
	} else if camera.DrawType == 1 {
		gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)
		camera.DrawFullWorld(world)
	}

	glfw.PollEvents()
	window.SwapBuffers()
}

func initGlfw(width, height int) *glfw.Window {
	if err := glfw.Init(); err != nil {
		panic(err)
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)

	window, err := glfw.CreateWindow(width, height, "GoCamera", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	return window
}

func initOpenGL() uint32 {
	if err := gl.Init(); err != nil {
		panic(err)
	}
	version := gl.GoStr(gl.GetString(gl.VERSION))
	log.Println("OpenGL version", version)

	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	prog := gl.CreateProgram()
	gl.AttachShader(prog, vertexShader)
	gl.AttachShader(prog, fragmentShader)
	gl.LinkProgram(prog)
	return prog
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logShader := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logShader))

		return 0, fmt.Errorf("failed to compile %v: %v", source, logShader)
	}

	return shader, nil
}
