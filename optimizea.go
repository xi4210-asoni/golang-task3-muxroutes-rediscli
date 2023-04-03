package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type setvalue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type server struct {
	database *redis.Client
	router   *mux.Router
}

func NewSetupClient() *server {

	//? new object for redis connect class
	redisConnect := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}) //! Returns an object of Redis Client type

	//TODO: Add a check if the connection is set properly

	routerConnect := mux.NewRouter()

	return &server{redisConnect, routerConnect} //passing the value in the struct
}

func (s *server) setRoutes() {
	s.router.HandleFunc("/set", s.setValue).Methods("POST")
	s.router.HandleFunc("/get", s.getValue).Methods("GET")
	s.router.HandleFunc("/del/{key}", s.delValue).Methods("DELETE")
}

// ! This method used for - **Dependency injections** property
// ** You can use additional parameters here and these parameters can be then added to the handler function for various other use cases: Authenication, any additional tasks
func (s *server) SetValue(abc string) http.HandlerFunc {

	// return func(w http.ResponseWriter, r *http.Request) {
	// fmt.Println(abc) //not using this parameter in the request here
	// }
	return s.setValue
}

func (s *server) setValue(w http.ResponseWriter, r *http.Request) {
	var req setvalue

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		fmt.Println(err) //! Console printing
		w.WriteHeader(http.StatusInternalServerError)
		// w.Write([]byte("Error in request"))
		demovar := map[string]string{
			"error": err.Error(),
		}

		bytes, _ := json.Marshal(demovar)
		// w.Write([]byte(demovar["error"]))
		w.Write(bytes)
	}

	//? Logic for single request
	err2 := s.database.Set(req.Key, req.Value, 0).Err()

	if err2 != nil {
		fmt.Println(err2) //! Console printing
		w.WriteHeader(http.StatusInternalServerError)
		//TODO: Create mssg json here and pass err message
		// w.Write([]byte({message: err2}))

		demovar := map[string]string{
			"error": err2.Error(),
		}

		bytes, _ := json.Marshal(demovar)
		// w.Write([]byte(demovar["error"]))
		w.Write(bytes)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Data Stored!"))
	}

}

func (s *server) GetValue(abc string) http.HandlerFunc {
	return s.getValue
}

func (s *server) getValue(w http.ResponseWriter, r *http.Request) {

	//? takes the query parameter from the url requested
	keys, ok := r.URL.Query()["key"]

	//? If query param is not available
	if !ok {
		keys, err := s.database.Keys("*").Result() //? Retrieves all the values from Redis Database
		if err != nil {                            //? Error Check common for all
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var resultList []setvalue  //? Create an array of objects to store the key value pairs
		for _, key := range keys { //? skip index and store the key values one by one using loops
			value, err := s.database.Get(key).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			resultList = append(resultList, setvalue{key, value}) //? Finally append all the data into the list
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resultList)
		return
	} else {
		value, err := s.database.Get(keys[0]).Result() //? Retrieving the first value from the array
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result := setvalue{keys[0], value}
		json.NewEncoder(w).Encode([]setvalue{result})

	}

}

func (s *server) DelValue(abc string) http.HandlerFunc {
	return s.delValue
}

func (s *server) delValue(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request) //? Retrieve the variables from the api request
	// fmt.Fprint(w, vars)
	key := vars["key"] //? Getting the key parameter
	// fmt.Fprint(w, key)

	err := s.database.Del(key).Err() //? Deleting the data from database
	// fmt.Println(num, err)
	if err != nil {
		// fmt.Println("No request to delete") //! For console only
		s.database.Close()
		fmt.Println("Gracefully Shutting down the connection")
		demovar := map[string]string{
			"error": err.Error(),
		}

		bytes, _ := json.Marshal(demovar)
		// w.Write([]byte(demovar["error"]))
		w.Write(bytes)
		return
	}

	w.Write([]byte("Data Deleted!"))
}

// // ? Create variable for client connection
// var client =

func main() {
	// var s server
	s := NewSetupClient() //? Setting up client details in server
	s.setRoutes()
	// s.router.
	http.ListenAndServe(":3000", s.router)
}
