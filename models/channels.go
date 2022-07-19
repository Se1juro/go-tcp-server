package models

type ChannelClients struct {
	Conn           string `json:"conn"`
	CurrentChannel int    `json:"currentChannel"`
	Status         string `json:"status"`
}

type Channels struct {
	ChannelId int              `json:"channelId"`
	Clients   []ChannelClients `json:"clients"`
}
