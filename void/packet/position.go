package packet

type Position struct {
	X float32
	Y float32
	Z float32
}

func (p Position) PositionX() float32 {
	return p.X
}
func (p Position) PositionY() float32 {
	return p.Y
}
func (p Position) PositionZ() float32 {
	return p.Z
}
func NewPosition(x, y, z float32) Position {
	return Position{x, y, z}
}
