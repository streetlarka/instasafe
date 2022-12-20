func AddTransaction(w http.ResponseWriter, r *http.Request) {
	// parse the request body
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// validate the transaction timestamp
	if transaction.Timestamp.After(time.Now()) {
		http.Error(w, "transaction timestamp is in the future", http.StatusUnprocessableEntity)
		return
	}

	// check if the transaction is older than 60 seconds
	if time.Now().Sub(transaction.Timestamp).Seconds() > 60 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// add the transaction to the database
	if err := db.InsertOne(context.TODO(), transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
