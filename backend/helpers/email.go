package helpers

import (
    "gopkg.in/gomail.v2"
)

func SendEmail(to string, subject string, body string) error {
    m := gomail.NewMessage()

    m.SetHeader("From", "webdevviscata@gmail.com")
    m.SetHeader("To", to)
    m.SetHeader("Subject", subject)
    m.SetBody("text/html", body)

    d := gomail.NewDialer("smtp.gmail.com", 465, "webdevviscata@gmail.com", "lpfy qhuo tptk iivt")
		d.SSL = true

    if err := d.DialAndSend(m); err != nil {
        return err
    }

    return nil
}
