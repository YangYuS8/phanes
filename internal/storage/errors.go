package storage

import "fmt"

type MissingTableError struct {
	Table string
}

func (e *MissingTableError) Error() string {
	return fmt.Sprintf("missing sqlite table %q", e.Table)
}
