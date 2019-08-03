package handler

import (
	"bytes"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"shuStudentCourse/infrastructure"
	"shuStudentCourse/service/token"
	"strings"
	"sync"
	"unicode"
)

func PingPongHandler(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("pong"))
}

func fetchPage(
	semesterId string,
	studentId string) *goquery.Document {
	response, _ := http.Get(os.Getenv("COURSE_SELECTION_ADDRESS_URL") + "?id=" + semesterId)
	body, _ := ioutil.ReadAll(response.Body)
	var courseSelectionUrlJson struct {
		Url string `json:"url"`
	}
	_ = json.Unmarshal(body, &courseSelectionUrlJson)
	data, _ := json.Marshal(struct {
		Url string `json:"url"`
	}{
		courseSelectionUrlJson.Url + "/StudentQuery/CtrlViewQueryCourseTable",
	})
	request, _ := http.NewRequest(
		"POST",
		os.Getenv("PROXY_ADDRESS")+"/get",
		bytes.NewBuffer(data))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+token.GenerateJWT(studentId))
	response, _ = http.DefaultClient.Do(request)
	doc, _ := goquery.NewDocumentFromReader(response.Body)
	return doc
}

func analyzePage(doc *goquery.Document, semesterId string, studentId string) []int64 {
	type courseTeacherPair struct {
		CourseId  string
		TeacherId string
	}
	var courseTeacherPairArray []courseTeacherPair
	var wg sync.WaitGroup
	doc.Find("tr").Each(func(i int, selection *goquery.Selection) {
		rowId := strings.TrimSpace(selection.Find("td").First().Text())
		if i >= 3 && len([]rune(rowId)) > 0 && unicode.IsUpper([]rune(rowId)[0]) {
			courseId := strings.TrimSpace(selection.Find("td:nth-child(2)").Text())
			teacherId := strings.TrimSpace(selection.Find("td:nth-child(4)").Text())
			wg.Add(1)
			courseTeacherPairArray = append(courseTeacherPairArray, courseTeacherPair{courseId, teacherId})
		}
	})
	var result []int64
	for _, pair := range courseTeacherPairArray {
		go func(courseId string, teacherId string) {
			defer wg.Done()
			response, _ := http.Get(os.Getenv("COURSE_INFO_URL") + "?semester_id=" + semesterId + "&course_id=" + courseId + "&in_course_teacher_id=" + teacherId)
			body, _ := ioutil.ReadAll(response.Body)
			responseStruct := struct {
				Id int64 `json:"id"`
			}{}
			_ = json.Unmarshal(body, &responseStruct)
			_, _ = infrastructure.DB.Exec(`
			INSERT INTO StudentCourse(semester_id, student_id, course_by_teacher_id) 
			VALUES ($1,$2,$3);
			`, semesterId, studentId, responseStruct.Id)
			result = append(result, responseStruct.Id)
		}(pair.CourseId, pair.TeacherId)
	}
	log.Println("student:", studentId, " semester:", semesterId, " fetched")
	wg.Wait()
	return result
}

func StudentCoursesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	tokenInHeader := r.Header.Get("Authorization")
	if len(tokenInHeader) <= 7 {
		w.WriteHeader(401)
		return
	}
	studentId := token.StudentIdForToken(tokenInHeader[7:])
	semesterId := r.URL.Query().Get("semester_id")
	row := infrastructure.DB.QueryRow(`
	SELECT array_agg(course_by_teacher_id)
	FROM studentCourse
	WHERE student_id=$1 AND semester_id=$2;
	`, studentId, semesterId)
	var result pq.Int64Array
	queryError := row.Scan(&result)
	if queryError != nil || len(result) == 0 {
		result = analyzePage(fetchPage(semesterId, studentId), semesterId, studentId)
	}
	body, _ := json.Marshal(result)
	_, _ = w.Write(body)
}
