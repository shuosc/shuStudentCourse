#!/usr/bin/env bash
export COURSE_SELECTION_URL="http://cloud.shu.xn--io0a7i:30000/api/course-selection-url"
export DB_ADDRESS="postgresql://localhost:5432/shuStudentCourse?sslmode=disable"
export PROXY_ADDRESS="http://cloud.shu.xn--io0a7i:30000/api/shu-course-proxy"
export COURSE_INFO_URL="http://cloud.shu.xn--io0a7i:30000/api/course"
export PORT="8001"
gin -p 8000 run main.go