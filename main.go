package main

import (
	"encoding/json"
	"fmt"
	"github.com/SCOTT-HAMILTON/Go-pathfinderalgo/AStar"
	"github.com/SCOTT-HAMILTON/Go-pathfinderalgo/Djikstra"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const nbw = 10
const nbh = 12

const DEBUGMODE = false

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//pixelgl.Run(run)
	os.Exit(run())
}

func saveToMapFile(map_walls []int) {
	b, err := json.Marshal(map_walls)
	print("Marshal init state : ")
	check(err)
	print("b : \n\"", b, "\"\n")
	err = ioutil.WriteFile("map.json", b, 0644)
	check(err)
}

func initMapFile() {
	var map_walls []int
	map_walls = make([]int, nbw*nbh)
	for i := 0; i < nbw*nbh; i++ {
		rand := rand.Intn(100)
		if rand < 90 {
			map_walls[i] = 0
		} else {
			map_walls[i] = 1
		}
	}
	saveToMapFile(map_walls)
}

func loadMap() (map_walls []int) {
	map_walls = make([]int, nbw*nbh)
	data, err := ioutil.ReadFile("map.json")
	print("\ndata first read : \n\"", data, "\"\n")
	if err != nil || len(data) == 0 {
		initMapFile()
		data, err = ioutil.ReadFile("map.json")
		print("\ndata to unmarshal 2nd read : \n\"", data, "\"\n")
		print("second try read state : ")
		check(err)
	}

	print("\ndata to unmarshal : \n\"", data, "\"\n")
	bytes := []byte(data)
	err = json.Unmarshal(bytes, &map_walls)
	print("unmarshal state : ")
	check(err)

	print("bytes : \n\"", bytes, "\"\n")

	return
}

func printD(text string, rect sdl.Rect, font *ttf.Font, render *sdl.Renderer) {
	solid, err := font.RenderUTF8Solid(text, sdl.Color{255, 255, 255, 255})
	check(err)
	defer solid.Free()
	var texture *sdl.Texture
	texture, err = render.CreateTextureFromSurface(solid)
	check(err)
	defer texture.Destroy()
	render.Copy(texture, nil, &rect)
} // debug font

func deleteMapCache(mapwalls *[]int) {
	for i, n := range *mapwalls {
		if n != 1 {
			(*mapwalls)[i] = 0
		}
	}
}

