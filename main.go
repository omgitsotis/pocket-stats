package main

import (
    "log"
    client "github.com/omgitsotis/pocket-stats/client"
)

var code string
var accessToken string

func main() {
    log.Fatal(client.ServeAPI())
}

// func healthcheck(w http.ResponseWriter, r *http.Request) {
//     fmt.Fprintln(w, "pocket-app healthy")
//     return
// }

// func recieveAuth(w http.ResponseWriter, r *http.Request) {
//     type Params struct {
//         ConsumerKey string `json:"consumer_key"`
//         Code string `json:"code"`
//     }
//
//     params := Params{"74935-9d486f66d2999047b61328f3", code}
//     b, err := json.Marshal(params)
//     if err != nil {
//         fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     req, err := http.NewRequest("POST", "https://getpocket.com/v3/oauth/authorize", bytes.NewBuffer(b))
// 	if err != nil {
// 		fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
//
// 	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	req.Header.Set("X-Accept", "application/json")
//
//     res, err := client.Do(req)
//     if err != nil {
//         fmt.Println("error doing request:", err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     defer res.Body.Close()
//     type Response struct {
//         AccessToken string `json:"access_token"`
//         Username string `json:"username"`
//     }
//
//     var response Response
//     err = json.NewDecoder(res.Body).Decode(&response)
//     if err != nil {
//         fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     fmt.Println(response.AccessToken)
//
//     type NewParam struct {
//         ConsumerKey string `json:"consumer_key"`
//         AccessToken string `json:"access_token"`
//         Tag string `json:"tag"`
//     }
//
//     newParam := NewParam{"74935-9d486f66d2999047b61328f3", response.AccessToken, "read now"}
//     b, err = json.Marshal(newParam)
//     if err != nil {
//         fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     req, err = http.NewRequest("POST", "https://getpocket.com/v3/get", bytes.NewBuffer(b))
// 	if err != nil {
// 		fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
//
// 	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	req.Header.Set("X-Accept", "application/json")
//
//     res, err = client.Do(req)
//     if err != nil {
//         fmt.Println("error doing request:", err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     defer res.Body.Close()
//     body, err := ioutil.ReadAll(res.Body)
// 	fmt.Println("post:\n", string(body))
// }
//
//
// func keepLines(s string, n int) string {
//     result := strings.Join(strings.Split(s, "\n")[:n], "\n")
//     return strings.Replace(result, "\r", "", -1)
// }
//
// func getAuth(w http.ResponseWriter, r *http.Request) {
//     type Params struct {
//         ConsumerKey string `json:"consumer_key"`
//         RedirectURI string `json:"redirect_uri"`
//     }
//
//     params := Params{"74935-9d486f66d2999047b61328f3", "http://localhost:8080/auth/recieved"}
//     b, err := json.Marshal(params)
//     if err != nil {
//         fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     req, err := http.NewRequest("POST", "https://getpocket.com/v3/oauth/request", bytes.NewBuffer(b))
// 	if err != nil {
// 		fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
// 	}
//
// 	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	req.Header.Set("X-Accept", "application/json")
//
//     res, err := client.Do(req)
//     if err != nil {
//         fmt.Println("error doing request:", err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     defer res.Body.Close()
//     type Response struct {
//         Code string `json:"code"`
//     }
//
//     var response Response
//     err = json.NewDecoder(res.Body).Decode(&response)
//     if err != nil {
//         fmt.Println(err.Error())
//         w.WriteHeader(http.StatusBadRequest)
// 		return
//     }
//
//     fmt.Println(response.Code)
//     code = response.Code
//     u := fmt.Sprintf(
//         "https://getpocket.com/auth/authorize?request_token=%s&redirect_uri=%s",
//         response.Code,
//         "http://localhost:8080/auth/recieved",
//     )
//     fmt.Fprintln(w, u)
// }
