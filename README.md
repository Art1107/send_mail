# SendGrid Webhook Event Handler

## 📝 รายละเอียดโปรเจค
ระบบรับและจัดการ webhook events จาก SendGrid พร้อมทั้งส่งการแจ้งเตือนผ่าน Lark เมื่อมีเหตุการณ์สำคัญ

### 🌟 คุณสมบัติหลัก
- รองรับการรับ webhook events จาก SendGrid
- ตรวจสอบความถูกต้องของลายเซ็น (signature verification)
- บันทึก logs ในรูปแบบไฟล์และ CSV
- ส่งการแจ้งเตือนผ่าน Lark webhook สำหรับเหตุการณ์สำคัญ
- รองรับเหตุการณ์หลัก: delivered, open, click, bounce, spam_report

### 🛠️ การติดตั้ง
1. Clone repository:
```bash
git clone <repository-url>
```

2. ติดตั้ง dependencies:
```bash
go mod download
```

3. สร้างไฟล์ .env และกำหนดค่าต่างๆ:
```env
SENDGRID_PUBLIC_KEY="your-public-key"
SERVER_PORT=":8080"
LARK_WEBHOOK_URL="your-lark-webhook-url"
LOG_FILE="sendgrid_events.log"
```

### 🚀 การใช้งาน
1. รันเซิร์ฟเวอร์:
```bash
go run cmd/main.go
```

2. เซิร์ฟเวอร์จะเริ่มทำงานที่พอร์ต 8080 (หรือตามที่กำหนดใน .env)

3. Endpoints:
- `/webhook` - สำหรับรับ SendGrid events
- `/test` - สำหรับทดสอบการทำงานของระบบ

### 📊 โครงสร้างโปรเจค
```
send_mail/
├── cmd/
│   └── main.go                 # จุดเริ่มต้นโปรแกรม
├── internal/
│   ├── adapters/
│   │   └── lark/              # การเชื่อมต่อกับ Lark
│   ├── core/                  # บิสเนสลอจิกหลัก
│   ├── domain/               # โมเดลและ errors
│   └── ports/                # interfaces
├── pkg/
│   ├── logger/               # ระบบบันทึกล็อก
│   └── verify/               # ตรวจสอบลายเซ็น
└── config/                   # การตั้งค่าระบบ
```

### 🔒 การตั้งค่าความปลอดภัย
- ระบบรองรับการตรวจสอบลายเซ็นของ SendGrid
- ต้องกำหนด SENDGRID_PUBLIC_KEY ใน .env file
- รองรับ HTTPS (ควรใช้ในการ deploy)

### 📝 การบันทึกข้อมูล
- บันทึกล็อกในรูปแบบไฟล์ปกติและ CSV
- ข้อมูลที่บันทึก: timestamp, ระดับความสำคัญ, ข้อความ, ประเภทเหตุการณ์, อีเมล, เวลาของเหตุการณ์

### 🤝 การสนับสนุน
หากพบปัญหาหรือต้องการเพิ่มฟีเจอร์ใหม่ กรุณาสร้าง issue ใน repository

### 📄 License
MIT License
