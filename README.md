# Shaman

## Introduction

Shaman is a support library to work on top of [Gost-DOM]. It helps writing tests
using a higher level of abstraction than the native DOM. E.g., to interact with
an input field, you may write.

[Gost-DOM]: https://github.com/gost-dom/browser

```Go
package my_test

import (
    "github.com/gost-dom/shaman"
    . "github.com/gost-dom/shaman/predicates"
)

func TestSomething(t *testing.T) {
    win := initWindow(t) // Return a gost-dom/browser/html.Window
    scope := shaman.WindowScope(win)
    form := scope.SubScope(ByRole(ariarole.Form))
    // Find a textbox with the accessibility name, "Email"
    form.Get(ByRole(ariarole.Textbox), ByName("Email")).Write("jd@example.com")
    // Find a password input field with the accessibility name, "Password"
    form.Get(ByRole(ariarole.PasswordText), ByName("Password")).Write("jd@example.com")
```

The "Accessibility name" is in fact, the label for an element. The code style
not only encourages building accessibility into your application, it allows the
visual design to change without breaking the test.

For example, the following two HTML snippets will both work with the previous
test.

```html
<label for="email-input">Email</label>
<input id="email-input" type="text" />
<label for="password-input">Password</label>
<input id="password-input" type="password" />
```

```html
<input aria-label="Email" type="text" placeholder="email" />
<input aria-label="Password" type="password" placeholder="******" />
```

This style of test supports changes to the UI, but helps promote good
accessibility, i.e., helping you to remember to add the `aria-label` for proper
screen reader support.

> [!NOTE]
> The ARIA role, `PasswordText` is not a real ARIA role. It is reported by
> Firefox's developer tools, but is helpful in test scenarios to locate password
> fields, as password input doesn't _have_ an ARIA role. 

> [!NOTE]
> Shaman is currently coupled to the interfaces exposed by Gost-DOM; but could
> be coupled to general interface, with Gost-DOM being just _one
> implementation; and other implementations could support Webdriver or other
> browser automation protocols. However, Shaman relies on very chatty
> communication, which would cause significant overhead using any kind of
> inter-process communication.
>
> Shaman is written to depend methods defined in the DOM and HTML DOM standards;
> adapted to Go idioms (errors as values, and Go naming conventions)
