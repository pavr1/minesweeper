package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func main() {
	//gate, err := gate.Start()

	// if err != nil {
	// 	fmt.Println(err.Error)
	// }

	// err = gate.CreateUser("Pablo", "Villalobos", "villaPab123!")

	// if err != nil {
	// 	fmt.Println(err.Error)
	// }

	http.HandleFunc("/", loadMain)

	log.Fatal(http.ListenAndServe(":8080", nil))

	fmt.Println("Server started at port 8080")

	select {}
}

func loadMain(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("./main_page.html")
	t.Execute(w, nil)
}
