package main

type Address struct {
	street string
	city   string
	state  string
	zip    string
}

type Person struct {
	Name    string
	Age     int
	Address Address
}

func (p *Person) Rename(newName string) *Person {
	if newName != p.Name {
		p.Name = newName
	}

	return p
}

func (p *Person) AgeInc() *Person {
	p.Age = p.Age + 1

	return p
}
