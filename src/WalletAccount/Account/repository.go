package Account

import (
	"fmt"
	"log"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)
//Repository ...
type Repository struct{}
// SERVER the DB server
const SERVER = "localhost:27017"
// DBNAME the name of the DB instance
const DBNAME = "WalletAccount"
// DOCNAME the name of the document
const DOCNAME = "account"

// GetAccounts returns the list of Accounts

func (r Repository) GetAccounts() Accounts {
	session, err := mgo.Dial(SERVER)
	if err != nil {
		fmt.Println("Failed to establish connection to Mongo server:", err)
	}
	defer session.Close()
	c := session.DB(DBNAME).C(DOCNAME)
	results := Accounts{}
	if err := c.Find(nil).All(&results); err != nil {
		fmt.Println("Failed to write results:", err)
	}
	return results
}
// AddAccount inserts an Accounts in the DB
func (r Repository) AddAccount(account Account) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	account.ID = bson.NewObjectId()
	session.DB(DBNAME).C(DOCNAME).Insert(account)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
// UpdateAccount updates an Account in the DB (not used for now)
func (r Repository) UpdateAccount(account Account) bool {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	account.ID = bson.NewObjectId()
	session.DB(DBNAME).C(DOCNAME).UpdateId(account.ID, account)
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}
// DeleteAccount deletes an Account (not used for now)
func (r Repository) DeleteAccount(id string) string {
	session, err := mgo.Dial(SERVER)
	defer session.Close()
	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(id) {
		return "NOT FOUND"
	}
	// Grab id
	oid := bson.ObjectIdHex(id)
	// Remove user
	if err = session.DB(DBNAME).C(DOCNAME).RemoveId(oid); err != nil {
		log.Fatal(err)
		return "INTERNAL ERR"
	}
	// Write status
	return "OK"
}
