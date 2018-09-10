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

type Instruction struct {
	Bytes               uint16
	Cycles              uint
	AddCycleOnPageCross bool
	AddressingMode      AddressingMode
	Assembly            string
	Opcode              byte
	Exec                func(*CPU, *InstructionContext)
}

type InstructionContext struct {
	PageCrossed    bool
	Address        uint16
	AddressingMode AddressingMode
}

var instructionMap = map[uint8]Instruction{
	0x00: Instruction{
		Bytes:               2,
		Cycles:              7,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "BRK",
		Opcode:              0x00,
		Exec:                BRK,
	},
	0x06: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "ASL",
		Opcode:              0x06,
		Exec:                ASL,
	},
	0x0A: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Accumulator,
		Assembly:            "ASL",
		Opcode:              0x0A,
		Exec:                ASL,
	},
	0x0E: Instruction{
		Bytes:               3,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "ASL",
		Opcode:              0x0E,
		Exec:                ASL,
	},
	0x10: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BPL",
		Opcode:              0x10,
		Exec:                BPL,
	},
	0x16: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "ASL",
		Opcode:              0x16,
		Exec:                ASL,
	},
	0x18: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "CLC",
		Opcode:              0x18,
		Exec:                CLC,
	},
	0x1E: Instruction{
		Bytes:               3,
		Cycles:              7,
		AddCycleOnPageCross: false,
		AddressingMode:      AbsoluteX,
		Assembly:            "ASL",
		Opcode:              0x1E,
		Exec:                ASL,
	},
	0x20: Instruction{
		Bytes:               3,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "JSR",
		Opcode:              0x20,
		Exec:                JSR,
	},
	0x21: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "AND",
		Opcode:              0x21,
		Exec:                AND,
	},
	0x24: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "BIT",
		Opcode:              0x24,
		Exec:                BIT,
	},
	0x25: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "AND",
		Opcode:              0x25,
		Exec:                AND,
	},
	0x29: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "AND",
		Opcode:              0x29,
		Exec:                AND,
	},
	0x2C: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "BIT",
		Opcode:              0x2C,
		Exec:                BIT,
	},
	0x2D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "AND",
		Opcode:              0x2D,
		Exec:                AND,
	},
	0x30: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BMI",
		Opcode:              0x30,
		Exec:                BMI,
	},
	0x31: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: true,
		AddressingMode:      IndirectIndexed,
		Assembly:            "AND",
		Opcode:              0x31,
		Exec:                AND,
	},
	0x35: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "AND",
		Opcode:              0x35,
		Exec:                AND,
	},
	0x39: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "AND",
		Opcode:              0x39,
		Exec:                AND,
	},
	0x3D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "AND",
		Opcode:              0x3D,
		Exec:                AND,
	},
	0x41: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "EOR",
		Opcode:              0x41,
		Exec:                EOR,
	},
	0x45: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "EOR",
		Opcode:              0x45,
		Exec:                EOR,
	},
	0x49: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "EOR",
		Opcode:              0x49,
		Exec:                EOR,
	},
	0x4A: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Accumulator,
		Assembly:            "LSR",
		Opcode:              0x4A,
		Exec:                LSR,
	},
	0x4C: Instruction{
		Bytes:               3,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "JMP",
		Opcode:              0x4C,
		Exec:                JMP,
	},
	0x4D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "EOR",
		Opcode:              0x4D,
		Exec:                EOR,
	},
	0x50: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BVC",
		Opcode:              0x50,
		Exec:                BVC,
	},
	0x51: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: true,
		AddressingMode:      IndirectIndexed,
		Assembly:            "EOR",
		Opcode:              0x51,
		Exec:                EOR,
	},
	0x55: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "EOR",
		Opcode:              0x50,
		Exec:                EOR,
	},
	0x58: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "CLI",
		Opcode:              0x58,
		Exec:                CLI,
	},
	0x59: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "EOR",
		Opcode:              0x59,
		Exec:                EOR,
	},
	0x5D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "EOR",
		Opcode:              0x5D,
		Exec:                EOR,
	},
	0x61: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "ADC",
		Opcode:              0x61,
		Exec:                ADC,
	},
	0x65: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "ADC",
		Opcode:              0x65,
		Exec:                ADC,
	},
	0x69: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "ADC",
		Opcode:              0x69,
		Exec:                ADC,
	},
	0x6C: Instruction{
		Bytes:               3,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      Indirect,
		Assembly:            "JMP",
		Opcode:              0x6C,
		Exec:                JMP,
	},
	0x6D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "ADC",
		Opcode:              0x6D,
		Exec:                ADC,
	},
	0x70: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BVS",
		Opcode:              0x70,
		Exec:                BVS,
	},
	0x71: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: true,
		AddressingMode:      IndirectIndexed,
		Assembly:            "ADC",
		Opcode:              0x71,
		Exec:                ADC,
	},
	0x75: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "ADC",
		Opcode:              0x75,
		Exec:                ADC,
	},
	0x79: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "ADC",
		Opcode:              0x79,
		Exec:                ADC,
	},
	0x7D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "ADC",
		Opcode:              0x7D,
		Exec:                ADC,
	},
	0x88: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "DEY",
		Opcode:              0x88,
		Exec:                DEY,
	},
	0x90: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BCC",
		Opcode:              0x90,
		Exec:                BCC,
	},
	0xA0: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "LDY",
		Opcode:              0xA0,
		Exec:                LDY,
	},
	0xA1: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "LDA",
		Opcode:              0xA1,
		Exec:                LDA,
	},
	0xA2: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "LDX",
		Opcode:              0xA2,
		Exec:                LDX,
	},
	0xA4: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "LDY",
		Opcode:              0xA4,
		Exec:                LDY,
	},
	0xA5: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "LDA",
		Opcode:              0xA5,
		Exec:                LDA,
	},
	0xA6: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "LDX",
		Opcode:              0xA6,
		Exec:                LDX,
	},
	0xA9: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "LDA",
		Opcode:              0xA9,
		Exec:                LDA,
	},
	0xAC: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "LDY",
		Opcode:              0xAC,
		Exec:                LDY,
	},
	0xAD: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "LDA",
		Opcode:              0xAD,
		Exec:                LDA,
	},
	0xAE: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "LDX",
		Opcode:              0xAE,
		Exec:                LDX,
	},
	0xB0: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BCS",
		Opcode:              0xB0,
		Exec:                BCS,
	},
	0xB1: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: true,
		AddressingMode:      IndirectIndexed,
		Assembly:            "LDA",
		Opcode:              0xB1,
		Exec:                LDA,
	},
	0xB4: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "LDY",
		Opcode:              0xB4,
		Exec:                LDY,
	},
	0xB5: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "LDA",
		Opcode:              0xB5,
		Exec:                LDA,
	},
	0xB6: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageY,
		Assembly:            "LDX",
		Opcode:              0xB6,
		Exec:                LDX,
	},
	0xB8: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "CLV",
		Opcode:              0xB8,
		Exec:                CLV,
	},
	0xB9: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "LDA",
		Opcode:              0xB9,
		Exec:                LDA,
	},
	0xBC: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "LDY",
		Opcode:              0xBC,
		Exec:                LDY,
	},
	0xBD: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "LDA",
		Opcode:              0xBD,
		Exec:                LDA,
	},
	0xBE: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "LDX",
		Opcode:              0xBE,
		Exec:                LDX,
	},
	0xC0: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "CPY",
		Opcode:              0xC0,
		Exec:                CPY,
	},
	0xC1: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "CMP",
		Opcode:              0xC1,
		Exec:                CMP,
	},
	0xC4: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "CPY",
		Opcode:              0xC4,
		Exec:                CPY,
	},
	0xC5: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "CMP",
		Opcode:              0xC5,
		Exec:                CMP,
	},
	0xC6: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "DEC",
		Opcode:              0xC6,
		Exec:                DEC,
	},
	0xC8: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "INY",
		Opcode:              0xC8,
		Exec:                INY,
	},
	0xC9: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "CMP",
		Opcode:              0xC9,
		Exec:                CMP,
	},
	0xCA: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "DEX",
		Opcode:              0xCA,
		Exec:                DEX,
	},
	0xCC: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "CPY",
		Opcode:              0xCC,
		Exec:                CPY,
	},
	0xCD: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "CMP",
		Opcode:              0xCD,
		Exec:                CMP,
	},
	0xCE: Instruction{
		Bytes:               3,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "DEC",
		Opcode:              0xCE,
		Exec:                DEC,
	},
	0xD0: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BNE",
		Opcode:              0xD0,
		Exec:                BNE,
	},
	0xD1: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: true,
		AddressingMode:      IndirectIndexed,
		Assembly:            "CMP",
		Opcode:              0x90,
		Exec:                CMP,
	},
	0xD5: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "CMP",
		Opcode:              0xD5,
		Exec:                CMP,
	},
	0xD6: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "DEC",
		Opcode:              0xD6,
		Exec:                DEC,
	},
	0xD9: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "CMP",
		Opcode:              0xD9,
		Exec:                CMP,
	},
	0xDD: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "CMP",
		Opcode:              0xDD,
		Exec:                CMP,
	},
	0xDE: Instruction{
		Bytes:               3,
		Cycles:              7,
		AddCycleOnPageCross: false,
		AddressingMode:      AbsoluteX,
		Assembly:            "DEC",
		Opcode:              0xDE,
		Exec:                DEC,
	},
	0xE0: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "CPX",
		Opcode:              0xE0,
		Exec:                CPX,
	},
	0xE4: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "CPX",
		Opcode:              0xE4,
		Exec:                CPX,
	},
	0xE6: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "INC",
		Opcode:              0xE6,
		Exec:                INC,
	},
	0xE8: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "INX",
		Opcode:              0xE8,
		Exec:                INX,
	},
	0xEC: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "CPX",
		Opcode:              0xEC,
		Exec:                CPX,
	},
	0xEE: Instruction{
		Bytes:               3,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "INC",
		Opcode:              0xEE,
		Exec:                INC,
	},
	0xF0: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Relative,
		Assembly:            "BEQ",
		Opcode:              0xF0,
		Exec:                BEQ,
	},
	0xFE: Instruction{
		Bytes:               3,
		Cycles:              7,
		AddCycleOnPageCross: false,
		AddressingMode:      AbsoluteX,
		Assembly:            "INC",
		Opcode:              0xFE,
		Exec:                INC,
	},
	0xF6: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "INC",
		Opcode:              0xF6,
		Exec:                INC,
	},
}

