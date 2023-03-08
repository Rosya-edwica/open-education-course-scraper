package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/Rosya-edwica/open-education-course-scraper/src/api"
	"github.com/Rosya-edwica/open-education-course-scraper/src/logger"
)

const GroupSize = 10

func main() {
	logger.Log.Println("Старт парсинга")
	start := time.Now().Unix()
	pages := groupPages()
	for _, page := range pages {
		var wg sync.WaitGroup
		wg.Add(len(page))
		for _, url := range page {
			go parsePage(url, &wg)
		}
		wg.Wait()
	}
	fmt.Println("Время выполнения программы:", time.Now().Unix()-start)
	logger.Log.Printf("Время выполнения программы: %d", time.Now().Unix()-start)
}

func parsePage(url string, wg *sync.WaitGroup) {
	logger.Log.Printf("Пагинация: %s", url)
	coursesLinks := api.GetCoursesList(url)
	for _, item := range coursesLinks {
		api.SaveCourse(item)
	}
	wg.Done()
}

func generatePageListLinks(count int) (links []string) {
	startPage := 1
	for i := startPage; i < count; i++ {
		links = append(links, fmt.Sprintf("https://courses.openedu.ru/api/courses/v1/courses/?page=%d", i))
	}
	return
}

func groupPages() (groups [][]string) {
	links := generatePageListLinks(990)
	count := len(links)
	var limit int
	for i := 0; i < count; i += GroupSize {
		limit += GroupSize
		if limit > count {
			limit = count
		}
		group := links[i:limit]
		groups = append(groups, group)
	}
	return
}
