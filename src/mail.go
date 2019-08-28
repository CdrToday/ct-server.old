package main

import (
	"fmt"
	"net/smtp"
)

func mailAuth() smtp.Auth {
	auth := conf().Mail.Auth

	return smtp.PlainAuth(
		auth.Ident,
		auth.User,
		auth.Pass,
		auth.Host,
	)
}

/// # Generate mail content
/// ```
/// From: cdr.today 👻 <cdr.today@foxmail.com>
/// To: reciver@example.com
/// Subject: Hello ✔︎
///
/// 8444aa5b-76ce-48de-b9d9-5108b5b39a13
/// ```
func _genMsg(to string, body string) []byte {
	msg := conf().Mail.Msg
	fmt.Println(msg)
	plain := "" +
		"From: " + msg.From + "\n" +
		"To: " + to + "\n" +
		"Subject: " + msg.Subject + " \n\n" +
		body

	return []byte(plain)
}

/// send uuid to author
func sendMail(to string, uuid string) bool {
	auth := mailAuth()

	err := smtp.SendMail(
		"smtp.qq.com:25", auth, "cdr.today@foxmail.com",
		[]string{to}, _genMsg(to, uuid),
	)

	if err != nil {
		return false
	} else {
		return true
	}
}

func main() {
	sendMail("udtrokia@163.com", "hello")
}
