package tests

import (
        "testing"
        "appengine/aetest"
        "net/http"
)

var inst aetest.Instance
var client = &http.Client{}

func TestHandles(t *testing.T) {
    instance, err := aetest.NewInstance(nil)
    inst = instance
    if err != nil {
            t.Fatalf("Failed to create instance: %v", err)
    }
    defer inst.Close()

    get(t, "/")
    get(t, "/auth")
    get(t, "/hashtag?hashtags=1,2,3")
    get(t, "/process")

}

func get(t *testing.T, location string) {
    req, err := inst.NewRequest("GET", location, nil)
    if err != nil {
            t.Fatalf("Failed to create req: %v", err)
    }
    resp, err := client.Do(req)
    t.Log(resp)
}
