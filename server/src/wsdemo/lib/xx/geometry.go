package xx

import (
	"math"
)

type Pos struct{ X, Y float64 }

func (p Pos) Add(q Pos) Pos {
	p.X += q.X
	p.Y += q.Y
	return p
}

func (p Pos) Sub(q Pos) Pos {
	p.X -= q.X
	p.Y -= q.Y
	return p
}

func (p Pos) Mul(s float64) Pos {
	p.X *= s
	p.Y *= s
	return p
}

func (p Pos) Dot(q Pos) float64 {
	return p.X*q.X + p.Y*q.Y
}

func (p Pos) Len() float64 {
	return math.Hypot(p.X, p.Y)
}

func (p Pos) Distance(q Pos) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// 判断两个方块是否重叠
func Overlap(p1, p2 Pos, r float64) bool {
	s := p1.Sub(p2)
	b := math.Abs(s.X) < r
	return b && math.Abs(s.Y) < r
}

func (dp Pos) Collide(p1, p2 Pos, r float64) Pos {
	p := p1.Add(dp)
	if Overlap(p, p2, r) {
		s := p1.Sub(p2)
		dp.X = (math.Abs(s.X) - r) * Sign(dp.X)
		dp.Y = (math.Abs(s.Y) - r) * Sign(dp.Y)
	}
	return dp
}

func (p Pos) Map() map[string]float64 {
	m := map[string]float64{
		"x": p.X, "y": p.Y,
	}
	return m
}

func Displace(dir string, dis float64) Pos {
	var x, y float64
	switch dir {
	case "up":
		y = dis
		break
	case "down":
		y = -dis
		break
	case "left":
		x = -dis
		break
	case "right":
		x = dis
		break
	}
	return Pos{X: x, Y: y}
}
