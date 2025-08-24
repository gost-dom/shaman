package shaman_test

// Contains helpers for constructing a DOM tree to test on.
//
// The current version may be unnecessarily complex for the problem at hand.

import (
	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
)

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
