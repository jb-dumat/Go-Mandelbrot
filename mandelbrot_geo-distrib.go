package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/cmplx"
	"os"
	"runtime"
	"sync"
)

const (
	maxEsc = 100
	rMin   = -2.
	rMax   = 1.
	iMin   = -1.
	iMax   = 1.
	width  = 750
	red    = 230
	green  = 235
	blue   = 255
)

/*
 * The mandelbrot function
 */
func mandelbrot(a complex128) float64 {
	i := 0
	for z := a; cmplx.Abs(z) < 2 && i < maxEsc; i++ {
		z = z*z + a
	}
	return float64(maxEsc-i) / maxEsc
}

/*
 * The main function
 */
func main() {
	maxProc := runtime.NumCPU()
	runtime.GOMAXPROCS(maxProc)

	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	bounds := image.Rect(0, 0, width, height)

	b := image.NewNRGBA(bounds)
	draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

	var wg = new(sync.WaitGroup)

	hIt := height / maxProc
	wIt := width / maxProc

	// We calculate the size of each subrectangle then we launch a goroutine.
	for w := 0; w < width; w += wIt {
		h := 0
		for ; h < height; h += hIt {
			wg.Add(1)

			// We launch the processing of the sub-rectangle
			go func(xStart int, yStart int, xEnd int, yEnd int) {
				defer wg.Done()
				for x := xStart; x < xEnd; x++ {
					for y := yStart; y < yEnd; y++ {
						fEsc := mandelbrot(complex(
							float64(x)/scale+rMin,
							float64(y)/scale+iMin))
						b.Set(x, y, color.NRGBA{uint8(red * fEsc),
							uint8(green * fEsc), uint8(blue * fEsc), 255})
					}
				}
			}(w, h, w+wIt, h+hIt)
		}
	}

	wg.Wait()

	f, err := os.Create("mandelbrot-geo.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = png.Encode(f, b); err != nil {
		fmt.Println(err)
	}
	if err = f.Close(); err != nil {
		fmt.Println(err)
	}
}
