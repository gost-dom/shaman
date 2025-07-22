package shaman

import (
	"fmt"

	"github.com/gost-dom/shaman/ariarole"

	"github.com/gost-dom/browser/dom"
)

// ByName is an [ElementPredicate] that matches elements by their accessibility
// name.
//
// See also: [ElementName]
type ByName string

func (n ByName) IsMatch(e dom.Element) bool { return GetName(e) == string(n) }

func (n ByName) String() string { return fmt.Sprintf("By accessibility name: %s", string(n)) }

// An [ElementPredicate] that matches elements by their [ARIA role].
//
// [ARIA role]: https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Roles
type ByRole ariarole.Role

func (r ByRole) IsMatch(
	e dom.Element,
) bool {
	return ariarole.GetElementRole(e) == ariarole.Role(r)
}

func (r ByRole) String() string { return fmt.Sprintf("By role: %s", string(r)) }

type byH1Predicate struct{}

// ByH1 is a predicate to find the <h1> element on the page. This predicate is
// intended to be used in a Get query that will fail if multiple elements match
// the predicate, to ensure the page has exactly one <h1> element.
//
// The <h1> element is for the "Page title", and a page should have only _one_
// title.
//
// See also: https://developer.mozilla.org/en-US/docs/Web/HTML/Reference/Elements/Heading_Elements#avoid_using_multiple_h1_elements_on_one_page
var ByH1 = byH1Predicate{}

func (s byH1Predicate) IsMatch(e dom.Element) bool { return e.TagName() == "H1" }

func (s byH1Predicate) String() string { return "Main heading (<h1>)" }
