package main
import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"	
	"strconv"
	"github.com/jackc/pgx/v5"
	
)

func roothandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Welcome to Go Server")
}
func healthHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintln(w, "Server is up and healthy")
}
func createUserHandler(w http.ResponseWriter, r *http.Request){
	
	var newUser User 
	err:=json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		fmt.Println(err)
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w,"Invalid Request Body")
	return
	}
	fmt.Println(newUser)
	// newUser.Id = len(users) +1
	// users = append(users, newUser)
	query := `
	insert into users (username,age,email)
	values ($1,$2,$3)
	returning id
	`
	 err =db.QueryRow(context.Background(),query,newUser.Name,newUser.Age,newUser.Email).Scan(&newUser.Id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"Could not create user")
		return
	}

	w.Header().Set("Content-Type","application/json")
	 w.WriteHeader(http.StatusCreated)
	 json.NewEncoder(w).Encode(newUser)
	

}
func getusersHandler(w http.ResponseWriter, r *http.Request){
	query := `select id,username,age,email from users`

	rows, err := db.Query(context.Background(),query)

	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w,"Could Not get All Users")
		return

	}
	defer rows.Close()
	var users [] User

	for rows.Next(){
		var user User
		err := rows.Scan(&user.Id, &user.Name, &user.Age,&user.Email)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintln(w,"Could not scan user")
			return
		}
		users = append(users, user)
	}

	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(users)
}

func getSingleUsersHandler(w http.ResponseWriter, r *http.Request){
	idParam := r.PathValue("id")
	// fmt.Printf("the value od id id %v and the type of the id is %T",idParam,idParam)

	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Invalid User Id")
	}
	for _, user := range users {
		if user.Id == id {
			w.Header().Set("Content-Type","application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w,"User not Found")



}

func updateUsersHandler(w http.ResponseWriter,r *http.Request){
	idParam := r.PathValue("id")

	id,err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Invalid user id")
		return 
	}
	 var updateUser User 

	  err = json.NewDecoder(r.Body).Decode(&updateUser)
	  if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Invalid Request Body")
		return
	  }
	  query := `
	  update users 
	  set username = $1, age = $2, email = $3
	  where id = $4 
	  returning id, username, age,email
	  `
	 err = db.QueryRow(context.Background(),query, updateUser.Name,
	 updateUser.Age,updateUser.Email,id).Scan(&updateUser.Id,&updateUser.Name,
	&updateUser.Age,&updateUser.Email)

	if err == pgx.ErrNoRows{
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w,"User Not Found")
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w,"Could not update user")
		return 
	}
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(updateUser)
	

}

func deleteUsersHandler(w http.ResponseWriter, r *http.Request){
	idParam := r.PathValue("id")

	id,err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Invalid user id")
		return 
	}
	// for idx,user := range users {
	// 	if user.Id == id {
	// 		// users = append(users[:idx],users[idx+1:]...)
	// 		users = slices.Delete(users,idx,idx+1)
	// 		w.WriteHeader(http.StatusNoContent)
	// 		return
	// 	}
	// }
	query := `delete from users where id = $1`
	cmdTag, err := db.Exec(context.Background(),query,id)
	if err != nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w,"Could not delete user")
		return 
	}

	if cmdTag.RowsAffected() != 1 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintln(w,"User not Found")
		return
	}

	w.WriteHeader(http.StatusNoContent)
	fmt.Println(w,"User deleted Successfully")
}