package main

import (
	"fmt"
)

type AddressingMode uint8

const (
	_               AddressingMode = iota
	Absolute                       // 1
	AbsoluteX                      // 2
	AbsoluteY                      // 3
	Accumulator                    // 4
	Immediate                      // 5
	Implied                        // 6
	IndexedIndirect                // 7
	Indirect                       // 8
	IndirectIndexed                // 9
	Relative                       // 10
	ZeroPage                       // 11
	ZeroPageX                      // 12
	ZeroPageY                      // 13
)

type InstructionContext struct {
	Mode        AddressingMode
	PageCrossed bool
	Address     uint16
}

type Instruction struct {
	Assembly string
	Exec     func(*CPU, *InstructionContext)
}

type InstructionMetadata struct {
	Cycles func(context *InstructionContext) uint
	Bytes  uint16
}

func fakeFunctionNeverCalled() {
	fmt.Println("So i can keep fmt imported, lol")
}
func branchRelative(cpu *CPU, context *InstructionContext) {
	var branchLocation uint16

	cpu.Cycles += 1
	// convert operand to signed offset
	relativeAddress := int8(cpu.Memory[context.Address])
	// convert signed offset to 16bit unsigned. If we don't convert
	// to a signed int8 first, we will not preserve the sign bits
	// meaning that addition will not correctly handle negative
	// operands
	branchLocation = cpu.PC + uint16(relativeAddress)

	if (branchLocation & 0xFF00) != (cpu.PC & 0xFF00) {
		// page crossed
		cpu.Cycles += 2
	}

	cpu.PC = branchLocation
}

func compare(cpu *CPU, register byte, operand byte) {
	result := register - operand

	cpu.CFlag = int8(result) >= 0
	cpu.setNegativeFlag(result)
	cpu.setZeroFlag(result)
}

var AND = Instruction{
	Assembly: "AND",
	Exec: func(cpu *CPU, context *InstructionContext) {
		cpu.A = (cpu.A & cpu.Memory[context.Address])
		cpu.setZeroFlag(cpu.A)
		cpu.setNegativeFlag(cpu.A)
	},
}

var ADC = Instruction{
	Assembly: "ADC",
	Exec: func(cpu *CPU, context *InstructionContext) {
		accumulator := cpu.A
		operand := cpu.Memory[context.Address]
		carry := cpu.flagToInt(cpu.CFlag)

		cpu.A = accumulator + operand + carry

		if cpu.A < accumulator {
			cpu.CFlag = true
		} else {
			cpu.CFlag = false
		}

		// Formula for setting the overflow flag taken from:
		// http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
		if ((accumulator ^ cpu.A) & (operand ^ cpu.A) & 0x80) != 0 {
			cpu.VFlag = true
		} else {
			cpu.VFlag = false
		}

		cpu.setZeroFlag(cpu.A)
		cpu.setNegativeFlag(cpu.A)
	},
}

var ASL = Instruction{
	Assembly: "ASL",
	Exec: func(cpu *CPU, context *InstructionContext) {
		var operand byte

		if cpu.A&0x80 != 0 {
			cpu.CFlag = true
		} else {
			cpu.CFlag = false
		}

		if context.Mode == Accumulator {
			operand = cpu.A
		} else {
			operand = cpu.Memory[context.Address]
		}

		cpu.A = operand << 1

		cpu.setZeroFlag(cpu.A)
		cpu.setNegativeFlag(cpu.A)
	},
}

var BCC = Instruction{
	Assembly: "BCC",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if !cpu.CFlag {
			branchRelative(cpu, context)
		}
	},
}

var BCS = Instruction{
	Assembly: "BCS",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if cpu.CFlag {
			branchRelative(cpu, context)
		}
	},
}

var BEQ = Instruction{
	Assembly: "BEQ",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if cpu.ZFlag {
			branchRelative(cpu, context)
		}
	},
}

var BIT = Instruction{
	Assembly: "BIT",
	Exec: func(cpu *CPU, context *InstructionContext) {
		operand := cpu.Memory[context.Address]
		if (cpu.A & operand) == 0 {
			cpu.ZFlag = true
		} else {
			cpu.ZFlag = false
		}

		cpu.NFlag = cpu.intToFlag(operand & 0x80)
		cpu.VFlag = cpu.intToFlag(operand & 0x40)
	},
}

var BMI = Instruction{
	Assembly: "BMI",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if cpu.NFlag {
			branchRelative(cpu, context)
		}
	},
}

var BNE = Instruction{
	Assembly: "BNE",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if !cpu.ZFlag {
			branchRelative(cpu, context)
		}
	},
}

var BPL = Instruction{
	Assembly: "BPL",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if !cpu.NFlag {
			branchRelative(cpu, context)
		}
	},
}

