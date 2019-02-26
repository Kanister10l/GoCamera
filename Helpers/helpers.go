package Helpers

import (
	"github.com/go-gl/gl/v4.1-compatibility/gl"
	"math"
)

func DegToRad(deg float32) float32 {
	return deg * math.Pi / 180.0
}

func MakeVao(points []float32) uint32 {
	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return vao
}

func NormalizePosition(x, y, maxX, maxY float32) (float32, float32) {
	return x/maxX, y/maxY
}
