# Monkey Language

This code follows the books [Writing An Interpreter In Go](https://interpreterbook.com/) and [Writing A Compiler In Go](https://compilerbook.com/) from Thorsten Ball. It also includes the Macro System implementation from the post [The Lost Chapter: A Macro System For Monkey](https://interpreterbook.com/lost/) from the same author.

This version differs from the original author's code in a few places that were convenient for me. For example, I've added a small utility that generates DOT files of the generated AST for easier visualization.

I highly recommend going through the books to learn the basics about programmiing language design, implementation, interpretation and compilation.

## Quick Review

Here's a small snippet that shows most of the languge features:

```rust
let name = "Monkey";
let age = 1;
let inspirations = ["Scheme", "Lisp", "JavaScript", "Clojure"];
let book = {
    "title": "Writing A Compiler In Go",
    "author": "Thorsten Ball",
    "prequel": "Writing An Interpreter In Go"
};

let printBookName = fn(book) {
    let title = book["title"];
    let author = book["author"];
    inspect(author + " - " + title);
};

printBookName(book);
// => prints: "Thorsten Ball - Writing A Compiler In Go"

let fibonacci = fn(x) {
    if (x == 0) {
        0
    } else {
        if (x == 1) {
            return 1;
        } else {
            fibonacci(x - 1) + fibonacci(x - 2);
        }
    }
};

let map = fn(arr, f) {
    let iter = fn(arr, accumulated) {
        if (len(arr) == 0) {
            accumulated
        } else {
            iter(rest(arr), push(accumulated, f(first(arr))));
        }
    };

    iter(arr, []);
};


let numbers = [1, 1 + 1, 4 - 1, 2 * 2, 2 + 3, 12 / 2];
map(numbers, fibonacci);
// => returns: [1, 1, 2, 3, 5, 8]
```

Features include:

* integers
* booleans
* strings
* arrays
* hashes
* prefix-, infix- and index operators
* conditionals
* global and local bindings
* first-class functions
* return statements
* closures
* _macros_ (TODO)