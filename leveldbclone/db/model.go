package db

type NotFoundError struct{}

func (m *NotFoundError) Error() string {
	return "value not found"
}