package source

import "time"

type TelegramChannel struct {
	Name          string
	ChannelID     int
	LastVisitedAt time.Time
	output        chan any
}

func NewTelegramChanel(name string, chanelId int) *TelegramChannel {
	return &TelegramChannel{
		Name:      name,
		ChannelID: chanelId,
		output:    make(chan any),
	}
}

func (chanel *TelegramChannel) Parse() chan string {
	panic("not implemented")
}
