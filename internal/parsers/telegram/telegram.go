package telegram

type Telegram struct {
	name      string
	channelId int
	output    chan string
}

func NewTelegram(name string, chanelId int) *Telegram {
	return &Telegram{
		name:      name,
		channelId: chanelId,
		output:    make(chan string),
	}
}

func (chanel *Telegram) Parse() chan string {
	panic("not implemented")
}
