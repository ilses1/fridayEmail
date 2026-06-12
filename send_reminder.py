import os
import smtplib
from email.mime.text import MIMEText
from email.utils import formataddr
from datetime import datetime, timezone, timedelta


BEIJING_TZ = timezone(timedelta(hours=8))


def build_message() -> MIMEText:
    subject = os.getenv("EMAIL_SUBJECT", "周报提醒：请写一下周报")
    body = os.getenv(
        "EMAIL_BODY",
        "今天是周五，记得写一下本周周报，整理本周工作进展、问题和下周计划。",
    )

    msg = MIMEText(body, "plain", "utf-8")
    msg["Subject"] = subject
    msg["From"] = formataddr((os.getenv("SENDER_NAME", "周报提醒"), os.environ["SMTP_USER"]))
    msg["To"] = os.environ["EMAIL_TO"]
    return msg


def send_email() -> None:
    smtp_host = os.environ["SMTP_HOST"]
    smtp_port = int(os.getenv("SMTP_PORT", "587"))
    smtp_user = os.environ["SMTP_USER"]
    smtp_password = os.environ["SMTP_PASSWORD"]
    email_to = os.environ["EMAIL_TO"]

    msg = build_message()

    with smtplib.SMTP(smtp_host, smtp_port) as server:
        server.ehlo()
        if smtp_port in (587, 25):
            server.starttls()
            server.ehlo()
        server.login(smtp_user, smtp_password)
        server.sendmail(smtp_user, [email_to], msg.as_string())


if __name__ == "__main__":
    now_beijing = datetime.now(BEIJING_TZ)
    print(f"[{now_beijing.isoformat()}] Sending weekly reminder email...")
    send_email()
    print("Reminder email sent successfully.")
