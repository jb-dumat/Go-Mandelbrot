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
	rMax   = .5
	iMin   = -1.
	iMax   = 1.
	width  = 750
	red    = 230
	green  = 235
	blue   = 255
)

// Tile contains the mandelbrot computation for one pixel at (x, y)
type Tile struct {
	r    float64
	x, y int
}

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
 * The worker Pool
 * Iterates over all tiles (which are being processing in the meanwhile)
 * to set the RGBA into the buffer.
 */
func worker(scale float64, b *image.NRGBA, tiles chan Tile, wg *sync.WaitGroup) {
	for tile := range tiles {
		b.Set(tile.x, tile.y, color.NRGBA{uint8(red * tile.r),
			uint8(green * tile.r), uint8(blue * tile.r), 255})
		wg.Done()
	}
}

/*
 * The main function
 */
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // Set the number of threads according to hardwarecm

	scale := width / (rMax - rMin)
	height := int(scale * (iMax - iMin))
	bounds := image.Rect(0, 0, width, height)
	b := image.NewNRGBA(bounds)
	draw.Draw(b, bounds, image.NewUniform(color.Black), image.ZP, draw.Src)

	var tiles = make(chan Tile, 100)
	var wg sync.WaitGroup

	// We launch N worker
	for t := 0; t < runtime.NumCPU(); t++ {
		go worker(scale, b, tiles, &wg)
	}

	// This must appears before the next statement
	wg.Add(width * height)

	// The main computation
	go func() {
		for x := 0; x < width; x++ {
			for y := 0; y < height; y++ {
				fEsc := mandelbrot(complex(
					float64(x)/scale+rMin,
					float64(y)/scale+iMin))
				tiles <- Tile{fEsc, x, y}
			}
		}
		close(tiles)
	}()

	// Wait for worker to end processing all goroutines.
	wg.Wait()

	f, err := os.Create("mandelbrot-pool.png")
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
