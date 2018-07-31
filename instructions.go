package main

import ()

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

func flagToInt(flag bool) uint8 {
	if flag {
		return 1
	} else {
		return 0
	}
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

		// Formula for overflow flag taken from:
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

var instructionMap = map[uint8]Instruction{
	0x21: AND,
	0x25: AND,
	0x29: AND,
	0x2D: AND,
	0x31: AND,
	0x35: AND,
	0x3D: AND,
	0x39: AND,
	0x65: ADC,
	0x69: ADC,
}

var addressingModeMap = map[uint8]AddressingMode{
	0x21: IndexedIndirect,
	0x25: ZeroPage,
	0x29: Immediate,
	0x2D: Absolute,
	0x31: IndirectIndexed,
	0x35: ZeroPageX,
	0x3D: AbsoluteX,
	0x39: AbsoluteY,
	0x65: ZeroPage,
	0x69: Immediate,
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

var instructionMetadata = map[uint8]InstructionMetadata{
	0x21: InstructionMetadata{
		Cycles: six,
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
	0x2D: InstructionMetadata{
		Cycles: four,
		Bytes:  3,
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
	0x65: InstructionMetadata{
		Cycles: three,
		Bytes:  2,
	},
	0x69: InstructionMetadata{
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

func six(context *InstructionContext) uint {
	return 6
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
