package graph

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDefaultEdge(t *testing.T) {
	Convey("DefaultEdge 测试", t, func() {
		So(DefaultEdge, ShouldNotBeNil)
		So(DefaultEdge.Zero, ShouldEqual, " ")
		So(DefaultEdge.X, ShouldEqual, " ")
		So(DefaultEdge.Y, ShouldEqual, " ")
		So(DefaultEdge.FC, ShouldEqual, "*")
		So(DefaultEdge.Scale, ShouldEqual, " ")
	})
}

func TestStandbyEdge(t *testing.T) {
	Convey("StandbyEdge 测试", t, func() {
		So(StandbyEdge, ShouldNotBeNil)
		So(StandbyEdge.Zero, ShouldEqual, "+")
		So(StandbyEdge.X, ShouldEqual, "-")
		So(StandbyEdge.Y, ShouldEqual, "|")
		So(StandbyEdge.FC, ShouldEqual, "*")
		So(StandbyEdge.Scale, ShouldEqual, "+")
	})
}

func TestEdgeStruct(t *testing.T) {
	Convey("Edge 结构测试", t, func() {
		e := &Edge{Zero: "0", X: "-", Y: "|", FC: "#", Scale: "+"}
		So(e.Zero, ShouldEqual, "0")
		So(e.FC, ShouldEqual, "#")
	})
}

func TestGraphDefaultEdge(t *testing.T) {
	Convey("Graph 默认Edge测试", t, func() {
		g := &Graph{XScale: 0.1, YScale: 0.1}
		So(g.Edge, ShouldBeNil)

		// Draw sets default edge if nil (may fail if not in terminal)
		g.Draw(func(x, y float64) bool { return false })
	})
}
