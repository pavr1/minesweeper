package main

import (
	"fmt"
	"log"
	"minesweeper/gate"
	"minesweeper/models"
	"net/http"
	"text/template"
)

var g *gate.Gate

func main() {
	var err error

	g, err = gate.Start()

	if err != nil {
		panic(err.Error)
	}

	http.HandleFunc("/", start)
	http.HandleFunc("/signup", Singup)

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server started at port 8080")

	select {}
}

func start(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("ui/main_page.html")
	t.Execute(w, nil)
}

func Singup(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user := models.User{}
		t, _ := template.ParseFiles("ui/sign_up.html")
		t.Execute(w, user)
	case "POST":
		name := r.FormValue("name")
		lastName := r.FormValue("lastName")
		password := r.FormValue("password")

		user := models.User{
			Name:     name,
			LastName: lastName,
			Password: password,
		}

		err := g.CreateUser(user)

		if err != nil {
			user.Message = err.Error()
		} else {
			user = models.User{}
			user.Message = "User created successfully!"
		}

		t, _ := template.ParseFiles("ui/sign_up.html")
		t.Execute(w, user)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
