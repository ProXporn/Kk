package modules

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/amarnathcjd/yoko/modules/db"
	tb "gopkg.in/telebot.v3"
)

var (
	HyperLink = regexp.MustCompile(`\[(.*?)\]\((.*?)\)`)
	Bold      = regexp.MustCompile(`\*(.*?)\*`)
	Italic    = regexp.MustCompile(`\_(.*?)\_`)
	Strike    = regexp.MustCompile(`\~(.*?)\~`)
	Underline = regexp.MustCompile(`•(.*?)•`)
	Spoiler   = regexp.MustCompile(`\|\|(.*?)\|\|`)
	Code      = regexp.MustCompile("`(.*?)`")
)

func Te(c tb.Context) error {
	r, _ := GetUser(c)
	b, _ := json.Marshal(r)
	fmt.Println(b)
	log.Print(b)
	c.Reply(string(fmt.Sprint(r)))
	return nil
}

func ParseMD(c *tb.Message) string {
	text := c.Text
	cor := 0
	for _, x := range c.Entities {
		offset, length := x.Offset, x.Length
		if x.Type == tb.EntityBold {
			text = string(text[:offset+cor]) + "<b>" + string(text[offset+cor:offset+cor+length]) + "</b>" + string(text[offset+cor+length:])
			cor += 7
		} else if x.Type == tb.EntityCode {
			text = string(text[:offset+cor]) + "<code>" + string(text[offset+cor:offset+cor+length]) + "</code>" + string(text[offset+cor+length:])
			cor += 13
		} else if x.Type == tb.EntityUnderline {
			text = string(text[:offset+cor]) + "<u>" + string(text[offset+cor:offset+cor+length]) + "</u>" + string(text[offset+cor+length:])
			cor += 7
		} else if x.Type == tb.EntityItalic {
			text = string(text[:offset+cor]) + "<i>" + string(text[offset+cor:offset+cor+length]) + "</i>" + string(text[offset+cor+length:])
			cor += 7
		} else if x.Type == tb.EntityStrikethrough {
			text = string(text[:offset+cor]) + "<s>" + string(text[offset+cor:offset+cor+length]) + "</s>" + string(text[offset+cor+length:])
			cor += 7
		} else if x.Type == "spoiler" {
			text = string(text[:offset+cor]) + "<tg-spoiler>" + string(text[offset+cor:offset+cor+length]) + "</tg-spoiler>" + string(text[offset+cor+length:])
			cor += 25
		}
	}
	for _, x := range HyperLink.FindAllStringSubmatch(text, -1) {
		if strings.Contains(x[2], "buttonurl") {
			continue
		}
		text = strings.Replace(text, x[0], fmt.Sprintf("<a href='%s'>%s</a>", x[2], x[1]), -1)
	}
	for _, x := range Bold.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<b>"+x[1]+"</b>", -1)

	}
	for _, x := range Italic.FindAllStringSubmatch(text, -1) {
		pattern, _ := regexp.Compile(`\_\_(.*?)\_\_`)
		if match := pattern.Match([]byte(x[0])); match {
			continue
		}
		text = strings.Replace(text, x[0], "<i>"+x[1]+"</i>", -1)

	}
	for _, x := range Strike.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<s>"+x[1]+"</s>", -1)

	}
	for _, x := range Underline.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<u>"+x[1]+"</u>", -1)

	}
	for _, x := range Spoiler.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<tg-spoiler>"+x[1]+"</tg-spoiler>", -1)

	}
	for _, x := range Code.FindAllStringSubmatch(text, -1) {
		text = strings.Replace(text, x[0], "<code>"+x[1]+"</code>", -1)
	}
	return text

}

