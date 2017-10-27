package main

import (
	"encoding/json"
	"log"
	"net/http"
	"gopkg.in/mux"
	"gopkg.in/mgo.v2"
	"time"
	"regexp"
	//"gopkg.in/mgo.v2/bson"
	"strings"
	"gopkg.in/mgo.v2/bson"
	"math/rand"
	"fmt"

	"strconv"
)

func ErrorWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("x-request-id", newUUID())
	w.Header().Set("datetime", time.Now().Format("2006-01-02 15:04:05+0700"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-roundtrip", "")
	w.Header().Set("x-job-id", "")

	//w.WriteHeader(code)
	//fmt.Fprintf(w, "{message: %q}", message)
	w.Write(json)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
	w.Header().Set("x-request-id", newUUID())
	w.Header().Set("datetime", time.Now().Format("2006-01-02 15:04:05+0700"))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("x-roundtrip", "")
	w.Header().Set("x-job-id", "")
	//w.WriteHeader(code)
	//w.Write(json)
	w.Write(json)
}

type WalletAccount struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	CitizenID     int     `json:"citizen_id" bson:"citizen_id"`
	FullName      string  `json:"full_name" bson:"full_Name"`
	WalletID      int     `json:"wallet_id" bson:"wallet_id"`
	OpenDateTime  string  `json:"open_datetime" bson:"open_datetime"`
	LedgerBalance float32 `json:"ledger_balance" bson:"ledger_balance"`
}

type MsgBody struct {
	RsBody RsBody `json:"rsBody"`
	Error ErrorList `json:"error"`
}

type RsBody struct {
	WalletID	int		`json:"wallet_id"`
	OpenDateTime  string  `json:"open_datetime"`
}

type Error struct {
	ErCode string 	`bson:"error code"`
	ErDesc string	`bson:"error description"`
}

type ErrorList struct {
	Error []Error
}

func main() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)
	ensureIndex(session)

	mux := mux.NewRouter()
	//mux.HandleFunc("/v1/accounts/{wallet_id}", getAccountByWalletID(session)).Methods("GET")
	mux.HandleFunc("/v1/accounts", createWallets(session)).Methods("POST")
	//http.ListenAndServe("localhost:5000", mux)
	log.Fatal(http.ListenAndServe("localhost:3333", mux))
}

func ensureIndex(s *mgo.Session) {
	session := s.Copy()
	defer session.Close()

	c := session.DB("wallets").C("accounts")

	index := mgo.Index{
		Key:		[]string{"citizen_id"},
		Unique:		true,
		DropDups:	true,
		Background:	true,
		Sparse:		true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func createWallets(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		var accounts WalletAccount
		var errorlst ErrorList

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&accounts)
		if err != nil {
			errorlst.Error = append(errorlst.Error,Error{"999", "Incorrect Body"})
		}
		if !LenCitizenId(accounts.CitizenID) {
			errorlst.Error = append(errorlst.Error,Error{"001", "Incorrect Citizen ID"})
		}

		if (!IsLetter(accounts.FullName)) || (!Len(accounts.FullName)) {
			errorlst.Error = append(errorlst.Error,Error{"003", "Incorrect Name"})
		}

		accounts.FullName=strings.ToUpper(accounts.FullName)
		accounts.WalletID=1234567890
		accounts.OpenDateTime = time.Now().Format("2006-01-02 15:04:05+0700")
		accounts.LedgerBalance = 0.00

		c := session.DB("wallets").C("accounts")

		err = c.Insert(accounts)
		if err != nil {
			if mgo.IsDup(err) {
				errorlst.Error = append(errorlst.Error,Error{"002", "Duplicate Citizen ID"})
			}

			//errorlst.Error = append(errorlst.Error,Error{"999", "Database Error"})
		}

		respbody := RsBody{
			OpenDateTime:accounts.OpenDateTime,
			WalletID:accounts.WalletID,
		}


		msgbody :=MsgBody{
			respbody,
			errorlst,
		}
		respBody, err := json.MarshalIndent(msgbody, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusCreated)

	}
}

/*
func getAccountByWalletID(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		session := s.Copy()
		defer session.Close()

		vars := mux.Vars(r)
		wallets := vars["wallet_id"]

		c := session.DB("wallets").C("accounts")

		var accounts WalletAccount
		err := c.Find(bson.M{"wallet_id": wallets}).One(&accounts)
		if err != nil {
			ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
			log.Println("Failed find wallet_id: ", err)
			return
		}

		if accounts.WalletID == nil {
			ErrorWithJSON(w, "Book not found", http.StatusNotFound)
			return
		}

		respBody, err := json.MarshalIndent(accounts, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		ResponseWithJSON(w, respBody, http.StatusOK)
	}
}
*/
var IsLetter = regexp.MustCompile(`^[a-zA-Z.,-]+( [a-zA-Z.,-]+)+$`).MatchString

func Len(s string) bool {
	if len(s)<=50 {
		return true
	}
	return false
}

func LenCitizenId(i int) bool {
	if len(strconv.Itoa(i))==13 {
		return true
	}
	return false
}

var walletId = randInt(0,9999999999)

func newUUID() (string) {
	uuid := make([]byte, 16)
	//n, err := io.ReadFull(rand.Reader, uuid)
	//if n != len(uuid) || err != nil {
	//	return "", err
	//}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}






