package model

import (
	"gopkg.in/mgo.v2/bson"
)

//Device - To represent the device collection
type Device struct {
	ID            bson.ObjectId `json:"id" bson:"_id"`
	DeviceUpdated string        `json:"devieUpdated" bson:"deviceUpdated"`
	Name          string        `json:"name" bson:"name" `
	Type          string        `json:"type" bson:"type"`
	Vers          string        `json:"vers" bson:"vers"`
	Group         string        `json:"group" bson:"group"`
	Net           struct {
		Intf []struct {
			Name string `json:"name,omitempty"`
			Mask string `json:"mask,omitempty"`
			Ipv4 string `json:"ipv4,omitempty"`
			Addr string `json:"addr,omitempty"`
		} `json:"intf" bson:"intf"`
	} `json:"net" bson:"net"`
	TicketID string `json:"ticketID" bson:"ticketID,omitempty"`
	//TicketCreated string `json:"ticketCreated" bson:"ticketCreated"`
}
