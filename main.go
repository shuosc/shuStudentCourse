package main

import (
	"log"
	"net/http"
	"os"
	"shuStudentCourse/handler"
)

func main() {
	http.HandleFunc("/ping", handler.PingPongHandler)
	http.HandleFunc("/student-courses", handler.StudentCoursesHandler)
	err := http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
