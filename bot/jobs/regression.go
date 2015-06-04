package jobs

import(
	"bot/session"
	"math"
	"log"
)
func LogisticRegression(s *session.Session){

    // Grab all people because of many iterations
    people := s.GetPeople()
    objective := Objective{
        People: people,
        Lambda: s.GetLambda(),
        Alpha: s.GetAlpha(),
        Size: float64(len(people)),
    }

    start := []float64{1, 1, 1, 0} // Bias, Following, Followers, Posts
    minimum := Minimize(objective, start)
    log.Println(minimum)

    s.SetTheta(minimum)
    s.StopProcessing()
}

func Minimize(objective Objective, thetas []float64) []float64{

    nthetas := objective.EvaluateGradient(thetas)
    next := objective.EvaluateFunction(nthetas)
    value := objective.EvaluateFunction(thetas)

    i := 0

    for (math.IsNaN(value) || value >= next) && i < 1000 {
        // log.Println(value)
        thetas = nthetas
        value = next
        nthetas = objective.EvaluateGradient(thetas)
        next = objective.EvaluateFunction(nthetas)
        i += 1
    }

    return thetas
}

type Objective struct{
    People []session.Person
    Lambda float64
    Alpha float64
    Size float64
}

// Cost Function
func (o Objective) EvaluateFunction(thetas []float64) float64{

    sum := 0.0
    for _, person := range o.People {
        sum += person.J(thetas)
    }
    //sum += -(o.Lambda/2) * thetas[3] * thetas[3]
    return -sum/o.Size
}

// Cost Function derivative
func (o Objective) EvaluateGradient(thetas []float64) (gradient []float64) {
    gradient = make([]float64, len(thetas))

    gradient[0] = o.CostD(0, thetas)
    gradient[1] = o.CostD(1, thetas)
    gradient[2] = o.CostD(2, thetas)
    //gradient[3] = o.CostD(3, thetas)

    return
}

func (o Objective) CostD(i int, thetas []float64) float64{

    sum := 0.0

    switch i {
    case 0:
        for _, person := range o.People {
            sum += (person.Sigmoid(thetas) - person.Y())
        }
        break;
    case 1:
        for _, person := range o.People {
            sum += (person.Sigmoid(thetas) - person.Y()) * person.Followers
        }
        break;
    case 2:
        for _, person := range o.People {
            sum += (person.Sigmoid(thetas) - person.Y()) * person.Following
        }
        break;
    case 3:
        for _, person := range o.People {
            sum += (person.Sigmoid(thetas) - person.Y()) * person.Posts
        }
        sum += o.Lambda * thetas[3]
        break;
    default:
        break;
    }
    return thetas[i] - o.Alpha * sum
}