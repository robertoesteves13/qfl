# Overview

QFL is a simple filtering language for queries designed to have just enough
functionality to provide a quite competent filtering scheme for your REST API
or something similar.

# Language

## Syntax
This is the syntax of how you define a filter:
```
comparator!value[,value...]
```

You can combine more than one filter by using the bar:
```
filter|filter
```


## Comparators
- eq: Equals (default if you pass only the value)
- lt: Less than
- gt: Greater than
- le: Less or equal
- ge: Greater or equal
- like: Searches for similar string
- is: Comma-separated list of possible options (similar to `IN ("a", "b")` in SQL)

## Symbols
- | (bar): combine filters from both sides
- , (comma): Separate elements in a list
- ! (mark): Indicate start of a value
- \ (backslash): escape the character in front of it (only affects symbols)

# Contributing

This is my first library, so things might not be optimal at the moment. I would
be grateful if you can contribute to the project and reduce the todo list items.

# TODO

- Better documentation
- More rigorous testing
- Maybe optimize

# Non-Goals

- Adding more features, this is just enough