func ParseString(t string, c tb.Context) (string, bool) {
	var Fillings = map[string]int{"{first}": 1, "{last}": 2, "{username}": 3, "{fullname}": 5, "{id}": 4, "{chatname}": 6, "{mention}": 7}
	q, preview := 0, true
	if strings.Contains(t, "{preview}") {
		preview = false
		t = strings.ReplaceAll(t, "{preview}", "")
	}
	for a, b := range Fillings {
		if strings.Contains(t, a) {
			q++
			t = strings.ReplaceAll(t, a, "%["+strconv.Itoa(b)+"]s")

		}
	}
	if strings.Contains(t, "{rules}") {
		t = strings.ReplaceAll(t, "{rules}", "[Rules](buttonurl://rules)")
	}
	first := c.Sender().FirstName
	last := c.Sender().LastName
	fullname := first
	if last != string("") {
		fullname += " " + last
	}
	username := c.Sender().Username
	id := strconv.Itoa(int(c.Sender().ID))
	mention := fmt.Sprintf("<a href='tg://user?id=%s'>%s</a>", id, first)
	if username == string("") {
		username = mention
	}
	chatname := c.Chat().Title
	if q != 0 {
		t = fmt.Sprintf(t, first, last, username, id, fullname, chatname, mention)
	}
	return t, preview

}

func ParseFile() {
}

type User struct {
	ID       int64
	Username string
	First    string
	Last     string
	Full     string
	DC       int64
	Mention  string
	Error    string
	Giga     bool
	Type     string
}

func GetObj(c tb.Context) (interface{}, string, error) {
	if c.Message().IsReply() {
		user := c.Message().ReplyTo.Sender
		if user.ID == int64(136817688) {
			user := c.Message().ReplyTo.SenderChat
			return *user, c.Message().Payload, nil
		}
		return *user, c.Message().Payload, nil
	} else if c.Message().Payload != string("") {
		Args := strings.SplitN(c.Message().Payload, " ", 1)
		if isInt(Args[0]) {
			id, _ := strconv.ParseInt(Args[0], 10, 64)
			user, err := c.Bot().ChatByID(id)
			if err != nil {
				return nil, "", err
			}
			if len(Args) > 1 {
				return *user, Args[1], err
			} else {
				return *user, "", err
			}

		} else {
			if len(Args) > 1 {
				return Args[0], Args[1], nil
			} else {
				return Args[0], "", nil
			}
		}
	} else {
		return nil, "", fmt.Errorf("you dont seem to be referring to a user or the ID specified is incorrect")
	}
}

func GetMention(id int64, name string) string {
	return fmt.Sprintf("<a href='tg://user?id=%d'>%s</a>", id, EscapeHTML(name))
}

func GetUser(c tb.Context) (User, string) {
	Obj, Payload, err := GetObj(c)
	if err != nil {
		c.Reply(err.Error())
		return User{}, ""
	}
	var user User
	switch Obj := Obj.(type) {
	case tb.User:
		user = User{
			ID:       Obj.ID,
			Username: "@" + Obj.Username,
			First:    EscapeHTML(Obj.FirstName),
			Last:     EscapeHTML(Obj.LastName),
			DC:       0,
			Mention:  GetMention(Obj.ID, Obj.FirstName),
			Giga:     false,
			Type:     "user",
		}
	case tb.Chat:
		if Obj.Title != string("") {
			user = User{
				ID:       Obj.ID,
				Username: "@" + Obj.Username,
				First:    EscapeHTML(Obj.Title),
				DC:       0,
				Mention:  EscapeHTML(Obj.Title),
				Giga:     false,
				Type:     "chat",
			}
		} else {
			user = User{
				ID:       Obj.ID,
				Username: "@" + Obj.Username,
				First:    EscapeHTML(Obj.FirstName),
				Last:     EscapeHTML(Obj.LastName),
				DC:       0,
				Mention:  GetMention(Obj.ID, Obj.FirstName),
				Giga:     false,
				Type:     "user",
			}
		}
	case string:
		user = ResolveUsername(Obj)
	}
	if user.Error != string("") {
		c.Reply(user.Error)
		return User{}, ""
	}
	return user, Payload

}

