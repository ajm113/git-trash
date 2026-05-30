# git-trash
*Simple trashbin functionality with git that supports wildcard matching and removing old branches for big projects.*

### Features

- Marks branches for deletion by appending a `git-trash/` prefix.
- Easy `empty` command that removes branches with the `git-trash/` prefix.
- Glob support for pattern matching branches you want to mark for deletion.
- Mark branches for deletion that are older then the specified number of days.


### Why?

Usually when working with large Git repos they start accomulating dead branches locally or remote.
I wrote this tool as a easy way to manage many branches without the fear of accidently permenly deleting em 
and make this ci tool easy to interface with CI/CD pipelines if need be. *More coming soon for GitHub workflows*


## Installing

*Requires Golang 1.26.3 to compile.*

- `git clone https://github.com/ajm113/git-trash.git`
- `cd git-trash`
- To compile, you may run `mage build` (if you have mage installed) or run `go build -o dist/git-trash`.

### For Linux/Unix or OSX

Simply `cp ./dist/git-trash ~/.local/bin` or which ever folder your OS recommends for user level binaries. 

### Windows

I haven't touched Windows in 10 years. TBD!

