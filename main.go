package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type SetValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// type ResponseResult struct {
// 	Key string `json:"key"`
// }

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

// values can be of multiple type
// type GetValue interface{}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/set", setValue).Methods("POST")
	// r.HandleFunc("/get?{key}", getValue).Methods("GET")  //need the key to change
	r.HandleFunc("/get", getValue).Methods("GET")
	r.HandleFunc("/del/{key}", delValue).Methods("DELETE") //need the key to delete

	http.ListenAndServe(":80", r)

}

//creating functions to delete the values

func delValue(w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request) //? Retrieve the variables from the api request
	// fmt.Fprint(w, vars)
	key := vars["key"] //? Getting the key parameter
	// fmt.Fprint(w, key)
	err := client.Del(key).Err() //? Deleting the data from database
	if err != nil {
		// fmt.Println("No request to delete") //! For console only
		client.Close()
		fmt.Println("Gracefully Shutting down the connection")
	}
}

func setValue(w http.ResponseWriter, r *http.Request) {
	// var req []SetValue //?array of the struct
	var req SetValue

	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		fmt.Println(err) //! Console printing
	}

	//? Logic for single request
	err2 := client.Set(req.Key, req.Value, 0).Err()

	if err2 != nil {
		fmt.Println(err2) //! Console printing
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Data Stored!"))

	// logic for getting multiple values
	// for _, kv := range req {
	// err := client.Set(kv.Key, kv.Value, 0).Err()
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		return
	// 	}
	// }

	// fmt.Fprint(w, value2)
	// fmt.Fprint(w, req.Key)
	// w.Write([]byte(req.Key))
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Fprint(w, req.Key, req.Value)
	// clienterr := client.Set(req.Key, req.Value, 0).Err()
	// if clienterr != nil {
	// }
	// fmt.Fprint(w, req.Key, req.Value)
}

func getValue(w http.ResponseWriter, r *http.Request) {

	//? takes the query parameter from the url requested
	keys, ok := r.URL.Query()["key"]

	//? If query param is not available
	if !ok {
		keys, err := client.Keys("*").Result() //? Retrieves all the values from Redis Database
		if err != nil {                        //? Error Check common for all
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var resultList []SetValue  //? Create an array of objects to store the key value pairs
		for _, key := range keys { //? skip index and store the key values one by one using loops
			value, err := client.Get(key).Result()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			resultList = append(resultList, SetValue{key, value}) //? Finally append all the data into the list
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resultList)
		return
	} else {
		value, err := client.Get(keys[0]).Result() //? Retrieving the first value from the array
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result := SetValue{keys[0], value}
		json.NewEncoder(w).Encode([]SetValue{result})

	}

}

// func getValue(w http.ResponseWriter, request *http.Request) {

// }

// func delResponseSet(response http.ResponseWriter, statusCode int, data interface{}) {
// 	result, _ := json.Marshal(data)
// 	response.Header().Set("Content-type", "application/json")
// 	response.WriteHeader(statusCode)
// 	response.Write(result)
// }

//Manual testing

//seting values
// err := client.Set("username", "user100", 0).Err()
// if err != nil {
// panic(err)
// fmt.Println(err)
// }
//getting values
// val, err2 := client.Get("username").Result()
// if err2 != nil {
// 	log.Fatal(err2)
// }
// fmt.Println(val)

// delete value
// client.Del("username").Result()
// fmt.Println("Deleting values")
//! Checking connection
// pong, err := client.Ping().Result()
// fmt.Println(pong, err)

// localhost connected
