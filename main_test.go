package main

import "testing"

func TestParsePorts(t *testing.T) {
	validInputs := []struct {
		single   []string
		multiple []string
	}{
		{single: []string{"22", "80", "8080"}, multiple: []string{"80,443,8080", "123,12345,54321"}},
	}

	t.Run("given single port", func(t *testing.T) {
		for _, validInput := range validInputs {
			for _, port := range validInput.single {
				port, err := parsePorts(port)
				if err != nil {
					t.Errorf("Error parsing single port: %v", err)
				}
				if len(port) != 1 {
					t.Errorf("Expected 1 port, got %d", len(port))
				}
				t.Log("Single port parsed successfully:", port[0])
			}
		}
	})

	t.Run("given multiple ports", func(t *testing.T) {
		for _, validInput := range validInputs {
			for _, ports := range validInput.multiple {
				ports, err := parsePorts(ports)
				if err != nil {
					t.Errorf("Error parsing multiple ports: %v", err)
				}
				if len(ports) < 2 {
					t.Errorf("Expected more than 1 port, got %d", len(ports))
				}
				t.Log("Multiple ports parsed successfully:", ports)
			}
		}
	})
}
