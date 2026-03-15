package heartbeat

import (
	"github.com/jotitan/monitor-pis/config"
	"sync"
)

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
	mapLocker  sync.Mutex
	alerts     map[string]*Alert
	mailSender MailSender
}

func (m *Manager) RunHeartbeat(hb Heartbeat) bool {
	success := hb.test()
	hb.manager.recordPoint(success, hb.url)

	return success
}

func (m *Manager) resetPoint(url string) {
	m.mapLocker.Lock()
	m.points[url] = 0
	m.mapLocker.Unlock()
}

func (m *Manager) increasePoint(url string) int {
	m.mapLocker.Lock()
	m.points[url] = m.points[url] + 1
	m.mapLocker.Unlock()
	return m.points[url]
}

func (m *Manager) recordPoint(success bool, url string) {
	a, exists := m.alerts[url]
	if !exists {
		return
	}
	if success {
		if a.alertSent && a.alertWhenUp {
			// Send email cause service is back
			m.mailSender.sendBackNormal(a.name, url, a.email)
		}
		m.resetPoint(url)
		a.alertSent = false
	} else {
		nb := m.increasePoint(url)
		if a.alertSent == false && nb >= a.threashold {
			m.mailSender.sendFail(a.name, url, a.email)
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
		mapLocker:  sync.Mutex{},
		alerts:     make(map[string]*Alert),
		mailSender: NewMailSender(emailConfig),
	}
}
