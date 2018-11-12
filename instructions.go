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
	0x01: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "ORA",
		Opcode:              0x01,
		Exec:                ORA,
	},
	0x05: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "ORA",
		Opcode:              0x05,
		Exec:                ORA,
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
	0x08: Instruction{
		Bytes:               1,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "PHP",
		Opcode:              0x08,
		Exec:                PHP,
	},
	0x09: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "ORA",
		Opcode:              0x09,
		Exec:                ORA,
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
	0x0D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "ORA",
		Opcode:              0x0D,
		Exec:                ORA,
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
	0x11: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: true,
		AddressingMode:      IndirectIndexed,
		Assembly:            "ORA",
		Opcode:              0x11,
		Exec:                ORA,
	},
	0x15: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "ORA",
		Opcode:              0x15,
		Exec:                ORA,
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
	0x19: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "ORA",
		Opcode:              0x19,
		Exec:                ORA,
	},
	0x1D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "ORA",
		Opcode:              0x1D,
		Exec:                ORA,
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
	0x26: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "ROL",
		Opcode:              0x26,
		Exec:                ROL,
	},
	0x28: Instruction{
		Bytes:               1,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "PLP",
		Opcode:              0x28,
		Exec:                PLP,
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
	0x2A: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Accumulator,
		Assembly:            "ROL",
		Opcode:              0x2A,
		Exec:                ROL,
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
	0x2E: Instruction{
		Bytes:               3,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "ROL",
		Opcode:              0x2E,
		Exec:                ROL,
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
	0x36: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "ROL",
		Opcode:              0x36,
		Exec:                ROL,
	},
	0x38: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "SEC",
		Opcode:              0x38,
		Exec:                SEC,
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
	0x3E: Instruction{
		Bytes:               3,
		Cycles:              7,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "ROL",
		Opcode:              0x3E,
		Exec:                ROL,
	},
	0x40: Instruction{
		Bytes:               1,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "RTI",
		Opcode:              0x40,
		Exec:                RTI,
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
	0x46: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "LSR",
		Opcode:              0x46,
		Exec:                LSR,
	},
	0x48: Instruction{
		Bytes:               1,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "PHA",
		Opcode:              0x48,
		Exec:                PHA,
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
	0x4E: Instruction{
		Bytes:               3,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "LSR",
		Opcode:              0x4E,
		Exec:                LSR,
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
	0x56: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "LSR",
		Opcode:              0x56,
		Exec:                LSR,
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
	0x5E: Instruction{
		Bytes:               3,
		Cycles:              7,
		AddCycleOnPageCross: false,
		AddressingMode:      AbsoluteX,
		Assembly:            "LSR",
		Opcode:              0x5E,
		Exec:                LSR,
	},
	0x60: Instruction{
		Bytes:               1,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "RTS",
		Opcode:              0x60,
		Exec:                RTS,
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
	0x66: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "ROR",
		Opcode:              0x66,
		Exec:                ROR,
	},
	0x68: Instruction{
		Bytes:               1,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "PLA",
		Opcode:              0x68,
		Exec:                PLA,
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
	0x6A: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Accumulator,
		Assembly:            "ROR",
		Opcode:              0x6A,
		Exec:                ROR,
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
	0x6E: Instruction{
		Bytes:               3,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "ROR",
		Opcode:              0x6E,
		Exec:                ROR,
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
	0x76: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "ROR",
		Opcode:              0x76,
		Exec:                ROR,
	},
	0x78: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "SEI",
		Opcode:              0x78,
		Exec:                SEI,
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
	0x7E: Instruction{
		Bytes:               3,
		Cycles:              7,
		AddCycleOnPageCross: false,
		AddressingMode:      AbsoluteX,
		Assembly:            "ROR",
		Opcode:              0x7E,
		Exec:                ROR,
	},
	0x81: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "STA",
		Opcode:              0x81,
		Exec:                STA,
	},
	0x84: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "STY",
		Opcode:              0x84,
		Exec:                STY,
	},
	0x85: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "STA",
		Opcode:              0x85,
		Exec:                STA,
	},
	0x86: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "STX",
		Opcode:              0x86,
		Exec:                STX,
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
	0x8A: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "TXA",
		Opcode:              0x8A,
		Exec:                TXA,
	},
	0x8C: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "STY",
		Opcode:              0x8C,
		Exec:                STY,
	},
	0x8D: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "STA",
		Opcode:              0x8D,
		Exec:                STA,
	},
	0x8E: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "STX",
		Opcode:              0x8E,
		Exec:                STX,
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
	0x91: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndirectIndexed,
		Assembly:            "STA",
		Opcode:              0x91,
		Exec:                STA,
	},
	0x94: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "STY",
		Opcode:              0x94,
		Exec:                STY,
	},
	0x95: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "STA",
		Opcode:              0x95,
		Exec:                STA,
	},
	0x96: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageY,
		Assembly:            "STA",
		Opcode:              0x96,
		Exec:                STX,
	},
	0x98: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "TYA",
		Opcode:              0x98,
		Exec:                TYA,
	},
	0x99: Instruction{
		Bytes:               3,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      AbsoluteY,
		Assembly:            "STA",
		Opcode:              0x99,
		Exec:                STA,
	},
	0x9A: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "TXS",
		Opcode:              0x9A,
		Exec:                TXS,
	},
	0x9D: Instruction{
		Bytes:               3,
		Cycles:              5,
		AddCycleOnPageCross: false,
		AddressingMode:      AbsoluteX,
		Assembly:            "STA",
		Opcode:              0x9D,
		Exec:                STA,
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
	0xA8: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "TAY",
		Opcode:              0xA8,
		Exec:                TAY,
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
	0xAA: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "TAX",
		Opcode:              0xAA,
		Exec:                TAX,
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
	0xBA: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "TSX",
		Opcode:              0xBA,
		Exec:                TSX,
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
	0xE1: Instruction{
		Bytes:               2,
		Cycles:              6,
		AddCycleOnPageCross: false,
		AddressingMode:      IndexedIndirect,
		Assembly:            "SBC",
		Opcode:              0xE1,
		Exec:                SBC,
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
	0xE5: Instruction{
		Bytes:               2,
		Cycles:              3,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPage,
		Assembly:            "SBC",
		Opcode:              0xE5,
		Exec:                SBC,
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
	0xE9: Instruction{
		Bytes:               2,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Immediate,
		Assembly:            "SBC",
		Opcode:              0xE9,
		Exec:                SBC,
	},
	0xEA: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "NOP",
		Opcode:              0xEA,
		Exec:                NOP,
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
	0xED: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      Absolute,
		Assembly:            "SBC",
		Opcode:              0xED,
		Exec:                SBC,
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
	0xF1: Instruction{
		Bytes:               2,
		Cycles:              5,
		AddCycleOnPageCross: true,
		AddressingMode:      IndirectIndexed,
		Assembly:            "SBC",
		Opcode:              0xF1,
		Exec:                SBC,
	},
	0xF5: Instruction{
		Bytes:               2,
		Cycles:              4,
		AddCycleOnPageCross: false,
		AddressingMode:      ZeroPageX,
		Assembly:            "SBC",
		Opcode:              0xF5,
		Exec:                SBC,
	},
	0xF8: Instruction{
		Bytes:               1,
		Cycles:              2,
		AddCycleOnPageCross: false,
		AddressingMode:      Implied,
		Assembly:            "SED",
		Opcode:              0xF8,
		Exec:                SED,
	},
	0xF9: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteY,
		Assembly:            "SBC",
		Opcode:              0xF9,
		Exec:                SBC,
	},
	0xFD: Instruction{
		Bytes:               3,
		Cycles:              4,
		AddCycleOnPageCross: true,
		AddressingMode:      AbsoluteX,
		Assembly:            "SBC",
		Opcode:              0xFD,
		Exec:                SBC,
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

// This function is used interally by the ADC and SBC instructions
var add = func(cpu *CPU, operand byte) {
	accumulator := cpu.A
	carry := cpu.flagToInt(cpu.CFlag)

	cpu.A = accumulator + operand + carry

	cpu.CFlag = cpu.A < accumulator

	// Formula for setting the overflow flag taken from:
	// http://www.righto.com/2012/12/the-6502-overflow-flag-explained.html
	cpu.VFlag = ((accumulator ^ cpu.A) & (operand ^ cpu.A) & 0x80) != 0

	cpu.setZeroAndNegativeFlags(cpu.A)
}

var AND = func(cpu *CPU, context *InstructionContext) {
	cpu.A = (cpu.A & cpu.Memory[context.Address])
	cpu.setZeroAndNegativeFlags(cpu.A)
}

var ADC = func(cpu *CPU, context *InstructionContext) {
	add(cpu, cpu.Memory[context.Address])
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

var NOP = func(cpu *CPU, context *InstructionContext) {}

var ORA = func(cpu *CPU, context *InstructionContext) {
	cpu.A = cpu.A | cpu.Memory[context.Address]
	cpu.setZeroAndNegativeFlags(cpu.A)
}

var PHA = func(cpu *CPU, context *InstructionContext) {
	cpu.stackPush(cpu.A)
}

var PHP = func(cpu *CPU, context *InstructionContext) {
	// Bits 5 and 4 are always set according to
	// https://wiki.nesdev.com/w/index.php/CPU_status_flag_behavior
	// though the flag status isn't changed
	cpu.stackPush(cpu.flagsToByte() | 0x30)
}

var PLA = func(cpu *CPU, context *InstructionContext) {
	cpu.A = cpu.stackPop()
}

var PLP = func(cpu *CPU, context *InstructionContext) {
	cpu.byteToFlags(cpu.stackPop())
}

var ROL = func(cpu *CPU, context *InstructionContext) {
	var operand *byte
	var flag bool
	if context.AddressingMode == Accumulator {
		operand = &cpu.A
	} else {
		operand = &cpu.Memory[context.Address]
	}

	flag = cpu.intToFlag(*operand & 0x80)
	*operand = *operand << 1
	*operand = *operand | cpu.flagToInt(cpu.CFlag)
	cpu.CFlag = flag
	cpu.setNegativeFlag(*operand)
}

var ROR = func(cpu *CPU, context *InstructionContext) {
	var operand *byte
	var flag bool

	if context.AddressingMode == Accumulator {
		operand = &cpu.A
	} else {
		operand = &cpu.Memory[context.Address]
	}

	flag = *operand > 0
	*operand = *operand >> 1
	if cpu.CFlag {
		*operand = *operand | 0x80
	}
	cpu.CFlag = flag
	cpu.setNegativeFlag(*operand)
}

var RTI = func(cpu *CPU, context *InstructionContext) {
	cpu.byteToFlags(cpu.stackPop())
}

var RTS = func(cpu *CPU, context *InstructionContext) {
	cpu.PC = cpu.stackPop16() + 1
}

var SBC = func(cpu *CPU, context *InstructionContext) {
	// Subtraction is the same as addition of the one's
	// complement. Here we flip the bits of the operand
	// (i.e. take the one's complement) and then run the
	// same logic as the ADC instruction
	operand := cpu.Memory[context.Address]
	add(cpu, ^operand)
}

var SEC = func(cpu *CPU, context *InstructionContext) {
	cpu.CFlag = true
}

var SED = func(cpu *CPU, context *InstructionContext) {
	cpu.DFlag = true
}

var SEI = func(cpu *CPU, context *InstructionContext) {
	cpu.IFlag = true
}

var STA = func(cpu *CPU, context *InstructionContext) {
	cpu.Memory[context.Address] = cpu.A
}

var STX = func(cpu *CPU, context *InstructionContext) {
	cpu.Memory[context.Address] = cpu.X
}

var STY = func(cpu *CPU, context *InstructionContext) {
	cpu.Memory[context.Address] = cpu.Y
}

var TAX = func(cpu *CPU, context *InstructionContext) {
	cpu.X = cpu.A
	cpu.setZeroAndNegativeFlags(cpu.X)
}

var TAY = func(cpu *CPU, context *InstructionContext) {
	cpu.Y = cpu.A
	cpu.setZeroAndNegativeFlags(cpu.Y)
}

var TSX = func(cpu *CPU, context *InstructionContext) {
	cpu.X = cpu.SP
	cpu.setZeroAndNegativeFlags(cpu.X)
}

var TXA = func(cpu *CPU, context *InstructionContext) {
	cpu.A = cpu.X
	cpu.setZeroAndNegativeFlags(cpu.A)
}

var TXS = func(cpu *CPU, context *InstructionContext) {
	cpu.SP = cpu.X
	cpu.setZeroAndNegativeFlags(cpu.SP)
}

var TYA = func(cpu *CPU, context *InstructionContext) {
	cpu.A = cpu.Y
	cpu.setZeroAndNegativeFlags(cpu.A)
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
