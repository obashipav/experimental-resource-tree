package main

import (
	"fmt"
	"github.com/OBASHITechnology/resourceList/web"
	"log"
	"net/http"
)

func main() {
	engine := web.Registration()

	fmt.Println("we are listening live at http://localhost:8080")
	if err := http.ListenAndServe(":8080", engine); err != nil {
		log.Fatal(err)
	}
}

