package shaman

import (
	"fmt"
	"iter"
	"strings"
	"testing"

	"github.com/gost-dom/shaman/ariarole"

	"github.com/gost-dom/browser/dom"
	"github.com/gost-dom/browser/html"
)

// An ElementPredicate is a type that checks if an element matches certain
// criteria, and is used to fine elements in the dom. E.g., finding the input
// element with the label "email".
//
// An implementation of ElementPredicate should also implement [fmt.Stringer],
// describing what the predicate is looking for. This provides better error
// messages for failed queries.
type ElementPredicate interface{ IsMatch(dom.Element) bool }

// An ElementPredicateFunc wraps a single function as a predicate to be used
// with [Scope.FindAll] or [Scope.Get].
//
// This type is intended for quick prototyping of test code.
//
// It is strongly suggested to create a new type for predicates that also implements
// [fmt.Stringer].
//
// See also [ElementPredicate]
type ElementPredicateFunc func(dom.Element) bool

func (f ElementPredicateFunc) IsMatch(e dom.Element) bool { return f(e) }

// predicates treats multiple predicates as one, simplifying the search for multiple
// predicates, as well as stringifying multiple predicates.
type predicates []ElementPredicate

func (o predicates) IsMatch(e dom.Element) bool {
	for _, o := range o {
		if !o.IsMatch(e) {
			return false
		}
	}
	return true
}

func (o predicates) String() string {
	names := make([]string, len(o))
	for i, o := range o {
		if s, ok := o.(fmt.Stringer); ok {
			names[i] = s.String()
		} else {
			names[i] = "Unknown predicate. No String()"
		}
	}
	return strings.Join(names, ", ")
}

// Scope represents a subset of a page, and can be used to find elements withing
// that scope.
type Scope struct {
	t    testing.TB
	root containerer
}

func (s Scope) container() dom.ElementContainer { return s.root.container() }

// containerer is the interface for the single method container that returns a
// [dom.ElementContainer]. When the scope relates to a window, this will return
// the current document.
type containerer interface {
	container() dom.ElementContainer
}

type windowContainerer struct{ win html.Window }

func (c windowContainerer) container() dom.ElementContainer { return c.win.Document() }

type simpleContainer struct{ c dom.ElementContainer }

func (s simpleContainer) container() dom.ElementContainer { return s.c }

// WindowScope create a new [Scope] that is bound to an [html.Window]. This
// scope will always reflect the current page displayed in the window.
func WindowScope(t testing.TB, win html.Window) Scope {
	return Scope{t, windowContainerer{win}}
}

// NewScope create a new [Scope] that is bound to a single [dom.Element].
// Holding on to a Scope returned from NewScope may return elements no longer
// if the original element is removed from the DOM.
func NewScope(t testing.TB, c dom.ElementContainer) Scope {
	return Scope{t: t, root: simpleContainer{c}}
}

// All returns an iterator over all elements in scope. If the scope is an
// element, the element itself will be included.
func (h Scope) All() iter.Seq[dom.Element] {
	return func(yield func(dom.Element) bool) {
		container := h.container()
		if self, ok := container.(dom.Element); ok {
			if !yield(self) {
				return
			}
		}
		for _, child := range container.Children().All() {
			func() {
				next, stop := iter.Pull(NewScope(h.t, child).All())
				defer stop()
				for {
					v, ok := next()
					if !ok || !yield(v) {
						return

					}
				}
			}()
		}
	}
}

// FindAll returns a sequence of all elements that match the specified options.
func (h Scope) FindAll(options ...ElementPredicate) iter.Seq[dom.Element] {
	opt := predicates(options)
	return func(yield func(dom.Element) bool) {
		next, done := iter.Pull(h.All())
		defer done()
		for {
			e, ok := next()
			if !ok {
				return
			}
			if opt.IsMatch(e) {
				if !yield(e) {
					return
				}
			}
		}
	}
}

// Find returns an element that matches the options if any. At most one element
// is expected to exist in the dom mathing the options. Returns nil if no
// element is found. If more than one is found, Fatalf is called on the
// specified [testing.TB] instance.
//
// Note: This must run in the same goroutine as the test case.
func (h Scope) Find(opts ...ElementPredicate) html.HTMLElement {
	h.t.Helper()
	next, stop := iter.Pull(h.FindAll(opts...))
	defer stop()
	if v, ok := next(); ok {
		if v2, ok := next(); ok {
			h.t.Fatalf(
				"At least two elements match options: %s\n1st match: %s\n2nd match: %s",
				predicates(opts),
				v.OuterHTML(),
				v2.OuterHTML(),
			)
		}
		return v.(html.HTMLElement)
	}
	return nil
}

// Get returns the element that matches the options. Exactly one element is
// expected to exist in the dom mathing the options. If zero, or more than one
// are found, a fatal error is generated.
func (h Scope) Get(opts ...ElementPredicate) html.HTMLElement {
	container := h.container()
	h.t.Helper()
	if !container.IsConnected() {
		h.t.Logf("WARN (shaman): Scope root element not connected to document")
	}
	if res := h.Find(opts...); res != nil {
		return res
	}
	h.t.Fatalf("No elements mathing options: %s", predicates(opts))
	return nil
}

// Query looks for one element that matches the options, and return it in return
// value e. Return ok tells whether an element was found. At most one element is
// expected to exist in the dom mathing the options. Of more than one are found,
// a fatal error is generated.
func (h Scope) Query(opts ...ElementPredicate) (e html.HTMLElement, ok bool) {
	h.t.Helper()
	res := h.Find(opts...)
	return res, res != nil
}

func (h Scope) Subscope(opts ...ElementPredicate) Scope {
	return NewScope(h.t, h.Get(opts...))
}

func (s Scope) Textbox(opts ...ElementPredicate) TextboxRole {
	opts = append(opts, ByRole(ariarole.Textbox))
	return TextboxRole{s.Get(opts...)}
}

func (s Scope) Checkbox(opts ...ElementPredicate) CheckboxRole {
	opts = append(opts, ByRole(ariarole.Checkbox))
	return CheckboxRole{s.Get(opts...)}
}

func (s Scope) PasswordText(opts ...ElementPredicate) TextboxRole {
	opts = append(opts, ByRole(ariarole.PasswordText))
	return TextboxRole{s.Get(opts...)}
}

// A helper to interact with "text boxes"
type TextboxRole struct {
	html.HTMLElement
}

func (tb TextboxRole) Value() string {
	v, _ := tb.HTMLElement.GetAttribute("value")
	return v
}

// Write is intended to simulate the user typing in. Currently it merely sets
// the value content attribute, making it only applicable to input elements, not
// custom implementations of the textbox aria role.
func (tb TextboxRole) Write(input string) { tb.SetAttribute("value", input) }

func (tb TextboxRole) Clear() { tb.SetAttribute("value", "") }

func (tb TextboxRole) ARIADescription() string {
	return GetDescription(tb)
}

type CheckboxRole struct {
	html.HTMLElement
}

func (cb CheckboxRole) Check()   { cb.setChecked(true) }
func (cb CheckboxRole) Uncheck() { cb.setChecked(false) }

func (cb CheckboxRole) setChecked(val bool) {
	input, ok := cb.HTMLElement.(html.HTMLInputElement)
	if !ok {
		// To support generic checkbox roles, the approach should probably be to
		// check for the presence of the aria-checked content attribute, and
		// call Click() on the element if it has the wrong state.
		panic("CheckboxRole.Check/Uncheck: only input elements are supported")
	}
	input.SetChecked(val)
}
