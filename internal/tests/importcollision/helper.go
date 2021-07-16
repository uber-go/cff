package importcollision

import (
	h "html/template"
	t "text/template"

	p "go.uber.org/cff/internal/tests/importcollision/package-with-dash"
)

// GetHTMLTemplate returns a nil html template.
func GetHTMLTemplate() *h.Template {
	return nil
}

// GetTextTemplate returns a nil text template.
func GetTextTemplate() *t.Template {
	return nil
}

// GetResult returns an empty string.
func GetResult(h1 *h.Template, h2 *t.Template, mt p.Foo) string {
	return ""
}

// GetFoo returns a Foo.
func GetFoo() p.Foo {
	return p.Foo{}
}
