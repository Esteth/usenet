package par2

import (
	"github.com/esteth/usenet/pkg/par2/gf"
)

type ConstantPool struct {
	currentPower int
	currentValue uint16
}

func NewConstantPool() ConstantPool {
	return ConstantPool{
		currentPower: 0,
		currentValue: 1,
	}
}

func (p *ConstantPool) Next() uint16 {
	for {
		p.currentValue = gf.Mul(p.currentValue, 2)
		p.currentPower++
		if p.currentPower%3 != 0 &&
			p.currentPower%5 != 0 &&
			p.currentPower%17 != 0 &&
			p.currentPower%257 != 0 {
			break
		}
	}
	return p.currentValue
}
