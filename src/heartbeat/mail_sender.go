package heartbeat

import (
	"fmt"
	"github.com/jotitan/monitor-pis/config"
	"log"
	"net/smtp"
)

type MailSender struct {
	config.EmailConfig
}

func NewMailSender(conf config.EmailConfig) MailSender {
	return MailSender{conf}
}

func (ms MailSender) sendBackNormal(name, url string) {
	subject := fmt.Sprintf("✅ Heartbeat service %s back to normal", name)
	message := fmt.Sprintf("Hello\n\nService %s is back to normal. \n\nUrl %s is fully available", name, url)
	ms.send(subject, message)
}

func (ms MailSender) sendFail(name, url string) {
	log.Println("Fail heartbeat", name, " => send mail")
	subject := fmt.Sprintf("⚠️ Heartbeat service %s fail", name)
	message := fmt.Sprintf("Hello\n\nService %s is not available. \n\nUrl %s is not responding, try to fix it.", name, url)
	ms.send(subject, message)
}

func (ms MailSender) send(subject, message string) {
	auth := smtp.PlainAuth("", ms.LoginSMTP, ms.PasswordSMTP, ms.HostSMTP)
	smtp.SendMail(ms.HostSMTP+":"+ms.PortSMTP, auth, ms.EmailSender, []string{ms.EmailRecipient}, []byte("Subject:"+subject+"\n\n"+message))
}
