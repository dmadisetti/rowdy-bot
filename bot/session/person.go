package session

import "math"

type Person struct{
    Following float64
    Followers float64
    Posts float64
    Follows bool // 1 or 0
    Followed bool
}

func (person *Person) Y() (y float64) {
    y = 1.0
    if person.Follows {
        y = 0.0
    }
    return
}

func (person *Person) Sigmoid(gradient []float64) float64{
	f := gradient[0] +  person.Followers * gradient[1] + person.Following * gradient[2] + person.Posts * gradient[3]
	return 1.0/(1.0 + math.Exp(-f))
}

func (person *Person) J(gradient []float64) float64{
    y := person.Y()
	h := person.Sigmoid(gradient)
	q := y * math.Log(h) + (1 - y) * math.Log(1 - h)
	return q
}