package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/faiface/pixel/imdraw"
	"math"
	"github.com/faiface/pixel/text"
  	"golang.org/x/image/font/basicfont"
	"image/color"
	"strconv"
)

const (
	BEZIER_WINDOW_HEIGHT = 768.0
	BEZIER_WINDOW_WIDTH = 1024.0
	CONTROL_WINDOW_HEIGHT = 669.0
	CONTROL_WINDOW_WIDTH = 500.0
	WINDOW_GAP = 10.0
	TOP_PADDING = 15.0
	USABLE_CONTROL_HEIGHT = CONTROL_WINDOW_HEIGHT - TOP_PADDING
	INIT_X = 69.0
	INIT_Y = 69.0
	POINT_RADIUS = 6.9
	MAX_POINTS = 4.0
	BASE_TEXT_HEIGHT = 12.0
	TEXT_HEIGHT =  20.0
	TEXT_SCALE = TEXT_HEIGHT / BASE_TEXT_HEIGHT
)

var (
	bezierWin *pixelgl.Window
	controlWin *pixelgl.Window
	state = NewState()
)

func init() {
	fmt.Println("Initializing Bezier Boy")
}

func main() {
	fmt.Println("Starting Bezier Boy")
	pixelgl.Run(Run)
}

func Run() {
	// Main Bezier Window
	bezierConfig := pixelgl.WindowConfig{
		Title:  "Bezier Boy",
		Bounds: pixel.R(0, 0, BEZIER_WINDOW_WIDTH, BEZIER_WINDOW_HEIGHT),
		VSync:  true,
		Resizable: false,
	}
	bezierWin, _ = pixelgl.NewWindow(bezierConfig)
	// Bezier Controls Window
	controlConfig := pixelgl.WindowConfig{
		Title:  "Bezier Boy Controls",
		Bounds: pixel.R(0, 0, CONTROL_WINDOW_WIDTH, CONTROL_WINDOW_HEIGHT),
		VSync:  true,
		Resizable: false,
	}
	controlWin, _ = pixelgl.NewWindow(controlConfig)

	bezierWin.SetPos(pixel.V(INIT_X, INIT_Y))
	controlWin.SetPos(pixel.V(INIT_X + BEZIER_WINDOW_WIDTH + WINDOW_GAP, INIT_Y))

	
	for !bezierWin.Closed() && !controlWin.Closed() {
		controlWin.Clear(colornames.Black)
		bezierWin.Clear(colornames.Black)

		HandleEvents()
		state.GenerateCurvePoints()
		state.TValue += 1
		if state.TValue > 100 {
			state.TValue = 0
		}
		Draw()

		controlWin.Update()
		bezierWin.Update()
	}
}

func Draw() {
	DrawControls()
	DrawBezier()
}

func DrawBezier() {
	imd := imdraw.New(nil)
	// Draw Bezier Points
	imd.Color = colornames.Pink
	// for _, point := range state.CurvePoints {
	// 	imd.Push(point)
	// 	imd.Circle(POINT_RADIUS / 6, 0)
	// }
	for i, _ := range state.CurvePoints {
		if i < len(state.CurvePoints) - 1 {
			imd.Push(state.CurvePoints[i], state.CurvePoints[i + 1])
			imd.Line(POINT_RADIUS / 3)
		}
	}
	// Draw Selected points
	imd.Color = colornames.Darkred
	for i, _ := range state.Points {
		if i < len(state.Points) - 1 {
			imd.Push(state.Points[i], state.Points[i + 1])
			imd.Line(POINT_RADIUS / 6)
		}
	}
	for i, point := range state.Points {
		imd.Color = colornames.Limegreen
		if i == state.DragPointIndex {
			imd.Color = colornames.Darkgreen
		}
		imd.Push(point)
		imd.Circle(POINT_RADIUS, 0)
	}
	// Draw T Points
	imd.Color = colornames.Lightgray
	mainTPoints := GetTPoints(state.Points)
	for i, _ := range mainTPoints {
		if i < len(mainTPoints) - 1 {
			imd.Push(mainTPoints[i], mainTPoints[i + 1])
			imd.Line(POINT_RADIUS / 6)
		}
	}
	imd.Color = colornames.White
	for _, point := range mainTPoints {
		imd.Push(point)
		imd.Circle(POINT_RADIUS, 0)
	}
	// Draw Bezier Curve
	imd.Color = colornames.Yellow
	bezierPoints := GetTPoints(mainTPoints)
	for i, _ := range bezierPoints {
		if i < len(bezierPoints) - 1 {
			imd.Push(bezierPoints[i], bezierPoints[i + 1])
			imd.Line(POINT_RADIUS / 6)
		}
	}
	imd.Color = colornames.Yellow
	for _, point := range bezierPoints {
		imd.Push(point)
		imd.Circle(POINT_RADIUS, 0)
	}
	finalPoints := GetTPoints(bezierPoints)
	imd.Color = colornames.Blue
	for _, point := range finalPoints {
		imd.Push(point)
		imd.Circle(POINT_RADIUS, 0)
	}

	imd.Draw(bezierWin)
}

