package currentstatus

import "strings"

type CurrentStatus int

const (
	SUCCESS CurrentStatus = iota
	FAILED
	PENDING
)

func (currStatus CurrentStatus) ToString() string {
	switch currStatus {
	case SUCCESS:
		return "SUCCESS"
	case FAILED:
		return "FAILED"
	case PENDING:
		return "PENDING"
	default:
		return "PENDING"
	}
}

func CurrentStatusParser(currentStatus string) CurrentStatus {
	switch strings.ToUpper(currentStatus) {
	case "SUCCESS":
		return SUCCESS
	case "PENDING":
		return PENDING
	case "FAILED":
		return FAILED
	default:
		return PENDING
	}
}
