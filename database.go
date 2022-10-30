package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func AddData(PersonData ...string) {
	db, err := sql.Open("sqlite3", "mail_list.db")
	check(err)

	defer db.Close()

	_, err = db.Exec("insert into mailing_list (mail, name, first_day, last_year) values ($1, $2, $3, $4)",
		PersonData[0], PersonData[1], PersonData[2], "")

	check(err)

}

func UpdateLastYear(lastYear string, PersonData []Person) {
	db, err := sql.Open("sqlite3", "mail_list.db")
	check(err)

	defer db.Close()
	for _, person := range PersonData {
		fmt.Println(person)
		_, err = db.Exec("UPDATE mailing_list SET last_year = $1 WHERE mail = $2",
			lastYear, person.Mail)

		check(err)
	}

}

func GetData() []Person {
	db, err := sql.Open("sqlite3", "mail_list.db")
	check(err)

	rows, err := db.Query("select * from mailing_list")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	var mailingList []Person

	for rows.Next() {
		person := Person{}
		err := rows.Scan(&person.Mail, &person.Name, &person.FirstDay, &person.LastYear)
		if err != nil {
			fmt.Println(err)
			continue
		}
		mailingList = append(mailingList, person)
	}

	return mailingList
}

func GoDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
