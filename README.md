# m4

This project is an Event Grid Web Hook written in Go that supports receiving event grid messages.

The handler can be "chained" as "middle ware".

## Developer Prerequisites

* Developer workstation will need
  * Go 1.11

## Getting Started

  * eventgrid
    * handler that parses the event grid "envelope" and handles validation events
  * logb
    * simple log wrapper for chaining requests

