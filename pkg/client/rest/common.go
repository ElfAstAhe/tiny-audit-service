package rest

const (
	AuthEventLogin   string = "login"
	AuthEventLogoff  string = "logoff"
	AuthEventRefresh string = "refresh"
	AuthEventRevoke  string = "revoke"
)

const (
	DataEventCreate string = "create"
	DataEventChange string = "change"
	DataEventRemove string = "remove"
)

const (
	AuditStatusSuccess string = "success"
	AuditStatusFail    string = "fail"
)
