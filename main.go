package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Position struct {
	x, y int
}
type Size struct {
	width, height int
}
type Snack struct {
	ar        []Position
	apel      Position
	direction string
	score     int
}
type Mode struct {
	line bool
}

const (
	fontPath  = "./assets/font/Roboto-Black.ttf"
	fontSize  = 40
	iconImage = "./assets/images/icon.png"
)

func main() {
	snackSize := 25
	gameSpeed := 100
	snackCount := 5
	var mode Mode = Mode{line: false}

	winWidth := 800
	winHeight := 800

	LEFT_BUTTON := 4
	RIGHT_BUTTON := 7
	TOP_BUTTON := 26
	BOTTOM_BUTTON := 22
	running := true
	tick := false

	var gameOver bool = false
	rand.Seed(time.Now().UnixNano())

	appelPosition := Position{
		x: int(math.Round(float64(rand.Intn(winWidth/snackSize-1) * snackSize))),
		y: int(math.Round(float64(rand.Intn(winHeight/snackSize-1) * snackSize))),
	}

	var snack Snack = Snack{
		ar:        []Position{},
		apel:      appelPosition,
		direction: "r",
		score:     0,
	}

	for i := 0; i < snackCount; i++ {
		snack.ar = append(snack.ar, Position{x: 0 - i*snackSize, y: 0})
	}

	if ttf.Init() != nil {
		fmt.Println("error to init ttf")
		return
	}
	defer ttf.Quit()

	surfaceIconImage, err := img.Load(iconImage)
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println(err)
		return
	}
	defer surfaceIconImage.Free()

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("snake game", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println(err)
		return
	}

	window.SetIcon(surfaceIconImage)
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	font, err := ttf.OpenFont(fontPath, fontSize)

	if err != nil {
		fmt.Println(err)
		return
	}
	defer font.Close()

	for running {
		if gameOver {
			renderer.SetDrawColor(0, 0, 0, 255)
			renderer.Clear()

			gameOverText, err := font.RenderUTF8Blended("game over", sdl.Color{R: 255, G: 0, B: 0, A: 255})
			if err != nil {
				fmt.Println(err)
				return
			}
			gameOverTexture, err := renderer.CreateTextureFromSurface(gameOverText)
			if err != nil {
				fmt.Println(err)
				return
			}

			renderer.Copy(gameOverTexture, nil, &sdl.Rect{X: (int32(winWidth / 2)) - (gameOverText.W / 2), Y: (int32(winHeight / 2)) - (gameOverText.H / 2), W: gameOverText.W, H: gameOverText.H})
			renderer.Present()

			gameOverText.Free()
			gameOverTexture.Destroy()

		} else {
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.Clear()
			if mode.line {
				verticalLineCount := (winHeight / snackSize)
				horizontalLineCount := (winWidth / snackSize)

				renderer.SetDrawColor(0, 0, 0, 255)
				for i := 0; i < verticalLineCount; i++ {
					renderer.DrawLine(0, int32(i*snackSize), int32(winHeight), int32(i*snackSize))
				}
				for i := 0; i < horizontalLineCount; i++ {
					renderer.DrawLine(int32(snackSize*i), 0, int32(i*snackSize), int32(winWidth))
				}
			}
			var scoreString string = strconv.FormatInt(int64(snack.score), 10)
			textSurface, err := font.RenderUTF8Blended(string("score:"+scoreString), sdl.Color{R: 0, G: 0, B: 0, A: 255})
			if err != nil {
				fmt.Println(err)
				return
			}
			texture, err := renderer.CreateTextureFromSurface(textSurface)
			if err != nil {
				fmt.Println(err)
				return
			}
			snakeHead := snack.ar[0]

			// check if this apel is eat
			if snakeHead.x == snack.apel.x && snakeHead.y == snack.apel.y {
				snack.score++
				appelPosition := Position{
					x: int(math.Round(float64(rand.Intn(winWidth/snackSize-1) * snackSize))),
					y: int(math.Round(float64(rand.Intn(winHeight/snackSize-1) * snackSize))),
				}
				snack.apel = appelPosition
			} else {
				snack.ar = snack.ar[:len(snack.ar)-1]
			}

			Apelrect := sdl.Rect{X: int32(snack.apel.x), Y: int32(snack.apel.y), W: int32(snackSize), H: int32(snackSize)}
			renderer.SetDrawColor(255, 0, 0, 255)
			renderer.FillRect(&Apelrect)

			firstSnake := snack.ar[0]
			if firstSnake.x >= (winWidth) || firstSnake.x <= -snackSize {
				gameOver = true
			} else if firstSnake.y >= (winHeight) || firstSnake.y <= -snackSize {
				gameOver = true
			}

			for i := 1; i < len(snack.ar); i++ {
				if firstSnake.x == snack.ar[i].x && firstSnake.y == snack.ar[i].y {
					gameOver = true
				}
			}

			// to update this direction
			if snack.direction == "l" {
				newPosition := Position{x: firstSnake.x - snackSize, y: firstSnake.y}
				snack.ar = append([]Position{newPosition}, snack.ar...)
			} else if snack.direction == "r" {
				newPosition := Position{x: firstSnake.x + snackSize, y: firstSnake.y}
				snack.ar = append([]Position{newPosition}, snack.ar...)
			} else if snack.direction == "t" {
				newPosition := Position{x: firstSnake.x, y: firstSnake.y - snackSize}
				snack.ar = append([]Position{newPosition}, snack.ar...)
			} else if snack.direction == "b" {
				newPosition := Position{x: firstSnake.x, y: firstSnake.y + snackSize}
				snack.ar = append([]Position{newPosition}, snack.ar...)
			}
			for i := 0; i < len(snack.ar); i++ {
				rect := sdl.Rect{X: int32(snack.ar[i].x), Y: int32(snack.ar[i].y), W: int32(snackSize), H: int32(snackSize)}
				if i != 0 {
					renderer.SetDrawColor(0, 0, 0, 255)
				} else {
					renderer.SetDrawColor(0, 255, 0, 255)
				}
				renderer.FillRect(&rect)
			}
			renderer.Copy(texture, nil, &sdl.Rect{X: int32(winWidth - 150), Y: 10, W: textSurface.W, H: textSurface.H})
			renderer.Present()

			texture.Destroy()
			textSurface.Free()
		}

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Scancode
				if !tick {
					tick = true

					if keyCode == sdl.Scancode(TOP_BUTTON) && snack.direction != "b" {
						snack.direction = "t"
					} else if keyCode == sdl.Scancode(RIGHT_BUTTON) && snack.direction != "l" {
						snack.direction = "r"
					} else if keyCode == sdl.Scancode(LEFT_BUTTON) && snack.direction != "r" {
						snack.direction = "l"
					} else if keyCode == sdl.Scancode(BOTTOM_BUTTON) && snack.direction != "t" {
						snack.direction = "b"
					} else if keyCode == 6 {
						mode.line = !mode.line
					} else if gameOver && keyCode == 21 {
						appelPosition := Position{
							x: int(math.Round(float64(rand.Intn(winWidth/snackSize-1) * snackSize))),
							y: int(math.Round(float64(rand.Intn(winHeight/snackSize-1) * snackSize))),
						}

						snack = Snack{
							ar:        []Position{},
							apel:      appelPosition,
							direction: "r",
							score:     0,
						}
						for i := 0; i < snackCount; i++ {
							snack.ar = append(snack.ar, Position{x: 0 - i*snackSize, y: 0})
						}
						gameOver = false

					}
				}
				break

			}
		}
		sdl.Delay(uint32(gameSpeed))
		tick = false
	}
}
