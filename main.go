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
	http.HandleFunc("/signup", singup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/menu", menu)
	http.HandleFunc("/creategame", createGame)

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server started at port 8080")

	select {}
}

func start(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("ui/main_page.html")
	t.Execute(w, nil)
}

func singup(w http.ResponseWriter, r *http.Request) {
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

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user := models.User{}
		t, _ := template.ParseFiles("ui/login.html")
		t.Execute(w, user)
	case "POST":
		name := r.FormValue("name")
		password := r.FormValue("password")

		user := models.User{
			Name:     name,
			Password: password,
		}

		valid, err := g.ValidateLogin(user)

		if err != nil {
			user.Message = err.Error()
			t, _ := template.ParseFiles("ui/login.html")
			t.Execute(w, user)
		} else {
			if valid {
				user.Message = "Welcome " + user.Name + "!"
				t, _ := template.ParseFiles("ui/menu.html")
				t.Execute(w, user)
			} else {
				user.Message = "Invalid user or password!"
				t, _ := template.ParseFiles("ui/login.html")
				t.Execute(w, user)
			}
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func menu(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user := models.User{}
		t, _ := template.ParseFiles("ui/menu.html")
		t.Execute(w, user)
	case "POST":
		name := r.FormValue("name")
		password := r.FormValue("password")

		user := models.User{
			Name:     name,
			Password: password,
		}

		valid, err := g.ValidateLogin(user)

		if err != nil {
			user.Message = err.Error()
		} else {
			user = models.User{}

			if valid {
				user.Message = "Welcome " + user.Name
				t, _ := template.ParseFiles("ui/menu.html")
				t.Execute(w, user)
			} else {
				user.Message = "Invalid user or password!"
				t, _ := template.ParseFiles("ui/login.html")
				t.Execute(w, user)
			}
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func createGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		user := models.User{}
		t, _ := template.ParseFiles("ui/create_game.html")
		t.Execute(w, user)
	case "POST":
		//to be worked
		// rows, _ := strconv.Atoi(r.FormValue("rows"))
		// columns, _ := strconv.Atoi(r.FormValue("columns"))
		// mines, _ := strconv.Atoi(r.FormValue("mines"))

		// game := models.Game{
		// 	UserId:       1, //????
		// 	TimeConsumed: 0,
		// 	Status:       "Pending",
		// 	Rows:         rows,
		// 	Columns:      columns,
		// 	Mines:        mines,
		// }

		// id, err := g.CreateGame(game)

		// if err != nil {
		// 	user.Message = err.Error()
		// } else {
		// 	user = models.User{}

		// 	if valid {
		// 		user.Message = "Welcome " + user.Name
		// 		t, _ := template.ParseFiles("ui/menu.html")
		// 		t.Execute(w, user)
		// 	} else {
		// 		user.Message = "Invalid user or password!"
		// 		t, _ := template.ParseFiles("ui/login.html")
		// 		t.Execute(w, user)
		// 	}
		// }
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
