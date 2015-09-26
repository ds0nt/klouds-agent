package main
// http interface
import (
  "net/http"
  "fmt"
	"github.com/gorilla/mux"
)


func startApp(w http.ResponseWriter, r *http.Request)  {
  err := r.ParseForm()
  	if err != nil {
  		//handle error http.Error() for example
  		return
  	}
    image := r.Form["image"][0]
  	image := r.Form["image"][0]


    client := redis.NewClient(&redis.Options{
      Addr:     "localhost:6379",
      Password: "", // no password set
      DB:       0,  // use default DB
  })

  pong, err := client.Ping().Result()
  fmt.Println(pong, err)
}



var API = []struct {
	path        string
	method      string
	handler     func(http.ResponseWriter, *http.Request)
	description string
}{
	{"/do", "POST", startApp, "runs an image"},
}

func main() {
  r := mux.NewRouter()

	for _, endpoint := range API {
		// extract the pieces
		path := endpoint.path
		method := endpoint.method
		handler := endpoint.handler
		description := endpoint.description
		// register the handler
		fmt.Printf("adding %v %v # %v\n", method, path, description)
		r.HandleFunc(path, handler).Methods(method)
	}

  http.Handle("/", r)
  http.ListenAndServe(":2010", nil)
}
