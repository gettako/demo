package main

import (
	"fmt"
	"charm.land/lipgloss/v2/tree"
)

func main() {
	t := tree.Root("Label")
	t.Child("Child 1").Child("Child 2")
	fmt.Println(t.String())
}
