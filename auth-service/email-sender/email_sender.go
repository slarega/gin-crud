package email_sender

//email_sender

import (
	"fmt"
	gomail "github.com/wneessen/go-mail"
	"log"
	"net/mail"
)

const (
	emailSender = "emailsender@gmail.com"
	smtpHost    = "smtp.mailtrap.io"
	smtpPort    = 2525
	smtpUser    = "your_mail@gmail.com"
	smtpPass    = "your_pass"
)

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func mailClient() *gomail.Client {
	client, err := gomail.NewClient(smtpHost,
		gomail.WithPort(smtpPort),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(smtpUser),
		gomail.WithPassword(smtpPass))
	if err != nil {
		log.Fatalf("Не удалось создать mail client: %s", err)
	}
	return client
}

func EmailSender(letterRecipient string) {
	//Проверка почты получателя
	validateEmail := validEmail(letterRecipient)
	if !validateEmail {
		fmt.Println("Адрес получателя недействителен")
		return
	}
	//Создание сообщения
	m := gomail.NewMsg()
	if err := m.From(emailSender); err != nil {
		fmt.Printf("Не удалось установить почту отправителя: %s", err)
	}
	if err := m.To(letterRecipient); err != nil {
		fmt.Printf("Не удалось установить почту получателя: %s", err)
	}
	m.Subject("email warning fake.com")
	m.SetBodyString(
		gomail.TypeTextPlain,
		"IP адрес вашего устройства изменился. Если это не вы - повторно авторизуйтесь на сайте",
	)

	client := mailClient()
	// Отправка письма
	if err := client.DialAndSend(m); err != nil {
		log.Fatalf("Не удалось отправить письмо: %s", err)
	}
	log.Printf("Письмо отправлено")
}
