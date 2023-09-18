package data

import "fmt"

type PersonalInfo struct {
	Field1 string
	Field2 int
	Field3 float32
}

func (pi PersonalInfo) TestMethod() {
	fmt.Println("Hello, World!")
}
