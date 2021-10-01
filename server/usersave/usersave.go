package usersave

import (
	"encoding/json"
	"fmt"
	"io"
)

type Cycle = string
type Tag = string

type JSONExpense struct {
	Name   string `json:"name"`
	Amount int32  `json:"amount"`
	Tag    Tag    `json:"tag"`
}

type JSONSavings struct {
	Name     string `json:"name"`
	Goal     int32  `json:"goal"`
	Amount   int32  `json:"amount"`
	Deadline int    `json:"deadline"`
}

type JSONUserSave struct {
	Cycle         Cycle         `json:"cycle"`
	Income        int32         `json:"income"`
	SavingsAmount int32         `json:"savingsAmount"`
	Savings       []JSONSavings `json:"savings"`
	Expenses      []JSONExpense `json:"expenses"`
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
