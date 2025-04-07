package domain

/*
Please note:

Any character with code between 1 and 126 inclusively can be escaped anywhere with a preceding '\' character, in which case it is treated as an ordinary character and not a part of the markup. This implies that '\' character usually must be escaped with a preceding '\' character.
Inside pre and code entities, all '`' and '\' characters must be escaped with a preceding '\' character.
Inside the (...) part of the inline link and custom emoji definition, all ')' and '\' must be escaped with a preceding '\' character.
In all other places characters '_', '*', '[', ']', '(', ')', '~', '`', '>', '#', '+', '-', '=', '|', '{', '}', '.', '!' must be escaped with the preceding character '\'.
In case of ambiguity between italic and underline entities __ is always greadily treated from left to right as beginning or end of an underline entity, so instead of ___italic underline___ use ___italic underline_**__, adding an empty bold entity as a separator.
A valid emoji must be provided as an alternative value for the custom emoji. The emoji will be shown instead of the custom emoji in places where a custom emoji cannot be displayed (e.g., system notifications) or if the message is forwarded by a non-premium user. It is recommended to use the emoji from the emoji field of the custom emoji sticker.
Custom emoji entities can only be used by bots that purchased additional usernames on Fragment.
*/

// general replies
const (
	ReplyStart         = "Hello %s\\!\nThis bot sends you reminders in a randomly timed manner\\."
	ReplyCreateNewUser = "Hello %s\\! Let's set up your user profile\\."
	ReplyHelp          = "Here are some commands use can use:\n" +
		"/start — Remindista startup\n" +
		"/help — Get instructions on how to use Remindista\n" +
		"/update\\_user — User profile settings\n" +
		"/add — Add new reminder\n" +
		"/list — List reminders\n" +
		"/update — Edit reminder parameters\n" +
		"/delete — Delete reminder\\(s\\)"
	ReplyUnkonwCommand  = "Unknown command 🤨\\."
	ReplyFailedFindUser = "Sorry, user profile data is not set 😕\\.\nUse /update_user update your profile\\."
	ReplyUnknown        = `🤨`
	ReplyCancel         = "Cancelled 👌"
	ReplyDone           = "Done ✅"
	ReplyYesNo          = `Say __yes__ or __no__\\.`
	ReplyCannotSkip     = "Sorry, cannot skip this step\\."
)

// error replies
const (
	ReplyErrorGettingReminder  = "Couldn't get reminder\\(s\\) 😢: %w"
	ReplyErrorCreatingReminder = "Couldn't create reminder 😐: %w"
	ReplyErrorUpdatingReminder = "Failed to update reminder 😐: %w"
	ReplyErrorDeletingReminder = "Failed to delete reminder\\(s\\) 😐: %w\\."
	ReplyErrorParsingLocation  = "Couldn't recognize location 😢: %w\\. Try one more time please\\."
	ReplyErrorParsingFrequency = "Couldn't parse frequency 😐: %w\\. Try again\\?"
	ReplyErrorParsingId        = "Couldn't parse reminder id 😢: %w\\. Try one more time\\."
	ReplyErrorParsingTag       = "Couldn't set reminder's tag 😢: %w\\. Try one more time\\."
	ReplyErrorParsingTime      = "Couldn't set time 😢: %w\\. Try one more time\\."
	ReplyErrorCreatingUser     = "Couldn't create user 😢: %w"
	ReplyErrorUpdatingUser     = "Couldn't update user 😢: %w"
)

// f-strings
const (
	ReplyReminderSet             = "Reminder set\\. Next reminder is _%s_\\."
	ReplySetLocation             = "Specify your location or send `skip` to leave _%s_\\."
	ReplySetWindowFloor          = "Set the lower time boundary for your notifications or send `skip` to leave _%s_\\."
	ReplySetWindowCeil           = "Set the upper time boundary for your notifications or send `skip` to leave _%s_\\."
	ReplyDeletedMultiple         = "Deleted %d reminder\\(s\\) ✅"
	ReplyUpdateReminderText      = "Send new reminder text or `skip` to keep:\n\n_%s_"
	ReplyUpdateReminderTag       = "Specify new reminder's tag or send `skip` to keep _%s_\\."
	ReplyUpdateReminderFrequency = "Set new reminder frequency or leave _%s_\\."
	ReplyUpdateReminderPrompt    = "Send new prompt text or `skip` to keep:\n\n_%s_"
	ReplyFrequencyUpdated        = "Reminder frequency updated\\. New frequency is _%s_\\."
	ReplyMaximumFrequency        = "Frequency is at its maximum of once per year\\."
	ReplyMinimumFrequency        = "Frequency is at its minimum of 1 minute\\."
)

// other replies
const (
	ReplySetReminderText      = "Send the text of your reminder\\."
	ReplySetReminderTag       = "Specify reminder's tag or send `skip`\\."
	ReplySetReminderPrompt    = "Send reminder prompt text or `skip`\\."
	ReplySetReminderFrequency = "Specify reminder frequency\\. Examples:\n2 days\n1 hour\n45 minutes"

	ReplyUserUpdated = "User profile updated\\. To change user setings, use /update\\_user\\."

	ReplyListReminders      = "Send tag name or `no\\_tag` to list reminders by tag\\. Say `all` to list all reminders, or `cancel` to exit\\."
	ReplyNoReminders        = "No reminders found\\. Use /add to create a reminder\\."
	ReplyNoRemindersWithTag = "No reminders found\\. Try another tag or list all reminders\\."
	ReplyListAnotherTag     = "Specify another tag, list all reminders, or send `cancel` to exit\\."

	ReplyModes              = "Send `id`, `tag`, or `all` to pick delete mode\\. To exit you can say `cancel`\\."
	ReplySetMode            = "You can delete reminders by id, tag, or just delete all of them\\. " + ReplyModes
	ReplySendId             = "Send rimender id, the one that looks like this: `0xfff`"
	ReplyNoSuchId           = "Couldn't find reminder with this id\\. Try another one\\?"
	ReplyDeleteMore         = "Delete more\\?"
	ReplySendTag            = "Send the tag you want to clear\\."
	ReplyNoSuchTag          = "Couldn't find reminders with this tag\\. Try another one\\?"
	ReplyConfirmDeletingAll = "Are you sure you want to delete all reminders\\? Answer yes or no\\."

	ReplyNoPromt = "(no prompt)"

	ReplyeConfirmDelete = "Are you sure you want to delete this reminder\\? Answer yes or no\\."
)
