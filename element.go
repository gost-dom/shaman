package shaman

import (
	"fmt"
	"strings"

	"github.com/gost-dom/browser/dom"
)

// ElementName returns the [accessibility name] of HTML element e. It returns
// empty string if e is nil.
//
// The accessibility name is a descriptive name an element, examples are:
//   - The name for <input id="name"> is the text content an associated <label>
//     element
//   - The name of a <button> is its text content.
//
// The default name can be overridden with either aria-label or aria-labelledby
// attributes.
//
// See also: https://developer.mozilla.org/en-US/docs/Web/Accessibility/Guides/Understanding_WCAG/Text_labels_and_names
//
// [accessibility name]: https://w3c.github.io/accname/#dfn-accessible-name
func ElementName(e dom.Element) string {
	// TODO: This should be exposed as IDL attributes
	if e == nil {
		return ""
	}
	doc := e.OwnerDocument()
	if l, ok := e.GetAttribute("aria-labelledby"); ok {
		ids := strings.Split(l, " ")
		labels := make([]string, 0, len(ids))
		for _, id := range ids {
			if labelElm := doc.GetElementById(id); labelElm != nil {
				labels = append(labels, labelElm.TextContent())
			}
		}
		return strings.Join(labels, " ")
	}
	if l, ok := e.GetAttribute("aria-label"); ok {
		return l
	}
	switch e.TagName() {
	case "INPUT":
		if id, ok := e.GetAttribute("id"); ok {
			if label, _ := doc.QuerySelector(fmt.Sprintf("label[for='%s']", id)); label != nil {
				return label.TextContent()
			}
		}
	case "A", "BUTTON", "LI": // How many more? Can it be calculated from webref?
		return e.TextContent()
	}
	return ""
}

// ElementDescription returns the [accessibility description] of an element. The
// description provides additional context to complement the name. E.g.,
// validation errors associated with an element would; or additional guidance
// about valid input.
//
// [accessibility description]: https://w3c.github.io/accname/#dfn-accessible-description
func ElementDescription(e dom.Element) string {
	if id, ok := e.GetAttribute("aria-describedby"); ok {
		ids := strings.Split(id, " ")
		ss := make([]string, 0, len(ids))
		for _, id := range ids {
			if e := e.OwnerDocument().GetElementById(id); e != nil {
				ss = append(ss, e.TextContent())
			}
		}
		return strings.Join(ss, "\n")
	}
	s, _ := e.GetAttribute("aria-description")
	return s
}

// Older named versions. The new "Element" prefix seems better, as they
// calculate some property of an element.

// Deprecated: Use [ElementName] instead. This function will be removed.
var GetName = ElementName

// Deprecated: Use [ElementDescription] instead. This function will be removed.
var GetDescription = ElementDescription