var AND = func(cpu *CPU, context *InstructionContext) {
	cpu.A = (cpu.A & cpu.Memory[context.Address])
	cpu.setZeroAndNegativeFlags(cpu.A)
}

var ADC = func(cpu *CPU, context *InstructionContext) {
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

	cpu.setZeroAndNegativeFlags(cpu.A)
}

var ASL = func(cpu *CPU, context *InstructionContext) {
	var operand byte

	if cpu.A&0x80 != 0 {
		cpu.CFlag = true
	} else {
		cpu.CFlag = false
	}

	if context.AddressingMode == Accumulator {
		operand = cpu.A
	} else {
		operand = cpu.Memory[context.Address]
	}

	cpu.A = operand << 1
	cpu.setZeroAndNegativeFlags(cpu.A)
}

var BCC = func(cpu *CPU, context *InstructionContext) {
	if !cpu.CFlag {
		branchRelative(cpu, context)
	}
}

var BCS = func(cpu *CPU, context *InstructionContext) {
	if cpu.CFlag {
		branchRelative(cpu, context)
	}
}

var BEQ = func(cpu *CPU, context *InstructionContext) {
	if cpu.ZFlag {
		branchRelative(cpu, context)
	}
}

var BIT = func(cpu *CPU, context *InstructionContext) {
	operand := cpu.Memory[context.Address]
	if (cpu.A & operand) == 0 {
		cpu.ZFlag = true
	} else {
		cpu.ZFlag = false
	}

	cpu.NFlag = cpu.intToFlag(operand & 0x80)
	cpu.VFlag = cpu.intToFlag(operand & 0x40)
}

