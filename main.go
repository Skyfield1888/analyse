package main

import (
	"net/http"

	"github.com/AltSimon/analyse/backend"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("frontend/")))
	//Table template endpoints
	http.HandleFunc("/table/add-row", backend.AddRowHandler)
	http.HandleFunc("/table/add-column", backend.AddColumnHandler)
	http.HandleFunc("/table/set-label", backend.RenameColumnHandler)
	http.HandleFunc("/table/save", backend.SavingHandler)
	http.HandleFunc("/table/new", backend.NewTableHandler)
	//use table testing without entire ui
	http.HandleFunc("/table/generate", backend.GenericTableHandler)

	http.ListenAndServe(":8081", nil)
}
