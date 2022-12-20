func GetStatistics(w http.ResponseWriter, r *http.Request) {
	// retrieve the transactions from the last 60 seconds
	sixtySecondsAgo := time.Now().Add(-time.Minute)
	cur, err := db.Find(context.TODO(), bson.M{"timestamp": bson.M{"$gte": sixtySecondsAgo}})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cur.Close(context.TODO())

	// calculate the statistics
	var sum, min, max float64
	count := 0
	for cur.Next(context.TODO()) {
		var transaction Transaction
		if err := cur.Decode(&transaction); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		amount := transaction.Amount
		sum += amount
		min = math.Min(min, amount)
		max = math.Max(max, amount)
		count++
	}

	// return the statistics
	stats := Statistics{Sum: sum, Min: min, Max: max, Count: count}
	if count > 0 {
		stats.Avg = sum / float64(count)
	}
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
