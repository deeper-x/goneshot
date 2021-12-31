```go

func register(w http.ResponseWriter, r *http.Request) {
	client := r.URL.Query().Get("client")
	validToken, err := GetJWT(client)
	if err != nil {
		fmt.Println("Failed to generate token")
	}

	fmt.Fprint(w, string(validToken))
}

func secretPage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Super Secret test Information")
}

func handleRequests() {
	http.HandleFunc("/", register)
	http.Handle("/check", isAuthorized(secretPage))
    
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	handleRequests()
}


```