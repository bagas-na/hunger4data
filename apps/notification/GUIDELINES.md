# Notification Service – Guidelines

## Tujuan

Notification service menangani:

- Konsumsi message RabbitMQ secara asynchronous
- Pengiriman notifikasi (email)
- Menyimpan log notifikasi -> ke database postgres

---

## Inbound Structure

- Menerima message melalui RabbitMQ
- Consumer berada di `transport/rabbitmq`

---

## Service Layer

- `service/notification.go`:
  - Menginterpretasikan event
  - Menentukan perilaku notifikasi

---

## Adapter

- `adapters/smtp`
  - Pengiriman email
- `adapters/repository`
  - Penyimpanan log notifikasi ke database postgres

---

## Messaging

- Service ini hanya mengonsumsi event
- Tidak mem-publish message

---

## Design Intent

- Tidak ada business workflow
- Pemanggilan external API khusus untuk pengiriman pesan
- Tidak ada logic autentikasi
