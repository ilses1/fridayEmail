package main

import (
	"crypto/tls"
	"fmt"
	"mime"
	"net/smtp"
	"os"
	"strconv"
	"time"
)

var (
	beijingTZ = time.FixedZone("CST", 8*60*60)

	mailHost       = envOr("MAIL_HOST", "SMTP_HOST", "smtp.qq.com")
	mailPort       = envPortOr("MAIL_PORT", "SMTP_PORT", 465)
	mailUser       = envOr("MAIL_USER", "SMTP_USER", "")
	mailPass       = envOr("MAIL_PASS", "SMTP_PASSWORD", "")
	mailReceiver   = envOr("MAIL_RECEIVER", "EMAIL_TO", "")
	mailSenderName = envOr("MAIL_SENDER_NAME", "", "周报提醒")
	mailSubject    = envOr("MAIL_SUBJECT", "", "周报提醒：请写一下周报")
	mailContent    = envOr("MAIL_CONTENT", "", "今天是周五，记得写一下本周周报，整理本周工作进展、问题和下周计划。")
)

func envOr(primary, fallback, def string) string {
	if v, ok := os.LookupEnv(primary); ok {
		return v
	}
	if fallback != "" {
		if v, ok := os.LookupEnv(fallback); ok {
			return v
		}
	}
	return def
}

func envPortOr(primary, fallback string, def int) int {
	s := envOr(primary, fallback, strconv.Itoa(def))
	port, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return port
}

func encodeHeader(s string) string {
	return mime.QEncoding.Encode("utf-8", s)
}

func buildMessage(fromName, fromAddr, to, subject, body string) []byte {
	from := encodeHeader(fromName) + " <" + fromAddr + ">"
	encodedSubject := encodeHeader(subject)
	return []byte(fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\nContent-Transfer-Encoding: 8bit\r\n\r\n%s",
		from, to, encodedSubject, body,
	))
}

func dialSMTP(host string, port int) (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%d", host, port)

	if port == 465 {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
		if err != nil {
			return nil, err
		}
		return smtp.NewClient(conn, host)
	}

	client, err := smtp.Dial(addr)
	if err != nil {
		return nil, err
	}

	if port == 587 || port == 25 {
		if err := client.StartTLS(&tls.Config{ServerName: host}); err != nil {
			client.Close()
			return nil, err
		}
	}

	return client, nil
}

func sendEmail(subject, content string) bool {
	if mailUser == "" || mailPass == "" || mailReceiver == "" {
		fmt.Println("未配置邮箱信息，跳过邮件发送")
		return false
	}

	client, err := dialSMTP(mailHost, mailPort)
	if err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}
	defer client.Close()

	auth := smtp.PlainAuth("", mailUser, mailPass, mailHost)
	if err := client.Auth(auth); err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}

	if err := client.Mail(mailUser); err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}

	if err := client.Rcpt(mailReceiver); err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}

	w, err := client.Data()
	if err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}

	msg := buildMessage(mailSenderName, mailUser, mailReceiver, subject, content)
	if _, err := w.Write(msg); err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}

	if err := w.Close(); err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}

	if err := client.Quit(); err != nil {
		fmt.Printf("❌ 邮件发送失败: %v\n", err)
		return false
	}

	fmt.Println("✅ 邮件发送成功")
	return true
}

func main() {
	nowBeijing := time.Now().In(beijingTZ)
	fmt.Printf("[%s] Sending weekly reminder email...\n", nowBeijing.Format("2006-01-02T15:04:05-07:00"))
	sendEmail(mailSubject, mailContent)
	fmt.Println("Reminder email sent successfully.")
}