func run() int {

	//initMapFile()
	map_walls := loadMap()

	star := AStar.NewAStar(nbw, nbh, 100, 19, &map_walls)
	djk := Djikstra.NewDjikstra(nbw, nbh, 100, 19, &map_walls)

	star.Init()
	djk.Init()

	const xMargin = 10
	const yMargin = 10
	const widthSpacing = 10
	const heightSpacing = 10
	const width = 55
	const height = 55
	windowWidth := int32(xMargin*2 + width*nbw + (nbw-1)*widthSpacing)
	windowHeight := int32(yMargin*2 + height*nbh + (nbh-1)*heightSpacing)
	mode := 0 //mode for debugging
	const nbMode = 3
	canUpdate := false
	guiMode := 0
	canRenderPath := true
	timer_update := time.Now()

	var window *sdl.Window
	var font *ttf.Font //debug font
	running := true
	var err error

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	err = ttf.Init()
	check(err)
	font, err = ttf.OpenFont("F25_Bank_Printer.ttf", 32)
	check(err)
	defer font.Close() //debug font
	window, err = sdl.CreateWindow("Go PathFinder", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	check(err)
	defer window.Destroy()

	var Render *sdl.Renderer

	Render, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_PRESENTVSYNC)
	check(err)
	Render.SetDrawColor(0, 0, 0, 255)
	defer Render.Destroy()

	renderPath := func(pathColor sdl.Color) {
		cx := int32(xMargin)
		cy := int32(yMargin)
		if canRenderPath {
			for _, n := range star.Path {
				if n == star.End {
					continue
				}
				x, y := star.ToCoord(n)
				x = xMargin + x*width + x*widthSpacing
				y = yMargin + y*height + y*heightSpacing
				rect := sdl.Rect{int32(x), int32(y), width, height}

				Render.SetDrawColor(66, 82, 90, 100)
				Render.FillRect(&rect)
			}

			for _, n := range djk.Path {
				if n == djk.End {
					continue
				}
				x, y := djk.ToCoord(n)
				x = xMargin + x*width + x*widthSpacing
				y = yMargin + y*height + y*heightSpacing
				rect := sdl.Rect{int32(x), int32(y), width, height}

				Render.SetDrawColor(224, 28, 85, 100)
				Render.FillRect(&rect)
			}
		}

		for i := 0; i < nbw*nbh; i++ {
			rect := sdl.Rect{cx, cy, width, height}
			if DEBUGMODE {
				if map_walls[i] == 2 {
					Render.SetDrawColor(255, 255, 0, 255)
					Render.DrawRect(&rect)
				} else if map_walls[i] == 3 {
					Render.SetDrawColor(255, 165, 0, 100)
					Render.FillRect(&rect)
				} else if map_walls[i] == 4 && i != star.Start {
					Render.SetDrawColor(51, 51, 51, 255)
					Render.FillRect(&rect)
				} else {
					Render.SetDrawColor(50, 205, 50, 255)
					Render.DrawRect(&rect)
				}

				if i == star.CurNode.Pos {
					Render.SetDrawColor(128, 0, 0, 255)
					Render.DrawRect(&rect)
				}
			}

			//modes
			if mode == 1 {
				x, y := star.ToCoord(i)
				str := strconv.Itoa(x) + "," + strconv.Itoa(y)
				printD(str, sdl.Rect{cx + 1, cy + 1, width - 2, height - 30}, font, Render)
			} else if mode == 2 {
				node := star.FindNei(i)
				if node.Pos != -1 {
					str := fmt.Sprintf("%.02f", node.GetF())
					printD(str, sdl.Rect{cx + 1, cy + 2, 40, 15}, font, Render)
					str = fmt.Sprintf("%.02f", node.GetG())
					printD(str, sdl.Rect{cx + 1, cy + 18, 40, 15}, font, Render)
					str = fmt.Sprintf("%.02f", node.GetH())
					printD(str, sdl.Rect{cx + 1, cy + 35, 40, 15}, font, Render)
				}
			} //debug font

			if (i+1) != 0 && (i+1)%(nbw) == 0 {
				cy += height + heightSpacing
				cx = xMargin
			} else {
				cx += width + widthSpacing
			}
		}
	}

	renderMap := func(drawCurPos bool) {
		cx := int32(xMargin)
		cy := int32(yMargin)
		for i := 0; i < nbw*nbh; i++ {
			rect := sdl.Rect{cx, cy, width, height}
			Render.DrawRect(&rect)

			if i == star.Start {
				Render.SetDrawColor(0, 0, 255, 255)
				Render.FillRect(&rect)
			}
			if i == star.End {
				Render.SetDrawColor(255, 0, 0, 255)
				Render.FillRect(&rect)
			}

			if map_walls[i] == 1 {
				Render.SetDrawColor(50, 205, 50, 255)
				Render.FillRect(&rect)
			} else {
				Render.SetDrawColor(50, 205, 50, 255)
				Render.DrawRect(&rect)
			}

			if drawCurPos && i == star.CurNode.Pos {
				Render.SetDrawColor(128, 0, 0, 255)
				Render.DrawRect(&rect)
			}

			if (i+1) != 0 && (i+1)%(nbw) == 0 {
				cy += height + heightSpacing
				cx = xMargin
			} else {
				cx += width + widthSpacing
			}
		}
	}

	runMode := func() {
		Render.SetDrawColor(0, 0, 0, 255)
		Render.Clear()

		renderMap(false)
		renderPath(sdl.Color{200, 200, 200, 10})
		Render.Present()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONUP {
					canUpdate = !canUpdate
					if !canRenderPath {
						canRenderPath = true
					}
					print("start pos : ", star.Start, ", end pos : ", star.End, "\n")
				}
				break
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					if t.Keysym.Sym == sdl.GetKeyFromName("Tab") {
						deleteMapCache(&map_walls)
						guiMode = 1 //go to editMode
						print("is in edit mode!!!\n")
					} else if t.Keysym.Sym == sdl.GetKeyFromName("Space") {
						mode++
						if mode > (nbMode - 1) {
							mode = 0
						}
					} else if t.Keysym.Sym == sdl.GetKeyFromName("R") {
						canRenderPath = false
						deleteMapCache(&map_walls)
						star.Init()
						djk.Init()
						canUpdate = false
					}
				}
				break
			}
		}

		if canUpdate && time.Now().Sub(timer_update) > 100000000 {
			timer_update = time.Now()
			starChan := make(chan bool)
			go star.Update(starChan)

			djkChan := make(chan bool)
			go djk.Update(djkChan)

			<-starChan
			<-djkChan

			{
				star.UpdateFinalPath()
				djk.UpdateFinalPath()

			}
		}

		sdl.Delay(16)
	}

	editMode := func() {
		Render.SetDrawColor(0, 0, 0, 255)
		Render.Clear()
		renderMap(true)
		Render.Present()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.MouseButtonEvent:
				if t.Type == sdl.MOUSEBUTTONUP {

				}
				break
			case *sdl.KeyboardEvent:
				if t.Type == sdl.KEYDOWN {
					if t.Keysym.Sym == sdl.GetKeyFromName("Tab") {
						guiMode = 0 //go to runMode
						saveToMapFile(map_walls)
						star.Init()
						djk.Init()
					} else if t.Keysym.Sym == sdl.GetKeyFromName("Down") {
						star.CurNode.Pos += nbw
						if star.CurNode.Pos > nbw*nbh-1 {
							star.CurNode.Pos -= nbw
						}
						djk.CurNode.Pos = star.CurNode.Pos
					} else if t.Keysym.Sym == sdl.GetKeyFromName("Up") {
						star.CurNode.Pos -= nbw
						if star.CurNode.Pos < 0 {
							star.CurNode.Pos += nbw
						}
						djk.CurNode.Pos = star.CurNode.Pos
					} else if t.Keysym.Sym == sdl.GetKeyFromName("Left") {
						star.CurNode.Pos--
						if star.CurNode.Pos < 0 {
							star.CurNode.Pos++
						}
						djk.CurNode.Pos = star.CurNode.Pos
					} else if t.Keysym.Sym == sdl.GetKeyFromName("Right") {
						star.CurNode.Pos++
						if star.CurNode.Pos > nbw*nbh-1 {
							star.CurNode.Pos--
						}
						djk.CurNode.Pos = star.CurNode.Pos
					} else if t.Keysym.Sym == sdl.GetKeyFromName("S") {
						star.Start = star.CurNode.Pos
						print("start new pos : ", star.Start, "\n")
						djk.Start = star.Start
					} else if t.Keysym.Sym == sdl.GetKeyFromName("E") {
						star.End = star.CurNode.Pos
						print("end new pos : ", star.End, "\n")
						djk.Start = star.Start
					} else if t.Keysym.Sym == sdl.GetKeyFromName("W") {
						if map_walls[star.CurNode.Pos] == 1 {
							map_walls[star.CurNode.Pos] = 0
						} else {
							map_walls[star.CurNode.Pos] = 1
						}
					}
				}
				break
			}
		}
		sdl.Delay(16)
	}

	for running {
		if guiMode == 0 {
			runMode()
			rect := sdl.Rect{0, 0, 50, 50}
			printD("run", rect, font, Render)
			Render.Present()
		} else {
			editMode()
			rect := sdl.Rect{0, 0, 67, 50}
			printD("edit", rect, font, Render)
			Render.Present()
		}
	}
	return 0
}
