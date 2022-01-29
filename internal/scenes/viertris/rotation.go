package viertris

type Rotation uint8

const (
	Rotation0   Rotation = iota
	Rotation90  Rotation = iota
	Rotation180 Rotation = iota
	Rotation270 Rotation = iota
	RotationMax Rotation = iota
)

func (r Rotation) RotateRight() Rotation {
	r++
	if r > Rotation270 {
		return Rotation0
	}
	return r
}

func (r Rotation) RotateLeft() Rotation {
	r--
	if r == 255 {
		return Rotation270
	}
	return r
}
