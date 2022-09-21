package internal

import "errors"

type PersonStatus int

const (
	ACTIVE PersonStatus = iota
	SUSPEND
)

func (i PersonStatus) String() string {
	return []string{"active", "suspended"}[i]
}
func ParsePersonStatus(v string) (interface{}, error) {
	result := ACTIVE
	switch v {
	case "active":
		result = ACTIVE
	case "suspended":
		result = SUSPEND
	default:
		return 0, errors.New("Unknown PersonStatus value: " + v)
	}
	return &result, nil
}
func SerializeTeamVisibilityType(values []PersonStatus) []string {
	result := make([]string, len(values))
	for i, v := range values {
		result[i] = v.String()
	}
	return result
}
