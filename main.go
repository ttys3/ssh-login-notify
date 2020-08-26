
// https://github.com/sendgrid/sendgrid-go
// https://linux.die.net/man/8/pam_exec

// /etc/pam.d/passwd
// password optional pam_exec.so seteuid /usr/bin/env SENDGRID_API_KEY="xxx" MAIL_FROM="" MAIL_TO="" /usr/local/bin/ssh-login-notify

// /etc/pam.d/sshd
// session optional pam_exec.so seteuid /usr/bin/env SENDGRID_API_KEY="xxx" MAIL_FROM="" MAIL_TO="" /usr/local/bin/ssh-login-notify
package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"

	"ssh-login-notify/pam"
)

var appName = "ssh-login-notify"
var version = "1.0.0"

var PAM *pam.PAMEnv
var Hostname string

func init() {
	PAM = pam.NewPAMEnv().Init()
	Hostname, _ = os.Hostname()
}

const mailTmpl = `
	------------------------------------
	User: {{ .PAM.PAM_USER }}
	Remote User: {{ .PAM.PAM_RHOST }}
	Remote Host: {{ .PAM.PAM_RHOST }}
	Service: {{ .PAM.PAM_SERVICE }}
	TTY: {{ .PAM.PAM_TTY }}
	Type: {{ .PAM.PAM_TYPE }}
	Date: {{ .Date }}
	Hostname: {{ .Hostname }}
	Reported By: {{ .AppName }} {{ .AppVer }} 
	-------------------------------------
`

type MailVars struct {
	PAM *pam.PAMEnv
	AppName string
	AppVer string
	Hostname string
	Date string
}

func NewV3MailInit(from *mail.Email, subject string, content ...*mail.Content) *mail.SGMailV3 {
	m := new(mail.SGMailV3)
	m.SetFrom(from)
	m.Subject = subject
	m.AddContent(content...)
	return m
}

func main() {

	subject := fmt.Sprintf("%s login on %s for account %s", PAM.PAM_SERVICE, Hostname, PAM.PAM_USER)

	t, err := template.New("ssh-notify").Parse(mailTmpl)
	if err != nil {
		panic(err)
	}
	tplData := MailVars{
		PAM: PAM,
		AppName: appName,
		AppVer: version,
		Hostname: Hostname,
		Date: time.Now().Format("2006-01-02 15:04:05"),
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, tplData)
	if err != nil {
		panic(err)
	}

	mailFrom :=  os.Getenv("MAIL_FROM")
	if mailFrom == "" {
		log.Fatalf("empty MAIL_FROM")
	}
	if !strings.Contains(mailFrom, "@") {
		log.Fatalf("invalid MAIL_FROM")
	}

	mailFromName :=  os.Getenv("MAIL_FROM_NAME")
	if mailFromName == "" {
		mailFromName = appName
	}

	from := mail.NewEmail(mailFromName, mailFrom)

	plainText := buf.String()
	log.Printf("mail content: %s", plainText)
	plainTextContent := mail.NewContent("text/plain", plainText)
	message := NewV3MailInit(from, subject, plainTextContent)


	mailTo :=  os.Getenv("MAIL_TO")
	if mailTo == "" {
		log.Fatalf("empty MAIL_TO")
	}
	if !strings.Contains(mailTo, "@") {
		log.Fatalf("invalid MAIL_TO")
	}
	p := mail.NewPersonalization()
	if strings.Contains(mailTo, ",") {
		tos := strings.Split(mailTo, ",")
		for _, to := range tos {
			to = strings.TrimSpace(to)
			if to != "" && strings.Contains(to, "@") {
				// add the email
				toEmail := mail.NewEmail("", to)
				p.AddTos(toEmail)
			}
		}
	} else {
		toEmail := mail.NewEmail("", mailTo)
		p.AddTos(toEmail)
	}

	sgApiKey := os.Getenv("SENDGRID_API_KEY")
	if sgApiKey == "" {
		log.Fatalf("empty SENDGRID_API_KEY")
	}
	message.AddPersonalizations(p)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)

	// https://sendgrid.api-docs.io/v3.0/how-to-use-the-sendgrid-v3-api/api-responses#status-codes
	if err != nil {
		log.Printf("request failed, err=%v", err)
	} else {
		if response.StatusCode == http.StatusAccepted {
			log.Println("success")
		} else {
			log.Printf("failed with StatusCode=%d", response.StatusCode)
		}
	}
}
