package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	backgroundImage string
	fontFile        string
	screenHeight    int32 = 480
	screenWidth     int32 = 640
)

func run() (err error) {
	var window *sdl.Window
	//var surface *sdl.Surface
	var pngImage *sdl.Surface
	//var otherImage *sdl.Surface

	err = ttf.Init()
	if err != nil {
		log.Fatalln("Error iniitializing font lib")
	}

	if err = sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Println(err)
		return
	}
	defer sdl.Quit()

	/*
		displayRect, err := sdl.GetDisplayBounds(0)
		if err != nil {
			log.Fatalln("Error creating window/renderer:", err)
		}
	*/

	// Create a window for us to draw the images on
	window, renderer, err := sdl.CreateWindowAndRenderer(screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalln("Error creating window/renderer:", err)
	}
	defer window.Destroy()

	// Load fonts
	font, err := ttf.OpenFont(fontFile, 64)
	if err != nil {
		log.Fatalln("can't open font:", err)
	}

	fontSurface, err := font.RenderUTF8Blended("Drag Up", sdl.Color{
		R: 255,
		G: 0,
		B: 0,
		A: 255,
	})
	if err != nil {
		log.Fatalln("can't render font:", err)
	}
	fontText, err := renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		log.Fatalln("Error creating font texture:", err)
	}
	fontSurface.Free()

	// Load a PNG image
	if pngImage, err = img.Load(backgroundImage); err != nil {
		log.Println(err)
		return err
	}
	defer pngImage.Free()

	imageText, err := renderer.CreateTextureFromSurface(pngImage)
	if err != nil {
		log.Println(err)
		return err
	}

	renderer.Copy(imageText, nil, &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight})

	// Calculate perfect center, based on font and screen size
	var textY int32 = screenHeight/2 - fontSurface.H/2
	var textX int32 = screenWidth/2 - fontSurface.W/2

	fontRect := &sdl.Rect{X: textX, Y: textY, W: fontSurface.W, H: fontSurface.H}
	renderer.Copy(fontText, nil, fontRect)
	renderer.Present()

	window.Show()

	// Run infinite loop until user closes the window
	running := true

	var touched bool

	for running {

		renderer.Clear()

		renderer.Copy(imageText, nil, &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight})

		fontSurface, err := font.RenderUTF8Solid(time.Now().Format("3:04:05PM"), sdl.Color{
			R: 0,
			G: 0,
			B: 0,
			A: 255,
		})
		if err != nil {
			log.Fatalln("can't render font:", err)
		}

		fontText, err := renderer.CreateTextureFromSurface(fontSurface)
		if err != nil {
			log.Fatalln("Error creating font texture:", err)
		}
		fontSurface.Free()

		//fontText.Update(fontRect, fontSurface2.Pixels(), int(fontSurface2.Pitch))

		renderer.Copy(fontText, nil, fontRect)
		fontText.Destroy()
		renderer.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.TouchFingerEvent:
				if t.Type == sdl.FINGERUP {
					// Draw original unlock button on finger up
					renderer.Clear()
					renderer.Copy(imageText, nil, &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight})

					fontRect.Y = textY
					renderer.Copy(fontText, nil, fontRect)

					renderer.Present()
					touched = false
				}
				if t.Type == sdl.FINGERMOTION {
					//log.Println("motion detected", t.DY)

					w, h := window.GetSize()
					xTouch := float32(w) * t.X
					yTouch := float32(h) * t.Y
					//log.Println(int32(yTouch), int32(xTouch))

					touchRect := &sdl.Rect{X: int32(xTouch), Y: int32(yTouch), W: 1, H: 1}

					// Check if our touch intersects with the unlock button
					if touchRect.HasIntersection(fontRect) {
						//log.Println("Touch and OMG intersect. Moving OMG!")

						/* Moving should be handled by the function below once dragging begins
						renderer.Clear()
						renderer.Copy(imageText, nil, &sdl.Rect{X: 0, Y: 0, W: 1280, H: 800})

						var textY int32 = int32(yTouch)
						renderer.Copy(fontText, nil, &sdl.Rect{X: 1280 / 2, Y: textY, W: 100, H: 50})
						renderer.Present()
						*/

						// Set this to true so the next check passes
						touched = true
					}

					// If unlock button was touched before and is being dragged upwards, perform the dragging
					if touched && float32(h)*t.DY < 0 {
						//log.Println("being dragged")
						//log.Println(textY)

						renderer.Clear()
						renderer.Copy(imageText, nil, &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight})

						var textY int32 = int32(yTouch)
						fontRect.Y = textY
						renderer.Copy(fontText, nil, fontRect)
						renderer.Present()

						touched = true

						// Once dragged halfway up, exit
						if textY < (screenHeight / 2) {
							err = renderer.Destroy()
							if err != nil {
								log.Fatalln("Error destroying renderer:", err)
							}

							err = window.Destroy()
							if err != nil {
								log.Fatalln("Error destroying window:", err)
							}
							os.Exit(0)
						}
					}

					/*
						renderer.SetDrawColor(255, 255, 255, 0)
						err := renderer.DrawPoint(int32(xTouch), int32(yTouch))
						if err != nil {
							log.Fatalln("can't destroy fonttext:", err)
						}
					*/

					/*
						_, h := window.GetSize()
						//log.Println("X touch:", float32(w)*t.DX)
						//log.Println("Y touch:", float32(h)*t.DY)

						// All X/Y coordinates are 'normalized', need to be multiplied by the resolution
						if (float32(h) * t.DY) < -30.0 {
							//log.Println("quit motion detected. quitting", t.DY)

							//otherImage.BlitScaled(nil, surface, &sdl.Rect{X: 0, Y: 0, W: 1280, H: 800})
							window.UpdateSurface()

							// Slight pause before exiting
							time.Sleep(250 * time.Millisecond)

							err = window.Destroy()
							if err != nil {
								log.Fatalln("Error destroying:", err)
							}
							os.Exit(0)
						}
					*/
				}
				//log.Println("Touch", t.Type, "moved by", t.DX, t.DY, "at", t.X, t.Y)
			// Exit on spacebar
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYUP && t.Keysym.Sym == sdl.GetKeyFromName("Space") {
					err = renderer.Destroy()
					if err != nil {
						log.Fatalln("Error destroying renderer:", err)
					}

					err = window.Destroy()
					if err != nil {
						log.Fatalln("Error destroying window:", err)
					}
					os.Exit(0)
				}
			}
		}

		sdl.Delay(16)
	}

	return
}

func main() {
	flag.StringVar(&backgroundImage, "bg", "bg.png", "Background image. PNG")
	flag.StringVar(&fontFile, "font", "sans.ttf", "Font file. TTF")

	flag.Parse()

	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
