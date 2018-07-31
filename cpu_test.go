package main

import (
	"testing"
)

func TestANDImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.Memory[0] = 0x29
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.A != 0x01 {
		t.Error("AND immediate (opcode 0x29) failed to give correct Accumulator value")
	}

	if cpu.ZFlag != false {
		t.Error("AND set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("AND set negative flag incorrectly")
	}

	if cpu.Cycles != 2 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("AND did not correctly update PC")
	}
}

func TestANDImmediateZero(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x01
	cpu.Memory[0] = 0x29
	cpu.Memory[1] = 0x00
	cpu.Exec()

	if cpu.A != 0x00 {
		t.Error("AND immediate (opcode 0x29) failed to give correct Accumulator value")
	}

	if cpu.ZFlag != true {
		t.Error("AND set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("AND set negative flag incorrectly")
	}

	if cpu.Cycles != 2 {
		t.Error("AND did not correctly set cycles flag")
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
		t.Error("AND zero page (opcode 0x25) failed to give correct Accumulator value")
	}

	if cpu.ZFlag != false {
		t.Error("AND set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("AND set negative flag incorrectly")
	}

	if cpu.Cycles != 3 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("AND did not correctly update PC")
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
		t.Error("AND zero page X (opcode 0x35) failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.ZFlag != false {
		t.Error("AND set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("AND set negative flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("AND did not correctly update PC")
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
		t.Error("AND zero page X (opcode 0x35) failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.ZFlag != false {
		t.Error("AND set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("AND set negative flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("AND did not correctly update PC")
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
		t.Error("AND zero page X (opcode 0x35) failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.ZFlag != false {
		t.Error("AND set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("AND set negative flag incorrectly")
	}

	if cpu.Cycles != 4 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("AND did not correctly update PC")
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
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("AND did not correctly update PC")
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
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("AND did not correctly set cycles flag")
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
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("AND did not correctly update PC")
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
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("AND did not correctly set cycles flag")
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
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("AND did not correctly update PC")
	}
}

func TestANDAIndexedIndirectWithOverflow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x0B
	cpu.Memory[0] = 0x21
	cpu.Memory[1] = 0xFF
	cpu.Memory[9] = 0x05
	cpu.Memory[0x0A] = 0x09
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("AND did not correctly set cycles flag")
	}
}

func TestANDAIndirectIndexed(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0x31
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0x05
	cpu.Memory[6] = 0x05
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("AND did not correctly update PC")
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
		t.Error("AND failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("AND did not correctly set cycles flag")
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

func TestADCImmediateWithCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0xFF
	cpu.Memory[0] = 0x69
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.A != 0x00 {
		t.Error("ADC failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 2 {
		t.Error("AND did not correctly set cycles flag")
	}

	if cpu.CFlag != true {
		t.Error("ADC set carry flag incorrectly")
	}

	if cpu.ZFlag != true {
		t.Error("ADC set zero flag incorrectly")
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
}
