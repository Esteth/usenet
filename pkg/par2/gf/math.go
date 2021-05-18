package gf

//go:generate go run gf_gen.go
//go:generate gofmt -w gf.go

const NW uint32 = 1 << 16

func Add(a int32, b int32) int32 {
	return a ^ b
}

func Mul(a int32, b int32) int32 {
	if a == 0 || b == 0 {
		return 0
	}
	sum_log := uint32(Gflog[a] + Gflog[b])
	if sum_log > NW-1 {
		sum_log -= NW - 1
	}
	return Gfilog[sum_log]
}

func Div(a int32, b int32) int32 {
	if a == 0 {
		return 0
	}
	if b == 0 {
		panic("division by zero")
	}
	diff_log := Gflog[a] - Gflog[b]
	if diff_log < 0 {
		diff_log += int32(NW - 1)
	}
	return Gfilog[diff_log]
}
