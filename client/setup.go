package client

import (
	"log"
	"net/http"

	"github.com/omgitsotis/pocket-stats/pocket"
	"github.com/omgitsotis/pocket-stats/pocket/dao/sqlite"
)

// func NewClient() *Client {
// 	p := pocket.NewPocket("74935-9d486f66d2999047b61328f3")
// 	return &Client{Pocket: p}
// }

func ServeAPI() error {
	// c := NewClient()
	// r := mux.NewRouter()
	//
	// corsObj := handlers.AllowedOrigins([]string{"*"})
	//
	// r.Methods("GET").Path("/auth").HandlerFunc(c.Retrieve)
	// r.Methods("GET").Path("/auth/recieved").HandlerFunc(c.Authorise)
	// r.Methods("GET").Path("/data").HandlerFunc(c.GetData)
	//
	// fmt.Println("Created router")
	// return http.ListenAndServe(":8082", handlers.CORS(corsObj)(r))

	sqlite, err := sqlite.NewSQLiteDAO("./database/pocket.db")
	if err != nil {
		return err
	}

	p := pocket.NewPocket(
		"74935-9d486f66d2999047b61328f3",
		&http.Client{},
		sqlite,
	)

	r := NewRouter(p)
	r.Handle("send auth", sendAuth)
	r.Handle("data init", initDB)
	r.Handle("auth cached", saveToken)
	r.Handle("data get", getStatistics)
	r.Handle("data update", updateDB)

	http.Handle("/", r)
	http.HandleFunc("/auth/recieved", r.RecievedAuth)

	log.Println("Serving application")
	return http.ListenAndServe(":4000", nil)

}
