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
	font            *ttf.Font
)

func getTimeSurface() *sdl.Surface {
	fontSurface, err := font.RenderUTF8Solid(time.Now().Format("03:04:05PM"), sdl.Color{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	})
	if err != nil {
		log.Fatalln("can't render font:", err)
	}
	return fontSurface
}

func getTimeTextureFromSurface(renderer *sdl.Renderer, surface *sdl.Surface) *sdl.Texture {
	fontText, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalln("Error creating font texture:", err)
	}
	surface.Free()
	return fontText
}

func getTimeTexture(renderer *sdl.Renderer) *sdl.Texture {
	surface := getTimeSurface()
	fontText, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalln("Error creating font texture:", err)
	}
	surface.Free()
	return fontText
}

func run() (err error) {
	var window *sdl.Window
	//var surface *sdl.Surface
	var pngImage *sdl.Surface
	//var otherImage *sdl.Surface

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

	// Initial surfaces to get size
	timeSurface := getTimeSurface()

	// Calculate perfect center, based on font and screen size
	var screenCenterY int32 = screenHeight/2 - timeSurface.H/2
	var screenCenterX int32 = screenWidth/2 - timeSurface.W/2

	// Define sctions of the screen
	fullRect := &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight}
	centerRect := &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: timeSurface.W, H: timeSurface.H}

	// Font surface to texture
	timeTexture := getTimeTextureFromSurface(renderer, timeSurface)

	// Load a PNG image
	if pngImage, err = img.Load(backgroundImage); err != nil {
		log.Println(err)
		return err
	}

	imageTexture, err := renderer.CreateTextureFromSurface(pngImage)
	if err != nil {
		log.Println(err)
		return err
	}
	pngImage.Free()

	// Copy background to the full screen, then text on top, to the center
	renderer.Copy(imageTexture, nil, fullRect)
	renderer.Copy(timeTexture, nil, centerRect)

	renderer.Present()

	window.Show()

	// Run infinite loop until user closes the window
	running := true

	var touched bool

	for running {

		renderer.Clear()

		timeTexture = getTimeTexture(renderer)

		renderer.Copy(imageTexture, nil, fullRect)
		renderer.Copy(timeTexture, nil, centerRect)
		timeTexture.Destroy()
		renderer.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.TouchFingerEvent:
				if t.Type == sdl.FINGERUP {
					touched = false
				}
				if t.Type == sdl.FINGERMOTION {
					//log.Println("motion detected", t.DY)

					// Calculate a rough estimate of the touch point
					w, h := window.GetSize()
					xTouch := float32(w) * t.X
					yTouch := float32(h) * t.Y
					touchRect := &sdl.Rect{X: int32(xTouch), Y: int32(yTouch), W: 1, H: 1}

					// Check if our touch intersects with the unlock button
					if touchRect.HasIntersection(centerRect) {

						// Set this to true so the next check passes
						touched = true
					}

					// If unlock button was touched before and is being dragged upwards, perform the dragging
					if touched && float32(h)*t.DY < 0 {
						//log.Println("being dragged")
						//log.Println(textY)

						var textY int32 = int32(yTouch)
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

	// Initialize font here, set to global font var
	var err error
	err = ttf.Init()
	if err != nil {
		log.Fatalln("Error iniitializing font lib")
	}
	font, err = ttf.OpenFont(fontFile, 64)
	if err != nil {
		log.Fatalln("can't open font:", err)
	}

	if err := run(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
