package graph

import (
	"math"

	"github.com/BYT0723/go-tools/funny/graph/functions"
	"github.com/BYT0723/go-tools/osx"
)

type Graph struct {
	XScale float64
	YScale float64
	// 打印字符的样式
	Edge *Edge
}

func (g *Graph) Draw(fc functions.BinaryFunction) error {
	w, h, err := osx.GetTermSize()
	if err != nil {
		return err
	}

	if g.Edge == nil {
		g.Edge = DefaultEdge
	}

	// 打印爱心形状
	for row := int(h) / 2; row > 0-int(h)/2; row-- {
		for col := 0 - int(w)/2; col < int(w)/2; col++ {
			var (
				y = float64(row) * g.YScale
				x = float64(col) * g.XScale
			)
			if fc(x, y) {
				print(g.Edge.FC)
			} else if x == 0 && y == 0 {
				print(g.Edge.Zero)
			} else if x == 0 {
				if y == math.Floor(y) {
					print(g.Edge.Scale)
				} else {
					print(g.Edge.Y)
				}
			} else if y == 0 {
				if x == math.Floor(x) {
					print(g.Edge.Scale)
				} else {
					print(g.Edge.X)
				}
			} else {
				print(" ")
			}
		}
	}
	return nil
}
