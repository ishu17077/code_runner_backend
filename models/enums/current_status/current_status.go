package currentstatus

import "strings"

type CurrentStatus int

const (
	SUCCESS CurrentStatus = iota
	FAILED
	PENDING
	TIME_LIMIT_EXCEEDED
	RESOURCE_LIMIT_EXCEEDED
)

func (currStatus CurrentStatus) ToString() string {
	switch currStatus {
	case SUCCESS:
		return "SUCCESS"
	case FAILED:
		return "FAILED"
	case PENDING:
		return "PENDING"
	case TIME_LIMIT_EXCEEDED:
		return "TIME_LIMIT_EXCEEDED"
	case RESOURCE_LIMIT_EXCEEDED:
		return "RESOURCE_LIMIT_EXCEEDED"
	default:
		return "FAILED"
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
	case "TIME_LIMIT_EXCEEDED":
		return TIME_LIMIT_EXCEEDED
	case "RESOURCE_LIMIT_EXCEEDED":
		return RESOURCE_LIMIT_EXCEEDED
	default:
		return FAILED
	}
}
