package main

import (
	"fmt"

	"strings"

	"io/ioutil"

	"errors"

	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/2takecarl/myopengl-go-project/opengl%201/src/myproject2"
)

type ShaderID uint32
type ProgramID uint32
type VAOID uint32
type VBOID uint32

// helper functions
func GetVersion() string {
	return gl.GoStr(gl.GetString(gl.VERSION))
}

type programInfo struct {
	path     string
	fragPath string
	modified time.Time
}

var loadedShaders []programInfo

func CheckShadersForChanges() {
	// for _, shaderInfo := range loadedShaders {
	// 	file, err := os.Stat(oriInfo.path)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	modTime := file.ModTime()
	// 	if !modTime.Equal(shaderInfo.modified) {
	// 		fmt.Println("shader modified")
	// 	}
	// }
}

func LoadShader(path string, shaderType uint32) (ShaderID, error) {
	shaderFile, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	shaderFileStr := string(shaderFile)
	shaderId, err := CreateShader(shaderFileStr, shaderType)
	if err != nil {
		return 0, err
	}

	return shaderId, nil
}
func CreateShader(shaderSource string, shaderType uint32) (ShaderID, error) {
	shaderId := gl.CreateShader(shaderType)
	shaderSource = shaderSource + "\x00"
	csource, free := gl.Strs(shaderSource)
	gl.ShaderSource(shaderId, 1, csource, nil)
	free()
	gl.CompileShader(shaderId)
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		fmt.Println("Shader... [FAILED]: \n" + log)
		return 0, errors.New("Trying To Compile Shader... [FAILED]")
	}
	fmt.Println("Shader... [PASSED]")
	return ShaderID(shaderId), nil
}
func CreateProgram(vertPath string, fragPath string) (ProgramID, error) {
	vert, err := LoadShader(vertPath, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}
	frag, err := LoadShader(fragPath, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, uint32(vert))
	gl.AttachShader(shaderProgram, uint32(frag))
	gl.LinkProgram(shaderProgram)

	var success int32
	gl.GetProgramiv(shaderProgram, gl.LINK_STATUS, &success)
	if success == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(shaderProgram, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(shaderProgram, logLength, nil, gl.Str(log))
		return 0, errors.New(("Program Link... [FAILED]: \n" + log))
	}
	fmt.Println("Shader Program... [PASSED]")
	gl.DeleteShader(uint32(vert))
	gl.DeleteShader(uint32(frag))
	// file, err := os.Stat(path)
	// if err != nil {
	// 	panic(err)
	// }
	// modTime := file.ModTime()
	// loadedShaders = append(loadedShaders, shaderInfo{path, modTime})
	return ProgramID(shaderProgram), nil
}
func GenBindBuffer(target uint32) VBOID {
	var VBO uint32
	gl.GenBuffers(1, &VBO)
	gl.BindBuffer(target, VBO)
	return VBOID(VBO)
}
func GenBindVertexArray() VAOID {
	var VAO uint32
	gl.GenVertexArrays(1, &VAO)
	gl.BindVertexArray(VAO)
	return VAOID(VAO)
}
func BindVertexArray(vaoID VAOID) {
	gl.BindVertexArray(uint32(vaoID))
}
func BufferDataFloat(target uint32, data []float32, usage uint32) {
	gl.BufferData(target, len(data)*4, gl.Ptr(data), usage)
}
func UseProgram(programID ProgramID) {
	gl.UseProgram(uint32(programID))
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 3)

	window, err := sdl.CreateWindow("test", 200, 200, 800, 600, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	window.GLCreateContext()
	defer window.Destroy()

	gl.Init()

	fmt.Println("version: ", GetVersion())

	shaderProgram, err := CreateProgram("shaders/main.vert", "shaders/main.frag")
	if err != nil {
		panic(err)
	}

	vertices := []float32{
		//first triangle
		-0.5, -0.5, 0.0,
		0.5, -0.5, 0.0,
		0.0, 0.5, 0.0}

	GenBindBuffer(gl.ARRAY_BUFFER)
	VAO := GenBindVertexArray()

	BufferDataFloat(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.BindVertexArray(0)

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				fmt.Println("SDL: QUITEVENT")
				println("quit")

			}
		}

		gl.ClearColor(0.1, 0.1, 0.1, 0.1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		UseProgram(shaderProgram)
		BindVertexArray(VAO)
		gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
		window.GLSwap()
		CheckShadersForChanges()

	}
}
