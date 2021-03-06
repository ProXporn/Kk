package modules

import (
	"fmt"
	"strings"

	tb "gopkg.in/telebot.v3"
)

var (
	menu        = &tb.ReplyMarkup{}
	help_button = sel.Data("HELP", "help_button")
	back_button = sel.Data("Back", "back_button")
)

func Start(c tb.Context) error {
	m := c.Message()
	if strings.HasPrefix(m.Payload, "notes") {
		return PrivateStartNote(c)
	} else if strings.HasPrefix(m.Payload, "allnotes") {
		return PrivateStartNotes(c)
	} else if strings.HasPrefix(m.Payload, "connect") {
		return PrivateConnect(c)
	} else if strings.HasPrefix(m.Payload, "rules") {
		return SendPrivateRules(c)
	}
	if m.Private() {
		menu.Inline(
			menu.Row(menu.URL("Support", "t.me/roseloverx_support"), menu.URL("Updates", "t.me/roseloverx_support")),
			menu.Row(menu.Data("Commands", "help_menu")),
			menu.Row(menu.URL("Add me to your group", fmt.Sprintf("https://t.me/%s?startgroup=true", BOT_USERNAME))))
		b.Send(m.Sender, fmt.Sprintf("Hey there! I am <b>%s</b>.\nIm an Anime themed Group Management Bot, feel free to add me to your groups!", BOT_NAME), menu)
		return nil
	}
	b.Reply(m, "Hey I'm Alive.")
	return nil
}

func Help_Menu(c tb.Context) error {
	if !c.Message().Private() {
		sel.Inline(sel.Row(sel.URL("Click here", fmt.Sprintf("https://t.me/%s?start=help", BOT_USERNAME))))
		c.Reply("Contact me at PM to get help.", sel)
	} else {
		gen_help_buttons(c, help_caption, true)
	}
	return nil
}

func gen_help_buttons(c tb.Context, text string, Reply bool) {
	sel.Inline(sel.Row(sel.Data("AFK", "help_button", "afk"), sel.Data("Admin", "help_button", "admin"), sel.Data("Bans", "help_button", "bans")), sel.Row(sel.Data("Chatbot", "help_button", "chatbot"), sel.Data("Feds", "help_button", "feds"), sel.Data("Greetings", "help_button", "greetings")), sel.Row(sel.Data("Inline", "help_button", "inline"), sel.Data("Locks", "help_button", "locks"), sel.Data("Misc", "help_button", "misc")), sel.Row(sel.Data("Notes", "help_button", "notes"), sel.Data("Pin", "help_button", "pin"), sel.Data("Stickers", "help_button", "stickers")), sel.Row(sel.Data("Warns", "help_button", "warns")))
	if Reply {
		c.Reply(text, sel)
	} else {
		c.Edit(text, sel)
	}
}

func HelpCB(c tb.Context) error {
	arg := c.Callback().Data
	text, ok := help[arg]
	sel.Inline(sel.Row(back_button))
	if ok {
		err := c.Edit(text.(string), &tb.SendOptions{ParseMode: "Markdown", ReplyMarkup: sel})
		check(err)
	}
	return nil
}

func back_cb(c tb.Context) error {
	gen_help_buttons(c, help_caption, false)
	return nil
}