func DrawControls() {
	// Display Point Coordinates
	atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	availableY := USABLE_CONTROL_HEIGHT
	for i := 0; i < MAX_POINTS ; i++ {
		y := float64(USABLE_CONTROL_HEIGHT - (WINDOW_GAP + (i * TEXT_HEIGHT)))
		availableY = y - WINDOW_GAP
		if len(state.Points) > i {
			basicTxt := text.New(pixel.V(WINDOW_GAP, y), atlas)
			basicTxt.Color = colornames.White
			fmt.Fprintln(basicTxt, strconv.Itoa(i) + " " + FormatVec(state.Points[i]))
			basicTxt.Draw(controlWin, pixel.IM.Scaled(basicTxt.Orig, TEXT_SCALE))
		}
	}
	DrawRect(controlWin, colornames.Gray, pixel.V(0, availableY), pixel.V(CONTROL_WINDOW_WIDTH, availableY - 2))
}

func HandleEvents() {
	ManageBezierInteractions()
}

func ManageBezierInteractions() {
	mousePos := bezierWin.MousePosition()
	if bezierWin.JustPressed(pixelgl.MouseButtonLeft) {
		selectedPoint := false
		// Check to see if clicked an exisitng point. If so, set dragging to be true.
		for i, point := range state.Points {
			if CheckCollision(mousePos, point, POINT_RADIUS) {
				state.DragPointIndex = i
				selectedPoint = true
				break
			}
		}
		// Add point if they're available.
		if !selectedPoint && state.PointLength() < MAX_POINTS {
			state.AddPoint(mousePos)
		}
	}
	if bezierWin.JustReleased(pixelgl.MouseButtonLeft) {
		state.DragPointIndex = -1
	}
	if state.DragPointIndex > -1 {
		state.Points[state.DragPointIndex].X = mousePos.X
		state.Points[state.DragPointIndex].Y = mousePos.Y 
	}
}

func CheckCollision(p1 pixel.Vec, p2 pixel.Vec, threshold float64) bool {
	return FindDistance(p1, p2) <= threshold
}

func FindDistance(p1 pixel.Vec, p2 pixel.Vec) float64 {
	yDelta := math.Abs(p2.Y - p1.Y)
	xDelta := math.Abs(p2.X - p1.X)

	return math.Sqrt(math.Pow(yDelta, 2) + math.Pow(xDelta, 2))
}

func GetTPoints(points []pixel.Vec) []pixel.Vec {
	tPoints := make([]pixel.Vec, 0)

	for i, _ := range points {
		if i < len(points) - 1 {
			p1 := points[i]
			p2 := points[i + 1]
			tPercent := float64(state.TValue) / 100.0

			yDelta := p2.Y - p1.Y
			xDelta := p2.X - p1.X
			newPoint := pixel.V(p1.X + (xDelta * tPercent), p1.Y + (yDelta * tPercent))

			tPoints = append(tPoints, newPoint)
		}
	}

	return tPoints
}

func DrawRect(target pixel.Target, c color.Color, p1 pixel.Vec, p2 pixel.Vec) {
	imd := imdraw.New(nil)
  
	imd.Color = c
	imd.Push(p1)
	imd.Push(pixel.V(p1.X, p2.Y))
	imd.Push(p2)
	imd.Push(pixel.V(p2.X, p1.Y))
	imd.Polygon(0)
  
	imd.Draw(target)
}

func FormatVec(vec pixel.Vec) string {
	return fmt.Sprintf("(%f, %f)", vec.X, vec.Y)
}