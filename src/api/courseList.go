package api

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

func GetCoursesList(url string) (links []string) {
	json := getJson(url)
	data := gjson.Get(json, "results").Array()
	for _, i := range data {
		link := createCourseAPIUrl(i.Get("id").String())
		links = append(links, link)
	}
	return
}

func createCourseAPIUrl(course_id string) (url string) {
	re := regexp.MustCompile(`:.*?\+.*?\+`) // Заберем ":eltech+ECON+" из строки "course-v1:eltech+ECON+spring_2021"
	finded := strings.ReplaceAll(re.FindString(course_id), ":", "")
	org_and_number := strings.Split(finded, "+")
	org := org_and_number[0]
	number := org_and_number[1]
	return fmt.Sprintf("https://openedu.ru/api/courses/export/%s/%s?format=json", org, number)
}
