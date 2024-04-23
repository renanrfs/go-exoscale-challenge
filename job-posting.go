package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"sync"
)

// Job represents a job posting.
type Job struct {
	ID          string `json:"id"`
	Company     string `json:"company"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

// JobsMap holds the in-memory job postings.
var (
	JobsMap = make(map[string]Job)
	mu      sync.Mutex
)

// getJobs handles GET requests and returns all jobs.
func getJobs(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()
	jobs := make([]Job, 0, len(JobsMap))
	for _, job := range JobsMap {
		jobs = append(jobs, job)
	}
	json.NewEncoder(w).Encode(jobs)
}

// checks if the position is open.
func isPositionOpen(title string) bool {
	for _, job := range JobsMap {
		if job.Title == title {
			return true
		}
	}
	return false
}

// addJob handles POST requests to add a new job.
func addJob(w http.ResponseWriter, r *http.Request) {
	var job Job
	if err := json.NewDecoder(r.Body).Decode(&job); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if job.Company == "" || job.Title == "" || job.Description == "" {
		http.Error(w, "Company, Title and Description are required.", http.StatusBadRequest)
		return
	}

	if isPositionOpen(job.Title) {
		http.Error(w, "Position is already open", http.StatusBadRequest)
		return
	}

	job.ID = uuid.NewString()
	mu.Lock()
	JobsMap[job.ID] = job
	mu.Unlock()
	getJobs(w, r)
}

// deleteJob handles DELETE requests to remove a job by ID.
func deleteJob(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/jobs/"):]
	mu.Lock()
	delete(JobsMap, id)
	mu.Unlock()
	getJobs(w, r)
}

// main sets up the HTTP server and routes.
func main() {
	http.HandleFunc("/jobs", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getJobs(w, r)
		case http.MethodPost:
			addJob(w, r)
		default:
			http.Error(w, "Unsupported request method.", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/jobs/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleteJob(w, r)
		} else {
			http.Error(w, "Unsupported request method.", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
