package main

func main() {
	service := GetSplitMoneyService()

	// Users
	user1 := NewUser("U1", "Joe", "joe@example.com")
	user2 := NewUser("U2", "Bob", "bob@example.com")
	user3 := NewUser("U3", "Dan", "dan@example.com")
	service.AddUser(user1)
	service.AddUser(user2)
	service.AddUser(user3)

	// groups and its users (Joe and Bob are friends, Bob and Dan live in the same flat)
	group1 := NewGroup("G1", "Friends")
	group2 := NewGroup("G2", "Flat")
	group1.AddMember(user1)
	group1.AddMember(user2)
	service.AddGroup(group1)
	group2.AddMember(user2)
	group2.AddMember(user3)
	service.AddGroup(group2)

	// expenses
	// Flat
	expense1 := NewExpense("E1", 300, "Rent", user2) // Bob paid the rent for both Bob and Dan
	expense1.AddSplit(NewEqualSplit(user2))
	expense1.AddSplit(NewEqualSplit(user3))
	service.AddExpense(group2.ID, expense1)
	expense2 := NewExpense("E2", 10.50, "Food", user3) // Dan paid for the burger and coffee of Bob
	expense2.AddSplit(NewPercentSplit(user2, 60))
	expense2.AddSplit(NewPercentSplit(user3, 40))
	service.AddExpense(group2.ID, expense2)

	// Friends
	expense3 := NewExpense("E3", 20, "Stuff", user1) // Joe paid on behalf of Bob
	expense3.AddSplit(NewPercentSplit(user2, 100))
	service.AddExpense(group1.ID, expense3)

	service.PrintBalances()
}