// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// NEHE Tutorial 06: Texture Mapping.
// http://nehe.gamedev.net/data/lessons/lesson.asp?lesson=06
package main

import (
	"errors"
	"fmt"
	"github.com/andrebq/gas"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"image"
	"image/png"
	"io"
	"os"
)

const (
	Title  = "Nehe 06"
	Width  = 640
	Height = 480
)

var (
	running      bool
	rotation     [3]float32
	texture      gl.Texture
	texturefiles [1]string
)

func init() {
	texturefiles[0], _ = gas.Abs("github.com/go-gl/examples/data/gopher.png")
}

func createTexture(r io.Reader) (gl.Texture, error) {
	img, err := png.Decode(r)
	if err != nil {
		return gl.Texture(0), err
	}

	rgbaImg, ok := img.(*image.NRGBA)
	if !ok {
		return gl.Texture(0), errors.New("texture must be an NRGBA image")
	}

	textureId := gl.GenTexture()
	textureId.Bind(gl.TEXTURE_2D)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	// flip image: first pixel is lower left corner
	imgWidth, imgHeight := img.Bounds().Dx(), img.Bounds().Dy()
	data := make([]byte, imgWidth*imgHeight*4)
	lineLen := imgWidth * 4
	dest := len(data) - lineLen
	for src := 0; src < len(rgbaImg.Pix); src += rgbaImg.Stride {
		copy(data[dest:dest+lineLen], rgbaImg.Pix[src:src+rgbaImg.Stride])
		dest -= lineLen
	}
	gl.TexImage2D(gl.TEXTURE_2D, 0, 4, imgWidth, imgHeight, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)

	return textureId, nil
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetFramebufferSizeCallback(reshape)
	window.SetKeyCallback(key)

	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	width, height := window.GetFramebufferSize()
	reshape(window, width, height)

	if err := initGL(); err != nil {
		fmt.Fprintf(os.Stderr, "init: %s\n", err)
		return
	}
	defer destroyGL()

	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGL() (err error) {

	goph, err := os.Open(texturefiles[0])
	if err != nil {
		panic(err)
	}
	defer goph.Close()

	texture, err = createTexture(goph)

	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LEQUAL)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
	gl.Enable(gl.DEPTH_TEST)
	gl.Enable(gl.TEXTURE_2D)
	return
}

// change view angle, exit upon ESC
func key(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press {
		return
	}

	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	default:
		return
	}
}

func destroyGL() {
	texture.Delete()
}

// new window size
func reshape(window *glfw.Window, width, height int) {

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-1.0, 1.0, -1.0, 1.0, 1, 20)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translated(0.0, 0.0, -20.0)
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity()

	gl.Translatef(0, 0, -5)

	gl.Rotatef(rotation[0], 1, 0, 0)
	gl.Rotatef(rotation[1], 0, 1, 0)
	gl.Rotatef(rotation[2], 0, 0, 1)

	rotation[0] += 0.3
	rotation[1] += 0.2
	rotation[2] += 0.4

	texture.Bind(gl.TEXTURE_2D)

	gl.Begin(gl.QUADS)
	// Front Face
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, 1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, 1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, 1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, 1) // Top Left Of The Texture and Quad
	// Back Face
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, -1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, -1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, -1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, -1) // Bottom Left Of The Texture and Quad
	// Top Face
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, 1, 1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, 1, 1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1) // Top Right Of The Texture and Quad
	// Bottom Face
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, -1, -1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, -1, -1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1) // Bottom Right Of The Texture and Quad
	// Right face
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(1, -1, -1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(1, 1, -1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(1, 1, 1) // Top Left Of The Texture and Quad
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(1, -1, 1) // Bottom Left Of The Texture and Quad
	// Left Face
	gl.TexCoord2f(0, 0)
	gl.Vertex3f(-1, -1, -1) // Bottom Left Of The Texture and Quad
	gl.TexCoord2f(1, 0)
	gl.Vertex3f(-1, -1, 1) // Bottom Right Of The Texture and Quad
	gl.TexCoord2f(1, 1)
	gl.Vertex3f(-1, 1, 1) // Top Right Of The Texture and Quad
	gl.TexCoord2f(0, 1)
	gl.Vertex3f(-1, 1, -1) // Top Left Of The Texture and Quad
	gl.End()
}
