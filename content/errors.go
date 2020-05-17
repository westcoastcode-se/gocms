package content

type NotFoundError struct {
	message string
}

func (p *NotFoundError) Error() string {
	return p.message
}
