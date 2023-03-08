package api

import (
	"github.com/Rosya-edwica/open-education-course-scraper/src/logger"
	"github.com/Rosya-edwica/open-education-course-scraper/src/models"
	"github.com/Rosya-edwica/open-education-course-scraper/src/postgres"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
	"github.com/tidwall/gjson"
)

func SaveCourse(url string) {
	var course models.Course
	json := getJson(url)
	if json == "" {
		logger.Log.Printf("Ошибка! Пустой json для страницы: %s", url)
		return
	}
	course.Title = gjson.Get(json, "title").String()
	course.Description = gjson.Get(json, "description").String()
	course.StartedAt = gjson.Get(json, "started_at").String()
	course.FinishedAt = gjson.Get(json, "finished_at").String()
	course.Skills = strings.Join(ClearSkills(gjson.Get(json, "results").String()), "|")
	course.Image = gjson.Get(json, "image").String()
	course.Requirements = getRequirements(json)
	course.Url = gjson.Get(json, "external_url").String()
	course.LecturesCount = int(gjson.Get(json, "lectures").Int())
	course.Program = gjson.Get(json, "content").String()
	course.DurationInWeek = int(gjson.Get(json, "duration.value").Int())
	course.HasCertificate = getCert(json)
	course.Price = getPrice(course.Url)
	parseTeachers(&course, json)
	if course.Title != "" && !strings.Contains(course.Url, "https://openedu.ru/course/123") && course.DurationInWeek != 0 {
		postgres.AddCourse(course)
	}

}

func getCert(json string) (certificate string) {
	cert := gjson.Get(json, "cert").Bool()
	if cert {
		return "Да"
	} else {
		return "Нет"
	}
}

func getRequirements(json string) (requirements string) {
	var reqs []string
	for _, item := range gjson.Get(json, "requirements").Array() {
		reqs = append(reqs, item.String())
	}
	return strings.Join(reqs, "|")
}

func getPrice(url string) (salary int) {
	c := colly.NewCollector()
	c.OnHTML("span.rub-box", func(h *colly.HTMLElement) {
		digitPattern := regexp.MustCompile(`\d+`)

		digit, err := strconv.Atoi(digitPattern.FindString(h.Text))
		checkErr(err)
		salary = digit
		return
	})
	c.OnHTML("div.card.card-certificate.card-itmo div.card-body", func(h *colly.HTMLElement) {
		salary = findSalaryInText(h.Text)
		if salary != 0 {
			return
		}
	})
	c.Visit(url)
	return
}

func clearDescription(text string) (desr string) {
	unTagsPattern := regexp.MustCompile(`<.*?>`)
	desr = unTagsPattern.ReplaceAllString(text, "")
	return
}

func ClearSkills(text string) (skills []string) {
	itemTagsPattern := regexp.MustCompile(`<li>.*?<\/li>`)
	tagsToReplacePattern := regexp.MustCompile(`<li>|<\/li>|;|·`)
	items := itemTagsPattern.FindAllString(text, -1)
	for _, item := range items {
		skills = append(skills, tagsToReplacePattern.ReplaceAllString(item, ""))
	}
	return
}

func parseTeachers(course *models.Course, json string) {
	data := gjson.Get(json, "teachers").Array()
	var names, descr, imgs []string
	for _, i := range data {
		names = append(names, i.Get("display_name").String())
		descr = append(descr, i.Get("description").String())
		imgs = append(imgs, i.Get("image").String())
	}
	course.TeachersDescriptions = strings.Join(descr, "|")
	course.TeachersName = strings.Join(names, "|")
	course.TeachersImages = strings.Join(imgs, "|")
}

func findSalaryInText(text string) (digit int) {
	digitPattern := regexp.MustCompile(`\d+`)
	salaryPattern := regexp.MustCompile(`\d+ ₽|\d+ руб`)

	salary := salaryPattern.FindString(text)
	digit, _ = strconv.Atoi(digitPattern.FindString(salary))
	return
}
