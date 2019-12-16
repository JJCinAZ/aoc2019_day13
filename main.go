package main

import (
	"cloud.google.com/aoc2019/day11/intcode"
	"fmt"
	"io/ioutil"
)

var (
	pixels                   []byte
	width, height int
	lastBallX, lastBallY,lastPaddleX int

)

func main() {
	data, err := ioutil.ReadFile("pgm.dat")
	if err != nil {
		panic(err)
	}
	width, height = 50,25
	pixels = make([]byte, width*height)
	for i := range pixels {
		pixels[i] = ' '
	}
	part2(string(data))
}

func getJoystick() int {
	joystick := 0
	if lastBallX != -1 {
		if lastPaddleX > lastBallX {
			joystick = -1
		} else if lastPaddleX < lastBallX {
			joystick = 1
		}
	}
	fmt.Printf("Paddle %d, Ball %dx%d, joystick %d\n",
		lastPaddleX, lastBallX, lastBallY, joystick)
	return joystick
}

func part2(sourceCode string) {
	outputChan := make(chan int)
	vm := intcode.NewVM(1, intcode.Compile(sourceCode), getJoystick, outputChan)
	vm.Pgm.SetMem(0, 2)
	go handleIO(outputChan)
	vm.Pgm.Debug(false)
	if err := vm.ExecPgm(); err != nil {
		panic(err)
	}
	printBitmap()
	blocks := 0
	for i := range pixels {
		if pixels[i] == '@' {
			blocks++
		}
	}
	fmt.Println(blocks)
}

/*
func part1(sourceCode string) {
	ioChan := make(chan int)
	vm := intcode.NewVM(1, intcode.Compile(sourceCode), ioChan)
	go handleIO(ioChan)
	vm.Pgm.Debug(false)
	if err := vm.ExecPgm(); err != nil {
		panic(err)
	}
	close(ioChan)
	printBitmap()
	blocks := 0
	for i := range pixels {
		if pixels[i] == '@' {
			blocks++
		}
	}
	fmt.Println(blocks)
}
*/

func printBitmap() {
	for y := 0; y < height; y++ {
		fmt.Println(string(pixels[y*width:y*width+width]))
	}
}

func handleIO(output chan int) {
	lastBallX = -1
	for {
		x, ok := <-output
		if !ok {
			return
		}
		y := <-output
		tileId := <-output
		if x == -1 {
			fmt.Println("score ", tileId)
			continue
		}
		if x < 0 || x >= width {
			panic("X beyond extents")
		}
		if y < 0 || y >= height {
			panic(fmt.Sprintf("Y beyond extents %d", y))
		}
		var c byte
		switch tileId {
		case 0:
			c = ' '
			if pixels[y*width+x] == '@' {
				fmt.Printf("hit at %dx%d\n", x, y)
			}
		case 1:
			c = '#'
		case 2:
			c = '@'
		case 3:
			lastPaddleX = x
			c = '-'
		case 4:
			lastBallX = x
			lastBallY = y
			c = '*'
		}
		pixels[y*width+x] = c
	}
}
