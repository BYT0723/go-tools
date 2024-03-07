package funny

import (
	"fmt"
	"math"
)

type BinaryFunction func(x, y float64) bool

type Graph struct {
	XScale float64
	YScale float64
	Edge   *Edge
}

type Edge struct {
	// 原点字符
	Zero string
	// X轴
	X string
	// Y轴
	Y string
	// 符合函数的字符
	FC string
	// 刻度线
	Scale string
}

var defaultEdge = &Edge{Zero: "+", X: "-", Y: "|", FC: "*", Scale: "+"}

func (g *Graph) Draw(fc BinaryFunction) error {
	size, err := GetTermSize()
	if err != nil {
		return err
	}

	if g.Edge == nil {
		g.Edge = defaultEdge
	}

	// 打印爱心形状
	for row := int(size.Row) / 2; row > 0-int(size.Row)/2; row-- {
		for col := 0 - int(size.Col)/2; col < int(size.Col)/2; col++ {
			var (
				y = float64(row) * g.YScale
				x = float64(col) * g.XScale
			)
			if fc(x, y) {
				fmt.Print("*")
			} else if x == 0 && y == 0 {
				fmt.Print(g.Edge)
			} else if x == 0 {
				if y == math.Floor(y) {
					fmt.Print(g.Edge.Scale)
				} else {
					fmt.Print(g.Edge.Y)
				}
			} else if y == 0 {
				if x == math.Floor(x) {
					fmt.Print(g.Edge.Scale)
				} else {
					fmt.Print(g.Edge.X)
				}
			} else {
				fmt.Print(" ")
			}
		}
	}
	return nil
}
