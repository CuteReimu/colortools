package colortools

import (
	"image"
	"image/color"
	"math"
)

// 径向渐变
type cervicalGradChgColorImage struct {
	c       []color.Color
	percent []float64
	line    image.Rectangle
	bound   image.Rectangle
	length  float64
}

func NewCervicalGradChgColorImage(bounds image.Rectangle, c []color.Color, percent []float64, line image.Rectangle) image.Image {
	syn := percent[len(percent)-1]
	p := make([]float64, len(percent))
	for i, v := range percent {
		p[i] = v / syn
	}
	lineWidth := float64(line.Max.X - line.Min.X)
	lineHeight := float64(line.Max.Y - line.Min.Y)
	return &cervicalGradChgColorImage{
		c:       c,
		percent: percent,
		line:    line,
		bound:   bounds,
		length:  math.Sqrt(lineWidth*lineWidth + lineHeight*lineHeight),
	}
}

func (c *cervicalGradChgColorImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *cervicalGradChgColorImage) Bounds() image.Rectangle {
	return c.bound
}

func (c *cervicalGradChgColorImage) At(x, y int) color.Color {
	xDist := c.line.Min.X - x
	yDist := c.line.Min.Y - y
	distance := math.Sqrt(float64(xDist*xDist + yDist*yDist))
	per := distance / c.length
	n, found := binarySearch(c.percent, per)
	if found {
		return c.c[n]
	}
	if n == 0 {
		return c.c[0]
	}
	length := len(c.percent)
	if n == length {
		return c.c[length-1]
	}
	per1 := (c.percent[n] - per) / (c.percent[n] - c.percent[n-1])
	per2 := 1.0 - per1
	R1, G1, B1, A1 := c.c[n-1].RGBA()
	R2, G2, B2, A2 := c.c[n].RGBA()
	return &color.RGBA{
		R: uint8((float64(R1)*per1 + float64(R2)*per2) / 256.0),
		G: uint8((float64(G1)*per1 + float64(G2)*per2) / 256.0),
		B: uint8((float64(B1)*per1 + float64(B2)*per2) / 256.0),
		A: uint8((float64(A1)*per1 + float64(A2)*per2) / 256.0),
	}
}

// 线性渐变
type lineGradChgColorImage struct {
	c       []color.Color
	percent []float64
	line    image.Rectangle
	bound   image.Rectangle
	angleX  float64
	angleY  float64
	lineLen float64
}

func NewLineGradChgColorImage(bounds image.Rectangle, c []color.Color, percent []float64, line image.Rectangle) image.Image {
	syn := percent[len(percent)-1]
	p := make([]float64, len(percent))
	for i, v := range percent {
		p[i] = v / syn
	}
	lineWidth := float64(line.Max.X - line.Min.X)
	lineHeight := float64(line.Max.Y - line.Min.Y)
	return &lineGradChgColorImage{
		c:       c,
		percent: percent,
		line:    line,
		bound:   bounds,
		angleX:  lineWidth,
		angleY:  lineHeight,
		lineLen: math.Sqrt(lineWidth*lineWidth + lineHeight*lineHeight),
	}
}

func (l *lineGradChgColorImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (l *lineGradChgColorImage) Bounds() image.Rectangle {
	return l.bound
}

func (l *lineGradChgColorImage) At(x, y int) color.Color {
	x -= l.line.Min.X
	y -= l.line.Min.Y
	if x == 0 && y == 0 {
		return l.c[0]
	}
	per := (float64(y)*l.angleY + float64(x)*l.angleX) / l.lineLen / l.lineLen
	if per < l.percent[0] {
		return l.c[0]
	}
	n, found := binarySearch(l.percent, per)
	if found {
		return l.c[n]
	}
	if n == 0 {
		return l.c[0]
	}
	length := len(l.percent)
	if n == length {
		return l.c[length-1]
	}
	per1 := (l.percent[n] - per) / (l.percent[n] - l.percent[n-1])
	per2 := 1.0 - per1
	R1, G1, B1, A1 := l.c[n-1].RGBA()
	R2, G2, B2, A2 := l.c[n].RGBA()
	return &color.RGBA{
		R: uint8((float64(R1)*per1 + float64(R2)*per2) / 256.0),
		G: uint8((float64(G1)*per1 + float64(G2)*per2) / 256.0),
		B: uint8((float64(B1)*per1 + float64(B2)*per2) / 256.0),
		A: uint8((float64(A1)*per1 + float64(A2)*per2) / 256.0),
	}
}

// 角度渐变
type taperedGradChgColorImage struct {
	c       []color.Color
	percent []float64
	line    image.Rectangle
	bound   image.Rectangle
	initAng float64
}

func NewTaperedGradChgColorImage(bounds image.Rectangle, c []color.Color, percent []float64, line image.Rectangle) image.Image {
	syn := percent[len(percent)-1]
	p := make([]float64, len(percent))
	for i, v := range percent {
		p[i] = v / syn
	}
	lineWidth := float64(line.Max.X - line.Min.X)
	lineHeight := float64(line.Max.Y - line.Min.Y)
	return &taperedGradChgColorImage{
		c:       c,
		percent: percent,
		line:    line,
		bound:   bounds,
		initAng: math.Atan2(lineWidth, lineHeight),
	}
}

func (t *taperedGradChgColorImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (t *taperedGradChgColorImage) Bounds() image.Rectangle {
	return t.bound
}

func (t *taperedGradChgColorImage) At(x, y int) color.Color {
	xDist := x - t.line.Min.X
	yDist := y - t.line.Min.Y
	per := (math.Atan2(float64(xDist), float64(yDist)) - t.initAng) / math.Pi / 2
	if per < 0 {
		per++
	}
	n, found := binarySearch(t.percent, per)
	if found {
		return t.c[n]
	}
	if n == 0 {
		return t.c[0]
	}
	length := len(t.percent)
	if n == length {
		return t.c[length-1]
	}
	per1 := (t.percent[n] - per) / (t.percent[n] - t.percent[n-1])
	per2 := 1.0 - per1
	R1, G1, B1, A1 := t.c[n-1].RGBA()
	R2, G2, B2, A2 := t.c[n].RGBA()
	return &color.RGBA{
		R: uint8((float64(R1)*per1 + float64(R2)*per2) / 256.0),
		G: uint8((float64(G1)*per1 + float64(G2)*per2) / 256.0),
		B: uint8((float64(B1)*per1 + float64(B2)*per2) / 256.0),
		A: uint8((float64(A1)*per1 + float64(A2)*per2) / 256.0),
	}
}

func binarySearch(sortedArray []float64, target float64) (int, bool) {
	low := 0
	height := len(sortedArray) - 1
	for low <= height {
		mid := low + (height-low)/2
		midValue := sortedArray[mid]
		if midValue == target {
			return mid, true
		} else if midValue > target {
			height = mid - 1
		} else {
			low = mid + 1
		}
	}
	return low, false
}
