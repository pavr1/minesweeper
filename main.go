package main

import (
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

	funcTemplate, err = template.New("ui/game.html").Funcs(template.FuncMap{
		"increase": func(i int) int {
			return -1
		},
	}).ParseFiles("ui/game.html")

	if err != nil {
		panic(err.Error)
	}

	http.HandleFunc("/", start)
	http.HandleFunc("/signup", singup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/menu", menu)
	http.HandleFunc("/creategame", createGame)
	http.HandleFunc("/LoadPendingGame", LoadPendingGame)

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

		t, _ := template.ParseFiles("ui/login.html")
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

		resultUser, err := g.ValidateLogin(user)

		if err != nil {
			user.Message = err.Error()
			t, _ := template.ParseFiles("ui/login.html")
			t.Execute(w, user)
		} else {
			if resultUser != nil {
				resultUser, err := user.ValidateUser(g.DbHandler)

				if err != nil {
					g.Cache.Set("USER_SESSION", user)
				}

				user = *resultUser

				g.Cache.Set("USER_SESSION", user)

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
		user := getLoggedinUser(w)

		t, _ := template.ParseFiles("ui/menu.html")
		t.Execute(w, user)
	case "POST":
		// name := r.FormValue("name")
		// password := r.FormValue("password")

		// user := models.User{
		// 	Name:     name,
		// 	Password: password,
		// }

		// valid, err := g.ValidateLogin(user)

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

		user, found := g.Cache.Get("USER_SESSION")

		if !found {
			//No user logged in
			user := models.User{}
			user.Message = "Session lost, please login again!"
			t, _ := template.ParseFiles("ui/login.html")
			t.Execute(w, user)
			return
		}

		game := &models.Game{
			UserId:       int64(user.(models.User).UserId),
			TimeConsumed: 0,
			Status:       "Pending",
			Rows:         rows,
			Columns:      columns,
			Mines:        mines,
		}

		newSpots, err := g.CreateGame(game)

		if err != nil {
			game.Message = err.Error()
			t, err := template.ParseFiles("ui/create_game.html")
			if err != nil {
				game.Message = err.Error()
			}

			t.Execute(w, game)
		} else {
			game.Message = "Game Created Successfully "
			game.Spots = *newSpots

			// t, err := template.ParseFiles("ui/game.html")
			// if err != nil {
			// 	game.Message = err.Error()
			// }

			// t.Execute(w, game)

			// t, err := funcTemplate.Clone()
			// t = t.Funcs(template.FuncMap{
			// 	"inc": func(i int) int {
			// 		i++

			// 		return i
			// 	},
			// })
			// if err != nil {
			// 	fmt.Println(err.Error())
			// 	game.Message = err.Error()
			// }
			//err = t.Execute(w, game)

			// tpl := template.Must(template.New("").Funcs(
			// 	template.FuncMap{
			// 		"inc": func() int {
			// 			return 1
			// 		},
			// 	},
			// ).ParseFiles("ui/game.html"))
			// tpl.Execute(w, game)

			err = template.Must(funcTemplate.Clone()).Funcs(template.FuncMap{
				"inc": func() int {
					return 1
				},
			}).Execute(w, game)
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func LoadPendingGame(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		gameId, _ := strconv.Atoi(r.FormValue("gameId"))
		game, err := g.GetSingleGame(int64(gameId))

		if err != nil {
			game.Message = err.Error()
		}

		t, _ := template.ParseFiles("ui/game.html")
		t.Execute(w, game)
	case "POST":
		rows, _ := strconv.Atoi(r.FormValue("rows"))
		columns, _ := strconv.Atoi(r.FormValue("columns"))
		mines, _ := strconv.Atoi(r.FormValue("mines"))

		user, found := g.Cache.Get("USER_SESSION")

		if !found {
			//No user logged in
			user := models.User{}
			user.Message = "Session lost, please login again!"
			t, _ := template.ParseFiles("ui/login.html")
			t.Execute(w, user)
			return
		}

		game := models.Game{
			UserId:       int64(user.(models.User).UserId),
			TimeConsumed: 0,
			Status:       "Pending",
			Rows:         rows,
			Columns:      columns,
			Mines:        mines,
		}

		newSpots, err := g.CreateGame(&game)

		if err != nil {
			game.Message = err.Error()
		} else {
			game.Message = "Game Created Successfully "
			game.Spots = *newSpots

			t, err := template.ParseFiles("ui/game.html")
			if err != nil {
				game.Message = err.Error()
			}

			t.Execute(w, game)
			// t, err := funcTemplate.Clone()
			// t = t.Funcs(template.FuncMap{
			// 	"increase": func(i int) int {
			// 		i++

			// 		return i
			// 	},
			// })
			// if err != nil {
			// 	fmt.Println(err.Error())
			// 	game.Message = err.Error()
			// }
			// err = t.Execute(w, game)
			// err = template.Must(funcTemplate.Clone()).Funcs(template.FuncMap{
			// 	"increase": func(i int) int {
			// 		i++

			// 		return i
			// 	},
			// }).Execute(w, game)

			if err != nil {
				fmt.Println(err.Error())
				game.Message = err.Error()
			}
		}
	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	g.Cache.Flush()

	getLoggedinUser(w)
}

func getLoggedinUser(w http.ResponseWriter) *models.User {
	obj, found := g.Cache.Get("USER_SESSION")

	if found {
		user := obj.(models.User)

		return &user
	}

	t, _ := template.ParseFiles("ui/login.html")
	t.Execute(w, models.User{
		Message: "Please loging",
	})

	return nil
}
