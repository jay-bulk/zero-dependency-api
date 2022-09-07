package main
import (
  "net/http"
  "encoding/json"
  "sync"
  "io/ioutil"
  "fmt"
  "time"
  "os"
  "strings"
)

//This is a "zero" dependency api
//Based on kubucation video on youtub https://www.youtube.com/watch?v=1v11Ym6Ct9Q
// RHETTB
// watchTime = 27:01 
// Music = Future - Mask Off (Aesthetic Remix) 
type Tester struct {
  Name string `json:"name"`
  Job string `json:"job"`
  ID string `json:"id"`
  Relation string  `json:"relation"`
  Station int `json:"station"`
}

type testHandlers struct {
  sync.Mutex
  store map[string]Tester
}
func (h *testHandlers) testers(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    h.get(w, r)
    return
  case "POST":
    h.post(w, r)
    return
  default:
    w.WriteHeader(http.StatusMethodNotAllowed)
    w.Write([]byte("Method not allowed"))
    return
  }
}
func (h *testHandlers) get(w http.ResponseWriter, r *http.Request) {
  testers := make([]Tester, len(h.store))
  h.Lock()
  i := 0
  for _, tester := range h.store {
    testers[i] = tester
    i++
  }
  h.Unlock()
  jsonBytes, err := json.Marshal(testers)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
  }
  w.Header().Add("content-type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(jsonBytes)
}
func (h *testHandlers) getTester(w http.ResponseWriter, r *http.Request) {
  parts := strings.Split(r.URL.String(), "/")
  if len(parts) != 3 {
    w.WriteHeader(http.StatusNotFound)
    return
  }
  
  h.Lock()
  tester, ok := h.store[parts[2]]
  h.Unlock()
  if !ok {
    w.WriteHeader(http.StatusNotFound)
    return
  }  
    jsonByters, err := json.Marshal(tester)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(err.Error()))
    }
    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonByters)
}

func (h *testHandlers) post(w http.ResponseWriter, r *http.Request) {
  bodyBytes, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  if err != nil {
     w.WriteHeader(http.StatusInternalServerError)
     w.Write([]byte(err.Error()))
     return
  }

  ct := r.Header.Get("content-type")
  if ct != "application/json" {
     w.WriteHeader(http.StatusUnsupportedMediaType)
     w.Write([]byte(fmt.Sprintf("need content-type 'application/json', but got '%s'", ct)))
     return
  }
   
  
  var tester Tester
  err = json.Unmarshal(bodyBytes, &tester)
  if err != nil {
     w.WriteHeader(http.StatusBadRequest)
     w.Write([]byte(err.Error()))
  }

  tester.ID = fmt.Sprintf("%d", time.Now().UnixNano())
  h.Lock()
  h.store[tester.ID] = tester
  defer h.Unlock()
  
}

func newTestHandlers() *testHandlers{
  return &testHandlers {
    store: map[string]Tester{ },
  }
}

type adminPortal struct {
  password string
}

func newAdminPortal() *adminPortal {
  password := os.Getenv("ADMIN_PASSWORD")
  if password == "" {
    panic("required env var ADMIN_PASSWORD not set")
  }
  return &adminPortal{password: password}
}

func (a adminPortal) handler(w http.ResponseWriter, r *http.Request) {
  user, pass, ok := r.BasicAuth()
  if !ok || user != "admin" || pass != a.password {
    w.WriteHeader(http.StatusUnauthorized)
    w.Write([]byte("401 - unauthorized"))
    return
  }
  w.Write([]byte("<html><h1>Super Secret admin portal</h1></html>"))
  
}

func main() {
  admin := newAdminPortal()
  testHandlers := newTestHandlers()
  http.HandleFunc("/testers", testHandlers.testers)
  http.HandleFunc("/testers/", testHandlers.getTester)
  http.HandleFunc("/admin", admin.handler)
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic(err)
  }
}
