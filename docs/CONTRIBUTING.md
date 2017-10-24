# Contributing

Issues, PR's or any other kind of contribution is in general very welcome.

This is currently developed using golang 1.8 and 1.9, so I don't know if any other version is working.

When you contribute code wise,
please make sure to run `make test`, that will execute tests, `go vet` and `gofmt` and will exit with a non zero status if failing. This will run in Travis anyway and will prevent you from committing again. :)


If you make changes which are changing the behavior of the program, please make sure to correctly adjust the read me 
according these changes.

Comments should be in the code if the code is not self explanatory, but you should also not clutter the code 
with to many comments.
