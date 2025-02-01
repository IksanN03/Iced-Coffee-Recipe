package utils

import (
	"crypto/tls"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

// 👇 Email template parser
func SendEmail(email string, link string) error {

	// Sender data.
	SMTP_HOST := os.Getenv("SMTP_HOST")
	SMTP_PORT, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	SMTP_SENDER_NAME := os.Getenv("SMTP_SENDER_NAME")
	SMTP_EMAIL := os.Getenv("SMTP_EMAIL")
	SMTP_PASSWORD := os.Getenv("SMTP_PASSWORD")

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", SMTP_SENDER_NAME)
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", `Verifikasi Email`)
	mailer.SetBody("text/html", ` 
<img alt="Logo" src="" style="width:20%"/><br/>
<font color="black"><strong>Hai `+email+`, </strong> selamat bergabung menjadi bagian dari kami.</font> 
<br /> <br /> Silahkan akses link di bawah untuk login<br/><br/>
<a style="background-color:#0084C8; color:white; padding:10px 20px; border-radius:30px; text-align:center;text-decoration:none" href="`+link+`">
     Verifikasi Email
</a><br/><br/><br/>`)

	d := gomail.NewDialer(
		SMTP_HOST,
		SMTP_PORT,
		SMTP_EMAIL,
		SMTP_PASSWORD,
	)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Send Email
	if err := d.DialAndSend(mailer); err != nil {
		return err
	}
	return nil
}
