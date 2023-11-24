package otpvalidation

import (
	"fmt"
	"net/smtp"
)

func main() {
	// Sender's email address and password
	from := "niki16809@gmail.com"
	password := "xagg yzsu kkvd jbrr"

	// Recipient's email address
	to := ""

	// SMTP server details
	smtpServer := "bhoomikasabalur@gmail.com"
	smtpPort := 587

	// Message content
	message := []byte("Subject: Test Email\n\nThis is a test email body.")

	// Authentication information
	auth := smtp.PlainAuth("", from, password, smtpServer)

	// SMTP connection
	smtpAddr := fmt.Sprintf("%s:%d", smtpServer, smtpPort)
	err := smtp.SendMail(smtpAddr, auth, from, []string{to}, message)
	if err != nil {
		fmt.Println("Error sending email:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
