package sendmail

import (
	"crypto/tls"
	//"flag"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"os"
	"time"
)

const (
	me       = "z.malinovskiy@beer-co.com"
	hostbeer = "mail.beer-co.com"
)

// Usage shows package usage
func Usage() string {
	ret := `
	msgsubj := flag.String("msgsubj", "", "message subject is a subject")
	msgbody := flag.String("msgbody", "", "message body is a msg body")
	senderalias := flag.String("senderalias", "", "sender alias is a From: header")

	flag.Parse()
	if len(flag.Args()) != 0 {
		fmt.Fprintf(os.Stderr, "%s", "no arguments\n")
		flag.PrintDefaults()
		return
	}
	if *msgbody == "" || *msgsubj == "" {
		fmt.Fprintf(os.Stderr, "%s", "switches required!\n")
		flag.PrintDefaults()
		return
	}

	newconn, conn, err := GetTLSConnection()
	if err != nil {
		log.Fatalf("Can't connect to mail server\n%s", err)
	}
	defer newconn.Close()
	defer conn.Close()
	c := Authenticate(newconn, "passhere")
	SendMailToMe(c, *msgsubj, *msgbody, *senderalias)
	`
	return ret
}

// GetTLSConnection gets connection objects
// needed defer newconn.Close()
// needed defer conn.Close()
func GetTLSConnection() (*tls.Conn, net.Conn, error) { //for pooling connections for several sending
	// conn, newconn escapes

	var newconn *tls.Conn

	conn, err := net.DialTimeout("tcp4", "["+hostbeer+"]:465", time.Second*20)

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return nil, nil, err
	}
	config := tls.Config{ServerName: hostbeer, InsecureSkipVerify: true} // to do not verify chain
	newconn = tls.Client(conn, &config)
	return newconn, conn, nil

}

// Authenticate passes a pass
func Authenticate(newconn *tls.Conn, pass string) *smtp.Client {
	//authentication
	c, err := smtp.NewClient(newconn, hostbeer)
	if err != nil {
		log.Fatalf("%s", err)
	}

	err = c.Hello("1cprogrammer")
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	a := smtp.PlainAuth("", me, pass, hostbeer)
	err = c.Auth(a)
	if err != nil {
		log.Fatalf("%s\n", err)
	}
	return c
}

// SendMailToMe send a msg from me
func SendMailToMe(c *smtp.Client, msgsubj, msgbody, senderalias string) {

	//actual sending
	err := c.Mail(me)
	if err != nil {
		log.Fatalf("%s", err)
	}
	err = c.Rcpt(me)
	if err != nil {
		log.Fatalf("%s", err)
	}

	wc, err := c.Data()
	if err != nil {
		log.Fatalf("%s", err)
	}
	fmt.Fprintf(wc, "Subject: %s\n", msgsubj)
	var _senderalias string
	if senderalias == "" {
		_senderalias = me

	} else {
		_senderalias = senderalias
	}
	fmt.Fprintf(wc, "From: %s\n", _senderalias)
	fmt.Fprintln(wc, "")
	fmt.Fprintf(wc, "%s", msgbody)
	wc.Close()
	err = c.Quit()
	if err != nil {
		log.Fatalf("%s", err)
	}

}

