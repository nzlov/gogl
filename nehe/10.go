package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"image"
	"image/png"
	"io"
	"math"
	"os"
)

const (
	Title        = "Nehe 10"
	Width        = 640
	Height       = 480
	piover180    = 0.0174532925
	numtriangles = 36
)

var (
	texture      gl.Texture
	texturefiles string
	worldfiles   string
	heading      float32
	xpos         float32
	zpos         float32

	yrot          float32
	walkbias      float32 = 0.0
	walkbiasangle float32 = 0.0
	lookupdown    float32 = 0.0
	z             float32 = 0.0

	sector1 sector
)

type vertex struct {
	x, y, z float32
	u, v    float32
}

type triangle struct {
	vertex [3]vertex
}

type sector struct {
	numtriangles int
	triangles    *[numtriangles]triangle
}

func init() {
	texturefiles = "Mud.png"
	worldfiles = "World.txt"
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

	glfw.WindowHint(glfw.DepthBits, 16)

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.SetFramebufferSizeCallback(reshape)
	window.SetKeyCallback(key)

	window.MakeContextCurrent()

	glfw.SwapInterval(1000 / 60)

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

	goph, err := os.Open(texturefiles)
	if err != nil {
		panic(err)
	}
	defer goph.Close()

	texture, err = createTexture(goph)

	gl.Enable(gl.TEXTURE_2D)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE) // 设置混色函数取得半透明效果
	gl.ClearColor(0, 0, 0, 0.0)
	gl.ClearDepth(1)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.DEPTH_TEST)
	gl.ShadeModel(gl.SMOOTH)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)

	SetupWorld()

	return
}

func SetupWorld() {
	var x, y, z, u, v float32
	worldf, err := os.Open(worldfiles)
	if err != nil {
		panic(err)
	}
	defer worldf.Close()
	sector1.triangles = new([numtriangles]triangle)
	sector1.numtriangles = numtriangles
	for loop := 1; loop < numtriangles; loop++ {
		for vert := 0; vert < 3; vert++ {
			fmt.Fscanf(worldf, "%f %f %f %f %f\n", &x, &y, &z, &u, &v)
			sector1.triangles[loop].vertex[vert].x = x
			sector1.triangles[loop].vertex[vert].y = y
			sector1.triangles[loop].vertex[vert].z = z
			sector1.triangles[loop].vertex[vert].u = u
			sector1.triangles[loop].vertex[vert].v = v
		}
	}
}

// change view angle, exit upon ESC
func key(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {

	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyUp:
		xpos -= float32(math.Sin(float64(heading*piover180)) * 0.05)
		zpos -= float32(math.Cos(float64(heading*piover180)) * 0.05)
		if walkbiasangle >= 359.0 {
			walkbiasangle = 0.0
		} else {
			walkbiasangle += 10
		}
		walkbias = float32(math.Sin(float64(walkbiasangle*piover180)) / 20.0)
	case glfw.KeyDown:
		xpos += float32(math.Sin(float64(heading*piover180)) * 0.05)
		zpos += float32(math.Cos(float64(heading*piover180)) * 0.05)
		if walkbiasangle <= 1.0 {
			walkbiasangle = 359.0
		} else {
			walkbiasangle -= 10
		}
		walkbias = float32(math.Sin(float64(walkbiasangle*piover180)) / 20.0)
	case glfw.KeyLeft:
		heading += 1.0
		yrot = heading
	case glfw.KeyRight:
		heading -= 1.0
		yrot = heading
	case glfw.KeyPageUp:
		z -= 0.02
		lookupdown -= 1.0
	case glfw.KeyPageDown:
		z += 0.02
		lookupdown += 1.0
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
	glu.Perspective(45.0, float32(width)/float32(height), 0.1, 100.0)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.LoadIdentity() // 重置当前矩阵

	texture.Bind(gl.TEXTURE_2D)

	var x_m, y_m, z_m, u_m, v_m float32
	xtrans := -xpos
	ztrans := -zpos
	ytrans := -walkbias - 0.25
	sceneroty := 360.0 - yrot

	var numtriangles int

	gl.Rotatef(lookupdown, 1.0, 0, 0)
	gl.Rotatef(sceneroty, 0, 1.0, 0)

	gl.Translatef(xtrans, ytrans, ztrans)

	numtriangles = sector1.numtriangles

	// Process Each Triangle
	for loop_m := 0; loop_m < numtriangles; loop_m++ {
		gl.Begin(gl.TRIANGLES)
		gl.Normal3f(0.0, 0.0, 1.0)
		x_m = sector1.triangles[loop_m].vertex[0].x
		y_m = sector1.triangles[loop_m].vertex[0].y
		z_m = sector1.triangles[loop_m].vertex[0].z
		u_m = sector1.triangles[loop_m].vertex[0].u
		v_m = sector1.triangles[loop_m].vertex[0].v
		gl.TexCoord2f(u_m, v_m)
		gl.Vertex3f(x_m, y_m, z_m)

		x_m = sector1.triangles[loop_m].vertex[1].x
		y_m = sector1.triangles[loop_m].vertex[1].y
		z_m = sector1.triangles[loop_m].vertex[1].z
		u_m = sector1.triangles[loop_m].vertex[1].u
		v_m = sector1.triangles[loop_m].vertex[1].v
		gl.TexCoord2f(u_m, v_m)
		gl.Vertex3f(x_m, y_m, z_m)

		x_m = sector1.triangles[loop_m].vertex[2].x
		y_m = sector1.triangles[loop_m].vertex[2].y
		z_m = sector1.triangles[loop_m].vertex[2].z
		u_m = sector1.triangles[loop_m].vertex[2].u
		v_m = sector1.triangles[loop_m].vertex[2].v
		gl.TexCoord2f(u_m, v_m)
		gl.Vertex3f(x_m, y_m, z_m)
		gl.End()
	}
}
