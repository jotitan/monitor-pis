package heartbeat

import "github.com/jotitan/monitor-pis/config"

type Alert struct {
	email       string
	name        string
	threashold  int
	alertWhenUp bool
	alertSent   bool
}

func newAlert(alert config.AlertConfig, name string) *Alert {
	return &Alert{alertSent: false, alertWhenUp: alert.AlertWhenUp, threashold: alert.Threashold, name: name, email: alert.Email}
}

// Manager Used to send alerts when heartbeat fail
type Manager struct {
	// For each url, store number of fail
	points     map[string]int
	alerts     map[string]*Alert
	mailSender MailSender
}

func (m *Manager) RunHeartbeat(hb Heartbeat) bool {
	success := hb.test()
	hb.manager.recordPoint(success, hb.url)

	return success
}

func (m *Manager) recordPoint(success bool, url string) {
	a, exists := m.alerts[url]
	if !exists {
		return
	}
	if success {
		if a.alertSent && a.alertWhenUp {
			// Send email cause service is back
			m.mailSender.sendBackNormal(a.name, url)
		}
		m.points[url] = 0
		a.alertSent = false
	} else {
		m.points[url] = m.points[url] + 1
		if a.alertSent == false && m.points[url] >= a.threashold {
			m.mailSender.sendFail(a.name, url)
			a.alertSent = true
		}
	}

}

func (m *Manager) RegisterAlert(hb *Heartbeat, alert config.AlertConfig) {
	hb.manager = m
	if alert.Threashold > 0 {
		m.points[hb.url] = 0
		m.alerts[hb.url] = newAlert(alert, hb.Name)
	}
}

func NewManager(emailConfig config.EmailConfig) *Manager {
	return &Manager{
		points:     make(map[string]int),
		alerts:     make(map[string]*Alert),
		mailSender: NewMailSender(emailConfig),
	}
}
