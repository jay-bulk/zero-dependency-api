package main


import "net/http"

type Tester struct {
  Name string `json:"name"`
  Job string `json:"job"`
  ID string `json:"id"`
  Relation string  `json:"relation"`
  Station int `json:"station"`
}

type testHandlers struct {
  store map[string]Tester
}
func (h *testHandlers) get(w http.ResponseWriter, r *http.Request) {

}

func newTestHandlers() *testHandlers{
  return &testHandlers {
    store: map[string]Tester{},
  }
}

func main() {
  testHandlers := newTestHandlers()
  http.HandleFunc("/testers", testHandlers.get)
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic(err)
  }
}
