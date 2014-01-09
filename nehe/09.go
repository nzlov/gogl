package main

import (
	"errors"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"image"
	"image/png"
	"io"
	"math/rand"
	"os"
	"time"
)

const (
	Title  = "Nehe 09"
	Width  = 640
	Height = 480
	num    = 50
)

type stars struct {
	r, g, b uint8   // 星星的颜色
	dist    float32 // 星星距离中心的距离
	angle   float32 // 当前星星所处的角度
}

var (
	texture      gl.Texture
	texturefiles [1]string
	twinkle      bool = false
	star         [num]stars
	spin         float32
	loop         int
	zoom         float32 = -15.0
	tilt         float32 = 90.0
	ztilt        float32 = 90.0
)

func init() {
	texturefiles[0] = "star.png"
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

	goph, err := os.Open(texturefiles[0])
	if err != nil {
		panic(err)
	}
	defer goph.Close()

	texture, err = createTexture(goph)

	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(0, 0, 0, 0.5)
	gl.ClearDepth(1)
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE) // 设置混色函数取得半透明效果
	gl.Enable(gl.BLEND)
	gl.Enable(gl.TEXTURE_2D)

	rand.Seed(time.Now().Unix())
	for loop = 0; loop < num; loop++ {
		star[loop].angle = 0.0
		star[loop].dist = float32(loop) / num * 5
		//fmt.Println(loop, star[loop].dist)
		star[loop].r = uint8(rand.Int() % 256)
		star[loop].g = uint8(rand.Int() % 256)
		star[loop].b = uint8(rand.Int() % 256)
	}

	return
}

// change view angle, exit upon ESC
func key(window *glfw.Window, k glfw.Key, s int, action glfw.Action, mods glfw.ModifierKey) {

	switch glfw.Key(k) {
	case glfw.KeyEscape:
		window.SetShouldClose(true)
	case glfw.KeyT:
		twinkle = !twinkle
	case glfw.KeyUp:
		tilt -= 0.5 // 屏幕向上倾斜
	case glfw.KeyDown:
		tilt += 0.5 // 屏幕向下倾斜
	case glfw.KeyLeft:
		ztilt -= 0.5 // 屏幕向上倾斜
	case glfw.KeyRight:
		ztilt += 0.5 // 屏幕向下倾斜
	case glfw.KeyPageUp:
		zoom -= 0.2 // 缩小
	case glfw.KeyPageDown:
		zoom += 0.2 // 放大
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
}

func drawScene() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	texture.Bind(gl.TEXTURE_2D)

	for loop = 0; loop < num; loop++ {
		gl.LoadIdentity()                            // 绘制每颗星星之前，重置模型观察矩阵
		gl.Translatef(0.0, 0.0, zoom)                // 深入屏幕里面
		gl.Rotatef(tilt, 1.0, 0.0, 0.0)              // 倾斜视角
		gl.Rotatef(ztilt, 0.0, 0.0, 1.0)             // 倾斜视角
		gl.Rotatef(star[loop].angle, 0.0, 1.0, 0.0)  // 旋转至当前所画星星的角度
		gl.Translatef(star[loop].dist, 0.0, 0.0)     // 沿X轴正向移动
		gl.Rotatef(-star[loop].angle, 0.0, 1.0, 0.0) // 取消当前星星的角度
		gl.Rotatef(-ztilt, 0.0, 0.0, 1.0)            // 取消屏幕倾斜
		gl.Rotatef(-tilt, 1.0, 0.0, 0.0)             // 取消屏幕倾斜
		if twinkle {                                 // 启用闪烁效果
			// 使用byte型数值指定一个颜色
			gl.Color4ub(star[(num-loop)-1].r, star[(num-loop)-1].g, star[(num-loop)-1].b, 255)
			gl.Begin(gl.QUADS) // 开始绘制纹理映射过的四边形
			gl.TexCoord2f(0.0, 0.0)
			gl.Vertex3f(-1.0, -1.0, 0.0)
			gl.TexCoord2f(1.0, 0.0)
			gl.Vertex3f(1.0, -1.0, 0.0)
			gl.TexCoord2f(1.0, 1.0)
			gl.Vertex3f(1.0, 1.0, 0.0)
			gl.TexCoord2f(0.0, 1.0)
			gl.Vertex3f(-1.0, 1.0, 0.0)
			gl.End() // 四边形绘制结束

		}
		gl.Rotatef(spin, 0.0, 0.0, 1.0) // 绕z轴旋转星星
		// 使用byte型数值指定一个颜色
		gl.Color4ub(star[loop].r, star[loop].g, star[loop].b, 255)
		gl.Begin(gl.QUADS) // 开始绘制纹理映射过的四边形
		gl.TexCoord2f(0.0, 0.0)
		gl.Vertex3f(-1.0, -1.0, 0.0)
		gl.TexCoord2f(1.0, 0.0)
		gl.Vertex3f(1.0, -1.0, 0.0)
		gl.TexCoord2f(1.0, 1.0)
		gl.Vertex3f(1.0, 1.0, 0.0)
		gl.TexCoord2f(0.0, 1.0)
		gl.Vertex3f(-1.0, 1.0, 0.0)
		gl.End() // 四边形绘制结束

		spin += 0.01                            // 星星的公转
		star[loop].angle += float32(loop) / num // 改变星星的自转角度
		star[loop].dist -= 0.01                 // 改变星星离中心的距离
		if star[loop].dist < 0.0 {              // 星星到达中心了么
			star[loop].dist += 5 // 往外移5个单位
			//fmt.Println(loop, star[loop].dist)
			star[loop].r = uint8(rand.Int() % 256) // 赋一个新红色分量
			star[loop].g = uint8(rand.Int() % 256) // 赋一个新绿色分量
			star[loop].b = uint8(rand.Int() % 256) // 赋一个新蓝色分量
		}
	}
}
