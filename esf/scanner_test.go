package main

import (
	"bytes"
	"log"
	"os"
)

var input = bytes.NewBufferString(`
# the MetaEventInfo is a special single level of inheritance of events.  Any
# fields in this Event are allowed in all other events.  Certain systems
# will set these fields, so the SenderIP, SenderPort, and ReceiptTime are
# set by all journallers and listeners, while the encoding is set by all
# emitters.
MetaEventInfo
{
  ip_addr SenderIP;    # IP address of Sender
  uint16  SenderPort;  # IP port of Sender
  int64   ReceiptTime; # time this event was received, in
                       # milliseconds since epoch
  int16   enc;         # encoding of strings in the event
  uint16  SiteID;      # id of site sending the event
}

UserLogin
{
  string  username;    # username of user
  uint64  password;    # unique hash of the users password
  ip_addr clientIP;    # client ip the user attempted to connect from
}
`)

func ExampleScanner() {
	scanner := NewScanner(input)
	esf, err := scanner.Scan()
	if err != nil {
		log.Println(err)
	}
	os.Stdout.Write(esf)

	// Output:
	// // fields in this Event are allowed in all other events.  Certain systems
	// // will set these fields, so the SenderIP, SenderPort, and ReceiptTime are
	// // set by all journallers and listeners, while the encoding is set by all
	// // emitters.
	// type MetaEventInfo struct {
	// 	// IP address of Sender
	// 	SenderIP net.IP
	// 	// IP port of Sender
	// 	SenderPort uint16
	// 	// time this event was received, in
	// 	// milliseconds since epoch
	// 	ReceiptTime int64
	// 	// encoding of strings in the event
	// 	enc int16
	// }
	//
	// type UserLogin struct {
	// 	MetaEventInfo
	// 	// username of user
	// 	username string
	// 	// unique hash of the users password
	// 	password uint64
	// }
	//
}
