# bartr/gohandlers

This project is a collection of Go http handlers.

* eventgrid
  * handler that parses the event grid "envelope" and handles validation events
  * this handler should be chained to a custom handler that does the actual work
  * see <https://github.com/bartr/m4> for a sample implementation
* rawrequest
  * this handler logs and displays the raw http request
  * this handler is a developer tool ONLY and is NOT designed for production!
  * see <https://github.com/bartr/gowac> for a sample implementation
* logb
  * simple log wrapper for chaining requests
  * both above samples use logb for logging

## Developer Prerequisites

* Go 1.11
