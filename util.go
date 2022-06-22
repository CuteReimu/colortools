package colortools

import (
	"fmt"
	"image"
	"image/color"
	"math"
)

// Equals 判断两个图像是否完全相同
func Equals(img1, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	if bounds1.Dx() != bounds2.Dx() || bounds1.Dy() != bounds2.Dy() {
		return false
	}
	dx, dy := bounds1.Min.X-bounds2.Min.X, bounds1.Min.Y-bounds2.Min.Y
	for i := bounds2.Min.X; i < bounds2.Max.X; i++ {
		for j := bounds2.Min.Y; j < bounds2.Max.Y; j++ {
			r1, g1, b1, a1 := bounds1.At(i+dx, j+dy).RGBA()
			r2, g2, b2, a2 := bounds2.At(i, j).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				return false
			}
		}
	}
	return true
}

// EqualsSub 判断img1在(x, y)位置的子图像与img2是否完全相同
func EqualsSub(img1 image.Image, x int, y int, img2 image.Image) bool {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	if bounds1.Dx()-x < bounds2.Dx() || bounds1.Dy()-y < bounds2.Dy() {
		return false
	}
	dx, dy := bounds1.Min.X+x-bounds2.Min.X, bounds1.Min.Y+y-bounds2.Min.Y
	for i := bounds2.Min.X; i < bounds2.Max.X; i++ {
		for j := bounds2.Min.Y; j < bounds2.Max.Y; j++ {
			r1, g1, b1, a1 := bounds1.At(i+dx, j+dy).RGBA()
			r2, g2, b2, a2 := bounds2.At(i, j).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				return false
			}
		}
	}
	return true
}

// Search 在img1中，从(x, y)位置开始搜索img2图像，返回值是相对于img1.Bounds().Min的相对位置
func Search(img1 image.Image, x int, y int, img2 image.Image) (point image.Point, ok bool) {
	bounds1 := img1.Bounds()
	bounds2 := img2.Bounds()
	width := bounds1.Dx() - bounds2.Dx() + 1
	height := bounds1.Dy() - bounds2.Dy() + 1
	for i := x; i < width; i++ {
		for j := y; j < height; j++ {
			if EqualsSub(img1, i, j, img2) {
				return image.Point{X: i, Y: j}, true
			}
		}
	}
	return image.Point{}, false
}

type screenImage struct {
	bounds image.Rectangle
	imgs   []image.Image
}

func (s *screenImage) ColorModel() color.Model {
	return s.imgs[0].ColorModel()
}

func (s *screenImage) Bounds() image.Rectangle {
	return s.bounds
}

func (s *screenImage) At(x, y int) color.Color {
	var r0, g0, b0, a0 uint32
	for _, img := range s.imgs {
		r, g, b, a := img.At(x, y).RGBA()
		r0 = 65535 - (65535-r)*(65535-r0)/65535
		g0 = 65535 - (65535-g)*(65535-g0)/65535
		b0 = 65535 - (65535-b)*(65535-b0)/65535
		a0 = 65535 - (65535-a)*(65535-a0)/65535
	}
	return color.RGBA{R: uint8(r0 >> 8), G: uint8(g0 >> 8), B: uint8(b0 >> 8), A: uint8(a0 >> 8)}
}

