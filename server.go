//This is a "zero" dependency api
//Based on kubucation
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

//Calls made to this API are held in memory and are not sent to any database


type Student struct {
  Name string `json:"student"`
  Job string `json:"role"`
  ID string `json:"student_id"`
  Class string  `json:"Class"`
  Professor int `json:"Professor"`
}

//Handlers are mutually exclusive (lockable/unlockable) can be locked/unlocked with .Lock(), .Unlock(); Can be used with defer (performs like an await in js)
type studentHandlers struct {
  sync.Mutex
  store map[string]Student
}

func (h *StudentHandlers) students(w http.ResponseWriter, r *http.Request) {
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

//Get function as used above (h.get(w, r))
func (h *studentHandlers) get(w http.ResponseWriter, r *http.Request) {
  students := make([]Student, len(h.store))
  h.Lock()
  i := 0
  for _, student := range h.store {
    students[i] = student
    i++
  }
  h.Unlock()
  
  //json marshaling is Go's JSON encoding/decoding Marshal is to encode it Unmarshall is the decode
  jsonBytes, err := json.Marshal(student)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte(err.Error()))
  }
  w.Header().Add("content-type", "application/json")
  w.WriteHeader(http.StatusOK)
  w.Write(jsonBytes)
}

//function for getting a random student from the array
func (h *studentHandlers) getRandomStudent(w http.ResponseWriter, r *http.request) {
  ids := make([]string, len(h.store))
  h.Lock()
  i := 0
  for id := range h.store {
    ids[i] == id
    i++
  }
  defer h.Unlock()
  var target string
  if len(ids) == 0 {
    w.WriteHeader(http.StatusNotFound)
    return
  } else if len(ids) == 1 {
      target = ids[0]
  } else {
     rand.Seed(time.Now().UnixNano())
     target = ids[rand.Intn(len(ids))]
  }
  w.Header().Add("location", fmt.Spring("/student/%s", target))
  w.WriteHeader(http.StatusFound)
}

//Get information on a specific student (default) return value if the api call doesn't have a path specified
func (h *studentHandlers) getStudent(w http.ResponseWriter, r *http.Request) {
  parts := strings.Split(r.URL.String(), "/")
  if len(parts) != 3 {
    w.WriteHeader(http.StatusNotFound)
    return
  }
  if parts[2] == "random" {
    h.getRandomStudent(w,r)
    return
  } 
  h.Lock()
  student, ok := h.store[parts[2]]
  h.Unlock()
  if !ok {
    w.WriteHeader(http.StatusNotFound)
    return
  }  
    jsonByters, err := json.Marshal(student)
    if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      w.Write([]byte(err.Error()))
    }
    w.Header().Add("content-type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(jsonByters)
}

//Create a new student object
func (h *studentHandlers) post(w http.ResponseWriter, r *http.Request) {
  bodyBytes, err := ioutil.ReadAll(r.Body)
  //Wait for the request Body Return
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
   
  var student Student
  err = json.Unmarshal(bodyBytes, &student)
  if err != nil {
     w.WriteHeader(http.StatusBadRequest)
     w.Write([]byte(err.Error()))
  }

  //Define a student id as a function of the epoc time that the function is called
  student.ID = fmt.Sprintf("%d", time.Now().UnixNano())
  h.Lock()
  h.store[student.ID] = student
  defer h.Unlock()
  
}

func newStudentHandlers() *studentHandlers{
  return &studentHandlers {
    store: map[string]Student{ },
  }
}

type adminPortal struct {
  password string
}

//Define admin login accessed via call to /admin
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
  studentHandlers := newStudentHandlers()
  http.HandleFunc("/students", studentHandlers.students)
  http.HandleFunc("/students/", studentHandlers.getStudent)
  http.HandleFunc("/admin", admin.handler)
  err := http.ListenAndServe(":8080", nil)
  if err != nil {
    panic(err)
  }
}
