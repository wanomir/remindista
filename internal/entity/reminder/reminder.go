package reminder

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/vedomirr/remindista/internal/domain"
)

type Reminder struct {
	Id           int
	UserId       int
	Text         string
	Tag          string
	Prompt       string
	Frequency    time.Duration
	NextReminder time.Time
}

type ReminderOption func(*Reminder)

func NewReminder(opts ...ReminderOption) (r Reminder) {
	for _, opt := range opts {
		opt(&r)
	}

	return r
}

func WithUserId(id int) ReminderOption {
	return func(r *Reminder) {
		r.UserId = id
	}
}

func (r *Reminder) String() string {
	var str strings.Builder
	str.WriteString(r.Text + "\n")

	if r.Tag != "" {
		str.WriteString(r.Tag + "\n")
	}

	if r.Prompt != "" {
		str.WriteString(r.Prompt + "\n")
	}

	str.WriteString(r.FreqeuncyString() + "\n")
	str.WriteString(r.IdToHex())

	return str.String()
}

func (r *Reminder) StringMdV2() string {
	var str strings.Builder
	str.WriteString(r.TextMdV2() + "\n")

	if r.Tag != "" {
		str.WriteString(r.TagMdV2() + "\n")
	}

	if r.Prompt != "" {
		str.WriteString("||" + r.PromptMdV2() + "||\n")
	}

	str.WriteString("`" + r.FreqeuncyString() + "`")

	return str.String()
}

func (r *Reminder) TextMdV2() string {
	return r.escapedMdV2(r.Text)
}

func (r *Reminder) TagMdV2() string {
	return r.escapedMdV2(r.Tag)
}

func (r *Reminder) PromptMdV2() string {
	return r.escapedMdV2(r.Prompt)
}

func (r *Reminder) FreqeuncyString() string {
	var str strings.Builder

	days := int(r.Frequency / (time.Hour * 24))
	if days > 0 {
		if days == 1 {
			str.WriteString("1 day ")
		} else {
			str.WriteString(fmt.Sprintf("%d days ", days))
		}
	}

	hours := int((r.Frequency - time.Duration(days)*time.Hour*24).Hours())
	if hours > 0 {
		if hours == 1 {
			str.WriteString("1 hour ")
		} else {
			str.WriteString(fmt.Sprintf("%d hours ", hours))
		}
	}

	mins := int((r.Frequency - time.Duration(days)*time.Hour*24 - time.Duration(hours)*time.Hour).Minutes())
	if mins > 0 {
		if mins == 1 {
			str.WriteString("1 minute")
		} else {
			str.WriteString(fmt.Sprintf("%d minutes", mins))
		}
	}

	result := str.String()
	if result != "" {
		result = "every " + result
	}

	return result
}

func (r *Reminder) SetTag(s string) error {
	if len(s) < 2 {
		return domain.ErrorShortTag
	}

	if s[0] != '#' {
		s = "#" + s
	}

	if s == "#no_tag" {
		r.Tag = ""
		return nil
	}

	r.Tag = strings.ToLower(s)

	return nil
}

func (r *Reminder) TagMatches(s string) bool {
	if len(s) < 2 {
		return false
	}

	if s[0] != '#' {
		s = "#" + s
	}

	if s == "#no_tag" {
		s = ""
	}

	return r.Tag == strings.ToLower(s)
}

func (r *Reminder) NextReminderString() string {
	return r.NextReminder.Format("on Jan _2 2006 at 15:04:05")
}

// TODO: remove this function
func (r *Reminder) HexToId(s string) error {
	if len(s) < 3 || s[:2] != "0x" {
		return errors.New("Must be a hexadecimal number with leading '0x'")
	}

	i, err := strconv.ParseInt(s[2:], 16, 64)
	if err != nil {
		return fmt.Errorf("Failed to parse hex: %w", err)
	}

	r.Id = int(i)

	return nil
}

// TODO: remove this function
func (r *Reminder) IdToHex() string {
	return fmt.Sprintf("%#x", r.Id)
}

func (r *Reminder) SetFrequency(s string) error {
	parts := strings.Split(s, " ")
	if len(parts) != 2 {
		return fmt.Errorf("should be 2 parts, got: %v", parts)
	}

	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return fmt.Errorf("N should natrual number, got: %v", parts[0])
	}

	t := strings.ToLower(strings.TrimSpace(parts[1]))
	switch {
	case strings.Contains(t, "minute"):
		r.Frequency = time.Duration(n) * time.Minute

	case strings.Contains(t, "hour"):
		r.Frequency = time.Duration(n) * time.Hour

	case strings.Contains(t, "day"):
		r.Frequency = time.Duration(n) * time.Hour * 24

	default:
		return fmt.Errorf("didn't recognize time unit: %v", t)
	}

	return nil
}

func (r *Reminder) UpdateNextReminder(userTime time.Time, floor, ceil time.Duration) {
	next := userTime.Add(r.RandomizedDuration(r.Frequency))
	nextTrunc := time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

	if next.Before(nextTrunc.Add(floor)) {
		r.NextReminder = nextTrunc.Add(floor).Add(r.RandomizedDuration(r.Frequency / 10))
		return
	}

	if next.After(nextTrunc.Add(ceil)) {
		r.NextReminder = nextTrunc.Add(24*time.Hour + floor).Add(r.RandomizedDuration(min(r.Frequency/5, time.Hour)))
		return
	}

	r.NextReminder = next
}

func (r *Reminder) RandomizedDuration(d time.Duration) time.Duration {
	x := d.Nanoseconds()      // toal duration in nanoseconds
	x += rand.Int63n(x) - x/4 // randomized duration by quarter distance

	return time.Duration(x)
}

func (r *Reminder) Keyboard() domain.Keyboard {
	return domain.Keyboard{[]domain.Item{
		{Key: "Delete", Val: fmt.Sprintf("%s %d", domain.CallbackDelete, r.Id)},
		{Key: "Update", Val: fmt.Sprintf("%s %d", domain.CallbackUpdate, r.Id)},
		{Key: "Freq ×2", Val: fmt.Sprintf("%s %d", domain.CallbackIncreaseFrequency, r.Id)},
		{Key: "Freq ÷2", Val: fmt.Sprintf("%s %d", domain.CallbackDecreaseFrequency, r.Id)},
	}}
}

/*
Inside pre and code entities, all '`' and '\' characters must be escaped with a preceding '\' character.
Inside the (...) part of the inline link and custom emoji definition, all ')' and '\' must be escaped with a preceding '\' character.
In all other places characters '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!' must be escaped with the preceding character '\'.
*/
func (r *Reminder) escapedMdV2(s string) string {
	specialChars := []string{"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"}

	for _, char := range specialChars {
		s = strings.Replace(s, char, "\\"+char, -1)
	}

	return s
}
