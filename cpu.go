package main

import (
	"fmt"
)

type CPU struct {
	PC     uint16
	SP     uint8
	A      uint8
	X      uint8
	Y      uint8
	F      byte // 8 bit status flags
	NFlag  bool // negative flag
	VFlag  bool // overflow flag
	UFlag  bool // unused
	BFlag  bool // technically unused (doesn't exist on hardware processor)
	DFlag  bool // decimal flag - no effect on NES
	IFlag  bool // interrupt disable flag
	ZFlag  bool // zero flag
	CFlag  bool // carry flag
	Cycles uint
	Memory [0x10000]byte
	Debug  bool
}

func (cpu *CPU) Print() {
	fmt.Printf("PC: %d SP: %d A: %d X: %d Y: %d, Cycles: %d\n", cpu.PC, cpu.SP, cpu.A, cpu.X, cpu.Y, cpu.Cycles)
}

func (cpu *CPU) PrintTest(instruction Instruction) {
	w0 := fmt.Sprintf("%02X", cpu.Memory[cpu.PC+0])
	w1 := fmt.Sprintf("%02X", cpu.Memory[cpu.PC+1])
	w2 := fmt.Sprintf("%02X", cpu.Memory[cpu.PC+2])
	if instruction.Bytes < 2 {
		w1 = "  "
	}
	if instruction.Bytes < 3 {
		w2 = "  "
	}
	fmt.Printf(
		"%4X  %s %s %s  %s %28s"+
			"A:%02X X:%02X Y:%02X P:%02X SP:%02X PPU:%3d\n",
		cpu.PC, w0, w1, w2, instruction.Assembly, "",
		cpu.A, cpu.X, cpu.Y, cpu.flagsToByte(), cpu.SP, (cpu.Cycles*3)%341)
}

func (cpu *CPU) Exec() {
	opcode := cpu.Memory[cpu.PC]
	context := context(cpu, opcode)

	if cpu.Debug {
		cpu.PrintTest(instructionMap[opcode])
	}

	instruction := instructionMap[opcode]
	cpu.PC += instruction.Bytes
	instruction.Exec(cpu, context)
	cpu.Cycles += instruction.Cycles
	if context.PageCrossed && instruction.AddCycleOnPageCross {
		cpu.Cycles += 1
	}
}

func (cpu *CPU) setZeroFlag(n byte) {
	// set zero flag if input is zero
	cpu.ZFlag = n == 0
}

func (cpu *CPU) setNegativeFlag(n byte) {
	// Set negative flag if input bit 7 is 1
	cpu.NFlag = n&0x80 > 0
}

func (cpu *CPU) setZeroAndNegativeFlags(n byte) {
	cpu.setZeroFlag(n)
	cpu.setNegativeFlag(n)
}
func (cpu *CPU) flagToInt(flag bool) uint8 {
	if flag {
		return 1
	} else {
		return 0
	}
}

func (cpu *CPU) intToFlag(n uint8) bool {
	return n != 0
}

func (cpu *CPU) flagsToByte() byte {
	return (cpu.flagToInt(cpu.NFlag) << 7) |
		(cpu.flagToInt(cpu.VFlag) << 6) |
		(cpu.flagToInt(cpu.UFlag) << 5) |
		(cpu.flagToInt(cpu.BFlag) << 4) |
		(cpu.flagToInt(cpu.DFlag) << 3) |
		(cpu.flagToInt(cpu.IFlag) << 2) |
		(cpu.flagToInt(cpu.ZFlag) << 1) |
		(cpu.flagToInt(cpu.CFlag) << 0)
}

func (cpu *CPU) byteToFlags(flags byte) {
	cpu.NFlag = cpu.intToFlag(flags & 0x80)
	cpu.VFlag = cpu.intToFlag(flags & 0x40)
	cpu.UFlag = cpu.intToFlag(flags & 0x20)
	cpu.BFlag = cpu.intToFlag(flags & 0x10)
	cpu.DFlag = cpu.intToFlag(flags & 0x08)
	cpu.IFlag = cpu.intToFlag(flags & 0x04)
	cpu.ZFlag = cpu.intToFlag(flags & 0x02)
	cpu.CFlag = cpu.intToFlag(flags & 0x01)
}

