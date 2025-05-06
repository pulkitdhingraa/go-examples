package main

import (
	"fmt"
	"sync"
)

type SplitMoneyService struct {
	users  map[string]*User
	groups map[string]*Group
}

var (
	instance *SplitMoneyService
	once     sync.Once
)

func GetSplitMoneyService() *SplitMoneyService {
	once.Do(func() {
		instance = &SplitMoneyService{
			users:  make(map[string]*User),
			groups: make(map[string]*Group),
		}
	})
	return instance
}

func (s *SplitMoneyService) AddUser(user *User) {
	s.users[user.ID] = user
}

func (s *SplitMoneyService) AddGroup(group *Group) {
	s.groups[group.ID] = group
}

func (s *SplitMoneyService) AddExpense(groupID string, expense *Expense) error {
	group, exists := s.groups[groupID]
	if !exists {
		return fmt.Errorf("Group doesn't exist")
	}
	group.AddExpense(expense)
	s.splitExpense(expense)
	s.updateBalances(expense)
	return nil
}

func (s *SplitMoneyService) splitExpense(expense *Expense) {
	totalAmount := expense.Amount
	totalSplits := len(expense.Splits)

	for _, split := range expense.Splits {
		switch v := split.(type) {
		case *EqualSplit:
			v.SetAmount(totalAmount / float64(totalSplits))
		case *PercentSplit:
			v.SetAmount(totalAmount * v.Percent / 100.0)
		}
	}
}

func (s *SplitMoneyService) updateBalances(expense *Expense) {
	for _, split := range expense.Splits {
		paidBy := expense.PaidBy
		user := split.GetUser()
		amount := split.GetAmount()
		if paidBy != user {
			s.updateBalance(paidBy, user, amount)
			s.updateBalance(user, paidBy, -amount)
		}
	}
}

func (s *SplitMoneyService) updateBalance(user1, user2 *User, amount float64) {
	key := user1.ID + ":" + user2.ID
	user1.Balances[key] += amount
}

func (s *SplitMoneyService) PrintBalances() {
	for _, user := range s.users {
		for key, balance := range user.Balances {
			fmt.Printf("Balance with %s: %.2f\n", key, balance)
		}
	}
}
