package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	screenWidth  = 30
	screenHeight = 20
	pipeGap      = 5
	pipeWidth    = 3
	gravity      = 1
	jumpStrength = -3
)

type Bird struct {
	y        int
	velocity int
}

type Pipe struct {
	x      int
	height int
}

var (
	bird  Bird
	pipes []Pipe
	score int
	mutex sync.Mutex
	quit  = make(chan struct{})
)

func initGame() {
	bird = Bird{y: screenHeight / 2, velocity: 0}
	pipes = append(pipes, Pipe{x: screenWidth, height: rand.Intn(screenHeight - pipeGap)})
}

func drawScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	for i := 0; i < screenHeight; i++ {
		for j := 0; j < screenWidth; j++ {
			if j == 5 && i == bird.y {
				fmt.Print("🐦")
			} else {
				isPipe := false
				for _, pipe := range pipes {
					if j >= pipe.x && j < pipe.x+pipeWidth {
						if i < pipe.height || i > pipe.height+pipeGap {
							fmt.Print("█")
							isPipe = true
						}
					}
				}
				if !isPipe {
					fmt.Print(" ")
				}
			}
		}
		fmt.Println()
	}
	fmt.Printf("Score: %d\n", score)
}

func updateGame() {
	for {
		mutex.Lock()
		bird.velocity += gravity
		bird.y += bird.velocity
		if bird.y < 0 {
			bird.y = 0
		}
		if bird.y >= screenHeight {
			fmt.Println("Game Over!")
			close(quit)
			mutex.Unlock()
			return
		}

		for i := range pipes {
			pipes[i].x--
		}
		if pipes[0].x+pipeWidth < 0 {
			pipes = pipes[1:]
			pipes = append(pipes, Pipe{x: screenWidth, height: rand.Intn(screenHeight - pipeGap)})
			score++
		}

		if pipes[0].x == 5 && (bird.y < pipes[0].height || bird.y > pipes[0].height+pipeGap) {
			fmt.Println("Game Over!")
			close(quit)
			mutex.Unlock()
			return
		}

		mutex.Unlock()
		drawScreen()
		time.Sleep(100 * time.Millisecond)
	}
}

func handleInput() {
	err := termbox.Init()
	if err != nil {
		fmt.Println("Failed to initialize termbox")
		return
	}

	defer termbox.Close()
	for {
		event := termbox.PollEvent()
		if event.Type == termbox.EventKey {
			if event.Key == termbox.KeySpace {
				mutex.Lock()
				bird.velocity = jumpStrength
				mutex.Unlock()
			}
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	initGame()
	go updateGame()
	go handleInput()

	<-quit
	fmt.Println("Thanks for playing!")
}
