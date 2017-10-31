package cachet

import (
    "strings"
    "github.com/mattn/go-xmpp"
    "github.com/Sirupsen/logrus"
)

// Investigating template
var defaultXMPPInvestigatingTpl = MessageTemplate{
	Subject: `{{ .Monitor.Name }} - {{ .SystemName }}`,
	Message: `{{ .Monitor.Name }} check **failed** (server time: {{ .now }})

{{ .FailReason }}`,
}

// Fixed template
var defaultXMPPFixedTpl = MessageTemplate{
	Subject: `{{ .Monitor.Name }} - {{ .SystemName }}`,
	Message: `**Resolved** - {{ .now }}

- - -

{{ .incident.Message }}`,
}

type XMPPMonitor struct {
	AbstractMonitor `mapstructure:",squash"`

	Username       string
	Password       string
}

// TODO: test
// TODO: xmpp library doesn't support connection timeout
func (monitor *XMPPMonitor) test() bool {
	var talk *xmpp.Client
	var err error
	options := xmpp.Options{Host: monitor.Target,
		User:          monitor.Username,
		Password:      monitor.Password,
		NoTLS:         true,
		StartTLS:      true,
	}

	talk, err = options.NewClient()
	if err != nil && strings.Index(err.Error(), "PLAIN ") < 0 {
		monitor.lastFailReason = err.Error()
		logrus.Errorf("XMPP error: %s", monitor.lastFailReason)
		return false
	}

	if talk != nil {
		talk.Close()
	}

	return true
}

// TODO: test
func (mon *XMPPMonitor) Validate() []string {
	mon.Template.Investigating.SetDefault(defaultXMPPInvestigatingTpl)
	mon.Template.Fixed.SetDefault(defaultXMPPFixedTpl)

	errs := mon.AbstractMonitor.Validate()

	if len(mon.Username) == 0 || len(mon.Password) == 0 {
		errs = append(errs, "Both 'username' and 'password' must be provided")
	}

	return errs
}

func (mon *XMPPMonitor) Describe() []string {
	features := mon.AbstractMonitor.Describe()
	features = append(features, "Username: "+mon.Username)

	return features
}
