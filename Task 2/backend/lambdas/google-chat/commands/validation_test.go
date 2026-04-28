package commands

import "testing"

func TestGoogleChatCommandFieldValidation(t *testing.T) {
	tests := []struct {
		name    string
		handler func() string
		want    string
	}{
		{
			name: "meal view invalid date",
			handler: func() string {
				return HandleMeal("view 2026-13-40", "12345").Text
			},
			want: `Error: invalid date "2026-13-40": expected YYYY-MM-DD`,
		},
		{
			name: "meal set invalid meal type",
			handler: func() string {
				return HandleMeal("set 2026-04-24 Lunch! YES", "12345").Text
			},
			want: `Error: invalid meal_type "Lunch!": use 1-32 lowercase letters, numbers, hyphens, or underscores`,
		},
		{
			name: "team meal set missing user id",
			handler: func() string {
				return HandleTeamMeal("set <> 2026-04-24 lunch YES", "12345").Text
			},
			want: "Error: user_id is required",
		},
		{
			name: "admin meal set invalid status",
			handler: func() string {
				return HandleAdminMeal("set users/99999 2026-04-24 lunch MAYBE", "12345").Text
			},
			want: `Error: invalid status "MAYBE": must be YES or NO`,
		},
		{
			name: "work location set invalid location",
			handler: func() string {
				return HandleWorkLocation("set 2026-04-24 HOME", "12345").Text
			},
			want: `Error: invalid location "HOME": must be OFFICE or WFH`,
		},
		{
			name: "day status special event note required",
			handler: func() string {
				return HandleDayStatus("set 2026-04-24 SPECIAL_EVENT", "12345").Text
			},
			want: "Error: note is required for SPECIAL_EVENT",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.handler(); got != test.want {
				t.Fatalf("response text = %q, want %q", got, test.want)
			}
		})
	}
}
