package bot

import(
    "log"
)

// Supervised learning -
//// Ideal person is a follower
//// Non ideal person is someone we follow, but doesn't follow back

func process(s *Session, users *Users, i int){
    for i >= 0 {

    	user  := GetUser(s, users.Data[i].Id)

		// TODO: grab media and create greedy hashtag logic
		log.Println(user)

		// TODO: Save data to aggregrate records

        // Decrement
        i--
    }
}

// Async set up multi calls
func listen(s *Session, next chan *Users, calls int) {

	// TODO: get optimal value for count given instagram limits at 5000 rph
	count := 3
    for {
        select {
	        case users := <-next:

			    i := len(users.Data) - 1
				log.Println(i)
				go process(s, users, i)

				// TODO: if end of line; save and start processing jobs, else store value for next time
	            if i == 49 && calls < count {
					log.Println("Got another batch")
	            	nxt  := make(chan *Users)
	            	batch := GetNext(s, users)
	            	go listen(s, nxt, calls + 1)
					nxt <- &batch
	            }
	            return
        }
    }
}

func Learn(s *Session) string{

	//TODO: Query Learning records to see where we need to go

    next  := make(chan *Users)
	batch := GetFollowers(s)
	go listen(s, next, 0)
	next <- &batch

	if(s.SetLearning()){
		return "Learning Jobs set"
	}
	return "Already learning"
}