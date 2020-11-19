package main

/*

Rectangle cheats:
X = width
Y = height

*/

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	_ "net/http/pprof"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type senseHatResponse struct {
	Temperature json.Number `json:"temperature"`
	Pressure    json.Number `json:"pressure"`
	Humidity    json.Number `json:"humidity"`
}

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
	senseHatTexture   *sdl.Texture
	senseRect         *sdl.Rect
	statTexture       *sdl.Texture
	statRect          *sdl.Rect
	daysSinceTexture  *sdl.Texture
	daysSinceRect     *sdl.Rect
	timeTexture       *sdl.Texture
	timeRect          *sdl.Rect
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

func getTimeTexture(renderer *sdl.Renderer) {
	timeSurface := newStringSurface(time.Now().Format("03:04:05PM"))
	timeRect = rectFromString(center, timeSurface, "large")
	timeTexture = newTextureFromSurface(renderer, timeSurface)
	timeSurface.Free()
	//return timeT, timeRect
}

func newStringSurface(s string) *sdl.Surface {
	fontSurface, err := font.RenderUTF8BlendedWrapped(s, fontColor, 1000)
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
		var screenCenterX int32 = screenWidth - surfaceWidth/2
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

func getSenseHatTexture(renderer *sdl.Renderer) {
	var senseData senseHatResponse
	resp, err := http.Get("http://raspberrypi.lan:8000/")
	if err != nil {
		log.Println("Error getting sense hat info:", err)
		//return nil, nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &senseData)
	if err != nil {
		log.Println("Error unmarshaling sense hat JSON:", err)
		//return nil, nil
	}

	senseSurface := newStringSurface(`
Temperature: ` + senseData.Temperature.String() + `
Pressure: ` + senseData.Pressure.String() + `
Humidity: ` + senseData.Humidity.String() + `%
	`)
	senseRect = rectFromString(upLeft, senseSurface, "small")
	senseHatTexture = newTextureFromSurface(renderer, senseSurface)
	senseSurface.Free()
	//return senseT, senseR
}

func getStatTexture(renderer *sdl.Renderer) {
	var ram runtime.MemStats
	runtime.ReadMemStats(&ram)

	ramSurface := newStringSurface(`
Allocated: ` + strconv.FormatUint(bToMb(ram.Alloc), 10) + `MB
System: ` + strconv.FormatUint(bToMb(ram.Sys), 10) + `MB
	`)
	statRect = rectFromString(upRight, ramSurface, "small")
	statTexture = newTextureFromSurface(renderer, ramSurface)
	ramSurface.Free()
	//return ramT, ramR
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}

func getDaysSinceTexture(renderer *sdl.Renderer) (*sdl.Texture, *sdl.Rect) {
	originalDate := time.Date(2019, time.October, 10, 11, 59, 0, 0, time.Local)
	daysSince := strconv.FormatFloat(time.Since(originalDate).Round(time.Hour).Hours()/24, 'f', 0, 64)
	daysSurface := newStringSurface(`Days Since Last Seizure:` + daysSince)
	daysR := rectFromString(lowerCenter, daysSurface, "small")
	daysT := newTextureFromSurface(renderer, daysSurface)
	daysSurface.Free()
	return daysT, daysR
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

		pngSurface, err = img.Load(backgroundImage)
		backgroundTexture, err = renderer.CreateTextureFromSurface(pngSurface)
		pngSurface.Free()

		renderer.Copy(backgroundTexture, nil, fullRect)
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

	/*
		timeSurface := getTimeSurface()
		timeRect = rectFromString(center, timeSurface, "large")
		timeSurface.Free()
	*/

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

	/*
		var newTexture *sdl.Texture
		var newSurface *sdl.Surface
		sdl.Do(func() {
			newSurface = newStringSurface(strconv.Itoa(int(time.Now().Unix())))
			newTexture = newTextureFromSurface(renderer, newSurface)
		})
		newRect := rectFromString(lowerRight, newSurface, "small")
	*/

	done := make(chan bool)

	var timerStarted bool

	for running {

		//timerM.Lock()

		sdl.Do(func() {
			renderer.Clear()

			timeTexture.Destroy()
			getTimeTexture(renderer)

			renderer.Copy(backgroundTexture, nil, fullRect)
			renderer.Copy(timeTexture, nil, timeRect)
			renderer.Copy(senseHatTexture, nil, senseRect)
			renderer.Copy(statTexture, nil, statRect)
			renderer.Copy(daysSinceTexture, nil, daysSinceRect)

			// Destroy textures (not sure if it's needed)
			//timeTexture.Destroy()
			//newTexture.Destroy()

		})

		if !timerStarted {
			timerStarted = true
			go func() {
				for {
					select {
					case <-done:
						return
					case <-senseHatTimer.C:
						sdl.Do(func() {
							senseHatTexture.Destroy()
							getSenseHatTexture(renderer)

							statTexture.Destroy()
							getStatTexture(renderer)

							daysSinceTexture.Destroy()
							getDaysSinceTexture(renderer)
						})
					}
				}
			}()
		}

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

	go func() {
		http.ListenAndServe("0.0.0.0:8080", nil)
	}()

	var exitcode int
	sdl.Main(func() {
		exitcode = run()
	})

	os.Exit(exitcode)
}
