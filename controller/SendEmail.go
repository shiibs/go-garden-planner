package controller

import (
	"bytes"
	"fmt"
	"net/smtp"
	"os"
	"text/template"
	"time"

	"github.com/shiibs/go-garden-planner/database"
	"github.com/shiibs/go-garden-planner/model"
)

func SendEmail(userEmail string, gardenLayout model.GardenLayout) error {
	// Sender data
	from := os.Getenv("EmailID")
	password := os.Getenv("EmailPassword")

	// Receiver data
	to := []string{
		userEmail,
	}

	// SMTP config
	smtpHost := os.Getenv("SmtpHost")
	smtpPort := os.Getenv("SmtpPort")

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Parsing HTML template
	tmpl, err := template.ParseFiles("static/garden.html")
	if err != nil {
		return fmt.Errorf("error parsing garden file, %w", err)
	}

	var body bytes.Buffer

	// MIME headers
	mimeHeaders := "MIME-version: 1.0;\nContent-Type: multipart/mixed; boundary=\"boundary\"\n\n"

	body.WriteString("Subject: Garden Schedule\n")
	body.WriteString(mimeHeaders)

	// HTML text part
	body.WriteString("--boundary\n")
	body.WriteString("Content-Type: text/html; charset=\"UTF-8\"\n\n")

	err = tmpl.Execute(&body, gardenLayout)
	if err != nil {
		return fmt.Errorf("error executing tmpl, %w", err)
	}

	body.WriteString("\n--boundary\n")
	body.WriteString("Content-Type: text/calendar; charset=\"UTF-8\"; method=REQUEST\n\n")

	// Create iCalendar content
	icalendarContent, err := createICalendar(gardenLayout)
	if err != nil {
		return fmt.Errorf("error creating iCalendar, %w", err)
	}
	body.WriteString(icalendarContent)

	body.WriteString("\n--boundary--\n")

	// Sending the email with attachments
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		return fmt.Errorf("error sending mail, %w", err)
	}

	return nil
}

func createICalendar(gardenLayout model.GardenLayout) (string, error) {
	var icalendar bytes.Buffer

	icalendar.WriteString("BEGIN:VCALENDAR\n")
	icalendar.WriteString("VERSION:2.0\n")
	icalendar.WriteString("PRODID:-//Your Organization//Your Product//EN\n")

	// Start date event
	icalendar.WriteString("BEGIN:VEVENT\n")
	icalendar.WriteString(fmt.Sprintf("UID:%d-start@yourdomain.com\n", gardenLayout.ID))
	icalendar.WriteString(fmt.Sprintf("DTSTAMP:%s\n", formatICalDate(time.Now())))
	icalendar.WriteString(fmt.Sprintf("DTSTART:%s\n", formatICalDate(gardenLayout.StartDate)))
	icalendar.WriteString(fmt.Sprintf("SUMMARY:Garden Start Date - %s\n", gardenLayout.Name))
	icalendar.WriteString("END:VEVENT\n")

	// Care dates events
	for i, date := range gardenLayout.CareDates {
		icalendar.WriteString("BEGIN:VEVENT\n")
		icalendar.WriteString(fmt.Sprintf("UID:%d-care-%d@yourdomain.com\n", gardenLayout.ID, i))
		icalendar.WriteString(fmt.Sprintf("DTSTAMP:%s\n", formatICalDate(time.Now())))
		icalendar.WriteString(fmt.Sprintf("DTSTART:%s\n", formatICalDate(date)))
		icalendar.WriteString(fmt.Sprintf("SUMMARY:Care Date - %s\n", gardenLayout.Name))
		icalendar.WriteString("END:VEVENT\n")
	}

	// Replanting schedules events
	var schedules []model.Schedule
	database.DBConn.Where("garden_id = ?", gardenLayout.ID).Find(&schedules)
	for _, schedule := range schedules {
		for _, date := range schedule.PlantingDates {
			icalendar.WriteString("BEGIN:VEVENT\n")
			icalendar.WriteString(fmt.Sprintf("UID:%d-replant-%s@yourdomain.com\n", gardenLayout.ID, schedule.PlantName))
			icalendar.WriteString(fmt.Sprintf("DTSTAMP:%s\n", formatICalDate(time.Now())))
			icalendar.WriteString(fmt.Sprintf("DTSTART:%s\n", formatICalDate(date)))
			icalendar.WriteString(fmt.Sprintf("SUMMARY:Replant Date for %s - %s\n", schedule.PlantName, gardenLayout.Name))
			icalendar.WriteString("END:VEVENT\n")
		}
	}

	icalendar.WriteString("END:VCALENDAR\n")

	return icalendar.String(), nil
}

func formatICalDate(t time.Time) string {
	return t.Format("20060102T150405Z")
}