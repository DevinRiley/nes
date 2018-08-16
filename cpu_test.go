package main

import (
	"testing"
)

func TestNewCPUSetsSP(t *testing.T) {
	cpu := NewCPU()
	if cpu.SP != 0xFF {
		t.Error("Stack Pointer not correctly initialized, got", cpu.SP)
	}
}

func TestStackPush(t *testing.T) {
	cpu := NewCPU()
	cpu.stackPush(0x17)

	if cpu.Memory[0x1FF] != 0x17 {
		t.Error("did not add expected value to top of stack")
	}

	if cpu.SP != 0xFE {
		t.Error("did not decrement stack pointer")
	}
}

func TestStackPop(t *testing.T) {
	cpu := NewCPU()
	cpu.SP = 0xFD
	cpu.Memory[0x1FE] = 0x50
	result := cpu.stackPop()

	if result != 0x50 {
		t.Error("did not return expected value, got", result)
	}

	if cpu.SP != 0xFE {
		t.Error("did not increment stack pointer")
	}
}

func TestFlagsToByte(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = false
	cpu.VFlag = true
	cpu.UFlag = false
	cpu.BFlag = true
	cpu.DFlag = false
	cpu.IFlag = true
	cpu.ZFlag = false
	cpu.CFlag = true
	result := cpu.flagsToByte()

	if result != 0x55 {
		t.Error("did not correctly set bit flags, got", result)
	}
}

func TestANDImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.Memory[0] = 0x29
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.A != 0x01 {
		t.Error("immediate (opcode 0x29) failed to give correct Accumulator value")
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestANDImmediateZero(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.Memory[0] = 0x29
	cpu.Memory[1] = 0x00
	cpu.Exec()

	if cpu.A != 0x00 {
		t.Error("immediate (opcode 0x29) failed to give correct Accumulator value")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestANDZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x02
	cpu.Memory[0] = 0x25
	cpu.Memory[1] = 0x09
	cpu.Memory[9] = 0x02
	cpu.Exec()

	if cpu.A != 0x02 {
		t.Error("zero page (opcode 0x25) failed to give correct Accumulator value")
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestANDZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x03
	cpu.X = 0x01
	cpu.Memory[0] = 0x35
	cpu.Memory[1] = 0x08
	cpu.Memory[9] = 0x03
	cpu.Exec()

	if cpu.A != 0x03 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestANDZeroPageXWithOverflow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x03
	cpu.X = 0xFF // should overflow the result and wraparound
	cpu.Memory[0] = 0x35
	cpu.Memory[1] = 0x0A
	cpu.Memory[9] = 0x03
	cpu.Exec()

	if cpu.A != 0x03 {
		t.Error("zero page X (opcode 0x35) failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestANDAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x03
	cpu.Memory[0] = 0x2D
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x03
	cpu.Exec()

	if cpu.A != 0x03 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestANDAbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.X = 0x01
	cpu.Memory[0] = 0x3D
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF02] = 0x04
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestANDAbsoluteXWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.X = 0x01
	cpu.Memory[0] = 0x3D
	cpu.Memory[1] = 0xFF
	cpu.Memory[2] = 0x00
	cpu.Memory[0x0100] = 0x04
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestANDAbsoluteY(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.Y = 0x01
	cpu.Memory[0] = 0x39
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF02] = 0x04
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestANDAbsoluteYWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0x39
	cpu.Memory[1] = 0xFF
	cpu.Memory[2] = 0x00
	cpu.Memory[0x0100] = 0x05
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestANDAIndexedIndirect(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x01
	cpu.Memory[0] = 0x21
	cpu.Memory[1] = 0xFE
	cpu.Memory[9] = 0x05
	cpu.Memory[0xFF] = 0x09
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestANDIndexedIndirectWithOverflow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x0B
	cpu.Memory[0] = 0x21
	cpu.Memory[1] = 0xFF
	cpu.Memory[9] = 0x05
	cpu.Memory[0x0A] = 0x09
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestANDIndirectIndexed(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0x31
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0x05
	cpu.Memory[6] = 0x05
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestANDAIndirectIndexedWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0x31
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0xFF
	cpu.Memory[0x100] = 0x05
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestADCImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.Memory[0] = 0x69
	cpu.Memory[1] = 0x80
	cpu.Exec()

	if cpu.A != 0x81 {
		t.Error("ADC failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 2 {
		t.Error("ADC did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}

	if cpu.NFlag != true {
		t.Error("ADC set negative flag incorrectly")
	}

	if cpu.CFlag != false {
		t.Error("ADC set carry flag incorrectly")
	}

	if cpu.VFlag != false {
		t.Error("ADC set overflow flag incorrectly")
	}

	if cpu.ZFlag != false {
		t.Error("ADC set zero flag incorrectly")
	}
}

func TestADCImmediateWithCarryIn(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x00
	cpu.CFlag = true
	cpu.Memory[0] = 0x69
	cpu.Memory[1] = 0x00
	cpu.Exec()

	if cpu.A != 0x01 {
		t.Error("ADC failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.CFlag != false {
		t.Error("ADC failed to clear the carry flag")
	}
}

func TestADCImmediateWithCarryOut(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0xFF
	cpu.Memory[0] = 0x69
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.A != 0x00 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.CFlag != true {
		t.Error("set carry flag incorrectly")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

}

func TestADCImmediateWithOverflow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x7F
	cpu.Memory[0] = 0x69
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.A != 0x80 {
		t.Error("ADC failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.VFlag != true {
		t.Error("ADC set overflow flag incorrectly")
	}
}

func TestADCZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.Memory[0] = 0x65
	cpu.Memory[1] = 0x09
	cpu.Memory[9] = 0x80
	cpu.Exec()

	if cpu.A != 0x81 {
		t.Error("ADC failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 3 {
		t.Error("ADC did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestADCZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.X = 0x01
	cpu.Memory[0] = 0x75
	cpu.Memory[1] = 0x08
	cpu.Memory[9] = 0x03
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("ADC failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("ADC did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestADCAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x03
	cpu.Memory[0] = 0x6D
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x03
	cpu.Exec()

	if cpu.A != 0x06 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestADCAbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.X = 0x01
	cpu.Memory[0] = 0x7D
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF02] = 0x04
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestADCAbsoluteXWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.X = 0x01
	cpu.Memory[0] = 0x7D
	cpu.Memory[1] = 0xFF
	cpu.Memory[2] = 0x00
	cpu.Memory[0x0100] = 0x04
	cpu.Exec()

	if cpu.A != 0x08 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestADCAbsoluteY(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.Y = 0x01
	cpu.Memory[0] = 0x79
	cpu.Memory[1] = 0xFE
	cpu.Memory[2] = 0x00
	cpu.Memory[0xFF] = 0x04
	cpu.Exec()

	if cpu.A != 0x08 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestADCAbsoluteYWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.Y = 0x01
	cpu.Memory[0] = 0x79
	cpu.Memory[1] = 0xFF
	cpu.Memory[2] = 0x00
	cpu.Memory[0x100] = 0x04
	cpu.Exec()

	if cpu.A != 0x08 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestADCIndexedIndirect(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x01
	cpu.Memory[0] = 0x61
	cpu.Memory[1] = 0xFE
	cpu.Memory[9] = 0x05
	cpu.Memory[0xFF] = 0x09
	cpu.Exec()

	if cpu.A != 0x0A {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestADCIndirectIndexed(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0x71
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0x05
	cpu.Memory[6] = 0x0A
	cpu.Exec()

	if cpu.A != 0x0F {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestADCIndirectIndexedWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0x71
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0xFF
	cpu.Memory[0x100] = 0x05
	cpu.Exec()

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestASLAccumulator(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.Memory[0] = 0x0A
	cpu.Exec()

	if cpu.A != 0x08 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles count")
	}

	if cpu.PC != 1 {
		t.Error("did not correctly set PC")
	}

	if cpu.CFlag != false {
		t.Error("did not correctly set carry flag")
	}

	if cpu.ZFlag != false {
		t.Error("did not correctly set zero flag")
	}

	if cpu.NFlag != false {
		t.Error("did not correctly set negative flag")
	}
}

func TestASLAccumulatorWithCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0xC0
	cpu.Memory[0] = 0x0A
	cpu.Exec()

	if cpu.A != 0x80 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.CFlag != true {
		t.Error("did not correctly set carry flag")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set negative flag")
	}

}

func TestASLZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x06
	cpu.Memory[1] = 0x09
	cpu.Memory[9] = 0x02
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("failed to give correct Accumulator value, got", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestASLZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.Memory[0] = 0x16
	cpu.Memory[1] = 0x08
	cpu.Memory[9] = 0x04
	cpu.Exec()

	if cpu.A != 0x08 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestASLAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x0E
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x07
	cpu.Exec()

	if cpu.A != 0x0E {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestASLAbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.Memory[0] = 0x1E
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF02] = 0x04
	cpu.Exec()

	if cpu.A != 0x08 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 7 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestBCC(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = false
	cpu.Memory[0] = 0x90
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x12 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBCCNoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = true
	cpu.Memory[0] = 0x90
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x2 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBCCNegativeRelativeAddress(t *testing.T) {
	cpu := NewCPU()
	cpu.PC = 0x01
	cpu.CFlag = false
	cpu.Memory[1] = 0x90
	cpu.Memory[2] = 0xFF // -1 in two's complement
	cpu.Exec()

	// the reason the expected value is 2 is that the offset specified in a branch instruction
	// is taken after the two bytes for the instruction and its operand are accounted for.
	// So here we start the PC at 1, add 2 to execute the BCC instruction, then branch back
	// 1 to land at 2.
	if cpu.PC != 2 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBCCWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = false
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0x90
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBCS(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = true
	cpu.Memory[0] = 0xB0
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x12 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBCSNoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = false
	cpu.Memory[0] = 0xB0
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x2 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBCSNegativeRelativeAddress(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = true
	cpu.PC = 0x10A
	cpu.Memory[0x10A] = 0xB0
	cpu.Memory[0x10B] = 0xF4 // -10
	cpu.Exec()

	if cpu.PC != 0x100 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBCSWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = true
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0xB0
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBEQ(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.Memory[0] = 0xF0
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x12 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBEQNoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.Memory[0] = 0xF0
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x2 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBEQWithNegativeRelativeAddress(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.PC = 0x10A
	cpu.Memory[0x10A] = 0xF0
	cpu.Memory[0x10B] = 0xF4 // -10
	cpu.Exec()

	if cpu.PC != 0x100 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBEQWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0xF0
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBITZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.VFlag = false
	cpu.A = 0x00
	cpu.PC = 0x01
	cpu.Memory[1] = 0x24
	cpu.Memory[2] = 0x03
	cpu.Memory[3] = 0xFF
	cpu.Exec()

	if cpu.ZFlag != true {
		t.Error("did not correctly set zero flag")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set negative flag")
	}

	if cpu.VFlag != true {
		t.Error("did not correctly set overflow flag")
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 0x3 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}
}

func TestBITZeroPageWithNonZeroValue(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.VFlag = false
	cpu.A = 0x01
	cpu.PC = 0x00
	cpu.Memory[0] = 0x24
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0x01
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("did not correctly set zero flag")
	}

	if cpu.NFlag != false {
		t.Error("did not correctly set negative flag")
	}

	if cpu.VFlag != false {
		t.Error("did not correctly set overflow flag")
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 0x2 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}
}

func TestBITAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.A = 0x00
	cpu.PC = 0x01
	cpu.Memory[1] = 0x2C
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x01
	cpu.Exec()

	if cpu.ZFlag != true {
		t.Error("did not correctly set zero flag")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 0x4 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}
}

func TestBMI(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.Memory[0] = 0x30
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x12 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBMINoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = false
	cpu.Memory[0] = 0x30
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x02 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBMIWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0x30
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBNE(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.Memory[0] = 0xD0
	cpu.Memory[1] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x11 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBNENoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.Memory[0] = 0xD0
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x02 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBNEWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0xD0
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x102 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBNEWithNegativeRelativeAddress(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.PC = 0x10A
	cpu.Memory[0x10A] = 0xD0
	cpu.Memory[0x10B] = 0xF4 // -10
	cpu.Exec()

	if cpu.PC != 0x100 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBPL(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = false
	cpu.Memory[0] = 0x10
	cpu.Memory[1] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x11 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBPLNoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.Memory[0] = 0x10
	cpu.Memory[1] = 0x10
	cpu.Exec()

	if cpu.PC != 0x02 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBPLWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = false
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0x10
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x102 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBPLWithNegativeRelativeAddress(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = false
	cpu.PC = 0x10A
	cpu.Memory[0x10A] = 0x10
	cpu.Memory[0x10B] = 0xF4 // -10
	cpu.Exec()

	if cpu.PC != 0x100 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVC(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = false
	cpu.Memory[0] = 0x50
	cpu.Memory[1] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x11 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVCNoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = true
	cpu.Memory[0] = 0x50
	cpu.Memory[1] = 0x50
	cpu.Exec()

	if cpu.PC != 0x02 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVCWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = false
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0x50
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x102 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVCWithNegativeRelativeAddress(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = false
	cpu.PC = 0x50A
	cpu.Memory[0x50A] = 0x10
	cpu.Memory[0x50B] = 0xF4 // -10
	cpu.Exec()

	if cpu.PC != 0x500 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVS(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = true
	cpu.Memory[0] = 0x70
	cpu.Memory[1] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x11 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVSNoBranch(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = false
	cpu.Memory[0] = 0x70
	cpu.Memory[1] = 0x50
	cpu.Exec()

	if cpu.PC != 0x02 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVSWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = true
	cpu.PC = 0xF1
	cpu.Memory[0xF1] = 0x70
	cpu.Memory[0xF2] = 0x0F
	cpu.Exec()

	if cpu.PC != 0x102 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBVSWithNegativeRelativeAddress(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = true
	cpu.PC = 0x50A
	cpu.Memory[0x50A] = 0x70
	cpu.Memory[0x50B] = 0xF4 // -10
	cpu.Exec()

	if cpu.PC != 0x500 {
		t.Error("failed to correctly set PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestBRK(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0xFFFE] = 0x0A
	cpu.Memory[0xFFFF] = 0x02
	cpu.PC = 0x00
	cpu.Exec()

	if cpu.PC != 0x020A {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.BFlag != true {
		t.Error("did not correctly set Break flag")
	}

	if cpu.Cycles != 7 {
		t.Error("did not correctly update cycles, got", cpu.Cycles)
	}
}

func TestCLC(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = true
	cpu.PC = 0x01
	cpu.Memory[1] = 0x18
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("did not clear C Flag")
	}
}

func TestCLCFlagUnset(t *testing.T) {
	cpu := NewCPU()
	cpu.CFlag = false
	cpu.PC = 0x01
	cpu.Memory[1] = 0x18
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("did not clear C Flag")
	}
}

func TestCLI(t *testing.T) {
	cpu := NewCPU()
	cpu.IFlag = true
	cpu.PC = 0x01
	cpu.Memory[1] = 0x58
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("did not clear I Flag")
	}
}

func TestCLIFlagUnset(t *testing.T) {
	cpu := NewCPU()
	cpu.IFlag = false
	cpu.PC = 0x01
	cpu.Memory[1] = 0x58
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("did not clear I Flag")
	}
}

func TestCLV(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = true
	cpu.PC = 0x01
	cpu.Memory[1] = 0xB8
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("did not clear V Flag")
	}
}

func TestCLVFlagUnset(t *testing.T) {
	cpu := NewCPU()
	cpu.VFlag = false
	cpu.PC = 0x01
	cpu.Memory[1] = 0xB8
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("did not clear V Flag")
	}
}

func TestCMPImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.CFlag = false
	cpu.A = 0x02
	cpu.PC = 0x01
	cpu.Memory[1] = 0xC9
	cpu.Memory[2] = 0x01
	cpu.Exec()

	if cpu.CFlag != true {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != false {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != false {
		t.Error("incorrectly set Negative flag")
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly update cycles, got", cpu.Cycles)
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPImmediateWithEqualInput(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.CFlag = false
	cpu.A = 0x01
	cpu.PC = 0x01
	cpu.Memory[1] = 0xC9
	cpu.Memory[2] = 0x01
	cpu.Exec()

	if cpu.CFlag != true {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != true {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != false {
		t.Error("incorrectly set Negative flag")
	}
}

func TestCMPImmediateWithOperandGreater(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.CFlag = false
	cpu.A = 0x01
	cpu.PC = 0x01
	cpu.Memory[1] = 0xC9
	cpu.Memory[2] = 0x02
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != false {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != true {
		t.Error("incorrectly set Negative flag")
	}

}

func TestCMPZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.CFlag = false
	cpu.A = 0x01
	cpu.PC = 0x01
	cpu.Memory[1] = 0xC5
	cpu.Memory[2] = 0x09
	cpu.Memory[9] = 0x02
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly update cycles, got", cpu.Cycles)
	}
}

func TestCMPZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x03
	cpu.X = 0x01
	cpu.Memory[0] = 0xD5
	cpu.Memory[1] = 0x08
	cpu.Memory[9] = 0x03
	cpu.Exec()

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Memory[0] = 0xCD
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x03
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.CFlag != true {
		t.Error("set carry flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPAbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x05
	cpu.Memory[0] = 0xDD
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF06] = 0x06
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.CFlag != false {
		t.Error("set carry flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPAbsoluteXWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0xFF
	cpu.Memory[0] = 0xDD
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFE
	cpu.Memory[0xFF00] = 0x0F
	cpu.Exec()

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPAbsoluteY(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x05
	cpu.Memory[0] = 0xD9
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF06] = 0x06
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.CFlag != false {
		t.Error("set carry flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPAbsoluteYWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0xFF
	cpu.Memory[0] = 0xD9
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFE
	cpu.Memory[0xFF00] = 0x0F
	cpu.Exec()

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPIndexedIndirect(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x01
	cpu.Memory[0] = 0xC1
	cpu.Memory[1] = 0xFE
	cpu.Memory[9] = 0x06
	cpu.Memory[0xFF] = 0x09
	cpu.Exec()

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPIndexedIndirectWithOverflow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x0B
	cpu.Memory[0] = 0xC1
	cpu.Memory[1] = 0xFF
	cpu.Memory[9] = 0x06
	cpu.Memory[0x0A] = 0x09
	cpu.Exec()

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPIndirectIndexed(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0xD1
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0x05
	cpu.Memory[6] = 0x06
	cpu.Exec()

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestCMPIndirectIndexedWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0xD1
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0xFF
	cpu.Memory[0x100] = 0x06
	cpu.Exec()

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestCPXImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x02
	cpu.Memory[0] = 0xE0
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.CFlag != true {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != false {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != false {
		t.Error("incorrectly set Negative flag")
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly update cycles, got", cpu.Cycles)
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestCPXImmediateWithEqualInput(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.PC = 0x01
	cpu.Memory[1] = 0xE0
	cpu.Memory[2] = 0x01
	cpu.Exec()

	if cpu.CFlag != true {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != true {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != false {
		t.Error("incorrectly set Negative flag")
	}
}

func TestCPXImmediateWithOperandGreater(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.CFlag = false
	cpu.A = 0x01
	cpu.PC = 0x01
	cpu.Memory[1] = 0xE0
	cpu.Memory[2] = 0x02
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != false {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != true {
		t.Error("incorrectly set Negative flag")
	}

}

func TestCPXZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.Memory[0] = 0xE4
	cpu.Memory[1] = 0x09
	cpu.Memory[9] = 0x02
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly update cycles, got", cpu.Cycles)
	}
}

func TestCPXAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x05
	cpu.Memory[0] = 0xEC
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x03
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.CFlag != true {
		t.Error("set carry flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestCPYImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x02
	cpu.Memory[0] = 0xE0
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.CFlag != true {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != false {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != false {
		t.Error("incorrectly set Negative flag")
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly update cycles, got", cpu.Cycles)
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestCPYImmediateWithEqualInput(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.PC = 0x01
	cpu.Memory[1] = 0xE0
	cpu.Memory[2] = 0x01
	cpu.Exec()

	if cpu.CFlag != true {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != true {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != false {
		t.Error("incorrectly set Negative flag")
	}
}

func TestCPYImmediateWithOperandGreater(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = false
	cpu.NFlag = false
	cpu.CFlag = false
	cpu.A = 0x01
	cpu.PC = 0x01
	cpu.Memory[1] = 0xE0
	cpu.Memory[2] = 0x02
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("incorrectly set Carry flag")
	}

	if cpu.ZFlag != false {
		t.Error("incorrectly set Zero flag")
	}

	if cpu.NFlag != true {
		t.Error("incorrectly set Negative flag")
	}

}

func TestCPYZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.Memory[0] = 0xE4
	cpu.Memory[1] = 0x09
	cpu.Memory[9] = 0x02
	cpu.Exec()

	if cpu.CFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly update cycles, got", cpu.Cycles)
	}
}

func TestCPYAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Y = 0x05
	cpu.Memory[0] = 0xCC
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x03
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.CFlag != true {
		t.Error("set carry flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestDECZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.NFlag = true
	cpu.Memory[0] = 0xC6
	cpu.Memory[1] = 0x0A
	cpu.Memory[10] = 0x02
	cpu.Exec()

	if cpu.Memory[10] != 0x01 {
		t.Error("failed to update memory value correclty, got", cpu.Memory[10])
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestDECZeroPageWithNegativeResult(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.NFlag = true
	cpu.Memory[0] = 0xC6
	cpu.Memory[1] = 0x0A
	cpu.Memory[10] = 0x00
	cpu.Exec()

	if cpu.Memory[10] != 0xFF {
		t.Error("failed to update memory value correclty, got", cpu.Memory[10])
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestDECZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.X = 0x1
	cpu.Memory[0] = 0xD6
	cpu.Memory[1] = 0x09
	cpu.Memory[10] = 0x01
	cpu.Exec()

	if cpu.Memory[10] != 0x00 {
		t.Error("failed to update memory value correclty, got", cpu.Memory[10])
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestDECAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.Memory[0] = 0xCE
	cpu.Memory[1] = 0x09
	cpu.Memory[2] = 0x09
	cpu.Memory[0x0909] = 0x02
	cpu.Exec()

	if cpu.Memory[0x0909] != 0x01 {
		t.Error("failed to update memory value correclty, got", cpu.Memory[10])
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestDECAbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.Memory[0] = 0xDE
	cpu.Memory[1] = 0x09
	cpu.Memory[2] = 0x09
	cpu.Memory[0x090A] = 0x02
	cpu.Exec()

	if cpu.Memory[0x090A] != 0x01 {
		t.Error("failed to update memory value correclty, got", cpu.Memory[10])
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 7 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestDEX(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x02
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.Memory[0] = 0xCA
	cpu.Exec()

	if cpu.X != 0x01 {
		t.Error("failed to update X register correctly, got", cpu.X)
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC")
	}
}
