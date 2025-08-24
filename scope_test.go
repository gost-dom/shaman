package shaman_test

import (
	"slices"
	"testing"

	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/gost-dom/shaman"
	"github.com/gost-dom/shaman/ariarole"
	"github.com/stretchr/testify/assert"
)

func TestScope_Find(t *testing.T) {
	t.Parallel()
	root := createRoot("div",
		child("a", textContent("Link 1")),
		child("div", textContent("Not link")),
		child("a", textContent("Link 2")),
	)

	t.Run("FindAll(ByRole(ariarole.Link))", func(t *testing.T) {
		t.Parallel()
		scope := shaman.NewScope(t, root)
		links := slices.Collect(scope.FindAll(shaman.ByRole(ariarole.Link)))
		if assert.Len(t, links, 2, "Length of found elements") {
			assert.Equal(t, "Link 1", links[0].TextContent())
			assert.Equal(t, "Link 2", links[1].TextContent())
		}
	})

	t.Run("All()", func(t *testing.T) {
		t.Parallel()
		scope := shaman.NewScope(t, root)
		all := slices.Collect(scope.All())
		if assert.Len(t, all, 4, "Length of found elements") {
			assert.Equal(t, "Link 1Not linkLink 2", all[0].TextContent())
			assert.Equal(t, "Link 1", all[1].TextContent())
			assert.Equal(t, "Not link", all[2].TextContent())
			assert.Equal(t, "Link 2", all[3].TextContent())
		}
	})
}

type container struct{ dom.ElementContainer }

type containerFunc func(dom.Element)

func createRoot(tagName string, opt ...containerFunc) dom.Element {
	w := html.NewWindow()
	d := html.NewHTMLDocument(w)
	e := d.CreateElement(tagName)
	for _, o := range opt {
		o(e)
	}
	return e
}

func child(tagName string, opt ...containerFunc) containerFunc {
	return func(e dom.Element) {
		d := e.OwnerDocument()
		child := d.CreateElement(tagName)
		for _, o := range opt {
			o(child)
		}
		e.AppendChild(child)
	}
}

func textContent(text string) containerFunc {
	return func(e dom.Element) {
		e.SetTextContent(text)
	}
}
