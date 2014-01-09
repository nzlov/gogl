// main.go
package main

import (
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

var rtri, rquad float32 = 0, 0 // 用于三角形的角度
func draw() {                  // 从这里开始进行所有的绘制

	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT) // 清除屏幕和深度缓存
	gl.LoadIdentity()                                   // 重置当前的模型观察矩阵
	gl.Translatef(-1.5, 0.0, -6.0)                      // 左移 1.5 单位，并移入屏幕 6.0

	gl.Rotatef(rtri, 0.0, 1.0, 0.0) // 绕Y轴旋转三角形

	gl.Begin(gl.TRIANGLES)       // 绘制三角形
	gl.Color3f(1.0, 0.0, 0.0)    // 设置当前色为红色
	gl.Vertex3f(0.0, 1.0, 0.0)   // 三角形的上顶点 (前侧面)
	gl.Color3f(0.0, 1.0, 0.0)    // 设置当前色为绿色
	gl.Vertex3f(-1.0, -1.0, 1.0) // 三角形的左下顶点 (前侧面)
	gl.Color3f(0.0, 0.0, 1.0)    // 设置当前色为蓝色
	gl.Vertex3f(1.0, -1.0, 1.0)  // 三角形的右下顶点 (前侧面)

	gl.Color3f(1.0, 0.0, 0.0)    // 设置当前色为红色
	gl.Vertex3f(0.0, 1.0, 0.0)   // 三角形的上顶点 (右侧面)
	gl.Color3f(0.0, 0.0, 1.0)    // 设置当前色为蓝色
	gl.Vertex3f(1.0, -1.0, 1.0)  // 三角形的左下顶点 (右侧面)
	gl.Color3f(0.0, 1.0, 0.0)    // 设置当前色为绿色
	gl.Vertex3f(1.0, -1.0, -1.0) // 三角形的右下顶点 (右侧面)

	gl.Color3f(1.0, 0.0, 0.0)     // 设置当前色为红色
	gl.Vertex3f(0.0, 1.0, 0.0)    // 三角形的上顶点 (后侧面)
	gl.Color3f(0.0, 0.0, 1.0)     // 设置当前色为蓝色
	gl.Vertex3f(1.0, -1.0, -1.0)  // 三角形的左下顶点 (后侧面)
	gl.Color3f(0.0, 1.0, 0.0)     // 设置当前色为绿色
	gl.Vertex3f(-1.0, -1.0, -1.0) // 三角形的右下顶点 (后侧面)

	gl.Color3f(1.0, 0.0, 0.0)     // 设置当前色为红色
	gl.Vertex3f(0.0, 1.0, 0.0)    // 三角形的上顶点 (左侧面)
	gl.Color3f(0.0, 0.0, 1.0)     // 设置当前色为蓝色
	gl.Vertex3f(-1.0, -1.0, -1.0) // 三角形的左下顶点 (左侧面)
	gl.Color3f(0.0, 1.0, 0.0)     // 设置当前色为绿色
	gl.Vertex3f(-1.0, -1.0, 1.0)  // 三角形的右下顶点 (左侧面)

	gl.End()                      // 三角形侧面绘制结束
	gl.Begin(gl.QUADS)            // 三角形底面绘制形
	gl.Vertex3f(1.0, -1.0, -1.0)  // 四边形的右上顶点 (底面)
	gl.Vertex3f(-1.0, -1.0, -1.0) // 四边形的左上顶点 (底面)
	gl.Vertex3f(-1.0, -1.0, 1.0)  // 四边形的左下顶点 (底面)
	gl.Vertex3f(1.0, -1.0, 1.0)   // 四边形的右下顶点 (底面)
	gl.End()                      // 三角形底面绘制结束

	gl.LoadIdentity()                // 重置当前的模型观察矩阵
	gl.Translatef(1.5, 0.0, -7.0)    // 右移1.5单位,并移入屏幕 6.0
	gl.Rotatef(rquad, 1.0, 1.0, 1.0) //  绕X轴旋转四边形
	gl.Begin(gl.QUADS)               //  绘制正方形

	gl.Color3f(1.0, 0.0, 0.0)    // 一次性将当前色设置为红色
	gl.Vertex3f(1.0, 1.0, -1.0)  // 四边形的右上顶点 (顶面)
	gl.Vertex3f(-1.0, 1.0, -1.0) // 四边形的左上顶点 (顶面)
	gl.Vertex3f(-1.0, 1.0, 1.0)  // 四边形的左下顶点 (顶面)
	gl.Vertex3f(1.0, 1.0, 1.0)   // 四边形的右下顶点 (顶面)

	gl.Color3f(1.0, 0.5, 0.0)    // 一次性将当前色设置为橙色
	gl.Vertex3f(1.0, 1.0, 1.0)   // 四边形的右上顶点 (前面)
	gl.Vertex3f(-1.0, 1.0, 1.0)  // 四边形的左上顶点 (前面)
	gl.Vertex3f(-1.0, -1.0, 1.0) // 四边形的左下顶点 (前面)
	gl.Vertex3f(1.0, -1.0, 1.0)  // 四边形的右下顶点 (前面)

	gl.Color3f(0.0, 0.0, 1.0)    // 一次性将当前色设置为绿色
	gl.Vertex3f(1.0, 1.0, -1.0)  // 四边形的右上顶点 (右面)
	gl.Vertex3f(1.0, 1.0, 1.0)   // 四边形的左上顶点 (右面)
	gl.Vertex3f(1.0, -1.0, 1.0)  // 四边形的左下顶点 (右面)
	gl.Vertex3f(1.0, -1.0, -1.0) // 四边形的右下顶点 (右面)

	gl.Color3f(0.0, 1.0, 0.0)     // 一次性将当前色设置为蓝色
	gl.Vertex3f(-1.0, 1.0, -1.0)  // 四边形的右上顶点 (后面)
	gl.Vertex3f(1.0, 1.0, -1.0)   // 四边形的左上顶点 (后面)
	gl.Vertex3f(1.0, -1.0, -1.0)  // 四边形的左下顶点 (后面)
	gl.Vertex3f(-1.0, -1.0, -1.0) // 四边形的右下顶点 (后面)

	gl.Color3f(1.0, 1.0, 0.0)     // 一次性将当前色设置为黄色
	gl.Vertex3f(-1.0, 1.0, 1.0)   // 四边形的右上顶点 (左面)
	gl.Vertex3f(-1.0, 1.0, -1.0)  // 四边形的左上顶点 (左面)
	gl.Vertex3f(-1.0, -1.0, -1.0) // 四边形的左下顶点 (左面)
	gl.Vertex3f(-1.0, -1.0, 1.0)  // 四边形的右下顶点 (左面)

	gl.Color3f(1.0, 1.0, 1.0)     // 一次性将当前色设置为白色
	gl.Vertex3f(1.0, -1.0, -1.0)  // 四边形的右上顶点 (底面)
	gl.Vertex3f(-1.0, -1.0, -1.0) // 四边形的左上顶点 (底面)
	gl.Vertex3f(-1.0, -1.0, 1.0)  // 四边形的左下顶点 (底面)
	gl.Vertex3f(1.0, -1.0, 1.0)   // 四边形的右下顶点 (底面)
	gl.End()

	rtri += 0.2 // 增加三角形的旋转变量

	rquad -= 0.15 // 减少四边形的旋转变量
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

	window, err := glfw.CreateWindow(300, 300, "nehe 05", nil, nil)
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
