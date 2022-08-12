package main


import "net/http"

type Tester struct {
  Name String `json:"name"`
  Job String `json:"job"`
  ID String `json:"id"`
  Relation String  `json:"relation"`
  Station int `json:"station"`
}
func testHandler(w http.ResponseWriter, r *http.Request) {

}
func main() {
  http.HandleFunc("/testers", testHandler)
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic(err)
  }
}
