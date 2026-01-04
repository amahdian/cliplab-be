package common

type Action string

const (
	ActionEndorse Action = "ENDORSE"
	ActionReject  Action = "REJECT"
	ActionApprove Action = "APPROVE"
)

func Actions() []Action {
	return []Action{ActionEndorse, ActionReject, ActionApprove}
}

func ApprovalActions() []string {
	return []string{string(ActionReject), string(ActionApprove)}
}

func IsAction(action string) bool {
	actions := Actions()
	for _, a := range actions {
		if string(a) == action {
			return true
		}
	}
	return false
}
