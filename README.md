# Friday Email Reminder

这是一个使用 Rust + GitHub Actions 的周报提醒项目（保留 Go / Python 版本供本地对照）。

## 功能

- 每周五北京时间 09:40 自动运行
- 通过 SMTP 发送邮件提醒你写周报
- 支持手动触发 `workflow_dispatch`

## 文件说明

- `src/main.rs`：发送提醒邮件的 Rust 脚本（GitHub Actions 使用）
- `Cargo.toml`：Rust 项目定义
- `cmd/send_reminder/main.go`：Go 版本（保留，可本地对照）
- `send_reminder.py`：Python 版本（保留，可本地对照）
- `.github/workflows/weekly-reminder.yml`：GitHub Actions 定时任务

## GitHub Secrets 配置

在仓库的 `Settings -> Secrets and variables -> Actions` 中添加：

- `MAIL_HOST`：SMTP 服务器地址，例如 `smtp.qq.com`
- `MAIL_PORT`：SMTP 端口，例如 `465` 或 `587`
- `MAIL_USER`：发件邮箱账号
- `MAIL_PASS`：SMTP 授权码或密码
- `MAIL_RECEIVER`：收件邮箱地址

兼容旧变量名：`SMTP_HOST`、`SMTP_PORT`、`SMTP_USER`、`SMTP_PASSWORD`、`EMAIL_TO`

## 定时说明

GitHub Actions 使用 UTC 时间。
北京时间 09:40 = UTC 01:40，因此 cron 写成：

```yaml
40 1 * * 5
```

其中 `5` 表示周五。

## 本地测试

先设置环境变量，再运行 Rust 版本：

```bash
cargo run --release
```

例如在 PowerShell 中：

```powershell
$env:SMTP_HOST="smtp.qq.com"
$env:SMTP_PORT="587"
$env:SMTP_USER="your@email.com"
$env:SMTP_PASSWORD="your_smtp_password"
$env:EMAIL_TO="target@email.com"
cargo run --release
```

也可使用 Go / Python 版本对照测试：

```powershell
go run ./cmd/send_reminder
python .\send_reminder.py
```

## 可自定义内容

你可以通过环境变量修改邮件标题和正文：

- `MAIL_SUBJECT`
- `MAIL_CONTENT`
- `MAIL_SENDER_NAME`
