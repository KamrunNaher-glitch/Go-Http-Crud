package main

import (
	"context"
	"fmt"
	"net/http"	
	"github.com/joho/godotenv"
)

func main(){
	var err error
	err = godotenv.Load()
	if err != nil{
		panic("env file not found")
	}
	
	connectDb()
	defer db.Close(context.Background())
	mux := http.NewServeMux()

	mux.HandleFunc("GET/",roothandler)
	mux.HandleFunc("GET /health",healthHandler)
	mux.HandleFunc("POST /createUser",createUserHandler)
	mux.HandleFunc("GET /users",getusersHandler)
	mux.HandleFunc("GET /users/{id}",getSingleUsersHandler)
	mux.HandleFunc("PUT /users/{id}",updateUsersHandler)
	mux.HandleFunc("DELETE /users/{id}",deleteUsersHandler)

	fmt.Println("Server is running at port 5000")

	err = http.ListenAndServe(":5000",mux)
	if err != nil {
		fmt.Println("Server Error",err)
	}

}

