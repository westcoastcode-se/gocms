package render

import (
	"net/http"
)

// Factory responsible for creating a renderer
type TemplateRendererFactory interface {
	// Create a new renderer
	NewRenderer(r *http.Request) TemplateRenderer
}
