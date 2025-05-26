package encoding_test

import (
	"testing"

	"bytes"

	"github.com/fischyn/omfetch/encoding"
	"github.com/fischyn/omfetch/sys/cpu"
)

func TestMarshalUnmarshalCPU(t *testing.T) {
	original := cpu.CPUResult{
		ProcessorName: "Intel(R) Core(TM) i7-9750H",
		Vendor:        "GenuineIntel",
		Identifier:    "Intel64 Family 6 Model 158",
		Mhz:           2600,
		CoresPhysical: 6,
		CoresLogical:  12,
		CoresActive:   6,
		Packages:      1,
	}

	data, err := encoding.MarshalCPUBinary(&original)

	if err != nil {
		t.Fatalf("MarshalCPU failed: %v", err)
	}

	result, err := encoding.UnMarshalCPU(data)

	if err != nil {
		t.Fatalf("UnMarshalCPU failed: %v", err)
	}

	if result.ProcessorName != original.ProcessorName {
		t.Errorf("ProcessorName mismatch: got %q, want %q", result.ProcessorName, original.ProcessorName)
	}
	if result.Vendor != original.Vendor {
		t.Errorf("Vendor mismatch: got %q, want %q", result.Vendor, original.Vendor)
	}
	if result.Identifier != original.Identifier {
		t.Errorf("Identifier mismatch: got %q, want %q", result.Identifier, original.Identifier)
	}
	if result.Mhz != original.Mhz {
		t.Errorf("Mhz mismatch: got %d, want %d", result.Mhz, original.Mhz)
	}
	if result.CoresPhysical != original.CoresPhysical {
		t.Errorf("CoresPhysical mismatch: got %d, want %d", result.CoresPhysical, original.CoresPhysical)
	}
	if result.CoresLogical != original.CoresLogical {
		t.Errorf("CoresLogical mismatch: got %d, want %d", result.CoresLogical, original.CoresLogical)
	}
	if result.CoresActive != original.CoresActive {
		t.Errorf("CoresActive mismatch: got %d, want %d", result.CoresActive, original.CoresActive)
	}
	if result.Packages != original.Packages {
		t.Errorf("Packages mismatch: got %d, want %d", result.Packages, original.Packages)
	}

	//Check if Marshal -> Unmarshal is idenpotent operation
	data2, err := encoding.MarshalCPUBinary(&result)

	if err != nil {
		t.Fatalf("MarshalCPU failed second time: %v", err)
	}

	if !bytes.Equal(data, data2) {
		t.Errorf("Binary output mismatch after re-marshal")
	}
}

var sampleCPU = cpu.CPUResult{
	ProcessorName: "Intel(R) Core(TM) i7-9750H",
	Vendor:        "GenuineIntel",
	Identifier:    "Intel64 Family 6 Model 158",
	Mhz:           2600,
	CoresPhysical: 6,
	CoresLogical:  12,
	CoresActive:   6,
	Packages:      1,
}

func BenchmarkMarshalCPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := encoding.MarshalCPUBinary(&sampleCPU)
		if err != nil {
			b.Fatalf("MarshalCPU failed: %v", err)
		}
	}
}

func BenchmarkUnmarshalCPU(b *testing.B) {
	data, err := encoding.MarshalCPUBinary(&sampleCPU)
	if err != nil {
		b.Fatalf("MarshalCPU failed: %v", err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := encoding.UnMarshalCPU(data)
		if err != nil {
			b.Fatalf("UnMarshalCPU failed: %v", err)
		}
	}
}
