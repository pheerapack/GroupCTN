package Account

import "gopkg.in/mgo.v2/bson"
//Account represents a Input Account
type Account struct {
	ID     		bson.ObjectId `bson:"_id"`
	FullName  	string        `json:"FullName"`
}
//Albums is an array of Album
type Accounts []Account
