package main

import "fmt"
import "time"
import "os"
import "github.com/veandco/go-sdl2/sdl"


func main(){
	emulatorMain()
}

func emulatorMain(){
	go startInputThread()
	go startUIThread()
	var chip chip8
	chip.initialize()
	chip.loadGame("PONG")
	for true{
		chip.cycle()
		if chip.drawFlag(){
			drawGraphics()
		}
		chip.setKeys()
		fmt.Println("cycle")
		//Sleep 1 second
		time.Sleep(1000000000)
		//return
	}
}

func startUIThread(){
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%T\n", err)
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	window.UpdateSurface()

	defer os.Exit(0)
	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
	}
}

func startInputThread(){

}

func drawGraphics(){

}