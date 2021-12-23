package colortools

import (
	"image"
	"image/color"
)

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