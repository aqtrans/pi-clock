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

func run() (err error) {
	fullRect := &sdl.Rect{X: 0, Y: 0, W: screenWidth, H: screenHeight}

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

	window, err := sdl.CreateWindow("SDL Clock", 0, 0, screenWidth, screenHeight, sdl.WINDOW_FULLSCREEN)
	if err != nil {
		log.Fatalln("Error creating window:", err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, 0)
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
	textTempRect := rectFromString(lowerLeft, tempSurface, "small")
	tempTexture := newTextureFromSurface(renderer, tempSurface)

	// Create a background texture to paint the background image, static text, and eventually time onto
	backgroundTexture, err = renderer.CreateTexture(sdl.PIXELFORMAT_RGB24, sdl.TEXTUREACCESS_TARGET, screenWidth, screenHeight)
	if err != nil {
		log.Fatalln("Error creating backgroundTexture:", err)
	}

	// Paint static items onto the background texture
	// With assistance from: https://stackoverflow.com/questions/40886350/how-to-connect-multiple-textures-in-the-one-in-sdl2
	renderer.SetRenderTarget(backgroundTexture)
	renderer.Copy(imageTexture, nil, fullRect)
	renderer.Copy(tempTexture, nil, textTempRect)
	renderer.SetRenderTarget(nil)
	renderer.Present()

	// Run infinite loop until user closes the window
	running := true

	//var touched bool

	// Define or calculate all the rectancles used to render
	// fullRect is the full size of the screen

	timeRect := rectFromString(center, getTimeSurface(), "large")

	window.Show()

	sdl.ShowCursor(sdl.DISABLE)

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

		sdl.Delay(100)
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