var BMI = func(cpu *CPU, context *InstructionContext) {
	if cpu.NFlag {
		branchRelative(cpu, context)
	}
}

var BNE = func(cpu *CPU, context *InstructionContext) {
	if !cpu.ZFlag {
		branchRelative(cpu, context)
	}
}

var BPL = func(cpu *CPU, context *InstructionContext) {
	if !cpu.NFlag {
		branchRelative(cpu, context)
	}
}

var BRK = func(cpu *CPU, context *InstructionContext) {
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
}

var BVC = func(cpu *CPU, context *InstructionContext) {
	if !cpu.VFlag {
		branchRelative(cpu, context)
	}
}

var BVS = func(cpu *CPU, context *InstructionContext) {
	if cpu.VFlag {
		branchRelative(cpu, context)
	}
}

var CLC = func(cpu *CPU, context *InstructionContext) {
	cpu.CFlag = false
}

var CLI = func(cpu *CPU, context *InstructionContext) {
	cpu.IFlag = false
}

var CLV = func(cpu *CPU, context *InstructionContext) {
	cpu.VFlag = false
}

var CMP = func(cpu *CPU, context *InstructionContext) {
	compare(cpu, cpu.A, cpu.Memory[context.Address])
}

var CPX = func(cpu *CPU, context *InstructionContext) {
	compare(cpu, cpu.X, cpu.Memory[context.Address])
}

