package components

import (
	"strconv"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

// Navbar returns a navigation bar with links to the home page,
// GitHub page and About page.
func Navbar(userID int, username string) Node {
	return Nav(Class("navbar"),
		Ol(
			NavbarItem("Home", "/"),
			NavbarItem("GitHub", "https://github.com/JneiraS/gTM"),
			NavbarItem("About", "/about"),
			ButtonClic(username, "/user/"+strconv.Itoa(userID)),
		),
	)
}

func NavbarItem(name, path string) Node {
	return Li(A(Href(path), Text(name)))
}

func ButtonClic(label, path string) Node {
	return Div(A(Href(path), Text(label)), ID("username"))
}
