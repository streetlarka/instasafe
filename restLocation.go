func ResetLocation(w http.ResponseWriter, r *http.Request) {
	// delete all locations from the database
	if _, err := db.DeleteMany(context.TODO(), bson.M{}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
