package main

import (
	"encoding/json"
	"math"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var scale float32

type WindowAnimations string

const (
	SINE_WAVE string = "sine wave"
	JUMPING   string = "jumping"
	STATIC    string = ""
)

type Animation struct {
	SpritesheetName      string  `json:"SpritesheetName"`
	Amount               int     `json:"Amount"`
	FrameHeight          int     `json:"FrameHeight"`
	FrameWidth           int     `json:"FrameWidth"`
	LengthSecs           float64 `json:"LengthSecs"`
	WindowAnimation      string  `json:"WindowAnimation"`
	WindowAnimationFlags float32 `json:"WindowAnimationFlags"`
	Loops                int     `json:"Loops"`
	clockSecs            float64
	spritesheet          rl.Texture2D
}

type AnimationFile struct {
	Name     string      `json:"Name"`
	Snippets []Animation `json:"Snippets"`
}

func DrawAnimation(animation Animation) {
	currentFrame := int(animation.clockSecs / animation.LengthSecs * float64(animation.Amount))
	sourceRect := rl.NewRectangle(
		float32(
			currentFrame*int(animation.FrameWidth)%int(animation.spritesheet.Width),
		),
		float32(
			int(currentFrame*int(animation.FrameWidth)/int(animation.spritesheet.Width))*int(animation.FrameHeight),
		),
		float32(animation.FrameWidth),
		float32(animation.FrameHeight),
	)
	rl.DrawTexturePro(
		animation.spritesheet,
		sourceRect,
		rl.NewRectangle(0, 0, float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())),
		rl.NewVector2(0, 0),
		0,
		rl.RayWhite,
	)
}

func MovingEvent() bool {
	return rl.IsMouseButtonDown(rl.MouseButtonLeft) && rl.IsWindowFocused()
}

var path string = "animation.json"
var CurrentSnippet Animation

var WIN_X, WIN_Y float32

func main() {
	animation := AnimationFile{}

	file, _ := os.ReadFile(path)
	json.Unmarshal(file, &animation)

	str, _ := json.Marshal(animation)
	println(string(str))
	rl.SetConfigFlags(rl.FlagWindowTransparent)

	rl.InitWindow(100, 100, animation.Name)
	rl.SetWindowState(rl.FlagWindowUndecorated | rl.FlagWindowTopmost)

	for i := range animation.Snippets {
		animation.Snippets[i].spritesheet = rl.LoadTexture(animation.Snippets[i].SpritesheetName)
		rl.SetTextureFilter(animation.Snippets[i].spritesheet, rl.FilterBilinear)
	}

	if len(animation.Snippets) == 0 {
		panic("Invalid Animation JSON, no valid snippets detected!")
	}

	CurrentSnippet = animation.Snippets[0]

	scale = float32(rl.GetMonitorHeight(rl.GetCurrentMonitor())) / 2000

	rl.SetWindowSize(int(float32(CurrentSnippet.FrameWidth)*scale), int(float32(CurrentSnippet.FrameHeight)*scale))
	defer rl.CloseWindow()

	rl.SetTargetFPS(144)

	WIN_X = float32(rl.GetWindowPosition().X)
	WIN_Y = float32(rl.GetWindowPosition().Y)

	lastMouse := rl.Vector2{X: 0, Y: 0}

	rl.InitAudioDevice()

	theme := rl.LoadMusicStream("theme.mp3")

	theme.Looping = true

	rl.PlayMusicStream(theme)

	AnimationNumber := 0
	for i := 0; !rl.WindowShouldClose(); {
		lastMouse = rl.GetMousePosition()
		rl.BeginDrawing()
		rl.ClearBackground(rl.Blank)

		if i > CurrentSnippet.Loops+1 {
			CurrentSnippet.clockSecs = 0
			AnimationNumber++
			AnimationNumber %= len(animation.Snippets)
			CurrentSnippet = animation.Snippets[AnimationNumber]
			CurrentSnippet.clockSecs = 0
			//rl.SetWindowSize(int(float32(CurrentSnippet.FrameWidth)*scale), int(float32(CurrentSnippet.FrameHeight)*scale))
			i = 0
		}

		DrawAnimation(CurrentSnippet)

		rl.EndDrawing()

		if !MovingEvent() {
			CurrentSnippet.clockSecs += float64(rl.GetFrameTime())
		}

		if CurrentSnippet.clockSecs > CurrentSnippet.LengthSecs {
			CurrentSnippet.clockSecs = math.Mod(CurrentSnippet.clockSecs, CurrentSnippet.LengthSecs)
			i++
		}

		if MovingEvent() {
			WIN_X += (rl.GetMousePosition().X - lastMouse.X)
			WIN_Y += rl.GetMousePosition().Y - lastMouse.Y
			rl.SetWindowPosition(int(WIN_X), int(WIN_Y))
			rl.UpdateMusicStream(theme)
			continue
		}

		switch CurrentSnippet.WindowAnimation {
		case SINE_WAVE:
			WINANI_SineWave()
		case JUMPING:
			WINANI_Jumping()
		case STATIC:
			fallthrough
		default:
			WINANI_Static()
		}

		rl.UpdateMusicStream(theme)
	}

	for i := range animation.Snippets {
		rl.UnloadTexture(animation.Snippets[i].spritesheet)
	}

}

func WINANI_SineWave() {
	offset := float32(math.Sin(CurrentSnippet.clockSecs/CurrentSnippet.LengthSecs*math.Pi*4)) * scale * CurrentSnippet.WindowAnimationFlags
	rl.SetWindowPosition(int(WIN_X), int(WIN_Y+offset))
}

func WINANI_Jumping() {
	offset := float32(math.Sin(CurrentSnippet.clockSecs/CurrentSnippet.LengthSecs*math.Pi*2)) * scale * CurrentSnippet.WindowAnimationFlags
	rl.SetWindowPosition(int(WIN_X), int(WIN_Y+offset))
}

func WINANI_Static() {
	rl.SetWindowPosition(int(WIN_X), int(WIN_Y))
}
