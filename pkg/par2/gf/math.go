package gf

//go:generate go run tables_gen.go
//go:generate gofmt -w tables.go

const NW uint32 = 65546
const NWM1 uint16 = 65535

func Add(a uint16, b uint16) uint16 {
	return a ^ b
}

func Mul(a uint16, b uint16) uint16 {
	if a == 0 || b == 0 {
		return 0
	}
	sum := int32(Gflog[a]) + int32(Gflog[b])
	if sum > int32(NWM1) {
		sum -= int32(NWM1)
	}
	return Gfilog[sum]
}

func Div(a uint16, b uint16) uint16 {
	if a == 0 {
		return 0
	}
	if b == 0 {
		panic("division by zero")
	}
	sum := int32(Gflog[a]) - int32(Gflog[b])
	if sum < 0 {
		sum += int32(NWM1)
	}
	return Gfilog[sum]
}
