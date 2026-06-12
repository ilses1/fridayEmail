import os
import smtplib
from datetime import datetime, timedelta, timezone
from email.header import Header
from email.mime.text import MIMEText


BEIJING_TZ = timezone(timedelta(hours=8))

MAIL_HOST = os.getenv("MAIL_HOST", os.getenv("SMTP_HOST", "smtp.qq.com"))
MAIL_PORT = int(os.getenv("MAIL_PORT", os.getenv("SMTP_PORT", "465")))
MAIL_USER = os.getenv("MAIL_USER", os.getenv("SMTP_USER"))
MAIL_PASS = os.getenv("MAIL_PASS", os.getenv("SMTP_PASSWORD"))
MAIL_RECEIVER = os.getenv("MAIL_RECEIVER", os.getenv("EMAIL_TO"))
MAIL_SENDER_NAME = os.getenv("MAIL_SENDER_NAME", "周报提醒")
MAIL_SUBJECT = os.getenv("MAIL_SUBJECT", "周报提醒：请写一下周报")
MAIL_CONTENT = os.getenv(
    "MAIL_CONTENT",
    "今天是周五，记得写一下本周周报，整理本周工作进展、问题和下周计划。",
)


def send_email(subject: str, content: str) -> bool:
    if not all([MAIL_USER, MAIL_PASS, MAIL_RECEIVER]):
        print("未配置邮箱信息，跳过邮件发送")
        return False

    try:
        msg = MIMEText(content, "plain", "utf-8")
        from_header = Header(MAIL_SENDER_NAME, "utf-8")
        from_header.append(f"<{MAIL_USER}>", "ascii")
        msg["From"] = from_header
        msg["To"] = MAIL_RECEIVER
        msg["Subject"] = Header(subject, "utf-8")

        if MAIL_PORT == 465:
            server = smtplib.SMTP_SSL(MAIL_HOST, MAIL_PORT)
        else:
            server = smtplib.SMTP(MAIL_HOST, MAIL_PORT)
            server.ehlo()
            if MAIL_PORT in (587, 25):
                server.starttls()
                server.ehlo()

        try:
            server.login(MAIL_USER, MAIL_PASS)
            server.sendmail(MAIL_USER, [MAIL_RECEIVER], msg.as_string())
        finally:
            server.quit()

        print("✅ 邮件发送成功")
        return True
    except Exception as e:
        print(f"❌ 邮件发送失败: {e}")
        return False


if __name__ == "__main__":
    now_beijing = datetime.now(BEIJING_TZ)
    print(f"[{now_beijing.isoformat()}] Sending weekly reminder email...")
    send_email(MAIL_SUBJECT, MAIL_CONTENT)
    print("Reminder email sent successfully.")
