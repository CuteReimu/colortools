package colortools

import (
	"image"
	"image/color"
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

// Screen 滤色
func Screen(imgs ...image.Image) *image.RGBA {
	maxSize := imgs[0].Bounds()
	for i := 1; i < len(imgs); i++ {
		bounds := imgs[i].Bounds()
		maxSize.Min.X = max(maxSize.Min.X, bounds.Min.X)
		maxSize.Min.Y = max(maxSize.Min.Y, bounds.Min.Y)
		maxSize.Max.X = min(maxSize.Max.X, bounds.Max.X)
		maxSize.Max.Y = min(maxSize.Max.Y, bounds.Max.Y)
	}
	img0 := image.NewRGBA(image.Rect(maxSize.Min.X, maxSize.Min.Y, maxSize.Max.X, maxSize.Max.Y))
	x0 := maxSize.Min.X
	x1 := maxSize.Max.X
	y0 := maxSize.Min.Y
	y1 := maxSize.Max.Y
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			var r0, g0, b0, a0 uint32
			for _, img := range imgs {
				r, g, b, a := img.At(x, y).RGBA()
				r0 = 65535 - (65535-r)*(65535-r0)/65535
				g0 = 65535 - (65535-g)*(65535-g0)/65535
				b0 = 65535 - (65535-b)*(65535-b0)/65535
				a0 = 65535 - (65535-a)*(65535-a0)/65535
			}
			img0.SetRGBA(x, y, color.RGBA{R: uint8(r0 >> 8), G: uint8(g0 >> 8), B: uint8(b0 >> 8), A: uint8(a0 >> 8)})
		}
	}
	return img0
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
	img2 := image.NewRGBA(rect)
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			c := color.RGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: uint8(a >> 8)}
			switch clockwise {
			case 1:
				img2.SetRGBA(rect.Max.X-j-1, i, c)
			case 2:
				img2.SetRGBA(rect.Max.X-i-1, rect.Max.Y-j-1, c)
			case 3:
				img2.SetRGBA(j, rect.Max.Y-i-1, c)
			}
		}
	}
	return img2
}
