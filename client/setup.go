package client

import (
	"net/http"
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

	r := NewRouter()
	r.Handle("send auth", sendAuth)
	r.Handle("data get", getData)
	r.Handle("auth cached", saveToken)

	http.Handle("/", r)
	http.HandleFunc("/auth/recieved", r.RecievedAuth)

	return http.ListenAndServe(":4000", nil)

}
