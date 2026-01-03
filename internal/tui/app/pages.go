package app

import "github.com/arthur404dev/dotts/pkg/vetru/messages"

const (
	PageDashboard messages.PageID = "dashboard"
	PageStatus    messages.PageID = "status"
	PageUpdate    messages.PageID = "update"
	PageDoctor    messages.PageID = "doctor"
	PageSettings  messages.PageID = "settings"
	PageWizard    messages.PageID = "wizard"
)

const (
	ActionSync   messages.ActionID = "sync"
	ActionUpdate messages.ActionID = "update"
	ActionDoctor messages.ActionID = "doctor"
)
