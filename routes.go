package routes

//? Single file for all the routes used

func (s *server) routes() {
	s.router.HandleFunc("/set", s.setValue()).Methods("POST")
	s.router.HandleFunc("/get", s.getValue()).Methods("GET")
	s.router.HandleFunc("/del/{key}", s.delValue()).Methods("DELETE")
}
