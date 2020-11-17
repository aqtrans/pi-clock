package main

/*

Rectangle cheats:
X = width
Y = height

*/

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

const (
	upLeft      = "upperLeft"
	upRight     = "upperRight"
	upCenter    = "upperCenter"
	center      = "center"
	centerLeft  = "centerLeft"
	centerRight = "centerRight"
	lowerLeft   = "lowerLeft"
	lowerRight  = "lowerRight"
	lowerCenter = "lowerCenter"
)

// Return only a surface with the current time
// Should only be used for calculating the clock position
func getTimeSurface() *sdl.Surface {
	return newStringSurface(time.Now().Format("03:04:05PM"))
}

func getTimeTexture(renderer *sdl.Renderer) *sdl.Texture {
	surface := getTimeSurface()
	fontText, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalln("Error creating time texture:", err)
	}
	surface.Free()
	return fontText
}

func newStringSurface(s string) *sdl.Surface {
	fontSurface, err := font.RenderUTF8Solid(s, sdl.Color{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	})
	if err != nil {
		log.Fatalln("Error creating string surface:", err)
	}
	return fontSurface
}

func newTextureFromSurface(renderer *sdl.Renderer, surface *sdl.Surface) *sdl.Texture {
	newTexture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalln("Error creating texture from surface:", err)
	}
	surface.Free()
	return newTexture
}

func newStringTexture(s string, renderer *sdl.Renderer) *sdl.Texture {
	surface := newStringSurface(s)
	stringTexture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		log.Fatalln("Error creating string texture:", err)
	}
	surface.Free()
	return stringTexture
}

// Using a given string, make a surface out of it, then create a rectangle using the surface bounds
func rectFromString(pos string, newSurface *sdl.Surface) *sdl.Rect {
	var rect *sdl.Rect

	switch pos {
	case center:
		var screenCenterY int32 = screenHeight/2 - newSurface.H/2
		var screenCenterX int32 = screenWidth/2 - newSurface.W/2
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: newSurface.W, H: newSurface.H}
	case centerLeft:
		var screenCenterY int32 = screenHeight/2 - newSurface.H/2
		rect = &sdl.Rect{X: 0, Y: screenCenterY, W: newSurface.W, H: newSurface.H}
	case centerRight:
		var screenCenterY int32 = screenHeight/2 - newSurface.H/2
		var screenCenterX int32 = screenWidth - newSurface.W
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: newSurface.W, H: newSurface.H}
	case upCenter:
		var screenCenterX int32 = screenWidth/2 - newSurface.W/2
		rect = &sdl.Rect{X: screenCenterX, Y: 0, W: newSurface.W, H: newSurface.H}
	case upLeft:
		rect = &sdl.Rect{X: 0, Y: 0, W: newSurface.W, H: newSurface.H}
	case upRight:
		var screenCenterX int32 = screenWidth - newSurface.W
		rect = &sdl.Rect{X: screenCenterX, Y: 0, W: newSurface.W, H: newSurface.H}
	case lowerCenter:
		var screenCenterX int32 = screenWidth/2 - newSurface.W/2
		var screenCenterY int32 = screenHeight - newSurface.H
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: newSurface.W, H: newSurface.H}
	case lowerLeft:
		var screenCenterY int32 = screenHeight - newSurface.H
		rect = &sdl.Rect{X: 0, Y: screenCenterY, W: newSurface.W, H: newSurface.H}
	case lowerRight:
		var screenCenterX int32 = screenWidth - newSurface.W
		var screenCenterY int32 = screenHeight - newSurface.H
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: newSurface.W, H: newSurface.H}
	default:
		rect = &sdl.Rect{X: 0, Y: 0, W: newSurface.W, H: newSurface.H}
	}

	newSurface.Free()

	return rect
}

func run() (err error) {

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

	window, err := sdl.CreateWindow("SDL Clock", 0, 0, screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalln("Error creating window:", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED+sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		log.Fatalln("Error creating renderer:", err)
	}

	// Load a PNG image
	pngSurface, err := img.Load(backgroundImage)
	if err != nil {
		log.Println("Error loading background image:", err)
		return err
	}

	imageTexture, err := renderer.CreateTextureFromSurface(pngSurface)
	if err != nil {
		log.Println(err)
		return err
	}
	pngSurface.Free()

	/// Try at combining textures
	tempSurface := newStringSurface("TEMPERATURE: 60f")
	textTempRect := &sdl.Rect{X: 0, Y: 0, W: 500, H: 100}
	tempTexture := newTextureFromSurface(renderer, tempSurface)
	tempSurface.Free()

	// Create a background texture to paint the background image, static text, and eventually time onto
	backgroundTexture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_TARGET, screenWidth, screenHeight)
	if err != nil {
		log.Fatalln("Error creating backgroundTexture:", err)
	}

	// Paint static items onto the background texture
	// With assistance from: https://stackoverflow.com/questions/40886350/how-to-connect-multiple-textures-in-the-one-in-sdl2
	renderer.SetRenderTarget(backgroundTexture)
	renderer.Copy(imageTexture, nil, &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight})
	renderer.Copy(tempTexture, nil, textTempRect)
	renderer.SetRenderTarget(nil)
	renderer.Present()

	// Run infinite loop until user closes the window
	running := true

	//var touched bool

	// Define or calculate all the rectancles used to render
	// fullRect is the full size of the screen
	fullRect := &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight}
	timeRect := rectFromString(center, getTimeSurface())

	window.Show()

	for running {

		renderer.Clear()

		timeTexture := getTimeTexture(renderer)
		//tempTexture := newStringTexture("TEMPERATURE: 60f", renderer)

		renderer.Copy(backgroundTexture, nil, fullRect)
		renderer.Copy(timeTexture, nil, timeRect)
		//renderer.Copy(tempTexture, nil, tempRect)
		renderer.Present()

		// Destroy textures (not sure if it's needed)
		//backgroundTexture.Destroy()
		timeTexture.Destroy()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {

			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.TouchFingerEvent:
				if t.Type == sdl.FINGERUP {
					//touched = false
				}
				if t.Type == sdl.FINGERMOTION {
					//log.Println("motion detected", t.DY)

					/* Unused, no touch stuff used yet, but leaving for reference:

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
