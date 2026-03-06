
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
)

// Student represents a student
type Student struct {
	ID   int
	Name string
}

// AttendanceRecord represents a student's attendance on a specific date
type AttendanceRecord struct {
	StudentID int
	Date      string
	Present   bool
}

var (
	students      []Student
	attendance    []AttendanceRecord
	templates     *template.Template
	studentsMux   sync.Mutex
	attendMux     sync.Mutex
	nextStudentID int
)

func main() {
	// Initialize some sample data
	students = []Student{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}
	nextStudentID = 4

	// Load templates
	templates = template.Must(template.ParseGlob("templates/*.html"))

	// Handlers
	http.HandleFunc("/", attendanceHandler)
	http.HandleFunc("/submit", submitHandler)
	http.HandleFunc("/add-student", addStudentHandler)

	fmt.Println("Server starting at port 8080")
	http.ListenAndServe(":8080", nil)
}

func attendanceHandler(w http.ResponseWriter, r *http.Request) {
	studentsMux.Lock()
	defer studentsMux.Unlock()

	err := templates.ExecuteTemplate(w, "index.html", students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func addStudentHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	studentName := r.FormValue("name")

	studentsMux.Lock()
	defer studentsMux.Unlock()

	students = append(students, Student{ID: nextStudentID, Name: studentName})
	nextStudentID++

	http.Redirect(w, r, "/", http.StatusFound)
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()

	attendMux.Lock()
	defer attendMux.Unlock()

	for _, student := range students {
		present := r.FormValue(fmt.Sprintf("student_%d", student.ID)) == "present"
		record := AttendanceRecord{
			StudentID: student.ID,
			Date:      r.FormValue("date"),
			Present:   present,
		}
		attendance = append(attendance, record)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
