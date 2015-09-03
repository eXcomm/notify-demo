package main
 
import (
	"fmt"
	"net/http"
	"sync"
)
 
type Server struct {
	mu      sync.Mutex
	clients map[chan string]bool
}
 
func (s *Server) Add() chan string {
	s.mu.Lock()
	ch := make(chan string)
	s.clients[ch] = true
	s.mu.Unlock()
	return ch
}
 
func (s *Server) Remove(ch chan string) {
	s.mu.Lock()
	delete(s.clients, ch)
	s.mu.Unlock()
}
 
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/broadcast" {
		s.Notify(w, r)
		return
	}
	ch := s.Add()
	defer s.Remove(ch)
	fmt.Fprintf(w, "msg: %s", <-ch)
}
 
func (s *Server) Notify(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	for ch := range s.clients {
		ch <- "BOOM!"
	}
	s.mu.Unlock()
	fmt.Fprintf(w, "ok")
}
 
func main() {
	http.Handle("/", &Server{
		clients: make(map[chan string]bool),
	})
	http.ListenAndServe(":9000", nil)
}
