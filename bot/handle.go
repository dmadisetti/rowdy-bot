package bot

import (
    "net/http"
    "fmt"
)

// Handler just to hold everything together
type Handler struct{
    handle func(http.ResponseWriter, *http.Request, *Session)
}

// Constructor
func NewHandler(path string, handle func(http.ResponseWriter, *http.Request, *Session)) {
    handler := &Handler{handle:handle}
    http.HandleFunc(path, handler.preHandle)
}

// Passed into http for all handlers
func (h *Handler)preHandle(w http.ResponseWriter, r *http.Request){

    // Create session
    session := NewSession(r)

    // get session from data store or create and prompt config
    if session.Load() {
        session.Save()
        fmt.Fprint(w, "Please Configure Settings")
        return
    }

    // Call handler set earlier
    h.handle(w, r, session)
}