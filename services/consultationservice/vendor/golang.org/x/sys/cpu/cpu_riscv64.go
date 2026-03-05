



//go:build riscv64

package cpu

const cacheLineSize = 64

func initOptions() {
	options = []option{
		{Name: "fastmisaligned", Feature: &RISCV64.HasFastMisaligned},
		{Name: "c", Feature: &RISCV64.HasC},
		{Name: "v", Feature: &RISCV64.HasV},
		{Name: "zba", Feature: &RISCV64.HasZba},
		{Name: "zbb", Feature: &RISCV64.HasZbb},
		{Name: "zbs", Feature: &RISCV64.HasZbs},
	}
}
