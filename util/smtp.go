package util

import (
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

type Message struct {
	To      []string
	From    string
	Subject string
	Body    string
	User    string
	Type    string
	Massive bool
	Info    string
}

// Content create mail content
func (m Message) Content() string {
	// set mail type
	contentType := "text/plain; charset=UTF-8"
	if m.Type == "html" {
		contentType = "text/html; charset=UTF-8"
	}

	// create mail content
	content := "From: " + m.User + "<" + m.From +
		">\r\nSubject: " + m.Subject + "\r\nContent-Type: " + contentType + "\r\n\r\n" + m.Body
	return content
}

type SmtpSendMail struct {
	MailHost     string
	MailAuthUser string
	MailAuthPass string
}

// SendMail Direct Send mail message
func (s *SmtpSendMail) SendMail(msg Message) (int, error) {
	host := strings.Split(s.MailHost, ":")

	// get message body
	content := msg.Content()

	auth := smtp.PlainAuth("", s.MailAuthUser, s.MailAuthPass, host[0])

	if len(msg.To) == 0 {
		return 0, fmt.Errorf("empty receive emails")
	}

	if len(msg.Body) == 0 {
		return 0, fmt.Errorf("empty email body")
	}

	if msg.Massive {
		// send mail to multiple emails one by one
		num := 0
		for _, to := range msg.To {
			body := []byte("To: " + to + "\r\n" + content)
			err := smtp.SendMail(s.MailHost, auth, msg.From, []string{to}, body)
			if err != nil {
				return num, err
			}
			num++
		}
		return num, nil
	}
	body := []byte("To: " + strings.Join(msg.To, ";") + "\r\n" + content)

	// send to multiple emails in one message
	err := smtp.SendMail(s.MailHost, auth, msg.From, msg.To, body)
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// AsyncSendMail Async Send mail message
func (s *SmtpSendMail) AsyncSendMail(msg Message) {
	// TODO may be need pools limit concurrent nums
	go func() {
		if num, err := s.SendMail(msg); err != nil {
			tos := strings.Join(msg.To, "; ")
			info := ""
			if len(msg.Info) > 0 {
				info = ", info: " + msg.Info
			}
			// log failed
			log.Printf("Async send email %d succeed, not send emails: %s%s err: %s", num, tos, info, err)
		}
	}()
}

// NewHTMLMessage Create html mail message
func NewHTMLMessage(To []string, From, Subject, Body string) Message {
	return Message{
		To:      To,
		From:    From,
		Subject: Subject,
		Body:    Body,
		Type:    "html",
	}
}
