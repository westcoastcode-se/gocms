package render

import "path/filepath"

type FactoryNotFound struct {
	view string
}

func (f *FactoryNotFound) Error() string {
	return "could not find factory for view: " + f.view
}

type TemplateRenderers struct {
	Renderers map[string]TemplateRendererFactory
}

// Add a new factory used when rendering templates
func (t *TemplateRenderers) AddFactory(suffix string, factory TemplateRendererFactory) {
	t.Renderers[suffix] = factory
}

// Figure out which factory to use when rendering the supplied view. This will return a FactoryNotFound error
// if no matching factory is found.
func (t *TemplateRenderers) FindFactory(view string) (TemplateRendererFactory, error) {
	suffix := filepath.Ext(view)
	if factory, ok := t.Renderers[suffix]; ok {
		return factory, nil
	}
	return nil, &FactoryNotFound{view}
}

func NewTemplateRenderers() *TemplateRenderers {
	return &TemplateRenderers{Renderers: make(map[string]TemplateRendererFactory)}
}
