package domain

// keyboards
var (
	KbCancel    = Keyboard{[]Item{{Key: "Cancel", Val: "cancel"}}}
	KbSkip      = Keyboard{[]Item{{Key: "Cancel", Val: "cancel"}, {Key: "Skip", Val: "skip"}}}
	KbAdd       = Keyboard{[]Item{{Key: "Add", Val: "/add"}}}
	KbLocations = Keyboard{
		[]Item{
			{Key: "Cancel", Val: "cancel"},
			{Key: "Skip", Val: "skip"},
			{Key: "Valid locations", Val: "https://en.wikipedia.org/wiki/List_of_tz_database_time_zones"},
		},
	}
	KbWindowFloor = Keyboard{
		[]Item{{Key: "4:00", Val: "4:00"}, {Key: "8:00", Val: "8:00"}},
		[]Item{{Key: "4:30", Val: "4:30"}, {Key: "8:30", Val: "8:30"}},
		[]Item{{Key: "5:00", Val: "5:00"}, {Key: "9:00", Val: "9:00"}},
		[]Item{{Key: "5:30", Val: "5:30"}, {Key: "9:30", Val: "9:30"}},
		[]Item{{Key: "6:00", Val: "6:00"}, {Key: "10:00", Val: "10:00"}},
		[]Item{{Key: "6:30", Val: "6:30"}, {Key: "10:30", Val: "10:30"}},
		[]Item{{Key: "7:00", Val: "7:00"}, {Key: "11:00", Val: "11:00"}},
		[]Item{{Key: "7:30", Val: "7:30"}, {Key: "11:30", Val: "11:30"}},
		[]Item{{Key: "Cancel", Val: "cancel"}, {Key: "Skip", Val: "skip"}},
	}
	KbWindowCeil = Keyboard{
		[]Item{{Key: "18:00", Val: "18:00"}, {Key: "21:00", Val: "21:00"}},
		[]Item{{Key: "18:30", Val: "18:30"}, {Key: "21:30", Val: "21:30"}},
		[]Item{{Key: "19:00", Val: "19:00"}, {Key: "22:00", Val: "22:00"}},
		[]Item{{Key: "19:30", Val: "19:30"}, {Key: "22:30", Val: "22:30"}},
		[]Item{{Key: "20:00", Val: "20:00"}, {Key: "23:00", Val: "23:00"}},
		[]Item{{Key: "20:30", Val: "20:30"}, {Key: "23:30", Val: "23:30"}},
		[]Item{{Key: "Cancel", Val: "cancel"}, {Key: "Skip", Val: "skip"}},
	}
	KbListReminders = Keyboard{[]Item{{"Cancel", "cancel"}, {"No tag", "no_tag"}, {"All", "all"}}}
	KbSetMode       = Keyboard{[]Item{{"Cancel", "cancel"}, {"ID", "id"}, {"Tag", "tag"}, {"All", "all"}}}
	KbYesNo         = Keyboard{[]Item{{"Yes", "yes"}, {"No", "no"}}}
)
