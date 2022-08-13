package main

import (
	"github.com/faiface/pixel"
)

type State struct {
	Points         []pixel.Vec
	TValue         int
	DragPointIndex int
	CurvePoints    []pixel.Vec
}

func NewState() *State {
	state := &State{}
	state.Points = make([]pixel.Vec, 0)
	state.TValue = 33
	state.DragPointIndex = -1

	return state
}

func (state *State) PointLength() int {
	return len(state.Points)
}

func (state *State) AddPoint(point pixel.Vec) {
	state.Points = append(state.Points, point)
}

func (state *State) GenerateCurvePoints() {
	originalTValue := state.TValue
	state.CurvePoints = make([]pixel.Vec, 0)
	for i :=0; i<=100; i++ {
		state.TValue = i
		mainTPoints := GetTPoints(state.Points)
		bezierPoints := GetTPoints(mainTPoints)
		finalPoints := GetTPoints(bezierPoints)
		state.CurvePoints = append(state.CurvePoints, finalPoints...)
	}
	state.TValue = originalTValue
}
