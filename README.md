# imdepth

Constructs an image from images with different focues.

usage:
  - `imdepth [flags] <dir_name>`

rules:
  - `<dir_name>` must constain a files named `<number>.<ext>`
  - this `<number>`s used as heights: `[0..255]`
  - `<ext>` should be jpg or png

flags:
  - `-r=2`: area around every point to get its sharp value

One can simply run this project by typing `go run main.go -r 5 wine` or `go build; ./imdepth -r 5 wine`.

You have to install Go 1.4 first.
