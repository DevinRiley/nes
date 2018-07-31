package main

import (
	"errors"
	"fmt"
	"io/ioutil"
)

type Mirroring int

const (
	Vertical Mirroring = iota
	Horizontal
)

type TVSystem int

const (
	NTSC TVSystem = iota
	PAL
)

type ROM struct {
	Mirroring       Mirroring
	CartridgeMemory bool
	Trainer         bool
	FourScreen      bool
	Mapper          uint8
	VSUnisystem     bool
	NES2Format      bool
	TVSystem        TVSystem
	PRGSize         uint
	CHRSize         uint
}

// ## Flags 6 #
//76543210
//||||||||
//|||||||+- Mirroring: 0: horizontal (vertical arrangement) (CIRAM A10 = PPU A11)
//|||||||              1: vertical (horizontal arrangement) (CIRAM A10 = PPU A10)
//||||||+-- 1: Cartridge contains battery-backed PRG RAM ($6000-7FFF) or other persistent memory
//|||||+--- 1: 512-byte trainer at $7000-$71FF (stored before PRG data)
//||||+---- 1: Ignore mirroring control or above mirroring bit; instead provide four-screen VRAM
//++++----- Lower nybble of mapper number

func parseFlags6Mirroring(flags byte) Mirroring {
	// Right-most bit position
	return Mirroring(flags & 0x01)
}

func parseFlags6CartridgeMemory(flags byte) bool {
	// Second bit position from the right
	return (uint8(flags>>1) & 0x01) == 1
}

func parseFlags6Trainer(flags byte) bool {
	// 512-byte trainer at $7000-$71FF (stored before PRG data)
	// third bit position from the right
	return (uint8(flags>>2) & 0x01) == 1
}

func parseFlags6FourScreen(flags byte) bool {
	// Ignore mirroring control or above mirroring bit; instead provide four-screen VRAM
	// fourth bit position from the right
	return (uint8(flags>>3) & 0x01) == 1
}

func parseFlags6MapperLowerNibble(flags byte) uint8 {
	// Lower nibble of mapper number
	return uint8(flags>>4) & 0x0F
}

// # Flags 7 #
// 76543210
// ||||||||
// |||||||+- VS Unisystem
// ||||||+-- PlayChoice-10 (8KB of Hint Screen data stored after CHR data)
// ||||++--- If equal to 2, flags 8-15 are in NES 2.0 format
// ++++----- Upper nybble of mapper number

func parseFlags7VSUnisystem(flags byte) bool {
	// Right-most bit position
	return (uint8(flags) & 0x01) == 1
}

func parseFlags7NES2RomFormat(flags byte) bool {
	// Third and fourth bit position from the right
	return uint8(flags>>2) == 2
}

func parseFlags7MapperUpperNibble(flags byte) uint8 {
	return uint8(flags >> 4)
}

// # Flags 9 #
// 76543210
//||||||||
//|||||||+- TV system (0: NTSC; 1: PAL)
//+++++++-- Reserved, set to zero

func parseFlags9TVSystem(flags byte) TVSystem {
	return TVSystem(uint8(flags) & 0x01)
}

func parsePrgRomSize(size byte) uint {
	// value is in 16kB blocks
	return uint(size) * 16384
}

func parseChrRomSize(size byte) uint {
	// value is in 8kB blocks
	return uint(size) * 8192
}

func validateHeader(header []byte) error {
	magicHeader := []byte{'N', 'E', 'S', 0x1A}

	for i, _ := range magicHeader {
		if header[i] != magicHeader[i] {
			return errors.New("Does not contain magic header!")
		}
	}

	return nil
}

func parseMapper(flags []byte) uint8 {
	lower := parseFlags6MapperLowerNibble(flags[0])
	upper := parseFlags7MapperUpperNibble(flags[1])

	return (upper << 4) | lower
}

func parseRom(data []byte) (*ROM, error) {
	var rom ROM

	err := validateHeader(data)
	if err != nil {
		fmt.Println(err)
		return &rom, err
	}

	rom.PRGSize = parsePrgRomSize(data[4])
	rom.CHRSize = parseChrRomSize(data[5])
	rom.Mirroring = parseFlags6Mirroring(data[6])
	rom.CartridgeMemory = parseFlags6CartridgeMemory(data[6])
	rom.Trainer = parseFlags6Trainer(data[6])
	rom.FourScreen = parseFlags6FourScreen(data[6])
	rom.Mapper = parseMapper(data[6:8])
	rom.VSUnisystem = parseFlags7VSUnisystem(data[7])
	rom.NES2Format = parseFlags7NES2RomFormat(data[7])
	rom.TVSystem = parseFlags9TVSystem(data[9])

	return &rom, err
}

func printRom(rom *ROM) {
	fmt.Printf("%+v\n", rom)
}

func main() {
	cpu := NewCPU()
	cpu.Memory[0] = 0x29
	cpu.Print()
	cpu.Exec()
	fmt.Printf("Memory: %d\n", uint(cpu.Memory[0]))
	cpu.Print()
	fmt.Println("Hello Nintendo World")

	// Slurp file into memory
	romFile, err := ioutil.ReadFile("nestest.nes")

	if err != nil {
		panic(fmt.Sprintf("Could not read rom file!"))
	}

	rom, err := parseRom(romFile)
	if err != nil {
		fmt.Println("Error parsing rom")
	} else {
		printRom(rom)
	}
}
