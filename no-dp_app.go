package main

import (
	//"database/sql"
	//"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type EventData struct {
	Type   string `json:"type"`
	Source string `json:"source"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}

/* Fill these details */
/*const (
	host     = "host"
	port     = 0000
	user     = "user"
	password = "password"
	dbname   = "postgres"
)*/

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func get_request(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Web Server Running!")
}

func post_request(w http.ResponseWriter, r *http.Request) {

	//Extracting data sent to /ingest
	var data EventData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	title := data.Title
	url := data.URL
	date := time.Now().Format("2006-01-02")

	//Inserting data into table.
	/*insert_query := `
	INSERT INTO jobs (title, url, timestamp) VALUES ($1, $2, $3)
	`
	_, err = db.Exec(insert_query, title, url, date)

	if err != nil {
		panic(err)
	}*/
	fmt.Printf(title, url, date)
	w.WriteHeader(http.StatusOK)
}

/*func exporter(db *sql.DB) {
	fname := "D:/job_list.csv"
	var err error
	isNewFile := false

	var f *os.File
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		f, err = os.Create(fname)
		isNewFile = true
	} else {
		f, err = os.OpenFile(fname, os.O_APPEND|os.O_WRONLY, 0644)
		info, _ := f.Stat()
		isNewFile = info.Size() == 0
	}

	if err != nil {
		panic(err)
	}
	defer f.Close()

	rows, err := db.Query("SELECT title, url, timestamp FROM jobs")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	w := csv.NewWriter(f)

	if isNewFile {
		w.Write([]string{"title", "url", "timestamp"})
	}

	for rows.Next() {
		var title, url string
		var timestamp time.Time

		err := rows.Scan(&title, &url, &timestamp)
		if err != nil {
			panic(err)
		}

		ts := timestamp.Format("2006-01-02 15:04:05")

		w.Write([]string{title, url, ts})
	}
	w.Flush()

	if err := w.Error(); err != nil {
		panic(err)
	}

	fmt.Println("Data Appended to CSV successfully")

}*/

/*func truncate(db *sql.DB) {
	db.Query("TRUNCATE TABLE jobs")
}*/

func main() {

	//Setting up connection with PostgreSQL
	/*psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}*/
	fmt.Println("Successfully connected!")

	//Trying to capture SIGINT
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)

	done := make(chan bool, 1)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		//exporter(db)
		//truncate(db)
		os.Exit(0)
	}()

	//Setting up router and calling functions based on GET and POST requests
	r := mux.NewRouter()
	r.HandleFunc("/", get_request).Methods("GET")
	r.HandleFunc("/ingest", func(w http.ResponseWriter, r *http.Request) {
		post_request(w, r)
	}).Methods("POST", "OPTIONS")

	fmt.Println("awaiting Signal")
	http.ListenAndServe(":8080", enableCORS(r))

	<-done
	fmt.Println("exiting")

}
