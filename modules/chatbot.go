package modules

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	tb "gopkg.in/tucnak/telebot.v3"
)

func Chat_bot(c tb.Context) error {
	is_chat := false
	if c.Message().IsReply() && c.Message().ReplyTo.Sender.ID == int64(5050904599) {
		is_chat = true
	} else if strings.Contains(c.Message().Text, "yoko") {
		is_chat = true
	}
	if !is_chat {
		return nil
	}
	text := strings.ReplaceAll(c.Message().Text, "yoko", "kuki")
	url_q := "https://icap.iconiq.ai/talk?&botkey=icH-VVd4uNBhjUid30-xM9QhnvAaVS3wVKA3L8w2mmspQ-hoUB3ZK153sEG3MX-Z8bKchASVLAo~&channel=7&sessionid=482070240&client_name=uuiprod-un18e6d73c-user-19433&id=true"
	req, err := http.PostForm(url_q, url.Values{"input": {text}})
	if err != nil {
		c.Reply(err.Error())
	}
	defer req.Body.Close()
	var resp mapType
	json.NewDecoder(req.Body).Decode(&resp)
	msg := resp["responses"].([]interface{})[0].(string)
	pattern := regexp.MustCompile(`<image>.+</image>`)
	media := pattern.FindAllStringSubmatch(msg, -1)
	if media != nil {
		if len(media) != 0 {
			file := strings.ReplaceAll(strings.ReplaceAll(media[0][0], "<image>", ""), "</image>", "")
			c.Reply(&tb.Animation{File: tb.FromURL(file)})
		}
	}
	chat := strings.SplitN(msg, "</image>", 2)
	var message string
	if len(chat) == 2 {
		message = chat[1]
	} else {
		message = chat[0]
	}
	c.Reply(message)
	return nil
}