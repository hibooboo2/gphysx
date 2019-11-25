package vector

import "time"

import "log"

import "strings"

type Vector struct {
	X float64
	Y float64
	Z float64
}

func (v Vector) AsIntPos() (int, int, int) {
	return int(v.X), int(v.Y), int(v.Z)
}

func (v Vector) Add(vec Vector) Vector {
	v.X += vec.X
	v.Y += vec.Y
	v.Z += vec.Z
	return v
}

func (v Vector) Transform(vec Vector, t func(val float64, val2 float64) float64) Vector {
	log.Printf("Before Transform: X: %f Y: %f Z: %f", v.X, v.Y, v.Z)
	v.X = t(v.X, vec.X)
	v.Y = t(v.Y, vec.Y)
	v.Z = t(v.Z, vec.Z)
	log.Printf("After Transform: X: %f Y: %f Z: %f", v.X, v.Y, v.Z)

	return v
}

func (v *Vector) Edge(origin Vector, end Vector) {
	if v.X >= end.X {
		v.X = end.X
	}
	if v.X < origin.X {
		v.X = origin.X
	}
	if v.Y >= end.Y {
		v.Y = end.Y
	}
	if v.Y < origin.Y {
		v.Y = origin.Y
	}

	if v.Z >= end.Z {
		v.Z = end.Z
	}
	if v.Z < origin.Z {
		v.Z = origin.Z
	}
}

type Particle struct {
	Name               string
	LA, A              Vector
	V                  Vector
	P                  Vector
	AirDragCoeffecient Vector
	Mass               float64
}

func (p *Particle) Next(t time.Duration, f Vector, origin Vector, end Vector) Vector {
	//Detect Coll

	if strings.TrimSpace(p.Name) != "" {
		log.Println("Force ORIG:", f)
	}

	f = f.Transform(f, func(val, val2 float64) float64 {
		return 9.8 * val2 / p.Mass
	})

	f = f.Add(p.AirDragCoeffecient.Transform(p.V, func(val float64, val2 float64) float64 {
		return val * (val2 * val2)
	}))
	if strings.TrimSpace(p.Name) != "" {
		log.Println("Force AFTER:", f)
	}
	timeStep := t.Seconds()
	p.LA = p.A
	p.P.X += (p.V.X * timeStep) + (0.5 * p.LA.X * (timeStep * timeStep))
	p.A.X = f.X / p.Mass
	xAvgAccel := (p.LA.X + p.A.X) / 2
	p.V.X += xAvgAccel * timeStep

	p.P.Y += (p.V.Y * timeStep) + (0.5 * p.LA.Y * (timeStep * timeStep))
	p.A.Y = f.Y / p.Mass
	yAvgAccel := (p.LA.Y + p.A.Y) / 2
	p.V.Y += yAvgAccel * timeStep

	p.P.Z += (p.V.Z * timeStep) + (0.5 * p.LA.Z * (timeStep * timeStep))
	p.A.Z = f.Z / p.Mass
	zAvgAccel := (p.LA.Z + p.A.Z) / 2
	p.V.Z += zAvgAccel * timeStep
	if strings.TrimSpace(p.Name) != "" {
		x, y, z := p.P.AsIntPos()
		log.Printf("Name: %s Drag: %v A: %v V: %v X: %v Y: %v Z: %v", p.Name, p.AirDragCoeffecient, p.A, p.V, x, y, z)
	}
	p.P.Edge(origin, end)
	return p.P
}
