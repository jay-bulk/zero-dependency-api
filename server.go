package main


import (
  "net/http"
  "encoding/json"
)

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
fun (h *testHandlers) testers(w http.ResponseWriter, r *htt.request) {
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
func (h *testHandlers) post(w http.ResponseWriter, r *http.Request) {
  bodyBytes, err := ioutil.ReadAll(r.Body)
  defer r.Body.Close()
  h.Lock()
  defer h.Unlock()
}

func newTestHandlers() *testHandlers{
  return &testHandlers {
    store: map[string]Tester{
      "id1": Tester{
        Name: "John Doe",
        Station: 89,
        ID: "113213",
        Relation: "founderson",
        Job: "technician",
      },
      "id2": Tester{
        Name: "exzibit",
        Station: 22,
        ID: "1",
        Relation: "none",
        Job: "entertainer",
      },
    },
  }
}

func main() {
  testHandlers := newTestHandlers()
  http.HandleFunc("/testers", testHandlers.testers)
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic(err)
  }
}