var CPY = func(cpu *CPU, context *InstructionContext) {
	compare(cpu, cpu.Y, cpu.Memory[context.Address])
}

var DEC = func(cpu *CPU, context *InstructionContext) {
	cpu.Memory[context.Address] = decrement(cpu, cpu.Memory[context.Address])
}

var DEX = func(cpu *CPU, context *InstructionContext) {
	cpu.X = decrement(cpu, cpu.X)
}

var DEY = func(cpu *CPU, context *InstructionContext) {
	cpu.Y = decrement(cpu, cpu.Y)
}

var EOR = func(cpu *CPU, context *InstructionContext) {
	operand := cpu.Memory[context.Address]
	cpu.A = cpu.A ^ operand
	cpu.setZeroAndNegativeFlags(cpu.A)
}

var INC = func(cpu *CPU, context *InstructionContext) {
	cpu.Memory[context.Address] = increment(cpu, cpu.Memory[context.Address])
}

var INX = func(cpu *CPU, context *InstructionContext) {
	cpu.X = increment(cpu, cpu.X)
}

var INY = func(cpu *CPU, context *InstructionContext) {
	cpu.Y = increment(cpu, cpu.Y)
}

var JMP = func(cpu *CPU, context *InstructionContext) {
	cpu.PC = context.Address
}

var JSR = func(cpu *CPU, context *InstructionContext) {
	cpu.stackPush16(cpu.PC - 1)
	cpu.PC = context.Address
}

var LDA = func(cpu *CPU, context *InstructionContext) {
	cpu.A = cpu.Memory[context.Address]
	cpu.setZeroAndNegativeFlags(cpu.A)
}

var LDX = func(cpu *CPU, context *InstructionContext) {
	cpu.X = cpu.Memory[context.Address]
	cpu.setZeroAndNegativeFlags(cpu.X)
}

var LDY = func(cpu *CPU, context *InstructionContext) {
	cpu.Y = cpu.Memory[context.Address]
	cpu.setZeroAndNegativeFlags(cpu.Y)
}

var LSR = func(cpu *CPU, context *InstructionContext) {
	var operand *byte

	if context.AddressingMode == Accumulator {
		operand = &cpu.A
	} else {
		operand = &cpu.Memory[context.Address]
	}

	cpu.CFlag = cpu.intToFlag(*operand & 0x80)
	*operand = *operand >> 1
	cpu.setZeroAndNegativeFlags(*operand)
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
	cpu.setZeroAndNegativeFlags(result)
}

func decrement(cpu *CPU, target byte) byte {
	result := target - 1
	cpu.setZeroAndNegativeFlags(result)

	return result
}

func increment(cpu *CPU, target byte) byte {
	result := target + 1
	cpu.setZeroAndNegativeFlags(result)

	return result
}
