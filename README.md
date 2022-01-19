# httperrorhandler
Error handling hook &amp; helper function to simplify writing API handler methods in Go.

# Example usage
```
import (
	httperr "github.com/peeveen/httperrorhandler"
)

...

func (s *server) errorHandler(w http.ResponseWriter, r *http.Request, e *httperrorhandler.Error) {
	var reason = e.Error.Error()
	if e.Message != "" {
		reason = fmt.Sprintf("%s (%s)", e.Message, reason)
	}
	logrus.Error(reason)
	w.WriteHeader(e.HTTPStatus)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(reason))
}

func (s *server) handleSomeAPIRequest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		httperr.Handle(w, r, func(w http.ResponseWriter, r *http.Request) *httperr.Error {
      // Do stuff
      if somethingHasGoneTerriblyWrong {
        return &httperr.Error{HTTPStatus: http.StatusInternalServerError, Error: errors.New("defluter valve blockage")}
      }
      err := thirdPartyDoodah.doSomething()
      if err!=nil {
        return &httperr.Error{HTTPStatus: http.StatusInternalServerError, Error: err, Message: "The doodah has failed"}
      }
    }, s.errorHandler)
  }
}

```