var BRK = Instruction{
	Assembly: "BRK",
	Exec: func(cpu *CPU, context *InstructionContext) {
		// push the PC onto the stack
		cpu.stackPush(byte(cpu.PC >> 8)) // high byte first
		cpu.stackPush(byte(cpu.PC))      // then the low byte
		// push the status flags onto the stack
		cpu.stackPush(cpu.flagsToByte())
		// load the interrupt address from $FFFE and $FFFF
		lo := uint16(cpu.Memory[0xFFFE])
		hi := uint16(cpu.Memory[0xFFFF])
		cpu.PC = (hi << 8) | lo
		// Set the break flag
		cpu.BFlag = true // When does this get unset?
	},
}

var BVC = Instruction{
	Assembly: "BVC",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if !cpu.VFlag {
			branchRelative(cpu, context)
		}
	},
}

var BVS = Instruction{
	Assembly: "BVS",
	Exec: func(cpu *CPU, context *InstructionContext) {
		if cpu.VFlag {
			branchRelative(cpu, context)
		}
	},
}

var CLC = Instruction{
	Assembly: "CLC",
	Exec: func(cpu *CPU, context *InstructionContext) {
		cpu.CFlag = false
	},
}

var CLI = Instruction{
	Assembly: "CLI",
	Exec: func(cpu *CPU, context *InstructionContext) {
		cpu.IFlag = false
	},
}

var CLV = Instruction{
	Assembly: "CLV",
	Exec: func(cpu *CPU, context *InstructionContext) {
		cpu.VFlag = false
	},
}

var CMP = Instruction{
	Assembly: "CMP",
	Exec: func(cpu *CPU, context *InstructionContext) {
		compare(cpu, cpu.A, cpu.Memory[context.Address])
	},
}

var CPX = Instruction{
	Assembly: "CPX",
	Exec: func(cpu *CPU, context *InstructionContext) {
		compare(cpu, cpu.X, cpu.Memory[context.Address])
	},
}

var CPY = Instruction{
	Assembly: "CPY",
	Exec: func(cpu *CPU, context *InstructionContext) {
		compare(cpu, cpu.Y, cpu.Memory[context.Address])
	},
}

var instructionMap = map[uint8]Instruction{
	0x00: BRK,
	0x06: ASL,
	0x0A: ASL,
	0x0E: ASL,
	0x10: BPL,
	0x16: ASL,
	0x18: CLC,
	0x1E: ASL,
	0x21: AND,
	0x24: BIT,
	0x25: AND,
	0x29: AND,
	0x2C: BIT,
	0x2D: AND,
	0x30: BMI,
	0x31: AND,
	0x35: AND,
	0x3D: AND,
	0x39: AND,
	0x50: BVC,
	0x58: CLI,
	0x61: ADC,
	0x65: ADC,
	0x69: ADC,
	0x6D: ADC,
	0x70: BVS,
	0x71: ADC,
	0x75: ADC,
	0x7D: ADC,
	0x79: ADC,
	0x90: BCC,
	0xB0: BCS,
	0xB8: CLV,
	0xC0: CPY,
	0xC1: CMP,
	0xC4: CPY,
	0xC5: CMP,
	0xC9: CMP,
	0xCC: CPY,
	0xCD: CMP,
	0xD0: BNE,
	0xD1: CMP,
	0xD5: CMP,
	0xD9: CMP,
	0xDD: CMP,
	0xE0: CPX,
	0xE4: CPX,
	0xEC: CPX,
	0xF0: BEQ,
}

var addressingModeMap = map[uint8]AddressingMode{
	0x00: Implied,
	0x06: ZeroPage,
	0x0A: Accumulator,
	0x0E: Absolute,
	0x10: Relative,
	0x16: ZeroPageX,
	0x18: Implied,
	0x1E: AbsoluteX,
	0x21: IndexedIndirect,
	0x24: ZeroPage,
	0x25: ZeroPage,
	0x29: Immediate,
	0x2C: Absolute,
	0x2D: Absolute,
	0x30: Relative,
	0x31: IndirectIndexed,
	0x35: ZeroPageX,
	0x3D: AbsoluteX,
	0x39: AbsoluteY,
	0x50: Relative,
	0x58: Implied,
	0x61: IndexedIndirect,
	0x65: ZeroPage,
	0x69: Immediate,
	0x6D: Absolute,
	0x70: Relative,
	0x71: IndirectIndexed,
	0x75: ZeroPageX,
	0x7D: AbsoluteX,
	0x79: AbsoluteY,
	0x90: Relative,
	0xB0: Relative,
	0xB8: Implied,
	0xC0: Immediate,
	0xC1: IndexedIndirect,
	0xC4: ZeroPage,
	0xC5: ZeroPage,
	0xC9: Immediate,
	0xCC: Absolute,
	0xCD: Absolute,
	0xD0: Relative,
	0xD1: IndirectIndexed,
	0xD5: ZeroPageX,
	0xD9: AbsoluteY,
	0xDD: AbsoluteX,
	0xE0: Immediate,
	0xE4: ZeroPage,
	0xEC: Absolute,
	0xF0: Relative,
}

