package encoding

import (
	"bytes"
	"encoding/binary"

	"github.com/fischyn/omfetch/sys/cpu"
)

type CPURS_FIELD uint8

const (
	PROCESSOR_NAME CPURS_FIELD = iota
	VENDOR
	IDENTIFIER
	MHZ
	CORES_PHYSICAL
	CORES_LOGICAL
	CORES_ACTIVE
	PACKAGES
)

func UnMarshalCPU(data []byte) (cpu.CPUResult, error) {

	var ret cpu.CPUResult

	offset := 0

	for offset < len(data) {
		field := CPURS_FIELD(data[offset])

		offset++

		fieldLen := binary.LittleEndian.Uint16(data[offset : offset+2])

		offset += 2

		value := data[offset : offset+int(fieldLen)]

		offset += int(fieldLen)

		switch field {
		case PROCESSOR_NAME:
			ret.ProcessorName = string(value)
		case VENDOR:
			ret.Vendor = string(value)
		case IDENTIFIER:
			ret.Identifier = string(value)
		case MHZ:
			ret.Mhz = binary.LittleEndian.Uint32(value)
		case CORES_PHYSICAL:
			ret.CoresPhysical = binary.LittleEndian.Uint16(value)
		case CORES_LOGICAL:
			ret.CoresLogical = binary.LittleEndian.Uint16(value)
		case CORES_ACTIVE:
			ret.CoresActive = binary.LittleEndian.Uint16(value)
		case PACKAGES:
			ret.Packages = binary.LittleEndian.Uint16(value)
		}
	}

	return ret, nil
}

func MarshalCPU(c *cpu.CPUResult) ([]byte, error) {
	var buf bytes.Buffer

	writeFiled(PROCESSOR_NAME, []byte(c.ProcessorName), &buf)
	writeFiled(VENDOR, []byte(c.Vendor), &buf)
	writeFiled(IDENTIFIER, []byte(c.Identifier), &buf)

	var buf4 [4]byte
	binary.LittleEndian.PutUint32(buf4[:], c.Mhz)
	writeFiled(MHZ, buf4[:], &buf)

	var buf2 [2]byte
	binary.LittleEndian.PutUint16(buf2[:], c.CoresPhysical)
	writeFiled(CORES_PHYSICAL, buf2[:], &buf)

	binary.LittleEndian.PutUint16(buf2[:], c.CoresLogical)
	writeFiled(CORES_LOGICAL, buf2[:], &buf)

	binary.LittleEndian.PutUint16(buf2[:], c.CoresActive)
	writeFiled(CORES_ACTIVE, buf2[:], &buf)

	binary.LittleEndian.PutUint16(buf2[:], c.Packages)
	writeFiled(PACKAGES, buf2[:], &buf)

	return buf.Bytes(), nil
}

func writeFiled(field CPURS_FIELD, data []byte, buf *bytes.Buffer) {
	// 1 b for field name + 2 byte for fieldValue length + n bytes for data
	//0 +----+--------+---------------+ N
	//  name + fLen   + data          +
	//0 +----+--------+---------------+ N
	buf.WriteByte(byte(field))

	var lenBuf [2]byte
	binary.LittleEndian.PutUint16(lenBuf[:], uint16(len(data)))
	buf.Write(lenBuf[:])

	buf.Write(data)
}
