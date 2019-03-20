# ReDoc golang middleware

Middleware that serves the swagger.json file and redoc ui from a specified path.

### Usage

```go
import (
	"git.holda.io/golang/redoc-middleware"
)

func main() {
	// Set the options for the middleware
	redocOpts := redocmiddle.Opts{
		// BasePath for the UI path, defaults to: /
		BasePath: "/",

		// Path combines with BasePath for the full UI path, defaults to: docs
		Path: "docs",

		// SpecPath the path to find the spec for
		SpecPath: "./docs/swagger.json",

		// RedocURL for the js that generates the redoc site,
		// defaults to: https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js
		RedocURL: "",

		// Title for the documentation site, default to: API documentation
		Title: "Example API documentation",
	}

	// Create the middleware with the specified options
	redocMiddleware := redocmiddle.NewHandler()

	router := chi.NewRouter()

	// Use the middleware in your router
	router.Use(redocMiddleware)
	
	http.ListenAndServe(port, router)
}
```
