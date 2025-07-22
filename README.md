# Shaman

## Introduction

Shaman is a support library to work on top of [Gost-DOM]. It helps writing tests
using a higher level of abstraction than the native DOM. This encourages
building accessibility in and results in code that is not only easier to read,
but also resilient to changes in UI that doesn't change semantics:

> [!Warning]
> This library is pre 0.1 state - breaking changes may be pushed with no
> warning (where feasible, old versions will live with a deprecation warning for
> a while)

[Gost-DOM]: https://github.com/gost-dom/browser

## Example

As an example, to write test to interact with an input field, you may write:

```go
package my_test

import (
    "github.com/gost-dom/shaman"
    . "github.com/gost-dom/shaman/predicates"
)

func TestSomething(t *testing.T) {
    win := initWindow(t) // Return a gost-dom/browser/html.Window
    // Find the <form> in the main landmark of the page.
    scope := shaman.WindowScope(t, win)
    mainContent := scope.SubScope(ariarole.Main)
    form := mainContent.SubScope(ByRole(ariarole.Form))
    // Find a textbox with the accessibility name, "Email"
    form.Textbox(ByName("Email")).Write("jd@example.com")
    // Find a password input field with the accessibility name, "Password"
    form.PasswordText(ByName("Password")).Write("very_secret")
    form.Get(ByName(ByRole(ariarole.Button), ByName("Submit"))).Click()
    // ...
}
```

## Easier to write

A common pattern is to assign `id` or `data-testid` attributes to elements just
to be able to find them in test cases. This practice does adds to mental load of
the developer:

- Which `id`-attributes to assign to which elements?
- When writing tests first, you have to deal with this _before_ even addressing
  the problem.

The _shaman style_ of thinking in terms of textboxes with labels, placed inside
a form, forces the developer to work in same level of abstraction as the
_problem domain_ itself: A user interacting with a page, identifying form
elements by the labels they have.

## Accessibility

The shaman style encourages writing tests that implcitly verify that elements
have the proper attributes to support accessibility, i.e., all input fields have
labels.

See the patterns section below for guidelines how to write tests that enforces a
higher level of accessibility.

## Resilient to changes in layout

If the layout of the application changes, but the _functionality_ remains, the
test is resilient to this change, if the new layout has the same semantics.

For example, the following 4 examples all have the same semantics, and provides
the same _functionality_ to the user. 


```html
<!-- A <label> is associated with an input field with the for-attribute -->
<label for="email-input">Email: </label>
<input id="email-input" type="text" />

<!-- An <input> field is a child of the label -->
<label>Email: <input id="email-input" type="text" /></label>

<!-- The input field references a random element using aria-labelledby -->
<input type="text" aria-labelledby="email-label" />
<span id="email-label">Your email address</span>

<!-- The label doesn't appear on the page, the input has an aria-label -->
<input type="text" aria-label="Email" placeholder="email" />
```

The shaman style of testing is not only resilient to the user interface
changing; it detects if you forget to add a label.

> [!NOTE]
> In my experience, 80% of all developers and UI designers are ignorant of
> accessibility. As a consequence any new project member are statistically very
> likely to break accessibility if not verified at design time. Fast
> developer-friendly tests is the best way to detect this early, preventing an
> unproductive path.
>
> (80% was a pretty conservative number. It's probably more like 95%)

## Good patterns

This is a collection of patterns that should help write more resilient tests
that also help achieve better accessibility.

### Verify exactly one `<h1>` heading

An `<h1>` is treated as a page title, so there should be exactly _one_. Using
the `Get` method on document scope fails when multiple elements match the
predicate, effectively verifying that exactly one `<h1>` exists with the
expected title:

```go
titleElm := shaman.WindowScope(t, win).Get(shaman.ByH1)
if got := titleElm.TextContent(); got != "Expected page title" {
    t.Error("Wrong title")
}
```

### Scope by landmarks

Scope by landmark relevant landmark. Most tests would generally verify behaviour
of the _main_ content, typically in a `<main>` element; so to enforce the
document structure.

```go
mainContent := windowScope.SubScope(ByRole(ariarole.Main))
```

Screen reader users typically rely landmarks to quickly find the relevant
content. By writing tests to enforce a main scope, it will have a dramatic
effect on the usability of the web application for that user base.

### Having input fields? Always scope them by form

```Go
loginForm := mainContent.SubScope(ByRole(ariarole.Form))
emailField := loginForm.Get(ByRole(ariarole.Textbox), ByName("Email"))
```

> [!NOTE]
> You _should_ add a `ByName` when finding a form, to ensure that it has a
> proper title. That isn't supported by shaman at the time of writing this.
> https://github.com/gost-dom/shaman/issues/2

## Can I use this with other libraries? (e.g., selenium, playwright)

Shaman is currently coupled to the interfaces exposed by Gost-DOM, but is 
written to depend on methods defined in the DOM and HTML DOM standards
adapted to Go idioms (errors as values, and Go naming conventions).

This could be coupled to general interface, with Gost-DOM being just _one_
implementation; other implementations could support Webdriver or other browser
automation protocols. 

_However_, Shaman relies on very chatty communication when processing the DOM
tree, which would cause significant overhead using any kind of inter-process
communication.
