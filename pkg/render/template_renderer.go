package render

import "io"

// Renderer responsible for rendering a view
type TemplateRenderer interface {
	// Render the supplied view using a model and fill the io.Writer with the resulting content.
	// Will return an error if one occurs.
	RenderView(writer io.Writer, view string, model interface{}) error
}
