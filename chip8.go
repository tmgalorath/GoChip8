package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
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
	rand.Seed(time.Now().UTC().UnixNano())
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
	chip.pc += 2

	// Decode opcode

	chip.decode(opcode)

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

func (chip chip8) decode(opcode uint16) {
	fmt.Printf("Exec: 0x%x\n", opcode)
	switch opcode & 0xF000 {
	case 0x0000:
		if opcode == 0x00E0 {
			//TODO disp_clear
		} else if opcode == 0x00EE {
			//TODO ret
		} else {
			//TODO call rca 1802
		}
		break
	case 0x1000:
		chip.pc = opcode & 0x0FFF
		break
	case 0x2000:
		chip.stack[chip.sp] = chip.pc
		chip.sp++
		chip.pc = opcode & 0x0FFF
		break
	case 0x3000:
		//will casting truncate left half.
		if chip.reg[(opcode&0x0F00)>>8] == uint8(opcode&0x00FF) {
			chip.pc += 2
		}
		break
	case 0x4000:
		if chip.reg[(opcode&0x0F00)>>8] != uint8(opcode&0x00FF) {
			chip.pc += 2
		}
		break
	case 0x5000:
		//contradicts itself
		if chip.reg[(opcode&0x0F00)>>8] == chip.reg[opcode&0x00FF] {
			chip.pc += 2
		}
		break
	case 0x6000:
		chip.reg[(opcode&0x0F00)>>8] = uint8(opcode & 0x00FF)
		break
	case 0x7000:
		chip.reg[(opcode&0x0F00)>>8] += uint8(opcode & 0x00FF)
		break
	case 0x8000:
		x := &chip.reg[(opcode&0x0F00)>>8]
		y := &chip.reg[(opcode&0x00F0)>>4]
		switch (opcode & 0x000F) {
		case 0x0000:
			*x = *y
			break
		case 0x0001:
			*x |= *y
			break
		case 0x0002:
			*x &= *y
			break
		case 0x0003:
			*x ^= *y
			break
		case 0x0004:
			*x += *y
			if *y > *x {
				chip.reg[0xF] = 1
			} else{
				chip.reg[0xF] = 0
			}
			break
		case 0x0005:
			*x -= *y
			if *y > *x {
				chip.reg[0xF] = 0
			} else{
				chip.reg[0xF] = 1
			}
			break
		case 0x0006:
			chip.reg[0xF] = *x & 0x1
			*x >>= 1
			break
		case 0x0007:
			*x = *y - *x
			if *x > *y {
				chip.reg[0xF] = 0
			} else{
				chip.reg[0xF] = 1
			}
			break
		case 0x000E:
			chip.reg[0xF] = *x & 0x80
			*x <<= 1
			break
		default:
			chip.decodeFail(opcode)
		}
		break
	case 0xA000:
		chip.ir = opcode & 0x0FFF
		break
	case 0xB000:
		chip.pc = uint16(chip.reg[0]) + (opcode & 0x0FFF)
		break
	case 0xC000:
		chip.reg[(opcode & 0x0F00) >> 8] = uint8(rand.Int31n(255)) & uint8(opcode & 0x00FF)
		break
	case 0xD000:
		//TODO display
		break
	case 0xE000:
		switch opcode & 0x00FF {
		case 0x009E:
			if chip.key[(opcode & 0x0F00) >> 8] {
				chip.pc += 2
			}
			break
		case 0x00A1:
			if chip.key[(opcode & 0x0F00) >> 8] {
				chip.pc += 2
			}
			break
		default:
			chip.decodeFail(opcode)
		}
		break
	case 0xF000:
		x := &chip.reg[(opcode & 0x0F00) >> 8]
		switch opcode & 0x00FF {
		case 0x0007:
			*x = chip.delayTimer
			break
		case 0x000A:
			*x = chip.waitKey()
			break
		case 0x0015:
			chip.delayTimer = *x
			break
		case 0x0018:
			chip.soundTimer = *x
			break
		case 0x001E:
			chip.ir += uint16(*x)
			break
		case 0x0029:
			//TODO chip.ir = address to sprite in *x
			break
		case 0x0033:
			chip.mem[chip.ir]     = chip.reg[(opcode & 0x0F00) >> 8] / 100;
			chip.mem[chip.ir + 1] = (chip.reg[(opcode & 0x0F00) >> 8] / 10) % 10;
			chip.mem[chip.ir + 2] = (chip.reg[(opcode & 0x0F00) >> 8] % 100) % 10;
			break
		case 0x0055:
			ix := (opcode & 0x0F00) >> 8
			for i := uint16(0); i < ix; i++ {
				chip.mem[chip.ir + i] = chip.reg[i]
			}
			break
		case 0x0065:
			ix := (opcode & 0x0F00) >> 8
			for i := uint16(0); i < ix; i++ {
				chip.reg[i] = chip.mem[chip.ir + i]
			}
			break
		default:
			chip.decodeFail(opcode)
		}
		break
	default:
		chip.decodeFail(opcode)
	}
}

func (chip chip8) decodeFail(opcode uint16){
	fmt.Printf("Unknown opcode: 0x%X\n", opcode);
}

func (chip chip8) drawFlag() bool {
	return false
}

func (chip chip8) setKeys() {

}

// Waits for key press then returns the key index pressed
func (chip chip8) waitKey() uint8 {
	return 0
}