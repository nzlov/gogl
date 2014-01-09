// main.go
package main

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

func draw() { // 从这里开始进行所有的绘制

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // 清除屏幕和深度缓存
	gl.LoadIdentity()                                   // 重置当前的模型观察矩阵
	gl.Translatef(-1.5, 0.0, -6.0)                      // 左移 1.5 单位，并移入屏幕 6.0

	gl.Begin(gl.TRIANGLES)       // 绘制三角形
	gl.Vertex3f(0.0, 1.0, 0.0)   // 上顶点
	gl.Vertex3f(-1.0, -1.0, 0.0) // 左下
	gl.Vertex3f(1.0, -1.0, 0.0)  // 右下
	gl.End()                     // 三角形绘制结束

	gl.Translatef(3.0, 0.0, 0.0) // 右移3单位

	gl.Begin(gl.QUADS)           //  绘制正方形
	gl.Vertex3f(-1.0, 1.0, 0.0)  // 左上
	gl.Vertex3f(1.0, 1.0, 0.0)   // 右上
	gl.Vertex3f(1.0, -1.0, 0.0)  // 左下
	gl.Vertex3f(-1.0, -1.0, 0.0) // 右下
	gl.End()
}

// new window size
func reshape(window *glfw.Window, width, height int) {
	h := float64(height) / float64(width)

	znear := 5.0
	zfar := 30.0
	xmax := znear * 0.5

	gl.Viewport(0, 0, width, height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Frustum(-xmax, xmax, -xmax*h, xmax*h, znear, zfar)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Translated(0.0, 0.0, -20.0)
}

// program & OpenGL initialization
func Init() {
	gl.ShadeModel(gl.SMOOTH)                           // 启用阴影平滑
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)                  // 黑色背景
	gl.ClearDepth(1.0)                                 // 设置深度缓存
	gl.Enable(gl.DEPTH_TEST)                           // 启用深度测试
	gl.DepthFunc(gl.LEQUAL)                            // 所作深度测试的类型
	gl.Hint(gl.PERSPECTIVE_CORRECTION_HINT, gl.NICEST) // 告诉系统对透视进行修正
}

func main() {
	if !glfw.Init() {
		panic("Failed to initialize GLFW")
	}

	defer glfw.Terminate()

	glfw.WindowHint(glfw.DepthBits, 16)

	window, err := glfw.CreateWindow(300, 300, "nehe 02", nil, nil)
	if err != nil {
		panic(err)
	}

	// Set callback functions
	window.SetFramebufferSizeCallback(reshape)

	window.MakeContextCurrent()
	glfw.SwapInterval(1)

	width, height := window.GetFramebufferSize()
	reshape(window, width, height)

	// Parse command-line options
	Init()

	// Main loop
	for !window.ShouldClose() {
		// Draw gears
		draw()

		// Swap buffers
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
