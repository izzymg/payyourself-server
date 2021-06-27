package usersave

import (
	"encoding/json"
	"fmt"
	"io"
)

type JSONCurrency struct {
	Cents int `json:"cents"`
}

type JSONExpense struct {
	Amount JSONCurrency `json:"amount"`
	Order  int          `json:"order"`
}

type JSONExpenseListItem struct {
	Name string      `json:"name"`
	Item JSONExpense `json:"item"`
}

type JSONSavingsGoal struct {
	Goal     JSONCurrency `json:"goal"`
	Amount   JSONCurrency `json:"amount"`
	Deadline int          `json:"deadline"`
	Order    int          `json:"order"`
}

type JSONSavingsGoalListItem struct {
	Name string          `json:"name"`
	Item JSONSavingsGoal `json:"item"`
}

type JSONUserSave struct {
	Income          JSONCurrency              `json:"income"`
	SavingsAmount   JSONCurrency              `json:"savingsAmount"`
	SavingsGoalList []JSONSavingsGoalListItem `json:"savingsGoalList"`
	ExpenseList     []JSONExpenseListItem     `json:"expenseList"`
}

func DecodeUserSave(jsonReader io.Reader) (*JSONUserSave, error) {
	decoder := json.NewDecoder(jsonReader)

	userSave := JSONUserSave{}
	err := decoder.Decode(&userSave)
	if err != nil {
		return nil, fmt.Errorf("failed to decode UserSave JSON: %w", err)
	}

	return &userSave, nil
}

func EncodeUserSave(userSave *JSONUserSave, w io.Writer) error {
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(userSave); err != nil {
		return fmt.Errorf("failed to encode user save %w", err)
	}
	return nil
}
