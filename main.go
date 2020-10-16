package main

import (
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"sync"
)

type store struct {
	data map[string]string
	m    sync.RWMutex
}

var s = &store{data: map[string]string{}, m: sync.RWMutex{}}

//BasicAuth provide basic auth method
func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func myshow(w http.ResponseWriter, r *http.Request) {

	s.m.RLock()
	defer s.m.RUnlock()

	fmt.Fprintf(w, "Get Methods Received\n")
	for k, v := range s.data {
		fmt.Fprintf(w, "Data  Name: %v, Value: %v\n", k, v)
	}

}

func myadd(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Post Methods Received\n")
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, " %v Invalid Post Method", http.ErrNotSupported)
		return
	}
	k := r.Form.Get("name")

	if k == "" {
		fmt.Fprintf(w, " %v Invalid Post Method", http.ErrNotSupported)
		return
	}
	v := r.Form.Get("val")
	s.m.Lock()
	s.data[k] = v
	s.m.Unlock()

	fmt.Fprintf(w, "Add Name:%v, Val:%v", k, v)

}

func myupdate(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Put Methods Received\n")
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, " %v Invalid Post Method", http.ErrNotSupported)
	}
	k := r.Form.Get("name")
	if k == "" {
		fmt.Fprintf(w, " %v Invalid Post Method", http.ErrNotSupported)
		return
	}
	v := r.Form.Get("val")
	if _, ok := s.data[k]; !ok {
		s.m.Lock()
		s.data[k] = v
		s.m.Unlock()
		fmt.Fprintf(w, "Add Name:%v, Val:%v", k, v)
	} else {
		s.m.Lock()
		s.data[k] = v
		s.m.Unlock()

		fmt.Fprintf(w, "Update Name:%v, Val:%v", k, v)
	}

}

func mydelete(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Delete Methods Received\n")
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		fmt.Fprintf(w, " %v Invalid Post Method", http.ErrNotSupported)
	}
	k := r.Form.Get("name")

	s.m.Lock()
	delete(s.data, k)
	s.m.Unlock()

	fmt.Fprintf(w, " Delete Name: %v", k)

}

func simpleMethod(w http.ResponseWriter, r *http.Request) {

	switch r.Method {

	case http.MethodGet:
		myshow(w, r)
	case http.MethodPost:
		myadd(w, r)
	case http.MethodPut:
		myupdate(w, r)
	case http.MethodDelete:
		mydelete(w, r)
	default:
		fmt.Fprintf(w, "%v Invalid request method. ", http.StatusMethodNotAllowed)
	}

}

func main() {

	http.HandleFunc("/", BasicAuth(simpleMethod, "admin", "123456", "Please enter your username and password for this site"))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
