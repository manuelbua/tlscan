package runner

import "strings"

func GetVersion() string {
	if len(strings.TrimSpace(version)) == 0 {
		return "(unversioned)"
	}
	return version
}
