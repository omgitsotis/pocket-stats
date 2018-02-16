package client

import (
	"net/http"
	"fmt"
	"os"
    "path/filepath"

	"github.com/gorilla/mux"
	pocket "github.com/omgitsotis/pocket-stats/pocket"
)

const (
	cdnReact           = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react.min.js"
	cdnReactDom        = "https://cdnjs.cloudflare.com/ajax/libs/react/15.5.4/react-dom.min.js"
	cdnBabelStandalone = "https://cdnjs.cloudflare.com/ajax/libs/babel-standalone/6.24.0/babel.min.js"
	cdnAxios           = "https://cdnjs.cloudflare.com/ajax/libs/axios/0.16.1/axios.min.js"
)

const indexHTML = `
<!DOCTYPE HTML>
<html>
  <head>
    <meta charset="utf-8">
    <title>Simple Go Web App</title>
  </head>
  <body>
    <div id='root'></div>
    <script src="` + cdnReact + `"></script>
    <script src="` + cdnReactDom + `"></script>
    <script src="` + cdnBabelStandalone + `"></script>
    <script src="` + cdnAxios + `"></script>
    <script src="/public/js/app.jsx" type="text/babel"></script>
  </body>
</html>
`

func NewClient() *Client {
	p := pocket.NewPocket("74935-9d486f66d2999047b61328f3")
	return &Client{Pocket:p}
}

func ServeAPI() error {
	c := NewClient()
	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer( http.Dir("./public") ) ))
	r.Methods("GET").Path("/").HandlerFunc(indexHandler)
	r.Methods("GET").Path("/auth").HandlerFunc(c.Retrieve)
	r.Methods("GET").Path("/auth/recieved").HandlerFunc(c.Authorise)
	r.Methods("GET").Path("/data").HandlerFunc(c.GetData)


	fmt.Println("Created router")
	return http.ListenAndServe(":8080", r)
}

func indexHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintf(w, indexHTML)
}
