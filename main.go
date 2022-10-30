package main

import (
	"bytes"
	mail "github.com/xhit/go-simple-mail/v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type DataSmtp struct {
	Host     string
	Port     int
	Username string
	Password string
}

type Person struct {
	Mail     string
	Name     string
	FirstDay string
	LastYear string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func sendMailMessage(persons []Person, mailData DataSmtp) {
	server := mail.NewSMTPClient()

	server.Host = mailData.Host
	server.Port = mailData.Port
	server.Username = mailData.Username
	server.Password = mailData.Password
	server.Encryption = mail.EncryptionSSL

	smtpClient, err := server.Connect()
	check(err)

	// Create email
	email := mail.NewMSG()

	emailHost := strings.Replace(mailData.Host, "smtp.", "@", 1)

	email.SetFrom(mailData.Username + emailHost)

	for _, person := range persons {
		tmpl, _ := template.ParseFiles("templates/mailing_letter.html")
		// create a new file

		var out bytes.Buffer
		tmpl.Execute(&out, person)

		htmlBody := out.String()

		email.SetBody(mail.TextHTML, htmlBody)
		email.AddTo(person.Mail)
		email.SetSubject("Happy Birthday!")
		// Send email
		err = email.Send(smtpClient)
		check(err)

		os.Remove("templates/changed_mailing_letter.html")
	}

	//email.AddAttachment("super_cool_file.png")

}

func sendingCongratulations(wg *sync.WaitGroup) {
	// Отправление поздравлений. С некой периодичностью обращается
	// к базе, и если сегодняшний день соответствует дню рождения, отправляет поздравление по почте

	defer wg.Done()

	var mailData DataSmtp

	mailData.Username = GoDotEnvVariable("Username")
	mailData.Password = GoDotEnvVariable("Password")
	mailData.Host = GoDotEnvVariable("smtpHost")
	mailData.Port, _ = strconv.Atoi(GoDotEnvVariable("smtpPort"))

	for {
		currentTime := time.Now()
		date := strconv.Itoa(currentTime.Day()) + "." + strconv.Itoa(int(currentTime.Month()))
		currentYear := strconv.Itoa(currentTime.Year())
		listPersons := GetData()

		var mailingList []Person
		for _, person := range listPersons {

			datePerson := person.FirstDay[0:5]
			lastYearPerson := person.LastYear

			if datePerson == date && currentYear != lastYearPerson {
				mailingList = append(mailingList, person)
			}
		}

		year := strconv.Itoa(currentTime.Year())
		if len(mailingList) != 0 {
			sendMailMessage(mailingList, mailData)
			UpdateLastYear(year, mailingList)

		}

		time.Sleep(10 * time.Second)
	}

}

func mailListHandler(writer http.ResponseWriter, request *http.Request) {
	//Обработчик веб-сервера, отвечающий за вывод списка людей, которым будут рассылаться
	// поздравления. Так же здесь есть возможность добавить нового человека

	if request.Method == "POST" {
		err := request.ParseForm()
		if err != nil {
			log.Println(err)
		}

		mailPerson := request.FormValue("mail")
		name := request.FormValue("name")
		firstDay := request.FormValue("first day")

		if len(mailPerson) != 0 && len(name) != 0 && len(firstDay) != 0 {
			AddData(mailPerson, name, firstDay)
		}

	}

	mailingList := GetData()
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(writer, mailingList)
}

func WebServer(wg *sync.WaitGroup) {
	//Веб-сервер

	defer wg.Done()
	mux := http.NewServeMux()
	mux.HandleFunc("/", mailListHandler)

	log.Println("Запуск веб-сервера на http://127.0.0.1:8000")
	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}

func main() {
	//Работа веб-сервера и отправка поздравлений работают одновременно
	// с помощью го-рутин

	var wg sync.WaitGroup

	wg.Add(2)

	go WebServer(&wg)
	go sendingCongratulations(&wg)

	wg.Wait()

}
