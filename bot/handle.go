package bot

import (
    "html/template"
    "net/http"
    "net/url"
    "fmt"
    "bot/session"
)

// Handler just to hold everything together
type Handler struct{
    handle func(http.ResponseWriter, *http.Request, *session.Session)
}

// Constructor
func NewHandler(path string, handle func(http.ResponseWriter, *http.Request, *session.Session)) {
    handler := &Handler{handle:handle}
    http.HandleFunc(path, handler.preHandle)
}

// Passed into http for all handlers
func (h *Handler)preHandle(w http.ResponseWriter, r *http.Request){

    // Create session
    s := session.NewSession(r)

    // get session from data store or create and prompt config
    if !s.LoadSettings() {
        if keys, ok := parseKeys(r); ok{
            s.InitAuth(keys["client_id"][0],keys["client_secret"][0],keys["callback"][0],keys["hash"][0])
        } else {
            s.Save()
            t, e := template.ParseGlob("templates/setup.html")
            if e != nil {
                fmt.Fprint(w, e)
                return
            }

            // render with records
            err := t.Execute(w, s)
            if err !=nil{
                panic(err)
            }
            return
        }
    }

    if !s.LoadMachine() {
        s.Save()
    }

    // Call handler set earlier
    h.handle(w, r, s)
}

func parseKeys(r *http.Request)(url.Values, bool){
    keys := []string{"client_id","client_secret","callback"}
    values := r.URL.Query()
    for i := 0; i < len(keys); i++ {    
        if _, suc := values["client_id"] ; !suc {
            return nil, false;
        }
    }
    return values, true;
}