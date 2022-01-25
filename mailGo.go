package rushtool

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/smtp"
	"strconv"
	"strings"
)

func TlsMail(sub, content, mailList, serverHost string, serverPort int, userName, passWord string) error {
	//host := "smtp.exmail.qq.com"
	//port := 465
	host := serverHost
	port := serverPort
	toEmail := mailList
	header := make(map[string]string)
	header["From"] = "ai_watch" + "<" + userName + ">"
	header["To"] = mailList
	header["Subject"] = sub
	header["Content-Type"] = "text/html; charset=UTF-8"
	body := content
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		userName,
		passWord,
		host,
	)
	err := sendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		userName,
		strings.Split(toEmail, ";"),
		[]byte(message),
	)
	if err != nil {
		return err
	}
	return nil
}

//return a smtp client
func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		LogConsole().Error("Dialing Error:", err)
		return nil, err
	}
	//分解主机端口字符串
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}
func someError(addr, err string) error {
	someA := errors.New(err + " ," + addr + " send failed")
	return someA
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时,smtp.NewClient()会卡住且不提示err
//len(to)>1时,to[1]开始提示是密送
func sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {
	//create smtp client
	c, err := dial(addr)
	if err != nil {
		LogConsole().Error("Create smpt client error:", err)
		return err
	}
	defer c.Close()
	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				//panic(err)
				log.Println("Error during AUTH", err)
				return err
			}
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			someP := someError(addr, err.Error())
			return someP
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
func MailText(sub, content, mailList, serverHost string, serverPort int, userName, passWord string) error {
	//host := "smtp.exmail.qq.com"
	//port := 465
	host := serverHost
	port := serverPort
	toEmail := mailList
	header := make(map[string]string)
	header["From"] = "ai_watch" + "<" + userName + ">"
	header["To"] = mailList
	header["Subject"] = sub
	header["Content-Type"] = "text/plain; charset=UTF-8"
	body := content
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body
	auth := smtp.PlainAuth(
		"",
		userName,
		passWord,
		host,
	)
	err := sendMailUsingTLS(
		fmt.Sprintf("%s:%d", host, port),
		auth,
		userName,
		strings.Split(toEmail, ";"),
		[]byte(message),
	)
	if err != nil {
		return err
	}
	return nil
}
func CueMail(sub, content, mailTo, serverAddr, user, password string, tls bool) error {
	host, port, _ := net.SplitHostPort(serverAddr)
	nPort, _ := strconv.Atoi(port)
	if tls {
		err := TlsMail(sub, content, mailTo, host, nPort, user, password)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return err
	} else if !tls {
		auth := myLoginAuth(user, password)
		to := strings.Split(mailTo, ";")
		msg := []byte(fmt.Sprintf("To: %s\r\n"+
			"From: %s\r\n"+
			"Subject: %s\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n", mailTo, user, sub, content))
		err := smtp.SendMail(serverAddr, auth, user, to, msg)
		if err != nil {
			//fmt.Println(err)
			return err
		}
		return err

	}
	return nil
}

type loginAuth struct {
	username, password string
}

func myLoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}
func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// return "LOGIN", []byte{}, nil
	return "LOGIN", []byte(a.username), nil
}
func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		}
	}
	return nil, nil
}
