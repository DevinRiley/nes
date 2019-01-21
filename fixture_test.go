package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"testing"
)

func setupCpu() *CPU {
	cpu := NewCPU()
	cpu.Debug = true
	loadTestRom(cpu)

	return cpu
}

func TestFixtureRom(t *testing.T) {
	// Read File
	logFile, err := os.Open("nestest.log")
	if err != nil {
		panic(err)
	}
	defer logFile.Close()

	scanner := bufio.NewScanner(logFile)

	cpu := setupCpu()
	for scanner.Scan() {
		var pc, a, x, y, sp, flags uint64
		var fail = false

		line := scanner.Text()

		if pc, err = strconv.ParseUint(line[0:4], 16, 64); err != nil {
			fmt.Println("Error parsing pc")
		}

		if a, err = strconv.ParseUint(line[50:52], 16, 64); err != nil {
			fmt.Println("Error parsing accumulator")
		}

		if x, err = strconv.ParseUint(line[55:57], 16, 64); err != nil {
			fmt.Println("Error parsing x register")
		}

		if y, err = strconv.ParseUint(line[60:62], 16, 64); err != nil {
			fmt.Println("Error parsing y register")
		}

		if flags, err = strconv.ParseUint(line[65:67], 16, 64); err != nil {
			fmt.Println("Error parsing flags values")
		}

		if sp, err = strconv.ParseUint(line[71:73], 16, 64); err != nil {
			fmt.Println("Error parsing stack pointer values")
		}

		if a != uint64(cpu.A) {
			fmt.Printf("Accumulator value wrong, expected %X but got %X\n", a, cpu.A)
			fail = true
		}

		if pc != uint64(cpu.PC) {
			fmt.Printf("PC value wrong, expected %X but got %X\n", pc, cpu.PC)
			fail = true
		}

		if x != uint64(cpu.X) {
			fmt.Printf("X value wrong, expected %X but got %X\n", x, cpu.X)
			fail = true
		}

		if y != uint64(cpu.Y) {
			fmt.Printf("Y value wrong, expected %X but got %X\n", y, cpu.Y)
			fail = true
		}

		if flags != uint64(cpu.flagsToByte()) {
			fmt.Printf("Flags value wrong, expected %X but got %X\n", flags, cpu.flagsToByte())
			fail = true
		}

		if sp != uint64(cpu.SP) {
			fmt.Printf("SP value wrong, expected %X but got %X\n", sp, cpu.SP)
			fail = true
		}

		if fail {
			os.Exit(1)
		}

		cpu.Exec()
	}
}
