package render

// Error that's raised if a requested template is not found
type TemplateNotFound struct {
	Path string
}

func (t *TemplateNotFound) Error() string {
	return "Could not find template: " + t.Path
}
