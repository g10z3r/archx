package main

type Go interface {
	Straight() error
	Left() error
	Right() error
	Back() error
}

type Transport interface {
	StartEngine() error
	StopEngine() error

	Go
}

type Address struct {
	street string
	city   string
	state  string
	zip    string
}

type Person struct {
	FirstName,
	LastName string
	Age       int
	Address   Address
	Transport Transport
	Skill     struct {
		Intelligence float32
		Speed        float32
	}
}

func (p *Person) ChangeFirstName(newFirstName string) *Person {
	if newFirstName != p.FirstName {
		p.FirstName = newFirstName
	}

	return p
}

func (p *Person) ChangeLastName(newLastName string) *Person {
	if newLastName != p.LastName {
		p.LastName = newLastName
	}

	return p
}

func (p *Person) AgeInc() *Person {
	p.Age = p.Age + 1

	return p
}
