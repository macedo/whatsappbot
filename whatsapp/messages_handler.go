package whatsapp

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/goodsign/monday"
	"github.com/macedo/whatsappbot/scheduler"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules"
	"github.com/olebedev/when/rules/br"
	"github.com/olebedev/when/rules/common"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func MessagesHandler(c *Client) whatsmeow.EventHandler {
	parser := when.New(&rules.Options{
		Distance: 10,
	})
	parser.Add(common.All...)
	parser.Add(br.All...)

	return func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			re := regexp.MustCompile(`^remember me|me lembra de ?(.*)?$`)
			matches := re.FindStringSubmatch(v.Message.GetConversation())
			if len(matches) == 0 {
				l.Printf("unhandled message - %s", v.Message.GetConversation())
				return
			}

			body := matches[1]

			l.Println(body)

			r, err := parser.Parse(body, time.Now())
			if err != nil {
				l.Printf("can't parser conversation - %s", err)
				return
			}

			l.Println(r)

			if r == nil {
				c.SendMessage(context.Background(), v.Info.Sender.ToNonAD(), &waProto.Message{
					Conversation: proto.String("Desculpe, não entendi direito. Eu sou facilmente confundido. Talvez tente as palavras em uma ordem diferente. Isso geralmente funciona: me lembra de [o que] [quando]"),
				})
			} else {
				reminder := strings.ReplaceAll(body, body[r.Index:r.Index+len(r.Text)], "")
				reminder = strings.TrimSpace(reminder)

				scheduler.ScheduleWithTime(func(ctx context.Context) {
					c.SendMessage(context.Background(), v.Info.Sender.ToNonAD(), &waProto.Message{
						Conversation: proto.String(reminder),
					})
				}, r.Time)

				if err == nil {
					l.Println(v.Info.Sender)
					c.SendMessage(context.Background(), v.Info.Sender.ToNonAD(), &waProto.Message{
						Conversation: proto.String(fmt.Sprintf("Vou te lembrar de %s as %s", reminder, monday.Format(r.Time, "15:04 de Monday dia 02 de January", monday.LocalePtBR))),
					})

					l.Printf("%s - %s", reminder, r.Time)
				}
			}
		}
	}
}
