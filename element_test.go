package shaman_test

import (
	"strings"
	"testing"

	. "github.com/gost-dom/shaman"

	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
	"github.com/stretchr/testify/assert"
)

func loadHTML(t *testing.T, h string) dom.Document {
	t.Helper()
	win, err := html.NewWindowReader(strings.NewReader(h))
	if err != nil {
		t.Fatalf("Error parsing HTML document")
	}

	return win.Document()
}

func TestElementName(t *testing.T) {
	t.Parallel()

	t.Run("<input>", func(t *testing.T) {
		doc := loadHTML(t, `
			<label for="input-1">Value 1</label><input id="input-1">
			<label for="input-2">Value 2</label><input id="input-2">
			<label for="input-3">Ignored value</label>
			<input id="input-3" aria-label="Value 3">
			<input id="input-4" aria-label="Ignored" aria-labelledby="label-4"><p id="label-4">Value 4</p>
			<input id="input-5" aria-labelledby="label-5a label-5b">
			<p id="label-5a">Value 5a</p>
			<p id="label-5b">Value 5b</p>
		`)
		assert.Equal(t, "Value 1", ElementName(doc.GetElementById("input-1")))
		assert.Equal(t, "Value 2", ElementName(doc.GetElementById("input-2")))
		assert.Equal(t,
			"Value 3", ElementName(doc.GetElementById("input-3")),
			"aria-label should win over a <label>",
		)
		assert.Equal(t,
			"Value 4", ElementName(doc.GetElementById("input-4")),
			"aria-labelledby should win over aria-label",
		)
		assert.Equal(t,
			"Value 5a Value 5b", ElementName(doc.GetElementById("input-5")),
			"aria-labelledby should accept multiple IDs",
		)
	})

	t.Run("<button>", func(t *testing.T) {
		doc := loadHTML(t, `<button id="btn">Click me!</button>`)
		assert.Equal(t, "Click me!", ElementName(doc.GetElementById("btn")))
	})

	t.Run("<a>", func(t *testing.T) {
		doc := loadHTML(t, `<a href="dummy" id="link">Click me!</button>`)
		assert.Equal(t, "Click me!", ElementName(doc.GetElementById("link")))
	})
}