func (cpu *CPU) stackPush16(value uint16) {
	cpu.stackPush(byte(value >> 8)) // high byte first
	cpu.stackPush(byte(value))      // then the low byte
}

func (cpu *CPU) stackPop16() uint16 {
	lo := cpu.stackPop()
	hi := cpu.stackPop()
	return uint16(hi)<<8 | uint16(lo)

}

func (cpu *CPU) stackPush(value byte) {
	cpu.Memory[0x100|uint16(cpu.SP)] = value
	cpu.SP -= 1
}

func (cpu *CPU) stackPop() byte {
	cpu.SP += 1
	return cpu.Memory[0x100|uint16(cpu.SP)]
}

func context(cpu *CPU, opcode byte) *InstructionContext {
	var address uint16
	var pageCrossed = false
	var mode = instructionMap[opcode].AddressingMode

	switch mode {
	case Immediate:
		address = cpu.PC + 1
	case Accumulator:
	case ZeroPage:
		address = uint16(cpu.Memory[cpu.PC+1]) & 0x00FF
	case ZeroPageX:
		address = uint16(cpu.Memory[cpu.PC+1]+cpu.X) & 0x00FF
	case ZeroPageY:
		address = uint16(cpu.Memory[cpu.PC+1]+cpu.Y) & 0x00FF
	case Absolute:
		address = uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])
	case AbsoluteX:
		address = (uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])) + uint16(cpu.X)
		if (address & 0x00FF) < uint16(cpu.X) {
			pageCrossed = true
		}
	case AbsoluteY:
		address = (uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])) + uint16(cpu.Y)
		if (address & 0x00FF) < uint16(cpu.Y) {
			pageCrossed = true
		}
	case IndexedIndirect:
		intermediateAddress := (uint8(cpu.Memory[cpu.PC+1]) + cpu.X)
		lo := cpu.Memory[intermediateAddress]
		hi := cpu.Memory[intermediateAddress+1]
		address = uint16(hi)<<8 | uint16(lo)
	case IndirectIndexed:
		intermediateAddress := cpu.Memory[cpu.PC+1]
		address = uint16(cpu.Memory[intermediateAddress]) + uint16(cpu.Y)
		if (address & 0x00FF) < uint16(cpu.Y) {
			pageCrossed = true
		}
	case Relative:
		address = cpu.PC + 1
	case Indirect:
		// NOTE: This addressing mode implements a hardware bug.
		// The operand of the instruction is an intermediate address.
		// This intermediate address contains the low byte of the JMP target address.
		// The high byte of the JMP target address is stored in the following address
		// (intermediate address + 1). However If the intermediate address falls on a
		// memory page boundary (i.e. the first byte is on $xxFF, where xx is any number),
		// it does not correctly look at the next page when reading the high byte of the
		// JMP target address. A concrete example: If the instruction has the operand $10FF,
		// it will read the LSB of the JMP address from $10FF, but will read the MSB of the JMP
		// address from $1000 instead of $1100.
		intermediateLo := uint16(cpu.Memory[cpu.PC+2])<<8 | uint16(cpu.Memory[cpu.PC+1])
		intermediateHi := (intermediateLo & 0xFF00) | ((intermediateLo + 1) & 0x00FF) // this is the bug
		address = uint16(cpu.Memory[intermediateHi])<<8 | uint16(cpu.Memory[intermediateLo])
	case Implied:
	}

	return &InstructionContext{
		PageCrossed:    pageCrossed,
		Address:        address,
		AddressingMode: mode,
	}
}

func NewCPU() *CPU {
	return &CPU{
		PC:     0x00,
		SP:     0xFD,
		A:      0x00,
		X:      0x00,
		Y:      0x00,
		F:      0x00,
		Cycles: 0,
	}
}
