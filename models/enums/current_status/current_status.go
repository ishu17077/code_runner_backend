package currentstatus

import "strings"

type CurrentStatus int

const (
	SUCCESS CurrentStatus = iota
	FAILED
	TIME_LIMIT_EXCEEDED
	RESOURCE_LIMIT_EXCEEDED
	RUNTIME_ERROR
	INTERNAL_ERROR
)

func (currStatus CurrentStatus) ToString() string {
	switch currStatus {
	case SUCCESS:
		return "SUCCESS"
	case FAILED:
		return "FAILED"
	case TIME_LIMIT_EXCEEDED:
		return "TIME_LIMIT_EXCEEDED"
	case RESOURCE_LIMIT_EXCEEDED:
		return "RESOURCE_LIMIT_EXCEEDED"
	case INTERNAL_ERROR:
		return "INTERNAL_ERROR"
	case RUNTIME_ERROR:
		return "RUNTIME_ERROR"
	default:
		return "FAILED"
	}
}

func CurrentStatusParser(currentStatus string) CurrentStatus {
	switch strings.ToUpper(currentStatus) {
	case "SUCCESS":
		return SUCCESS
	case "FAILED":
		return FAILED
	case "TIME_LIMIT_EXCEEDED":
		return TIME_LIMIT_EXCEEDED
	case "RESOURCE_LIMIT_EXCEEDED":
		return RESOURCE_LIMIT_EXCEEDED
	case "INTERNAL_ERROR":
		return INTERNAL_ERROR
	case "RUNTIME_ERROR":
		return RUNTIME_ERROR
	default:
		return FAILED
	}
}
