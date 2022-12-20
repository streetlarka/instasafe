func main() {
	// connect to the database
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.TODO())
	db = client.Database("transactions")

	// set up the routes
	r := mux.NewRouter()
	r.HandleFunc("/transactions", AddTransaction).Methods("POST")
	r.HandleFunc("/statistics", GetStatistics).Methods("GET")
	r.HandleFunc("/transactions", DeleteTransactions).Methods("DELETE")
	r.HandleFunc("/location", SetLocation).Methods("POST")
	r.HandleFunc("/location", ResetLocation).Methods("DELETE")

	// start the server
	log.Println("Starting server on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
