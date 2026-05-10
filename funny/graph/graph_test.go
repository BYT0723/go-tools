package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultEdge(t *testing.T) {
	t.Run("DefaultEdge 测试", func(t *testing.T) {
		assert.NotNil(t, DefaultEdge)
		assert.Equal(t, " ", DefaultEdge.Zero)
		assert.Equal(t, " ", DefaultEdge.X)
		assert.Equal(t, " ", DefaultEdge.Y)
		assert.Equal(t, "*", DefaultEdge.FC)
		assert.Equal(t, " ", DefaultEdge.Scale)
	})
}

func TestStandbyEdge(t *testing.T) {
	t.Run("StandbyEdge 测试", func(t *testing.T) {
		assert.NotNil(t, StandbyEdge)
		assert.Equal(t, "+", StandbyEdge.Zero)
		assert.Equal(t, "-", StandbyEdge.X)
		assert.Equal(t, "|", StandbyEdge.Y)
		assert.Equal(t, "*", StandbyEdge.FC)
		assert.Equal(t, "+", StandbyEdge.Scale)
	})
}

func TestEdgeStruct(t *testing.T) {
	t.Run("Edge 结构测试", func(t *testing.T) {
		e := &Edge{Zero: "0", X: "-", Y: "|", FC: "#", Scale: "+"}
		assert.Equal(t, "0", e.Zero)
		assert.Equal(t, "#", e.FC)
	})
}

func TestGraphDefaultEdge(t *testing.T) {
	t.Run("Graph 默认Edge测试", func(t *testing.T) {
		g := &Graph{XScale: 0.1, YScale: 0.1}
		assert.Nil(t, g.Edge)

		// Draw sets default edge if nil (may fail if not in terminal)
		g.Draw(func(x, y float64) bool { return false })
	})
}
