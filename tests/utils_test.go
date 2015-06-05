package tests

import(
	"../bot/utils"
	"testing"
)

func TestComputeHmac256(t *testing.T){
	// From instagram example
	// http://instagram.com/developer/restrict-api-requests/
	if utils.ComputeHmac256("200.15.1.1","6dc1787668c64c939929c17683d7cb74") != "7e3c45bc34f56fd8e762ee4590a53c8c2bbce27e967a85484712e5faa0191688"{
		t.Fatalf("Single IP Hash Broken")
	}
	if utils.ComputeHmac256("200.15.1.1,131.51.1.35","6dc1787668c64c939929c17683d7cb74") != "13cb27eee318a5c88f4456bae149d806437fb37ba9f52fac0b1b7d8c234e6cee"{
		t.Fatalf("Multi IP Hash Broken")
	}
}

func TestFollowerDecay(t *testing.T){
	// Set some test
	testFollowerDecay(t, 0.75, 1000)
	testFollowerDecay(t, 0.49, 3000)
}

func testFollowerDecay(t *testing.T, magic float64, target int64){
	followed_by := target
	follows := int64(float64(target) * magic)

	// By definition of Magic and Target, this should be true
	if x := utils.FollowerDecay(followed_by,follows,magic,float64(target)); x != 0 {
		t.Fatalf("Follower Decay is Broken",x,follows,followed_by)
	}
}

func TestLimit(t *testing.T){
	bound  := 100
	amount := 0
	interval := 0
	utils.Limit(&amount,interval,bound)

	bound    = 100
	amount   = 10
	interval = 0
	utils.Limit(&amount,interval,bound)
}