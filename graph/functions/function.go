package functions

import "math"

type BinaryFunction func(x, y float64) bool

// 爱心
func Love(x, y float64) bool {
	return math.Pow(math.Pow(x, 2)+math.Pow(y, 2)-1, 3)-math.Pow(x, 2)*math.Pow(y, 3) <= 0
}

// 循环爱心
func CircularLove(x, y float64) bool {
	return y <= math.Sqrt(math.Abs(math.Sin(x)))+math.Sqrt(math.Abs(math.Cos(x))) && y >= math.Sqrt(math.Abs(math.Sin(x)))-math.Sqrt(math.Abs(math.Cos(x)))
}

func RoseLine(x, y float64) bool {
	theta := math.Atan2(y, x)
	r := 2 * math.Sin(4*theta)
	xp := r * math.Cos(theta)
	yp := r * math.Sin(theta)
	return math.Abs(x) <= math.Abs(xp) && math.Abs(y) <= math.Abs(yp)
}
