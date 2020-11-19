package main

/*

Rectangle cheats:
X = width
Y = height

*/

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var (
	backgroundImage string
	fontFile        string
	screenHeight    int32 = 480
	screenWidth     int32 = 800
	font            *ttf.Font
	fontColor       = sdl.Color{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}
	backgroundTexture *sdl.Texture
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
	fontSurface, err := font.RenderUTF8Solid(s, fontColor)
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
func rectFromString(pos string, newSurface *sdl.Surface, size string) *sdl.Rect {
	var rect *sdl.Rect

	surfaceHeight := newSurface.H
	surfaceWidth := newSurface.W

	if size == "large" {
		surfaceHeight = surfaceHeight * 2
		surfaceWidth = surfaceWidth * 2
	}

	if size == "small" {
		surfaceHeight = surfaceHeight / 2
		surfaceWidth = surfaceWidth / 2
	}

	switch pos {
	case center:
		var screenCenterY int32 = screenHeight/2 - surfaceHeight/2
		var screenCenterX int32 = screenWidth/2 - surfaceWidth/2
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: surfaceWidth, H: surfaceHeight}
	case centerLeft:
		var screenCenterY int32 = screenHeight/2 - surfaceHeight/2
		rect = &sdl.Rect{X: 0, Y: screenCenterY, W: surfaceWidth, H: surfaceHeight}
	case centerRight:
		var screenCenterY int32 = screenHeight/2 - surfaceHeight/2
		var screenCenterX int32 = screenWidth - surfaceWidth
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: surfaceWidth, H: surfaceHeight}
	case upCenter:
		var screenCenterX int32 = screenWidth/2 - surfaceWidth/2
		rect = &sdl.Rect{X: screenCenterX, Y: 0, W: surfaceWidth, H: surfaceHeight}
	case upLeft:
		rect = &sdl.Rect{X: 0, Y: 0, W: surfaceWidth, H: surfaceHeight}
	case upRight:
		var screenCenterX int32 = screenWidth - surfaceWidth
		rect = &sdl.Rect{X: screenCenterX, Y: 0, W: surfaceWidth, H: surfaceHeight}
	case lowerCenter:
		var screenCenterX int32 = screenWidth/2 - surfaceWidth/2
		var screenCenterY int32 = screenHeight - surfaceHeight
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: surfaceWidth, H: surfaceHeight}
	case lowerLeft:
		var screenCenterY int32 = screenHeight - surfaceHeight
		rect = &sdl.Rect{X: 0, Y: screenCenterY, W: surfaceWidth, H: surfaceHeight}
	case lowerRight:
		var screenCenterX int32 = screenWidth - surfaceWidth
		var screenCenterY int32 = screenHeight - surfaceHeight
		rect = &sdl.Rect{X: screenCenterX, Y: screenCenterY, W: surfaceWidth, H: surfaceHeight}
	default:
		rect = &sdl.Rect{X: 0, Y: 0, W: surfaceWidth, H: surfaceHeight}
	}

	return rect
}

func run() int {
	var window *sdl.Window
	var renderer *sdl.Renderer
	var backgroundTexture *sdl.Texture
	var err error

	fullRect := &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight}

	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Println(err)
		return 1
	}

	defer func() {
		sdl.Do(func() {
			sdl.Quit()
		})
	}()

	/*
		displayRect, err := sdl.GetDisplayBounds(0)
		if err != nil {
			log.Fatalln("Error creating window/renderer:", err)
		}
	*/

	sdl.Do(func() {
		window, err = sdl.CreateWindow("SDL Clock", 0, 0, screenWidth, screenHeight, sdl.WINDOW_FULLSCREEN)

	})

	if err != nil {
		log.Println("Error creating window:", err)
		return 1
	}

	sdl.Do(func() {
		renderer, err = sdl.CreateRenderer(window, -1, 0)

	})

	if err != nil {
		log.Println("Error creating renderer:", err)
		return 1
	}

	sdl.Do(func() {
		var pngSurface *sdl.Surface
		var imageTexture *sdl.Texture

		pngSurface, err = img.Load(backgroundImage)
		imageTexture, err = renderer.CreateTextureFromSurface(pngSurface)
		pngSurface.Free()

		/// Try at combining textures
		tempSurface := newStringSurface("TEMPERATURE: 60f")
		textTempRect := rectFromString(lowerLeft, tempSurface, "small")
		tempTexture := newTextureFromSurface(renderer, tempSurface)

		// Create a background texture to paint the background image, static text, and eventually time onto
		backgroundTexture, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_TARGET, screenWidth, screenHeight)

		// Paint static items onto the background texture
		// With assistance from: https://stackoverflow.com/questions/40886350/how-to-connect-multiple-textures-in-the-one-in-sdl2
		renderer.SetRenderTarget(backgroundTexture)
		renderer.Copy(imageTexture, nil, fullRect)
		renderer.Copy(tempTexture, nil, textTempRect)
		renderer.SetRenderTarget(nil)
		renderer.Present()
	})

	if err != nil {
		log.Println("Error creating backgroundTexture:", err)
		return 1
	}

	// Run infinite loop until user closes the window
	running := true

	//var touched bool

	// Define or calculate all the rectancles used to render
	// fullRect is the full size of the screen

	timeRect := rectFromString(center, getTimeSurface(), "large")

	defer func() {
		sdl.Do(func() {
			renderer.Destroy()
		})
	}()

	sdl.Do(func() {
		window.Show()
		sdl.ShowCursor(sdl.DISABLE)
	})

	senseHatTimer := time.NewTicker(5 * time.Second)
	//var timerM sync.Mutex

	//omg := "yeah"
	newSurface := newStringSurface(strconv.Itoa(int(time.Now().Unix())))
	newRect := rectFromString(lowerRight, newSurface, "small")
	newTexture := newTextureFromSurface(renderer, newSurface)

	done := make(chan bool)

	for running {

		//timerM.Lock()

		sdl.Do(func() {
			renderer.Clear()

			timeTexture := getTimeTexture(renderer)

			renderer.Copy(backgroundTexture, nil, fullRect)
			renderer.Copy(timeTexture, nil, timeRect)
			renderer.Copy(newTexture, nil, newRect)

			// Destroy textures (not sure if it's needed)
			timeTexture.Destroy()
			//newTexture.Destroy()

		})

		go func() {
			for {
				select {
				case <-done:
					return
				case t := <-senseHatTimer.C:
					sdl.Do(func() {
						newTexture = newStringTexture(strconv.Itoa(int(t.Unix())), renderer)
					})
					fmt.Println("Tick at", t)

				}
			}
		}()

		sdl.Do(func() {
			renderer.Present()
			sdl.Delay(100)
		})
	}

	return 0
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

	var exitcode int
	sdl.Main(func() {
		exitcode = run()
	})

	os.Exit(exitcode)
}