func ResolveUsername(u string) User {
	resp, err := Client.Get(ResolveURL + u)

	if err != nil {
		log.Println(err)
		return User{Error: "ResolveUsernameRequestError"}
	}
	defer resp.Body.Close()
	var data mapType
	json.NewDecoder(resp.Body).Decode(&data)
	if err, ok := data["error"]; ok {
		return User{Error: err.(string)}
	}
	var user User
	if _type, ok := data["type"]; ok && _type == "user" {
		user.Type = "user"
		if id, ok := data["id"]; ok {
			user.ID = int64(id.(float64))
		}
		if username, ok := data["username"]; ok {
			user.Username = "@" + username.(string)
		}
		if first, ok := data["first_name"]; ok {
			user.First = EscapeHTML(first.(string))
			user.Mention = GetMention(int64(data["id"].(float64)), first.(string))
		}
		if last, ok := data["last_name"]; ok {
			user.Last = EscapeHTML(last.(string))
		}
		if user.Last != "" {
			user.Full = user.First + " " + user.Last
		} else {
			user.Full = user.First
		}
		if dc, ok := data["dc_id"]; ok {
			user.DC = int64(dc.(float64))
		}
	} else if chat, ok := data["type"]; ok && chat == "channel" {
		user.Type = "channel"
		if id, ok := data["id"]; ok {
			user.ID = int64(id.(float64))
		}
		if username, ok := data["username"]; ok {
			user.Username = "@" + username.(string)
		}
		if first, ok := data["title"]; ok {
			user.First = EscapeHTML(first.(string))
			user.Mention = EscapeHTML(first.(string))
		}
		if giga, ok := data["gigagroup"]; ok {
			user.Giga = giga.(bool)
		}
		if dc, ok := data["dc_id"]; ok {
			user.DC = int64(dc.(float64))
		}
	}
	return user

}

func (user *User) User() *tb.User {
	return &tb.User{
		ID:           user.ID,
		FirstName:    user.First,
		LastName:     user.Last,
		Username:     user.Username,
		LanguageCode: "en",
		IsBot:        false,
	}
}

func (user *User) Chat() *tb.Chat {
	return &tb.Chat{
		ID:       user.ID,
		Title:    user.First,
		Username: user.Username,
	}
}

func (user *User) Approved(chatID int64) bool {
	if db.IsApproved(chatID, user.ID) {
		return true
	}
	if user.ID == OWNER_ID || user.ID == BOT_ID {
		return true
	}
	return false
}

func GetReason(r string) string {
	if r != string("") {
		return ", " + "\n<b>Reason:</b> " + r
	}
	return "."
}

func EscapeHTML(s string) string {
	for x, y := range map[string]string{"<": "&lt;", ">": "&gt;", "&": "&amp;"} {
		s = strings.ReplaceAll(s, x, y)
	}
	return s
}

func GetForwardID(c tb.Context) (int64, string, string) {
	Message := c.Message().ReplyTo
	var ID int64
	var FirstName string
	var Type = "user"
	if Message.OriginalSender != nil {
		if Message.OriginalSender.ID != 0 {
			ID = Message.OriginalSender.ID
		}
		if Message.OriginalSender.FirstName != string("") {
			FirstName = Message.OriginalSender.FirstName
		}
	} else if Message.OriginalChat != nil {
		if Message.OriginalChat.ID != 0 {
			ID = Message.OriginalChat.ID
		}
		if Message.OriginalChat.Title != string("") {
			FirstName = Message.OriginalChat.Title
		}
		Type = "chat"
	} else if Message.OriginalSignature != string("") {
		FirstName = Message.OriginalSignature
	} else if Message.OriginalSenderName != string("") {
		FirstName = Message.OriginalSenderName
	}
	return ID, FirstName, Type
}

func GetArgs(c tb.Context) string {
	if c.Message().IsReply() {
		if c.Message().ReplyTo.Text != "" {
			return c.Message().ReplyTo.Text
		} else {
			return c.Message().ReplyTo.Caption
		}
	} else {
		Args := strings.SplitN(c.Message().Text, " ", 2)
		if len(Args) > 1 {
			return Args[1]
		}
	}
	return ""
}

func ParseCountry(s string) string {
	for x, y := range COUNTRY_CODES {
		if strings.EqualFold(s, x) {
			return strings.ToUpper(y.(string))
		}
		if strings.EqualFold(s, y.(string)) {
			return strings.ToUpper(y.(string))
		}
	}
	return "US"
}
