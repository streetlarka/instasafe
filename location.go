type Location struct {
	City string `json:"city"`
}

func SetLocation(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var location Location
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// save the location to the database
	if err := db.InsertOne(context.TODO(), location); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
