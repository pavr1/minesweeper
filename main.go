package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"minesweeper/gate"
	"minesweeper/models"
	"net/http"
	"strconv"
	"text/template"
)

var g *gate.Gate
var funcTemplate *template.Template
var db *sql.DB

func main() {
	var err error

	g, err = gate.Start()

	if err != nil {
		panic(err.Error)
	}

	defer g.DbHandler.DB.Close()

	if err != nil {
		panic(err.Error)
	}

	// s := &http.Server{
	// 	Addr:           ":8080",
	// 	Handler:        nil,
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	http.HandleFunc("/", start)
	http.HandleFunc("/signup", singup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/menu", menu)
	http.HandleFunc("/creategame", createGame)
	http.HandleFunc("/loadGame", loadGame)
	http.HandleFunc("/processSpot", processSpot)

	//log.Fatal(s.ListenAndServe())

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server started at port 8080")

	select {}
}

func start(w http.ResponseWriter, r *http.Request) {
	ui, err := getLoggedinUser(w)

	if err != nil {
		ui.Message = err.Error()
		t, _ := template.ParseFiles("ui/main_page.html")
		t.Execute(w, ui)
	} else {
		t, _ := template.ParseFiles("ui/menu.html")
		t.Execute(w, ui)
	}
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

		ui, err := g.CreateUser(name, lastName, password)

		if err != nil {
			ui.Message = err.Error()
		} else {
			ui = models.UIUser{}
			ui.Message = "User created successfully!"
		}

		t, _ := template.ParseFiles("ui/login.html")
		t.Execute(w, ui)
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

		ui, err := g.ValidateLogin(name, password)

		if err != nil {
			ui.Message = err.Error()
			t, _ := template.ParseFiles("ui/login.html")
			t.Execute(w, ui)
		} else {
			ui.Message = "Hello " + ui.Name + "!"
			g.Cache.Set("USER_SESSION", ui)

			t, _ := template.ParseFiles("ui/menu.html")
			t.Execute(w, ui)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func menu(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		ui, err := getLoggedinUser(w)

		if err != nil {
			ui.Message = err.Error()
			t, _ := template.ParseFiles("ui/login.html")
			t.Execute(w, ui)
		} else {
			t, _ := template.ParseFiles("ui/menu.html")
			t.Execute(w, ui)
		}
	case "POST":
		//no actions
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
		rows, _ := strconv.Atoi(r.FormValue("rows"))
		columns, _ := strconv.Atoi(r.FormValue("columns"))
		mines, _ := strconv.Atoi(r.FormValue("mines"))

		ui, err := getLoggedinUser(w)

		if err != nil {
			ui.Message = err.Error()
			t, _ := template.ParseFiles("ui/login.html")
			t.Execute(w, nil)
			return
		}

		ctx := context.Background()
		uigame, err := g.CreateGame(ctx, ui.UserId, rows, columns, mines)
		if err != nil {
			ui.Message = err.Error()
			t, _ := template.ParseFiles("ui/create_game.html")
			t.Execute(w, ui)
			return
		}

		t, err := template.ParseFiles("ui/game.html")
		if err != nil {
			uigame.Message = err.Error()
		}

		t.Execute(w, uigame)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func loadGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		gameId, _ := strconv.Atoi(r.FormValue("gameId"))
		game, err := g.GetSingleGame(int64(gameId))

		if err != nil {
			game.Message = err.Error()
		}

		t, _ := template.ParseFiles("ui/game.html")
		t.Execute(w, game)
	default:
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	g.Cache.Flush()

	getLoggedinUser(w)

	user := models.User{
		Name:     "",
		Password: "",
	}

	t, _ := template.ParseFiles("ui/login.html")
	t.Execute(w, user)
}

func getLoggedinUser(w http.ResponseWriter) (models.UIUser, error) {
	obj, found := g.Cache.Get("USER_SESSION")

	if found {
		user := obj.(models.UIUser)

		return user, nil
	}

	return models.UIUser{}, fmt.Errorf("Session lost, please log back in!")
}

func processSpot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		gameId, _ := strconv.Atoi(r.FormValue("gameId"))
		spotId, _ := strconv.Atoi(r.FormValue("spotId"))
		status := r.FormValue("status")

		getLoggedinUser(w)

		ui, err := g.ProcessSpot(int64(gameId), int64(spotId), status)

		if err != nil {
			ui.Message = err.Error()
		}

		t, _ := template.ParseFiles("ui/game.html")
		t.Execute(w, ui)
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}
