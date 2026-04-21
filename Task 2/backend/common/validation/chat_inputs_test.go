package validation

import "testing"

func TestDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "valid", input: "2026-04-21"},
		{name: "empty", input: "", wantErr: true},
		{name: "wrong format", input: "21-04-2026", wantErr: true},
		{name: "invalid day", input: "2026-02-30", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Date(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("Date(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			}
		})
	}
}

func TestMealType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "Lunch", want: "lunch"},
		{name: "underscore", input: "late_snack", want: "late_snack"},
		{name: "empty", input: "", wantErr: true},
		{name: "space", input: "late snack", wantErr: true},
		{name: "symbol", input: "snack!", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := MealType(test.input)
			if (err != nil) != test.wantErr {
				t.Fatalf("MealType(%q) error = %v, wantErr %v", test.input, err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("MealType(%q) = %q, want %q", test.input, got, test.want)
			}
		})
	}
}

func TestEnums(t *testing.T) {
	if _, err := Status("yes"); err != nil {
		t.Fatalf("Status validation failed: %v", err)
	}
	if _, err := Location("wfh"); err != nil {
		t.Fatalf("Location validation failed: %v", err)
	}
	if _, err := DayStatusType("special_event"); err != nil {
		t.Fatalf("DayStatusType validation failed: %v", err)
	}
	if _, err := Note("   ", true); err == nil {
		t.Fatal("expected note validation error for required blank note")
	}
}
