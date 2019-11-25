package main

import (
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell"
	"github.com/hibooboo2/physics/vector"
)

func main() {
	f, _ := os.Create("log.txt")
	log.SetOutput(f)
	log.SetFlags(log.Lshortfile)
	pts := []*vector.Particle{}
	pts = append(pts, &vector.Particle{P: vector.Vector{X: -5, Y: -5}, Mass: 1})
	pts = append(pts, &vector.Particle{P: vector.Vector{X: 5, Y: -5}, Mass: 1})
	pts = append(pts, &vector.Particle{P: vector.Vector{X: 5, Y: 5}, Mass: 1})
	pts = append(pts, &vector.Particle{P: vector.Vector{X: -5, Y: 5}, Mass: 1})

	s, err := tcell.NewTerminfoScreen()
	if err != nil {
		log.Fatal("Failed to start screen", err)
	}
	err = s.Init()
	if err != nil {
		log.Fatal("Failed to start screen", err)
	}
	defer s.Fini()
	origin := &vector.Vector{20, 10, 0}
	player := &vector.Particle{Mass: 1, Name: "player"}
	player.AirDragCoeffecient.Y = -0.000605
	// player.V.X = 1
	gravity := vector.Vector{Y: 1}
	done := monitorKeys(s, []*vector.Vector{&player.P})
	last := time.Now()
	timeMoved := 0
	for {

		timeMoved++
		time.Sleep(time.Millisecond * 1000 / 10)
		s.Clear()
		select {
		case <-done:
			s.Fini()
			os.Exit(0)
		default:
		}

		frames := 45
		if timeMoved%frames > frames/2 {
			s.SetContent(3, 3, '🔴', nil, tcell.StyleDefault)
		} else {
			s.SetContent(3, 3, '💚', nil, tcell.StyleDefault)
		}

		now := time.Since(last)
		last = time.Now()
		// now = time.Millisecond

		OX, OY, _ := origin.AsIntPos()
		for _, pt := range pts {
			// pt.Next(now, gravity)
			x, y, _ := pt.P.AsIntPos()
			s.SetContent(x+OX, y+OY, '*', nil, tcell.StyleDefault)
		}
		player.Next(now, gravity, vector.Vector{-10, -10, -5}, vector.Vector{10, 10, 5})
		PX, PY, _ := player.P.AsIntPos()
		s.SetContent(OX+PX, OY+PY, '😃', nil, tcell.StyleDefault)
		s.SetContent(40, 5, '0'+rune(int(player.A.X)%10), nil, tcell.StyleDefault)
		s.SetContent(40, 6, '0'+rune(int(player.V.X)%10), nil, tcell.StyleDefault)
		s.Show()
	}

}

func monitorKeys(s tcell.Screen, vectors []*vector.Vector) chan struct{} {
	done := make(chan struct{})
	go func() {
		for {
			evt := s.PollEvent()
			force := vector.Vector{}
			forceMod := 200.00
			switch e := evt.(type) {
			case *tcell.EventKey:
				switch e.Key() {
				case tcell.KeyUp:
					force.Y += forceMod
				case tcell.KeyDown:
					force.Y -= forceMod
				case tcell.KeyLeft:
					force.X -= forceMod
				case tcell.KeyRight:
					force.X += forceMod
				case tcell.KeyEsc:
					done <- struct{}{}
				}
			}
			for _, v := range vectors {
				v.Add(force)
			}
		}
	}()

	return done
}
