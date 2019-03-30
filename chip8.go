package main

import (
	"fmt"
	"io/ioutil"
)

type chip8 struct {
	opcode     uint8
	mem        [4096]uint8
	reg        [16]uint8
	ir         uint16
	pc         uint16
	gfx        [64][32]bool
	delayTimer uint8
	soundTimer uint8
	stack      [16]uint16
	sp         uint16
	key        [16]bool
}

func (chip chip8) initialize() {
	// Initialize registers and memory once
	chip.pc = 0x200 // Program counter starts at 0x200
	chip.opcode = 0 // Reset current opcode
	chip.ir = 0     // Reset index register
	chip.sp = 0     // Reset stack pointer

	// Clear display
	// Clear stack
	// Clear registers V0-VF
	// Clear memory

	// Load fontset
	chip8FontSet := [80]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}
	for i := 0; i < 80; i++ {
		chip.mem[i] = chip8FontSet[i]
	}

	// Reset timers

}

func (chip chip8) loadGame(name string) {
	buffer, err := ioutil.ReadFile("roms/" + name)
	if err == nil {
		fmt.Println(len(buffer))
		for i := 0; i < len(buffer); i++ {
			chip.mem[i+512] = buffer[i]
		}
	}
}

func (chip chip8) cycle() {

	// Fetch opcode
	op1 := uint16(chip.mem[chip.pc])
	op2 := uint16(chip.mem[chip.pc+1])
	opcode := op1<<8 | op2

	// Decode opcode
	switch opcode & 0xF000 {
	// Some opcodes //

	case 0xA000: // ANNN: Sets I to the address NNN
		// Execute opcode
		chip.ir = opcode & 0x0FFF
		chip.pc += 2
		break

	// More opcodes //

	default:
		fmt.Printf("Unknown opcode: 0x%X\n", opcode);
	}

	// Update timers
	if chip.delayTimer > 0 {
		chip.delayTimer--
	}
	if chip.soundTimer > 0 {
		if chip.soundTimer == 1 {
			fmt.Printf("BEEP!\n")
		}

		chip.soundTimer--;
	}
}

func (chip chip8) drawFlag() bool {
	return false
}

func (chip chip8) setKeys() {

}