// Screen 滤色
func Screen(imgs ...image.Image) image.Image {
	if len(imgs) == 0 {
		panic("need at least 1 image")
	}
	rect := imgs[0].Bounds()
	for i := 1; i < len(imgs); i++ {
		bounds := imgs[i].Bounds()
		rect.Min.X = max(rect.Min.X, bounds.Min.X)
		rect.Min.Y = max(rect.Min.Y, bounds.Min.Y)
		rect.Max.X = min(rect.Max.X, bounds.Max.X)
		rect.Max.Y = min(rect.Max.Y, bounds.Max.Y)
	}
	if rect.Min.X >= rect.Max.X || rect.Min.Y >= rect.Max.Y {
		panic("the images have no common area")
	}
	return &screenImage{
		bounds: rect,
		imgs:   imgs,
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type rotateImage struct {
	img       image.Image
	clockwise int8
	rect      image.Rectangle
}

func (img *rotateImage) ColorModel() color.Model {
	return img.img.ColorModel()
}

func (img *rotateImage) Bounds() image.Rectangle {
	return img.rect
}

func (img *rotateImage) At(x, y int) color.Color {
	switch img.clockwise {
	case 0:
		return img.img.At(x, y)
	case 1:
		return img.img.At(y, img.rect.Max.X-x-1)
	case 2:
		return img.img.At(img.rect.Max.X-x-1, img.rect.Max.Y-y-1)
	case 3:
		return img.img.At(img.rect.Max.Y-y-1, x)
	}
	panic("unreachable")
}

// Rotate 旋转图片，clockwise参数为1表示顺时针，2表示180°，-1表示逆时针
func Rotate(img image.Image, clockwise int8) image.Image {
	clockwise %= 4
	if clockwise == 0 {
		return img
	}
	if clockwise < 0 {
		clockwise += 4
	}
	rect := img.Bounds()
	if clockwise%2 == 1 {
		rect.Min.X, rect.Min.Y = rect.Min.Y, rect.Min.X
		rect.Max.X, rect.Max.Y = rect.Max.Y, rect.Max.X
	}
	return &rotateImage{
		img:       img,
		clockwise: clockwise,
		rect:      rect,
	}
}

type zoomGraphicsImage struct {
	img        image.Image
	imgRect    image.Rectangle
	rect       image.Rectangle
	background color.Color
}

func (z *zoomGraphicsImage) ColorModel() color.Model {
	return z.img.ColorModel()
}

func (z *zoomGraphicsImage) Bounds() image.Rectangle {
	return z.rect
}

func (z *zoomGraphicsImage) At(x, y int) color.Color {
	if (image.Point{X: x, Y: y}).In(z.imgRect) {
		return z.img.At(x, y)
	}
	return z.background
}

// ZoomGraphics 放大或缩小画布，img-原图，rect-新画布的位置，background-空余部分的背景色
func ZoomGraphics(img image.Image, rect image.Rectangle, background color.Color) image.Image {
	return &zoomGraphicsImage{
		img:        img,
		imgRect:    img.Bounds(),
		rect:       rect,
		background: background,
	}
}

type mosaicImage struct {
	img   image.Image
	rect  image.Rectangle
	px    int
	cache [][]color.RGBA
}

func (m *mosaicImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (m *mosaicImage) Bounds() image.Rectangle {
	return m.img.Bounds()
}

func (m *mosaicImage) At(x, y int) color.Color {
	if !(image.Point{X: x, Y: y}).In(m.rect) {
		return m.img.At(x, y)
	}
	return m.cache[(x-m.rect.Min.X)/m.px][(y-m.rect.Min.Y)/m.px]
}

// Mosaic 马赛克，img-原图，rect-马赛克的区域，px-马赛克的像素
func Mosaic(img image.Image, rect image.Rectangle, px int) image.Image {
	if px < 0 {
		panic(fmt.Sprint("illegal param: ", px, " px"))
	} else if px <= 1 {
		return img
	}
	cache := make([][]color.RGBA, (rect.Dx()-1)/px+1)
	for i := 0; i*px < rect.Dx(); i++ {
		cache[i] = make([]color.RGBA, (rect.Dy()-1)/px+1)
		for j := 0; j*px < rect.Dy(); j++ {
			var r, g, b, a uint32
			count := 0
			for xx := 0; xx < px; xx++ {
				xi := i*px + xx
				if xi >= rect.Max.X {
					break
				}
				for yy := 0; yy < px; yy++ {
					yi := j*px + yy
					if yi >= rect.Max.Y {
						break
					}
					r1, g1, b1, a1 := img.At(xi, yi).RGBA()
					r += r1 >> 8
					g += g1 >> 8
					b += b1 >> 8
					a += a1 >> 8
					count++
				}
			}
			cache[i][j] = color.RGBA{
				R: uint8(math.Round(float64(r) / float64(count))),
				G: uint8(math.Round(float64(g) / float64(count))),
				B: uint8(math.Round(float64(b) / float64(count))),
				A: uint8(math.Round(float64(a) / float64(count))),
			}
		}
	}
	return &mosaicImage{
		img:   img,
		rect:  rect,
		px:    px,
		cache: cache,
	}
}
