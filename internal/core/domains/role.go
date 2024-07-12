package domains

import "fmt"

type Role struct {
	ID          uint
	Name        string
	Description string
	Weight      int
}

type RoleW int

const (
	Owner RoleW = iota
	Maintainer
	Editor
	Viewer
)

var roleNames = []string{
	"Owner",
	"Maintainer",
	"Editor",
	"Viewer",
}

func (r RoleW) String() string {
	if r < Owner || r > Viewer {
		return "Unknown"
	}
	return roleNames[r]
}

func ParseRole(s string) (RoleW, error) {
	for i, name := range roleNames {
		if name == s {
			return RoleW(i), nil
		}
	}
	return -1, fmt.Errorf("invalid role: %s", s)
}
