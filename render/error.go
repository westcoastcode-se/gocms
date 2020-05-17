package render

type TemplateNotFound struct {
	Path string
}

func (t *TemplateNotFound) Error() string {
	return "Could not find template: " + t.Path
}

