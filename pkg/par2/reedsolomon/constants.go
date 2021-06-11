package reedsolomon

import (
	"github.com/esteth/usenet/pkg/par2/gf"
)

type constantPool struct {
	currentPower int
	currentValue uint16
}

func newConstantPool() constantPool {
	return constantPool{
		currentPower: 0,
		currentValue: 1,
	}
}

func (p *constantPool) Next() uint16 {
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
