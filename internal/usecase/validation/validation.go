package validation

import "taskapi/internal/domain"

func IsValidStatus(s domain.Status) bool {
	switch s {
	case domain.StatusTodo, domain.StatusInProgress, domain.StatusDone:
		return true
	default:
		return false
	}
}

func StatusString(s *domain.Status) string {
	if s == nil {
		return ""
	}
	return string(*s)
}

func ErrString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
