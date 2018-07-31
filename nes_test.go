package main

import (
	"testing"
)

func TestReturErrorOnBadHeader(t *testing.T) {
	badHeader := []byte{'B', 'A', 'D', 0x1A}
	err := validateHeader(badHeader)
	if err == nil {
		t.Error("bad header did not return error")
	}
}

func TestNoErrorOnBadHeader(t *testing.T) {
	badHeader := []byte{'N', 'E', 'S', 0x1A}
	err := validateHeader(badHeader)
	if err != nil {
		t.Error("good header is returning an error")
	}
}

func TestGetPrgRomSize(t *testing.T) {
	header := byte(0xFF)
	size := parsePrgRomSize(header)

	if size != (16384 * 0xFF) {
		t.Error("Incorrectly read PRG ROM size")
	}
}

func TestGetChrRomSize(t *testing.T) {
	header := byte(0xBA)
	size := parseChrRomSize(header)

	if size != (8192 * 0xBA) {
		t.Error("Incorrectly read CHR ROM size")
	}
}

func TestParseFlags6MirroringVertical(t *testing.T) {
	flags := byte(0x00)

	mirroring := parseFlags6Mirroring(flags)

	if mirroring != Vertical {
		t.Error("Incorrectly parsed vertical mirroring")
	}
}

func TestParseFlags6MirroringHorizontal(t *testing.T) {
	flags := byte(0x01)

	mirroring := parseFlags6Mirroring(flags)

	if mirroring != Horizontal {
		t.Error("Incorrectly parsed horizontal mirroring")
	}
}

func TestParseFlags6CartridgeDoesNotHaveMemory(t *testing.T) {
	flags := byte(0x00)

	onboardMemory := parseFlags6CartridgeMemory(flags)

	if onboardMemory != false {
		t.Error("Incorrectly parsed cartridge onboard memory")
	}
}

func TestParseFlags6CartridgeHasMemory(t *testing.T) {
	flags := byte(0x02)

	onboardMemory := parseFlags6CartridgeMemory(flags)

	if onboardMemory != true {
		t.Error("Incorrectly parsed cartridge onboard memory")
	}
}

func TestParseFlags6Trainer(t *testing.T) {
	flags := byte(0x00)

	trainer := parseFlags6Trainer(flags)

	if trainer != false {
		t.Error("Incorrectly parsed presence of trainer data")
	}
}

func TestParseFlags6NoTrainer(t *testing.T) {
	flags := byte(0xFF)

	trainer := parseFlags6Trainer(flags)

	if trainer != true {
		t.Error("Incorrectly parsed presence of trainer data")
	}
}

func TestParseFlags6DoNotFourScreen(t *testing.T) {
	flags := byte(0x00)

	fourScreen := parseFlags6FourScreen(flags)

	if fourScreen != false {
		t.Error("Incorrectly parsed four screen flag")
	}
}

func TestParseFlags6FourScreen(t *testing.T) {
	flags := byte(0x08)

	fourScreen := parseFlags6FourScreen(flags)

	if fourScreen != true {
		t.Error("Incorrectly parsed four screen flag")
	}
}

func TestParseFlags6MapperLowerNibble(t *testing.T) {
	flags := byte(0x08 << 4)

	lowerNibble := parseFlags6MapperLowerNibble(flags)

	if lowerNibble != 0x08 {
		t.Error("Incorrectly parsed mapper lower nibble")
	}
}

func TestParseFlags7NotVSUnisystem(t *testing.T) {
	flags := byte(0x00)

	vsUnisystem := parseFlags7VSUnisystem(flags)

	if vsUnisystem != false {
		t.Error("Incorrectly parsed VS Unisystem flag")
	}
}

func TestParseFlags7IsVSUnisystem(t *testing.T) {
	flags := byte(0x01)

	vsUnisystem := parseFlags7VSUnisystem(flags)

	if vsUnisystem != true {
		t.Error("Incorrectly parsed VS Unisystem flag")
	}
}

func TestParseFlags7NotNES2RomFormat(t *testing.T) {
	flags := byte(0x00)

	nes2RomFormat := parseFlags7NES2RomFormat(flags)

	if nes2RomFormat != false {
		t.Error("Incorrectly parsed NES 2.0 ROM format")
	}
}

func TestParseFlags7IsNES2RomFormat(t *testing.T) {
	flags := byte(0x08)

	nes2RomFormat := parseFlags7NES2RomFormat(flags)

	if nes2RomFormat != true {
		t.Error("Incorrectly parsed NES 2.0 ROM format")
	}
}

func TestParseFlags7MapperUpperNibble(t *testing.T) {
	flags := byte(0x0C << 4)

	nibble := parseFlags7MapperUpperNibble(flags)

	if nibble != 0x0C {
		t.Error("Incorrectly parsed mapper upper nibble")
	}
}

func TestParseFlags9TVSystemNTSC(t *testing.T) {
	flags := byte(0x00)

	isNTSC := parseFlags9TVSystem(flags)

	if isNTSC != 0x00 {
		t.Error("Incorrectly parsed TV System NTSC")
	}
}

func TestParseFlags9TVSystemPAL(t *testing.T) {
	flags := byte(0x01)

	isNTSC := parseFlags9TVSystem(flags)

	if isNTSC != 0x01 {
		t.Error("Incorrectly parsed TV System NTSC")
	}
}

func TestParseMapper(t *testing.T) {
	// pass in flags 6 and 7, since mapper is spread across both
	flagSix := byte(0x01 << 4)
	flagSeven := byte(0x01 << 4)
	flags := []byte{flagSix, flagSeven}

	mapper := parseMapper(flags)

	if mapper != 0x11 {
		t.Error("Incorrectly parsed full mapper byte")
	}
}

func TestParseRom(t *testing.T) {
	romData := []byte{
		'N',  // Magic Header
		'E',  // Magic Header
		'S',  // Magic Header
		0x1A, // Magic Header
		0x0F, // PRG Size
		0x0A, // CHR Size
		0xFF, // Flags 6
		0xFF, // Flags 7
		0x00, // PRG RAM Size
		0x01, // Flags 9
		0x00, // Flags 10
		0x00, // Zero-filled
		0x00, // Zero-filled
		0x00, // Zero-filled
		0x00, // Zero-filled
		0x00, // Zero-filled
		0x00, // Zero-filled
		0x00, // Zero-filled
	}

	rom, err := parseRom(romData)
	if err != nil {
		t.Error("Failed to parse ROM!")
	}

	if rom.PRGSize != 245760 {
		t.Error("Incorrect PRGSize")
	}

	if rom.CHRSize != 0x0A*8192 {
		t.Error("Incorrect CHRSize")
	}

	if rom.Mirroring != Horizontal {
		t.Error("Incorrect Mirroring")
	}

	if rom.CartridgeMemory != true {
		t.Error("Incorrect CartridgeMemory")
	}

	if rom.Trainer != true {
		t.Error("Incorrect Trainer")
	}

	if rom.FourScreen != true {
		t.Error("Incorrect FourScreen")
	}

	if rom.Mapper != 0xFF {
		t.Error("Incorrect Mapper")
	}

	if rom.VSUnisystem != true {
		t.Error("Incorrect VSUnisystem")
	}

	if rom.NES2Format != false {
		t.Error("Incorrect NES 2.0 Format")
	}

	if rom.TVSystem != PAL {
		t.Error("Incorrect TVSystem (NTSC or PAL)")
	}
}
