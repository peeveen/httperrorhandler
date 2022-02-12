# httperrorhandler

Error handling hook &amp; helper function to simplify writing API handler methods in Go.

Tries to follow the [RFC-7807](https://datatracker.ietf.org/doc/html/rfc7807) recommendation for HTTP errors.

# Example usage

```
import (
	httperr "github.com/peeveen/httperrorhandler"
)

...

func (s *server) handleSomeAPIRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httperr.Handle(w, r, func(w http.ResponseWriter, r *http.Request) *httperr.Error {
			// Do stuff
			if somethingHasGoneTerriblyWrong {
				return &httperr.Error{Type: "http://myapp/valve/blockage", Status: http.StatusInternalServerError, Detail: "There has been a defluter valve blockage."}
			}
			err := thirdPartyDoodah.doSomething()
			if err != nil {
				return httperr.Wrap(err, &httperr.Error{Type: "http://myapp/internal", Status: http.StatusInternalServerError, Detail: "The doodah has failed!"})
			}
		}, httperr.DefaultErrorHandler)
	}
}

```

# Error handling

There is a default error handler implementation available (`httperrorhandler.DefaultErrorHandler`), or you can provide your own.

The default implementation will write the HTTP status code and the JSON representation of the error to the response as the `application/problem+json` content type.