var instructionMetadata = map[uint8]InstructionMetadata{
	0x00: InstructionMetadata{
		Cycles: seven,
		Bytes:  2,
	},
	0x06: InstructionMetadata{
		Cycles: five,
		Bytes:  2,
	},
	0x0A: InstructionMetadata{
		Cycles: two,
		Bytes:  1,
	},
	0x0E: InstructionMetadata{
		Cycles: six,
		Bytes:  3,
	},
	0x10: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0x16: InstructionMetadata{
		Cycles: six,
		Bytes:  2,
	},
	0x18: InstructionMetadata{
		Cycles: two,
		Bytes:  1,
	},
	0x1E: InstructionMetadata{
		Cycles: seven,
		Bytes:  3,
	},
	0x21: InstructionMetadata{
		Cycles: six,
		Bytes:  2,
	},
	0x24: InstructionMetadata{
		Cycles: three,
		Bytes:  2,
	},
	0x25: InstructionMetadata{
		Cycles: three,
		Bytes:  2,
	},
	0x29: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0x2C: InstructionMetadata{
		Cycles: four,
		Bytes:  3,
	},
	0x2D: InstructionMetadata{
		Cycles: four,
		Bytes:  3,
	},
	0x30: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0x31: InstructionMetadata{
		Cycles: fiveIncrementOnPageCross,
		Bytes:  2,
	},
	0x35: InstructionMetadata{
		Cycles: four,
		Bytes:  2,
	},
	0x3D: InstructionMetadata{
		Cycles: fourIncrementOnPageCross,
		Bytes:  3,
	},
	0x39: InstructionMetadata{
		Cycles: fourIncrementOnPageCross,
		Bytes:  3,
	},
	0x50: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0x58: InstructionMetadata{
		Cycles: two,
		Bytes:  1,
	},
	0x65: InstructionMetadata{
		Cycles: three,
		Bytes:  2,
	},
	0x61: InstructionMetadata{
		Cycles: six,
		Bytes:  2,
	},
	0x69: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0x6D: InstructionMetadata{
		Cycles: four,
		Bytes:  3,
	},
	0x70: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0x71: InstructionMetadata{
		Cycles: fiveIncrementOnPageCross,
		Bytes:  2,
	},
	0x75: InstructionMetadata{
		Cycles: four,
		Bytes:  2,
	},
	0x7D: InstructionMetadata{
		Cycles: fourIncrementOnPageCross,
		Bytes:  3,
	},
	0x79: InstructionMetadata{
		Cycles: fourIncrementOnPageCross,
		Bytes:  3,
	},
	0x90: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0xB0: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0xB8: InstructionMetadata{
		Cycles: two,
		Bytes:  1,
	},
	0xC0: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0xC1: InstructionMetadata{
		Cycles: six,
		Bytes:  2,
	},
	0xC4: InstructionMetadata{
		Cycles: three,
		Bytes:  2,
	},
	0xC5: InstructionMetadata{
		Cycles: three,
		Bytes:  2,
	},
	0xC9: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0xCC: InstructionMetadata{
		Cycles: four,
		Bytes:  3,
	},
	0xCD: InstructionMetadata{
		Cycles: four,
		Bytes:  3,
	},
	0xD0: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0xD1: InstructionMetadata{
		Cycles: fiveIncrementOnPageCross,
		Bytes:  2,
	},
	0xD5: InstructionMetadata{
		Cycles: four,
		Bytes:  2,
	},
	0xD9: InstructionMetadata{
		Cycles: fourIncrementOnPageCross,
		Bytes:  3,
	},
	0xDD: InstructionMetadata{
		Cycles: fourIncrementOnPageCross,
		Bytes:  3,
	},
	0xE0: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
	0xE4: InstructionMetadata{
		Cycles: three,
		Bytes:  2,
	},
	0xEC: InstructionMetadata{
		Cycles: four,
		Bytes:  3,
	},
	0xF0: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
	},
}

func two(context *InstructionContext) uint {
	return 2
}

func three(context *InstructionContext) uint {
	return 3
}

func four(context *InstructionContext) uint {
	return 4
}

func five(context *InstructionContext) uint {
	return 5
}

func six(context *InstructionContext) uint {
	return 6
}

func seven(context *InstructionContext) uint {
	return 7
}

func fourIncrementOnPageCross(context *InstructionContext) uint {
	if context.PageCrossed {
		return 5
	} else {
		return 4
	}
}

func fiveIncrementOnPageCross(context *InstructionContext) uint {
	if context.PageCrossed {
		return 6
	} else {
		return 5
	}
}

// TODO: Potential refactor of instruction struct, metadata, and map.
//       Would allow the use of a single map to lookup instructions
//       and their metadata.
//
// type InstructionRefactor struct {
// 	Cycles func(context *InstructionContext) uint
// 	Bytes  uint16
//  AddCycleOnPageCross bool
//  AddressingMode AddressingMode
//  Assembly string
//  Opcode   byte
// 	Exec     func(*InstructionContext)
// }
//
// var instructionMap = map[uint8]InstructionRefactor{
// 	0x21: Instruction{
// 		Cycles: 6,
// 		Bytes:  2,
// 		AddCycleOnPageCross: true,
//      AddressingMode: Immediate,
//      Assembly: "AND",
//      Opcode: 0x21
// 		Exec: AND,
// 	},
// }
