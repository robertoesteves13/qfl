[![GoDoc](https://godoc.org/github.com/robertoesteves13/qfl?status.png)](https://godoc.org/github.com/robertoesteves13/qfl)
# Overview

This repository contains a parser for the QFL language that transforms into
a filter data structure that stores a collection of rules for a given record,
and a SQL Query builder that can convert these rules into WHERE statements.

# Library
Checkout [documentation](https://pkg.go.dev/github.com/robertoesteves13/qfl) for examples on how to use the library.

# Language
QFL is a simple filtering language for queries designed to have just enough
functionality to provide a quite competent filtering scheme for your REST API
or something similar.

## Syntax
This is the syntax of how you define a filter:
```
[comparator!]value
```

If only the value is passed, it will use the `eq` comparator. Additionally, if
you want to filter for some data that is equal to one of the specified values,
you can specify a list of values separated with commas. Note that this is only
supported on the `eq` comparator:
```
eq!value[,value...]
```

You can combine more than one filter by using the bar:
```
filter|filter
```

## Comparators
- eq: Equals
- lt: Less than
- gt: Greater than
- le: Less or equal
- ge: Greater or equal
- lk: Searches for similar string

## Symbols
- | (bar): combine filters from both sides
- , (comma): Separate elements in a list
- ! (mark): Indicate start of a value
- \ (backslash): escape the character in front of it (only affects symbols)

# Contributing
If you found a bug, missing documentation or have any improvements for performance,
you can contribute to the codebase by opening a PR addressing the problem.
If you want to add a feature, open an issue describing why it is necessary and
to be sure we both agree that this needs to be implemented, so you don't waste
time implementing something that won't be accepted.

This library should be considered stable, so there shouldn't be any breaking
changes. If however there is a feature that most users would benefit to, a
new version of the module might be made for it. Any breaking changes will be
marked as a new major version of the library. You can also just copy the source
code into your project since it's in public domain.

