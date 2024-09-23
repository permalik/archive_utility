package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load .env", err)
	}
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	from := os.Getenv("FROM")
	y, m, d := time.Now().Date()
	today := fmt.Sprintf("%d %s %d", y, m, d)
	recipient := "tymalik@pm.me"
	to := fmt.Sprintf("To: %s\r\n", recipient)
	subject := fmt.Sprintf("Subject: Atlas Redeployment for %s:\r\n", today)
	body := fmt.Sprintf("- archive_ui\n- archive_resume\r\n")
	msg := []byte(to + subject + "\r\n" + body)
	recipients := []string{recipient}

	hostname := os.Getenv("HOSTNAME")
	auth := smtp.PlainAuth("", username, password, hostname)

	err = smtp.SendMail(hostname+":587", auth, from, recipients, msg)
	if err != nil {
		log.Fatal(err)
	}
}
