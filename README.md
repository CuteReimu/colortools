# Color Tools
![](https://img.shields.io/github/languages/top/CuteReimu/colortools "Language")
[![](https://img.shields.io/github/workflow/status/CuteReimu/colortools/Go)](https://github.com/CuteReimu/colortools/actions/workflows/golangci-lint.yml "Analysis")

Useful color tools.

## Install

```
go get github.com/CuteReimu/colortools
```

## Usage

```go
package main

import (
	"github.com/CuteReimu/colortools"
	"image/color"
	"image/png"
	"os"
)

func main() {
	f, _ := os.Open("1.png")
	defer f.Close()
	img, _ := png.Decode(f)

	c := make([]color.Color, 361)
	p := make([]float64, 361)
	for i := 0; i <= 360; i++ {
		c[i] = &colortools.HSV{H: float64(i), S: 1.0, V: 0.5}
		p[i] = float64(i) / 360.0
	}
	img1 := colortools.NewLineGradChgColorImage(img.Bounds(), c, p, img.Bounds())
	img2 := colortools.Screen(img1, img)
	
	f2, _ := os.Create("2.png")
	defer f2.Close()
	_ = png.Encode(f2, img2)
}
```
