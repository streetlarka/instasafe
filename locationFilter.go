func GetStatistics(w http.ResponseWriter, r *http.Request) {
	// check if a location is set
	var location Location
	if err := db.FindOne(context.TODO(), bson.M{}).Decode(&location); err != nil {
		if err == mongo.ErrNoDocuments {
			// no location is set, allow access from any location
			location.City = ""
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// get the statistics for the given location
	var stats Statistics
	if err := db.Aggregate(context.TODO(), []bson.M{
		{"$match": bson.M{"timestamp": bson.M{"$gte": time.Now().Add(-time.Minute)}}},
		{"$group": bson.M{"_id": "$location.city", "sum": bson.M{"$sum": "$amount"}, "count": bson.M{"$sum": 1}}},
		{"$sort": bson.M{"sum": -1}},
		{"$limit": 1},
	}).Decode(&stats); err != nil {
		if err == mongo.ErrNoDocuments {
			// no transactions in the last 60 seconds
			stats.Sum = 0
			stats.Count = 0
			stats.Avg = 0
			stats.Min = 0
			stats.Max = 0
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// check if the location matches the set location
	if location.City != "" && location.City != stats.Location {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// return the statistics
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
