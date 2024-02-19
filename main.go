package main

import (
	"fmt"
	"math"
	"time"

	"github.com/inancgumus/screen"
)

// Render space height
const HEIGHT = 40

// Render space width
const WIDTH = 40

// Torus body radius
const R1 = 2

// Torus rotational radius
const R2 = 3

// Distance of donut from the viewer
const K2 = 10

// Donut scale
const K1 = WIDTH * K2 * 3 / (8 * (R1 + R2))

// Theta rotation step
const DELTAT = 0.07

// Phi rotation step
const DELTAP = 0.02

// Rendering characters
const RENCH = ".,-~:;=!*#$@"

// FPS
const FPS = 144

func renderFrame(angleA float64, angleB float64) [HEIGHT][WIDTH]rune {
	cosA := math.Cos(angleA)
	sinA := math.Sin(angleA)
	cosB := math.Cos(angleB)
	sinB := math.Sin(angleB)

	var theta, phi float64

	output := [HEIGHT][WIDTH]rune{}
	for i, row := range output {
		for j := range row {
			output[i][j] = ' '
		}
	}

	zBuf := [HEIGHT][WIDTH]float64{}

	for theta = 0.0; theta < 2*math.Pi; theta += DELTAT {
		cosT := math.Cos(theta)
		sinT := math.Sin(theta)

		for phi = 0.0; phi < 2*math.Pi; phi += DELTAP {
			cosP := math.Cos(phi)
			sinP := math.Sin(phi)
			circleX := R2 + R1*cosT
			circleY := R1 * sinT

			x := circleX*(cosB*cosP+sinA*sinB*sinP) - circleY*cosA*sinB
			y := circleX*(sinB*cosP-sinA*cosB*sinP) + circleY*cosA*cosB
			z := K2 + cosA*circleX*sinP + circleY*sinA

			ooz := 1 / z

			xProj := int(WIDTH/2 + K1*ooz*x)
			yProj := int(HEIGHT/2 - K1*ooz*y)

			lum := cosP*cosT*sinB - cosA*cosT*sinP - sinA*sinT + cosB*(cosA*sinT-cosT*sinA*sinP)
			if lum > 0 {
				if ooz > zBuf[yProj][xProj] {
					zBuf[yProj][xProj] = ooz
					lumIndex := int(math.Floor(lum * (float64(len(RENCH)-1) / math.Sqrt2)))
					output[yProj][xProj] = rune(RENCH[lumIndex])
				}
			}
		}
	}

	return output
}

func printFrame(frame [HEIGHT][WIDTH]rune) {
	output := ""
	for i := 0; i < HEIGHT; i++ {
		for j := 0; j < WIDTH; j++ {
			output += string(frame[i][j])
		}
		output += "\n"
	}
	fmt.Print(output)
}

func main() {
	A := 0.0
	B := math.Pi / 2

	// Not sure why, but the control character approach doesn't work for me, had to resort to a third party import
	screen.Clear()

	for {
		renderStart := time.Now()

		frame := renderFrame(A, B)

		screen.MoveTopLeft()
		printFrame(frame)

		A += DELTAT
		B += DELTAP

		if FPS > 0 {
			renderDuration := time.Since(renderStart)

			delay := time.Duration((time.Second.Nanoseconds() - FPS*renderDuration.Nanoseconds()) / FPS)
			time.Sleep(delay)
		}
	}
}
