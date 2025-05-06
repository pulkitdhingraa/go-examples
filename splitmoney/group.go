package main

type Group struct {
	ID string
	Name string
	Members []*User
	Expenses []*Expense
}

func NewGroup(id, name string) *Group {
	return &Group{
		ID: id,
		Name: name,
		Members: []*User{},
		Expenses: []*Expense{},
	}
}

func (g *Group) AddMember(member *User) {
	g.Members = append(g.Members, member)
}

func (g *Group) AddExpense(expense *Expense) {
	g.Expenses = append(g.Expenses, expense)
}