package main
import "fmt"


func main(){
	emulatorMain()
}

func emulatorMain(){
	setupGraphics()
	setupInput()
	var chip chip8
	chip.initialize()
	chip.loadGame("PONG")
	for true{
		chip.cycle()
		if chip.drawFlag(){
			drawGraphics()
		}
		chip.setKeys()
		//temporary return until we make a way to exit
		fmt.Println("cycle")
		return
	}
}

func setupGraphics(){

}

func setupInput(){

}

func drawGraphics(){

}