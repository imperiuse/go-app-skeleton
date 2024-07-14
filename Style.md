# Style Guide

This document contains suggestions (not a rule). 

However, if you would like to have clean and maintainable code, please, follow them.
Let's make our code simple, clean and beauty!

## Style
* Configure `go fmt`, `go imports`, `golangci-lint` in your IDE (e.g. IDE -> add Watchers -> go fmt)
* golangci-lint's project settings  you can find here -> `.golangci-lint.yml`
* In general we suppose to use [Uber Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
* Also please [Go concurrency checklist](https://github.com/golang/go/wiki/CodeReviewComments) and [Code Review Comments](https://github.com/code-review-checklists/go-concurrency)
* We expect that you have own profession experience as a software developer, and if you need do some extraordinary thing, you can do them. 

## File naming
* Not use this symbols: `-`, `_`, `:` (except some extra ordinary situations)
* Not user UPPER symbols

## Package dividing
* We strongly suggest use de-facto [typical go project structure](https://github.com/golang-standards/project-layout)

## Variable naming 
* Consts - `CamelCase`
* Vars - `lowCamelCase`
* Func - `CamelCase` || `lowCamelCase`
* context always ctx
* Reqctx context always `rctx` 
* Request context always `reqCtx` (for `fasthttp` only)
* Request context always `с` (for `gin.Context` only)
* Preferable name of the interface ends with suffix `-er`  (example: `Closer`, `Writer`, `Reader` and etc.)

## Coding rules
* Must comment non-standard or tricky places
* Try use less comments, try use meaningfully naming for vars, const, struct, func 
* Use empty line to keep code for keep code spacy
* New line at the end of each file

## Limits for functions
* `ctx context.Context` should be first param in every func (if ctx is necessary)
* Inline complexity not more than (5)
* Line lengths to many than (120)
* Length of defer fun not more than (3 line)
* Func length not more than 70 line
* Number of params not more than 7
* All params should have own type  ~~f(a,b string)~~ --> f(a string, b string)
* For DI purposes all obj/structure should have own "constructors"
* Methods with panic should begin with Must
* Not use `else` if you can use `return`

## Other moments
* DON'T PANIC! Do not use panic approach 
* NO INIT (not use init func)
* Use DI Pattern (we use uber.fx)
* Use strategy pattern 
* Decouple code strategy, use interface, it's help easy integrate mocks and do unit tests. 

## Errors
* Do not ignore any errors -> Logging or Metrics

## Unit tests
* Every pull request should not decrease test coverage

Feel free for add more rule, if you think there are necessary.
Have a nice codding!

MORE INFO HERE -> Google style guide -> https://google.github.io/styleguide/go/index
