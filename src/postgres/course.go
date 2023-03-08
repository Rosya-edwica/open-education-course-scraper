package postgres

import (
	"database/sql"
	"fmt"
	"github.com/Rosya-edwica/open-education-course-scraper/src/logger"
	"github.com/Rosya-edwica/open-education-course-scraper/src/models"
	"strings"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "postgres"
	password = "admin"
	dbname   = "open_education"
)

func connect() *sql.DB {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", conn)
	checkErr(err)
	return db
}

func AddCourse(course models.Course) {
	db := connect()
	defer db.Close()

	smt := `INSERT INTO course (Название_курса, Описание_курса, Сертификат, Продолжительность_в_неделях, Ссылка_на_страницу, Дата_начала, Дата_окончания, Картинка_курса, Цена, Навыки, Что_нужно_знать, Фото_преподавателя, ФИО_преподавателя, О_преподавателе, Программа) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)`
	tx, _ := db.Begin()
	_, err := db.Exec(smt, course.Title, course.Description, course.HasCertificate, course.DurationInWeek, course.Url, strings.ReplaceAll(course.StartedAt, "-", "."), strings.ReplaceAll(course.FinishedAt, "-", "."), course.Image, course.Price, course.Skills, course.Requirements, course.TeachersImages, course.TeachersName, course.TeachersDescriptions, course.Program)
	if err != nil {
		if err.Error() != `pq: повторяющееся значение ключа нарушает ограничение уникальности "course_pkey"` {
			fmt.Println(err)
			er := err.(*pq.Error)
			logger.Log.Fatal(fmt.Sprintf("Code:%s\tHint:%s\tMessage:%s\tWhere:%s", er.Code, er.Hint, er.Message, er.Where))
			panic(fmt.Sprintf("[%s]", err.Error()))
		} else {
			tx.Commit()
			db.Close()
			logger.Log.Printf("Не смогли сохранить курс: %s. Причина: %s", course.Title, err)
		}
	} else {
		tx.Commit()
		db.Close()
		logger.Log.Printf("Успех! Добавили курс [%s]", course.Title)
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
