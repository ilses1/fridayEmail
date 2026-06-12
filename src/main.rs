use chrono::{FixedOffset, Utc};
use lettre::message::header::ContentType;
use lettre::message::Mailbox;
use lettre::transport::smtp::authentication::Credentials;
use lettre::transport::smtp::client::{Tls, TlsParameters};
use lettre::{Message, SmtpTransport, Transport};
use std::env;

fn env_or(primary: &str, fallback: Option<&str>, default: &str) -> String {
    if let Ok(v) = env::var(primary) {
        return v;
    }
    if let Some(fb) = fallback {
        if let Ok(v) = env::var(fb) {
            return v;
        }
    }
    default.to_string()
}

fn env_port_or(primary: &str, fallback: Option<&str>, default: u16) -> u16 {
    let s = env_or(primary, fallback, &default.to_string());
    s.parse().unwrap_or(default)
}

struct MailConfig {
    host: String,
    port: u16,
    user: String,
    pass: String,
    receiver: String,
    sender_name: String,
    subject: String,
    content: String,
}

impl MailConfig {
    fn from_env() -> Self {
        Self {
            host: env_or("MAIL_HOST", Some("SMTP_HOST"), "smtp.qq.com"),
            port: env_port_or("MAIL_PORT", Some("SMTP_PORT"), 465),
            user: env_or("MAIL_USER", Some("SMTP_USER"), ""),
            pass: env_or("MAIL_PASS", Some("SMTP_PASSWORD"), ""),
            receiver: env_or("MAIL_RECEIVER", Some("EMAIL_TO"), ""),
            sender_name: env_or("MAIL_SENDER_NAME", None, "周报提醒"),
            subject: env_or("MAIL_SUBJECT", None, "周报提醒：请写一下周报"),
            content: env_or(
                "MAIL_CONTENT",
                None,
                "今天是周五，记得写一下本周周报，整理本周工作进展、问题和下周计划。",
            ),
        }
    }
}

fn build_transport(host: &str, port: u16, user: &str, pass: &str) -> Result<SmtpTransport, lettre::transport::smtp::Error> {
    let creds = Credentials::new(user.to_string(), pass.to_string());

    if port == 465 {
        let tls = TlsParameters::new(host.to_string())?;
        return Ok(
            SmtpTransport::relay(host)?
                .port(port)
                .tls(Tls::Wrapper(tls))
                .credentials(creds)
                .build(),
        );
    }

    let mut builder = SmtpTransport::relay(host)?.port(port).credentials(creds);

    if port == 587 || port == 25 {
        let tls = TlsParameters::new(host.to_string())?;
        builder = builder.tls(Tls::Required(tls));
    }

    Ok(builder.build())
}

fn send_email(config: &MailConfig) -> bool {
    if config.user.is_empty() || config.pass.is_empty() || config.receiver.is_empty() {
        println!("未配置邮箱信息，跳过邮件发送");
        return false;
    }

    let from_addr = match config.user.parse() {
        Ok(addr) => addr,
        Err(e) => {
            println!("❌ 邮件发送失败: {e}");
            return false;
        }
    };
    let to_addr = match config.receiver.parse() {
        Ok(addr) => addr,
        Err(e) => {
            println!("❌ 邮件发送失败: {e}");
            return false;
        }
    };

    let from = Mailbox::new(Some(config.sender_name.clone()), from_addr);

    let email = match Message::builder()
        .from(from)
        .to(to_addr)
        .subject(&config.subject)
        .header(ContentType::TEXT_PLAIN)
        .body(config.content.clone())
    {
        Ok(msg) => msg,
        Err(e) => {
            println!("❌ 邮件发送失败: {e}");
            return false;
        }
    };

    let mailer = match build_transport(&config.host, config.port, &config.user, &config.pass) {
        Ok(m) => m,
        Err(e) => {
            println!("❌ 邮件发送失败: {e}");
            return false;
        }
    };

    match mailer.send(&email) {
        Ok(_) => {
            println!("✅ 邮件发送成功");
            true
        }
        Err(e) => {
            println!("❌ 邮件发送失败: {e}");
            false
        }
    }
}

fn main() {
    let beijing_tz = FixedOffset::east_opt(8 * 3600).expect("invalid timezone offset");
    let now_beijing = Utc::now().with_timezone(&beijing_tz);
    println!(
        "[{}] Sending weekly reminder email...",
        now_beijing.format("%Y-%m-%dT%H:%M:%S%:z")
    );

    let config = MailConfig::from_env();
    send_email(&config);

    println!("Reminder email sent successfully.");
}
