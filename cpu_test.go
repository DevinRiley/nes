package main

import (
	"testing"
)

type CpuTestHarness struct {
	Cpu    *CPU
	Opcode byte
}

func (h *CpuTestHarness) SetupAccumulator() {
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.A = 0x08
}

func (h *CpuTestHarness) SetupZeroPage() {
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x17
	h.Cpu.Memory[0x17] = 0x07
}

func (h *CpuTestHarness) SetupZeroPageX() {
	h.Cpu.X = 0x0F
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x80
	h.Cpu.Memory[0x8F] = 0x07
}

func (h *CpuTestHarness) SetupZeroPageY() {
	h.Cpu.Y = 0x0F
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x80
	h.Cpu.Memory[0x8F] = 0x07
}

func (h *CpuTestHarness) SetupAbsolute() {
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x80
	h.Cpu.Memory[2] = 0x80
	h.Cpu.Memory[0x8080] = 0x07
}

func (h *CpuTestHarness) SetupAbsoluteX() {
	h.Cpu.X = 0x01
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x80
	h.Cpu.Memory[2] = 0xFF
	h.Cpu.Memory[0xFF81] = 0x07
}

func (h *CpuTestHarness) SetupAbsoluteXPageCross() {
	h.Cpu.X = 0x01
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0xFF
	h.Cpu.Memory[2] = 0xF0
	h.Cpu.Memory[0xF100] = 0x07
}

func (h *CpuTestHarness) SetupAbsoluteY() {
	h.Cpu.Y = 0x01
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x80
	h.Cpu.Memory[2] = 0xFF
	h.Cpu.Memory[0xFF81] = 0x07
}

func (h *CpuTestHarness) SetupAbsoluteYPageCross() {
	h.Cpu.Y = 0x01
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0xFF
	h.Cpu.Memory[2] = 0xF0
	h.Cpu.Memory[0xF100] = 0x07
}

func (h *CpuTestHarness) SetupIndexedIndirect() {
	h.Cpu.X = 0x01
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0xFE
	h.Cpu.Memory[9] = 0x07
	h.Cpu.Memory[0xFF] = 0x09
}

func (h *CpuTestHarness) SetupIndirectIndexed() {
	h.Cpu.Y = 0x01
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x02
	h.Cpu.Memory[2] = 0x05
	h.Cpu.Memory[6] = 0x07
}

func (h *CpuTestHarness) SetupIndirectIndexedPageCross() {
	h.Cpu.Y = 0x01
	h.Cpu.Memory[0] = h.Opcode
	h.Cpu.Memory[1] = 0x02
	h.Cpu.Memory[2] = 0xFF
	h.Cpu.Memory[0x100] = 0x07
}

func (h *CpuTestHarness) Run() {
	h.Cpu.Exec()
}

func NewCpuTestHarness() *CpuTestHarness {
	return &CpuTestHarness{Cpu: NewCPU()}
}

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

func TestDEY(t *testing.T) {
	cpu := NewCPU()
	cpu.Y = 0x02
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.Memory[0] = 0x88
	cpu.Exec()

	if cpu.Y != 0x01 {
		t.Error("failed to update Y register correctly, got", cpu.Y)
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

func TestEORImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.NFlag = true
	cpu.A = 0x06
	cpu.Memory[0] = 0x49
	cpu.Memory[1] = 0x05
	cpu.Exec()

	if cpu.A != 0x03 {
		t.Error("did not correclty set accumulator, got", cpu.A)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}

	if cpu.ZFlag != false {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != false {
		t.Error("set negative flag incorrectly")
	}
}

func TestEORZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x00
	cpu.Memory[0] = 0x45
	cpu.Memory[1] = 0x05
	cpu.Memory[5] = 0x00
	cpu.Exec()

	if cpu.A != 0x00 {
		t.Error("did not correclty set accumulator, got", cpu.A)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}
}

func TestEORZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x80
	cpu.X = 0x01
	cpu.Memory[0] = 0x55
	cpu.Memory[1] = 0x04
	cpu.Memory[5] = 0x01
	cpu.Exec()

	if cpu.A != 0x81 {
		t.Error("did not correclty set accumulator, got", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}

	if cpu.NFlag != true {
		t.Error("set zero flag incorrectly")
	}
}

func TestEORAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x03
	cpu.Memory[0] = 0x4D
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF01] = 0x00
	cpu.Exec()

	if cpu.A != 0x03 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestEORAbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.X = 0x01
	cpu.Memory[0] = 0x5D
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF02] = 0x01
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

func TestEORAbsoluteXWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.X = 0x01
	cpu.Memory[0] = 0x5D
	cpu.Memory[1] = 0xFF
	cpu.Memory[2] = 0x00
	cpu.Memory[0x0100] = 0x01
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestEORAbsoluteY(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.Y = 0x01
	cpu.Memory[0] = 0x59
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0xFF
	cpu.Memory[0xFF02] = 0x01
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

func TestEORAbsoluteYWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.Y = 0x01
	cpu.Memory[0] = 0x59
	cpu.Memory[1] = 0xFF
	cpu.Memory[2] = 0x00
	cpu.Memory[0x0100] = 0x01
	cpu.Exec()

	if cpu.A != 0x05 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestEORIndexedIndirect(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.X = 0x01
	cpu.Memory[0] = 0x41
	cpu.Memory[1] = 0xFE
	cpu.Memory[9] = 0x01
	cpu.Memory[0xFF] = 0x09
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestEORIndirectIndexed(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x04
	cpu.Y = 0x01
	cpu.Memory[0] = 0x51
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0x05
	cpu.Memory[6] = 0x05
	cpu.Exec()

	if cpu.A != 0x01 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestEORIndirectIndexedWithPageCross(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x05
	cpu.Y = 0x01
	cpu.Memory[0] = 0x51
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0xFF
	cpu.Memory[0x100] = 0x01
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("failed to give correct Accumulator value, gave", cpu.A)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestINCZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.NFlag = true
	cpu.Memory[0] = 0xE6
	cpu.Memory[1] = 0x0A
	cpu.Memory[10] = 0x02
	cpu.Exec()

	if cpu.Memory[10] != 0x03 {
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

func TestINCZeroPageWithNegativeResult(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.Memory[0] = 0xE6
	cpu.Memory[1] = 0x0A
	cpu.Memory[10] = 0x7F
	cpu.Exec()

	if cpu.Memory[10] != 0x80 {
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

func TestINCZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.X = 0x1
	cpu.Memory[0] = 0xF6
	cpu.Memory[1] = 0x09
	cpu.Memory[10] = 0x01
	cpu.Exec()

	if cpu.Memory[10] != 0x02 {
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

	if cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}
}

func TestINCAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.ZFlag = false
	cpu.Memory[0] = 0xEE
	cpu.Memory[1] = 0x09
	cpu.Memory[2] = 0x09
	cpu.Memory[0x0909] = 0xFF
	cpu.Exec()

	if cpu.Memory[0x0909] != 0x00 {
		t.Error("failed to update memory value correclty, got", cpu.Memory[0x0909])
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

	if cpu.PC != 3 {
		t.Error("did not correctly update PC")
	}
}

func TestINCAbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x01
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.Memory[0] = 0xFE
	cpu.Memory[1] = 0x09
	cpu.Memory[2] = 0x09
	cpu.Memory[0x090A] = 0x02
	cpu.Exec()

	if cpu.Memory[0x090A] != 0x03 {
		t.Error("failed to update memory value correclty, got", cpu.Memory[0x090A])
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

func TestINX(t *testing.T) {
	cpu := NewCPU()
	cpu.X = 0x02
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.Memory[0] = 0xE8
	cpu.Exec()

	if cpu.X != 0x03 {
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

func TestINY(t *testing.T) {
	cpu := NewCPU()
	cpu.Y = 0x02
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.Memory[0] = 0xC8
	cpu.Exec()

	if cpu.Y != 0x03 {
		t.Error("failed to update Y register correctly, got", cpu.Y)
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

func TestJMPAbsolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x4C
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0x01
	cpu.Exec()

	if cpu.PC != 0x0101 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestJMPIndirect(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x6C
	cpu.Memory[1] = 0x01
	cpu.Memory[2] = 0x02
	cpu.Memory[0x0201] = 0x07
	cpu.Memory[0x0202] = 0x01
	cpu.Exec()

	if cpu.PC != 0x0107 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestJMPIndirectWithHardwareBug(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x6C
	cpu.Memory[1] = 0xFF
	cpu.Memory[2] = 0x02
	cpu.Memory[0x02FF] = 0x07
	cpu.Memory[0x0200] = 0x01
	cpu.Exec()

	if cpu.PC != 0x0107 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestJSR(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x20
	cpu.Memory[1] = 0x02
	cpu.Memory[2] = 0x02
	cpu.Exec()

	if cpu.PC != 0x0202 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	returnAddress := cpu.stackPop16()
	if returnAddress != 0x02 {
		t.Error("did not push the expected return address onto the stack, got:", returnAddress)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA9
	cpu.Memory[1] = 0x07
	cpu.Exec()

	if cpu.A != 0x07 {
		t.Error("did not correctly load accumulator, got", cpu.A)
	}

	if cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAImmediateSetsZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA9
	cpu.Memory[1] = 0x00
	cpu.Exec()

	if cpu.ZFlag != true {
		t.Error("did not correctly set ZFlag")
	}

	if cpu.NFlag != false {
		t.Error("did not correctly set NFlag")
	}
}

func TestLDAImmediateSetsNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA9
	cpu.Memory[1] = 0x80
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}
}

func TestLDAZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xA5
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xB5
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xAD
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xBD
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAAbsoluteXWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xBD
	harness.SetupAbsoluteXPageCross()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAAbsoluteY(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xB9
	harness.SetupAbsoluteY()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAAbsoluteYWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xB9
	harness.SetupAbsoluteYPageCross()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAIndexedIndirect(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xA1
	harness.SetupIndexedIndirect()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("failed to give correct Accumulator value, gave", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAIndirectIndexed(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xB1
	harness.SetupIndirectIndexed()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("failed to give correct Accumulator value, gave", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDAIndirectIndexedWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xB1
	harness.SetupIndirectIndexedPageCross()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("failed to give correct Accumulator value, gave", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC")
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDXImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA2
	cpu.Memory[1] = 0x07
	cpu.Exec()

	if cpu.X != 0x07 {
		t.Error("did not correctly load accumulator, got", cpu.X)
	}

	if cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDXImmediateSetsZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA2
	cpu.Memory[1] = 0x00
	cpu.Exec()

	if cpu.ZFlag != true {
		t.Error("did not correctly set ZFlag")
	}

	if cpu.NFlag != false {
		t.Error("did not correctly set NFlag")
	}
}

func TestLDXImmediateSetsNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA2
	cpu.Memory[1] = 0x80
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}
}

func TestLDXZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xA6
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.X != 0x07 {
		t.Error("did not correctly set X register, got: ", harness.Cpu.X)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDXZeroPageY(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xB6
	harness.SetupZeroPageY()
	harness.Run()

	if harness.Cpu.X != 0x07 {
		t.Error("did not correctly set X register, got: ", harness.Cpu.X)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDXAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xAE
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.X != 0x07 {
		t.Error("did not correctly set X register, got: ", harness.Cpu.X)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDXAbsoluteY(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xBE
	harness.SetupAbsoluteY()
	harness.Run()

	if harness.Cpu.X != 0x07 {
		t.Error("did not correctly set X register, got: ", harness.Cpu.X)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDXAbsoluteYWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xBE
	harness.SetupAbsoluteYPageCross()
	harness.Run()

	if harness.Cpu.X != 0x07 {
		t.Error("did not correctly set X register, got: ", harness.Cpu.X)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDYImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA0
	cpu.Memory[1] = 0x07
	cpu.Exec()

	if cpu.Y != 0x07 {
		t.Error("did not correctly load Y register, got", cpu.Y)
	}

	if cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDYImmediateSetsZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA0
	cpu.Memory[1] = 0x00
	cpu.Exec()

	if cpu.ZFlag != true {
		t.Error("did not correctly set ZFlag")
	}

	if cpu.NFlag != false {
		t.Error("did not correctly set NFlag")
	}
}

func TestLDYImmediateSetsNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xA0
	cpu.Memory[1] = 0x80
	cpu.Exec()

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}
}

func TestLDYZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xA4
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.Y != 0x07 {
		t.Error("did not correctly set Y register, got: ", harness.Cpu.Y)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDYZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xB4
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.Y != 0x07 {
		t.Error("did not correctly set Y register, got: ", harness.Cpu.Y)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDYAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xAC
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.Y != 0x07 {
		t.Error("did not correctly set Y register, got: ", harness.Cpu.Y)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDYAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xBC
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.Y != 0x07 {
		t.Error("did not correctly set Y register, got: ", harness.Cpu.Y)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLDYAbsoluteXWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0xBC
	harness.SetupAbsoluteXPageCross()
	harness.Run()

	if harness.Cpu.Y != 0x07 {
		t.Error("did not correctly set Y register, got: ", harness.Cpu.Y)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLSRAccumulator(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x08
	cpu.CFlag = true
	cpu.NFlag = true
	cpu.Memory[0] = 0x4A
	cpu.Exec()

	if cpu.A != 0x04 {
		t.Error("Did not correctly set the accumulator, got: ", cpu.A)
	}

	if cpu.CFlag != false {
		t.Error("Did not correctly set carry flag")
	}

	if cpu.NFlag != false {
		t.Error("Did not correctly set negative flag")
	}
}

func TestLSRAccumulatorWithCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x80
	cpu.Memory[0] = 0x4A
	cpu.Exec()

	if cpu.A != 0x40 {
		t.Error("Did not correctly set the accumulator, got: ", cpu.A)
	}

	if cpu.CFlag != true {
		t.Error("Did not correctly set carry flag")
	}
}

func TestLSRZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x46
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.Memory[0x17] != 0x03 {
		t.Error("did not correctly set Memory, got: ", harness.Cpu.Memory[0x17])
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLSRZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x56
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.Memory[0x8F] != 0x03 {
		t.Error("did not correctly set Memory, got: ", harness.Cpu.Memory[0x8F])
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLSRAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x4E
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.Memory[0x8080] != 0x03 {
		t.Error("did not correctly set Memory, got: ", harness.Cpu.Memory[0x8080])
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestLSRAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x5E
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.Memory[0xFF81] != 0x03 {
		t.Error("did not correctly set Memory, got: ", harness.Cpu.Memory[0xFF81])
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 7 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestNOP(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xEA
	cpu.Exec()

	if cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAImmediate(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.NFlag = true
	cpu.A = 0x03
	cpu.Memory[0] = 0x09
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.A != 0x03 {
		t.Error("failed to give correct Accumulator value, got: ", cpu.A)
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

func TestORAZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x05
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x15
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x0D
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x1D
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAAbsoluteXWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x1D
	harness.SetupAbsoluteXPageCross()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAAbsoluteY(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x19
	harness.SetupAbsoluteY()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAAbsoluteYWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x19
	harness.SetupAbsoluteYPageCross()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAIndexedIndirect(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x01
	harness.SetupIndexedIndirect()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAIndirectIndexed(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x11
	harness.SetupIndirectIndexed()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestORAIndirectIndexedWithPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x11
	harness.SetupIndirectIndexedPageCross()
	harness.Run()

	if harness.Cpu.A != 0x07 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestPHA(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0xD3
	cpu.Memory[0] = 0x48
	cpu.Exec()

	stack := cpu.stackPop()

	if stack != 0xD3 {
		t.Error("did not correctly push A onto stack, got: ", stack)
	}

	if cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestPHP(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.VFlag = true
	cpu.UFlag = false
	cpu.BFlag = false
	cpu.DFlag = true
	cpu.IFlag = true
	cpu.ZFlag = true
	cpu.CFlag = true
	cpu.Memory[0] = 0x08
	cpu.Exec()

	stack := cpu.stackPop()

	if stack&0x30 != 0x30 {
		t.Error("did not set bits 5 and 4")
	}

	if stack != 0xFF {
		t.Error("did not correctly push flags onto stack, got: ", stack)
	}

	if cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 3 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestPLA(t *testing.T) {
	cpu := NewCPU()
	cpu.NFlag = true
	cpu.ZFlag = true
	cpu.stackPush(0x07)
	cpu.Memory[0] = 0x68
	cpu.Exec()

	if cpu.A != 0x07 {
		t.Error("did not correctly push flags onto stack, got: ", cpu.A)
	}

	if cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestPLP(t *testing.T) {
	cpu := NewCPU()
	cpu.stackPush(0xFF)
	cpu.Memory[0] = 0x28
	cpu.Exec()

	if cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 4 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}
	if cpu.VFlag != true {
		t.Error("set overflow flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}
}

func TestROL(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x2A
	harness.SetupAccumulator()
	harness.Cpu.CFlag = true
	harness.Cpu.NFlag = true
	harness.Run()

	if harness.Cpu.A != 0x11 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.CFlag != false {
		t.Error("did not correctly set CFlag, got: ", harness.Cpu.CFlag)
	}

	if harness.Cpu.NFlag != false {
		t.Error("did not correctly set NFlag, got: ", harness.Cpu.NFlag)
	}

	if harness.Cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}
}

func TestROLZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x26
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.Memory[0x17] != 0x0E {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0x17])
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles")
	}
}

func TestROLZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x36
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.Memory[0x8F] != 0x0E {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0x8F])
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestROLAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x2E
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.Memory[0x8080] != 0x0E {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0x8080])
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestROLAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x3E
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.Memory[0xFF81] != 0x0E {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0xFF81])
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 7 {
		t.Error("did not correctly set cycles")
	}
}

func TestROR(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x6A
	harness.SetupAccumulator()
	harness.Cpu.CFlag = false
	harness.Cpu.NFlag = true
	harness.Run()

	if harness.Cpu.A != 0x04 {
		t.Error("did not correctly set Accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.CFlag != true {
		t.Error("did not correctly set CFlag, got: ", harness.Cpu.CFlag)
	}

	if harness.Cpu.NFlag != false {
		t.Error("did not correctly set NFlag, got: ", harness.Cpu.NFlag)
	}

	if harness.Cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}
}

func TestRORZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x66
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.Memory[0x17] != 0x03 {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0x17])
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles")
	}
}

func TestRORZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x76
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.Memory[0x8F] != 0x03 {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0x8F])
	}

	if harness.Cpu.PC != 0x02 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestRORAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x6E
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.Memory[0x8080] != 0x03 {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0x8080])
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestRORAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Opcode = 0x7E
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.Memory[0xFF81] != 0x03 {
		t.Error("did not correctly set memory value, got: ", harness.Cpu.Memory[0xFF81])
	}

	if harness.Cpu.PC != 0x03 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 7 {
		t.Error("did not correctly set cycles")
	}
}

func TestRTI(t *testing.T) {
	cpu := NewCPU()
	cpu.stackPush(0xFF)
	cpu.Memory[0] = 0x40
	cpu.Exec()

	if cpu.PC != 0x01 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}
	if cpu.VFlag != true {
		t.Error("set overflow flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}

	if cpu.ZFlag != true {
		t.Error("set zero flag incorrectly")
	}

	if cpu.NFlag != true {
		t.Error("set negative flag incorrectly")
	}
}

func TestRTS(t *testing.T) {
	cpu := NewCPU()
	cpu.stackPush(0xFA)
	cpu.Memory[0] = 0x60
	cpu.Exec()

	if cpu.PC != 0xFB {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 6 {
		t.Error("did not correctly set cycles flag")
	}
}

func TestSBC(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x0A
	cpu.CFlag = true
	cpu.Memory[0] = 0xE9
	cpu.Memory[1] = 0x09
	cpu.Exec()

	if cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", cpu.A)
	}

	if cpu.CFlag != true {
		t.Error("did not correctly set carry, got: ", cpu.CFlag)
	}

	if cpu.VFlag != false {
		t.Error("did not correctly set overflow, got: ", cpu.VFlag)
	}

	if cpu.PC != 2 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCWithOverflow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x50
	cpu.CFlag = true
	cpu.Memory[0] = 0xE9
	cpu.Memory[1] = 0xB0
	cpu.Exec()

	if cpu.A != 160 {
		t.Error("did not correctly set accumulator, got: ", cpu.A)
	}

	if cpu.CFlag != false {
		t.Error("did not correctly set carry, got: ", cpu.CFlag)
	}

	if cpu.VFlag != true {
		t.Error("did not correctly set overflow, got: ", cpu.VFlag)
	}
}

func TestSBCWithOverflowWithoutBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0xD0
	cpu.CFlag = true
	cpu.Memory[0] = 0xE9
	cpu.Memory[1] = 0x70
	cpu.Exec()

	if cpu.A != 96 {
		t.Error("did not correctly set accumulator, got: ", cpu.A)
	}

	if cpu.CFlag != true {
		t.Error("did not correctly set carry, got: ", cpu.CFlag)
	}

	if cpu.VFlag != true {
		t.Error("did not correctly set overflow, got: ", cpu.VFlag)
	}
}

func TestSBCWithBorrow(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x02
	cpu.Memory[0] = 0xE9
	cpu.Memory[1] = 0x01
	cpu.Exec()

	if cpu.A != 0 {
		t.Error("did not correctly set accumulator, got: ", cpu.A)
	}

	if cpu.CFlag != true {
		t.Error("did not correctly set carry, got: ", cpu.CFlag)
	}

	if cpu.VFlag != false {
		t.Error("did not correctly set overflow, got: ", cpu.VFlag)
	}
}

func TestSBCZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xE5
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xF5
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xED
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xFD
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCAbsoluteXPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xFD
	harness.SetupAbsoluteXPageCross()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCAbsoluteY(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xF9
	harness.SetupAbsoluteY()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCAbsoluteYPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xF9
	harness.SetupAbsoluteYPageCross()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCIndexedIndirect(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xE1
	harness.SetupIndexedIndirect()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCIndirectIndexed(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xF1
	harness.SetupIndirectIndexed()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles")
	}
}

func TestSBCIndirectIndexedPageCross(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.CFlag = true
	harness.Cpu.A = 0x08
	harness.Opcode = 0xF1
	harness.SetupIndirectIndexedPageCross()
	harness.Run()

	if harness.Cpu.A != 1 {
		t.Error("did not correctly set accumulator, got: ", harness.Cpu.A)
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestSEC(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x38
	cpu.Exec()

	if cpu.CFlag != true {
		t.Error("did not set carry")
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}
}

func TestSED(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0xF8
	cpu.Exec()

	if cpu.DFlag != true {
		t.Error("did not set decimal flag")
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}
}

func TestSEI(t *testing.T) {
	cpu := NewCPU()
	cpu.Memory[0] = 0x78
	cpu.Exec()

	if cpu.IFlag != true {
		t.Error("did not set interrupt disable")
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTAZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.A = 0xDD
	harness.Opcode = 0x85
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.Memory[0x17] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x17])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTAZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.A = 0xDD
	harness.Opcode = 0x95
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.Memory[0x8F] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x8F])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTAAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.A = 0xDD
	harness.Opcode = 0x8D
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.Memory[0x8080] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x8080])
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTAAbsoluteX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.A = 0xDD
	harness.Opcode = 0x9D
	harness.SetupAbsoluteX()
	harness.Run()

	if harness.Cpu.Memory[0xFF81] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0xFF81])
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTAAbsoluteY(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.A = 0xDD
	harness.Opcode = 0x99
	harness.SetupAbsoluteY()
	harness.Run()

	if harness.Cpu.Memory[0xFF81] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0xFF81])
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 5 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTAIndexedIndirect(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.A = 0xDD
	harness.Opcode = 0x81
	harness.SetupIndexedIndirect()
	harness.Run()

	if harness.Cpu.Memory[0x09] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x09])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTAIndirectIndexed(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.A = 0xDD
	harness.Opcode = 0x91
	harness.SetupIndirectIndexed()
	harness.Run()

	if harness.Cpu.Memory[0x06] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x09])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 6 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTXZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.X = 0xDD
	harness.Opcode = 0x86
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.Memory[0x17] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x17])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTXZeroPageY(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.X = 0xDD
	harness.Opcode = 0x96
	harness.SetupZeroPageY()
	harness.Run()

	if harness.Cpu.Memory[0x8F] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x8F])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTXAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.X = 0xDD
	harness.Opcode = 0x8E
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.Memory[0x8080] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x8080])
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTYZeroPage(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.Y = 0xDD
	harness.Opcode = 0x84
	harness.SetupZeroPage()
	harness.Run()

	if harness.Cpu.Memory[0x17] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x17])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 3 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTYZeroPageX(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.Y = 0xDD
	harness.Opcode = 0x94
	harness.SetupZeroPageX()
	harness.Run()

	if harness.Cpu.Memory[0x8F] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x8F])
	}

	if harness.Cpu.PC != 2 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestSTYAbsolute(t *testing.T) {
	harness := NewCpuTestHarness()
	harness.Cpu.Y = 0xDD
	harness.Opcode = 0x8C
	harness.SetupAbsolute()
	harness.Run()

	if harness.Cpu.Memory[0x8080] != 0xDD {
		t.Error("did not correctly set memory, got: ", harness.Cpu.Memory[0x8080])
	}

	if harness.Cpu.PC != 3 {
		t.Error("did not correctly update PC, got", harness.Cpu.PC)
	}

	if harness.Cpu.Cycles != 4 {
		t.Error("did not correctly set cycles")
	}
}

func TestTAX(t *testing.T) {
	cpu := NewCPU()
	cpu.A = 0x87
	cpu.ZFlag = true
	cpu.Memory[0] = 0xAA
	cpu.Exec()

	if cpu.X != 0x87 {
		t.Error("did not correctly set X register, got: ", cpu.X)
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}
}

func TestTAY(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.A = 0x87
	cpu.Memory[0] = 0xA8
	cpu.Exec()

	if cpu.Y != 0x87 {
		t.Error("did not correctly set Y register, got: ", cpu.Y)
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}
}

func TestTSX(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.SP = 0x87
	cpu.Memory[0] = 0xBA
	cpu.Exec()

	if cpu.X != 0x87 {
		t.Error("did not correctly set X register, got: ", cpu.X)
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}
}

func TestTXA(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.X = 0x87
	cpu.Memory[0] = 0x8A
	cpu.Exec()

	if cpu.A != 0x87 {
		t.Error("did not correctly set Y register, got: ", cpu.A)
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}
}

func TestTXS(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.X = 0x87
	cpu.Memory[0] = 0x9A
	cpu.Exec()

	if cpu.SP != 0x87 {
		t.Error("did not correctly set SP register, got: ", cpu.SP)
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}
}

func TestTYA(t *testing.T) {
	cpu := NewCPU()
	cpu.ZFlag = true
	cpu.Y = 0x87
	cpu.Memory[0] = 0x98
	cpu.Exec()

	if cpu.A != 0x87 {
		t.Error("did not correctly set Y register, got: ", cpu.A)
	}

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}

	if cpu.NFlag != true {
		t.Error("did not correctly set NFlag")
	}

	if cpu.ZFlag != false {
		t.Error("did not correctly set ZFlag")
	}
}

func TestCLD(t *testing.T) {
	cpu := NewCPU()
	cpu.DFlag = true
	cpu.Memory[0] = 0xD8
	cpu.Exec()

	if cpu.PC != 1 {
		t.Error("did not correctly update PC, got", cpu.PC)
	}

	if cpu.Cycles != 2 {
		t.Error("did not correctly set cycles")
	}

	if cpu.DFlag != false {
		t.Error("did not clear DFlag")
	}
}
