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
	fmt.Println("So i can keep fmt imported")
}
func flagToInt(flag bool) uint8 {
	if flag {
		return 1
	} else {
		return 0
	}
}

func intToFlag(n uint8) bool {
	if n == 0 {
		return false
	} else {
		return true
	}
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

var AND = Instruction{
	Assembly: "AND",
	Exec: func(cpu *CPU, context *InstructionContext) {
		cpu.A = (cpu.A & cpu.Memory[context.Address])
		cpu.setZeroFlag()
		cpu.setNegativeFlag()
	},
}

var ADC = Instruction{
	Assembly: "ADC",
	Exec: func(cpu *CPU, context *InstructionContext) {
		accumulator := cpu.A
		operand := cpu.Memory[context.Address]
		carry := flagToInt(cpu.CFlag)

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

		cpu.setZeroFlag()
		cpu.setNegativeFlag()
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

		cpu.setZeroFlag()
		cpu.setNegativeFlag()
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

		cpu.NFlag = intToFlag(operand & 0x80)
		cpu.VFlag = intToFlag(operand & 0x40)
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

var instructionMap = map[uint8]Instruction{
	0x06: ASL,
	0x0A: ASL,
	0x0E: ASL,
	0x10: BPL,
	0x16: ASL,
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
	0xD0: BNE,
	0xF0: BEQ,
}

var addressingModeMap = map[uint8]AddressingMode{
	0x06: ZeroPage,
	0x0A: Accumulator,
	0x0E: Absolute,
	0x10: Relative,
	0x16: ZeroPageX,
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
	0xD0: Relative,
	0xF0: Relative,
}

var instructionMetadata = map[uint8]InstructionMetadata{
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
	0xD0: InstructionMetadata{
		Cycles: two,
		Bytes:  2,
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
// type InstructionMetadata struct {
// 	Cycles func(context *InstructionContext) uint
// 	Bytes  uint16
//  Assembly string
//  AddCycleOnPageCross bool
//  AddressingMode AddressingMode
// }
//
// type Instruction struct {
// 	Metadata InstructionMetadata
// 	Exec     func(*CPU, *InstructionContext)
// }
//
// var instructionMap = map[uint8]InstructionRefactor{
// 	0x21: Instruction{
// 		Metadata: {
// 			Cycles: 6,
// 			Bytes:  2,
// 			AddCycleOnPageCross: true,
//          AddressingMode: Immediate,
// 		},
//      Assembly: "AND",
// 		Exec: AND,
// 	},
// }
