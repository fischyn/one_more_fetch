package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/fischyn/wfetch/sys/cpu"
)

type FIELD uint8

const (
	//CPU fields
	PROCESSOR_NAME FIELD = iota
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

		if offset+3 > len(data) {
			return ret, fmt.Errorf("unexpected end of data at position %d", offset)
		}

		field := FIELD(data[offset])

		offset++

		fieldLen := binary.LittleEndian.Uint16(data[offset : offset+2])

		offset += 2

		if offset+int(fieldLen) > len(data) {
			return ret, fmt.Errorf("not enough data for field %d", field)
		}

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
			if len(value) != 4 {
				return ret, fmt.Errorf("invalid length for MHZ")
			}
			ret.Mhz = binary.LittleEndian.Uint32(value)
		case CORES_PHYSICAL:
			if len(value) != 2 {
				return ret, fmt.Errorf("invalid length for CORES_PHYSICAL")
			}
			ret.CoresPhysical = binary.LittleEndian.Uint16(value)
		case CORES_LOGICAL:
			if len(value) != 2 {
				return ret, fmt.Errorf("invalid length for CORES_LOGICAL")
			}
			ret.CoresLogical = binary.LittleEndian.Uint16(value)
		case CORES_ACTIVE:
			if len(value) != 2 {
				return ret, fmt.Errorf("invalid length for CORES_ACTIVE")
			}
			ret.CoresActive = binary.LittleEndian.Uint16(value)
		case PACKAGES:
			if len(value) != 2 {
				return ret, fmt.Errorf("invalid length for PACKAGES")
			}
			ret.Packages = binary.LittleEndian.Uint16(value)
		default:
			return ret, fmt.Errorf("unknown field id %d", field)
		}
	}

	return ret, nil
}

func MarshalCPUBinary(c *cpu.CPUResult) ([]byte, error) {
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

func writeFiled(field FIELD, data []byte, buf *bytes.Buffer) {
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
