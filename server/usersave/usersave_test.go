package usersave

import (
	"io"
	"os"
	"testing"
)

func readValidJSON() io.ReadCloser {
	file, err := os.Open("examples/valid.json")
	if err != nil {
		panic(err)
	}

	return file
}

func TestDecodeUserSave(t *testing.T) {
	validJSON := readValidJSON()
	defer validJSON.Close()

	userSave, err := DecodeUserSave(validJSON)
	if err != nil {
		t.Error(err)
	}

	if userSave.Cycle != "Fortnightly" {
		t.Errorf("expected Fortnightly cycle, got %s", userSave.Cycle)
	}

	if userSave.Income != 900033 {
		t.Errorf("expected 900033 cents, got %d", userSave.Income)
	}

	if userSave.SavingsAmount != 405015 {
		t.Errorf("expected 405015 cents, got %d", userSave.SavingsAmount)
	}

	if size := len(userSave.Savings); size != 1 {
		t.Errorf("expected savings list of len 1, got %d", size)
	}

	if size := len(userSave.Expenses); size != 2 {
		t.Errorf("expected expense list of len 3, got %d", size)
	}

	if userSave.Expenses[0].Tag != "Housing" {
		t.Errorf("expected Housing tag, got %s", userSave.Expenses[0].Tag)
	}
}
