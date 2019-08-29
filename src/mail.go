package main

import (
	"net/smtp"
)

func mailAuth() smtp.Auth {
	t := conf()
	ident := t.Get("mail.auth.ident").(string)
	user := t.Get("mail.auth.user").(string)
	pass := t.Get("mail.auth.pass").(string)
	host := t.Get("mail.auth.host").(string)

	return smtp.PlainAuth(
		ident, user, pass, host,
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
	t := conf()
	cc := t.Get("mail.msg.subject").(string)
	from := t.Get("mail.msg.from").(string)

	plain := "" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		cc + " \n\n" +
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
