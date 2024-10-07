package main

import (
	"iguana/src"
)

func main() {
	iguana.GenerateHTML("./rsc/index.html", "iguana.html")

	// ! Loosing leading '<' of parent file
	// iguana.GenerateHTML("./rsc/index_oneline.html", "iguana.html")
}