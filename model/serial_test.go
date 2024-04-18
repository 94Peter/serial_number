package model_test

import (
	"fmt"
	"testing"

	"github.com/94peter/serial_number/model"
)

func TestGetSerial(t *testing.T) {
	s := model.NewSerial(&model.Config{
		PersistanceFile: "test.txt",
	})
	s.CreateSerial("abc", 100)
	// Testing for a prefix that exists in the map
	t.Run("Existing Prefix", func(t *testing.T) {
		result, err := s.GetSerial("abc")
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if result != 101 {
			t.Errorf("Expected serial number to be 101, got %d", result)
		}
	})

	// Testing for a prefix that does not exist in the map
	t.Run("Non-Existing Prefix", func(t *testing.T) {
		_, err := s.GetSerial("xyz")
		if err != model.Err_PrefixNotFound {
			t.Errorf("Expected error Err_PrefixNotFound, got %v", err)
		}
	})
}

func Test_Persistance(t *testing.T) {
	s := model.NewSerial(&model.Config{
		PersistanceFile: "test.txt",
	})
	s.CreateSerial("abc", 100)
	fmt.Println(s.GetSerial("abc"))
	fmt.Println(s.GetSerial("abc"))
	err := s.Persistance()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	ss := model.NewSerial(&model.Config{
		PersistanceFile: "test.txt",
	})
	fmt.Println(ss.GetSerial("abc"))
	fmt.Println(ss.GetSerial("abc"))
	t.Error("bbb")
}
