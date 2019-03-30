package main

type chip8 struct {
	opcode		uint8
	mem			[4096]uint8
	reg			[16]uint8
	ir			uint16
	pc			uint16
	gfx			[64][32]bool
	delayTimer 	uint8
	soundTimer 	uint8
	stack		[16]uint16
	sp			uint16
	key			[16]bool
}

func (chip chip8) initialize(){

}

func (chip chip8) loadGame(name string){

}

func (chip chip8) cycle(){

}

func (chip chip8) drawFlag() bool {
	return false
}

func (chip chip8) setKeys(){

}