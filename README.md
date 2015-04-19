# polcmp

polcmp is a program that uses a git repository as a backend, via git-html5.

**but wait, that's a JS library! why?** this is a client side application that
clones the target repo, reads over the markdown files, and generates a page
dynamically that compares them side by side.

It's called polcmp because I'll be using this for
[polcmp](https://github.com/polcmp).

By client-side application, we mean a client-side gopherjs library that's
embedded into an index.html page.

## how it works
* navigate to index.html
* loads dependent JS files
* loads polcmp
* requests a html5 filesystem
* clones the specified repository into the html5 filesystem (depth 1)
* does stuff with the files.
    * in my case, reads each file, parses markdown, renders new div

## how to build
`gopherjs build && mv -fv polcmp.js static && mv -fv polcmp.js.map static`

## but wait, CORS!
Currently we're proxying through
[cors-anywhere](https://github.com/Rob--W/cors-anywhere), but obviously that's
not a permanent solution. But especially for [polcmp], it's vital that it
remains on an open repository that is easily accessible and merge requests can
be seen, and is not on my servers or another private one so owner bias isn't an
issue.

[relevant issue here](https://github.com/isaacs/github/issues/263).

## do you even bower bro
i'm not a frontend dev ;-;
