# Payment Service – Guidelines

## Tujuan

Payment service memiliki tanggung jawab atas:

- Eksekusi pembayaran
- Integrasi dengan payment provider
- Menyimpan log pembayaran -> ke database postgres

---

## Inbound Structure

- Menerima request melalui gRPC
- Handler mengadaptasi gRPC call ke service method

---

## Service Layer

- `service/payment.go` mengorkestrasi alur pembayaran
- Mengontrol urutan:
  - Pemanggilan external API
  - Persistence data
  - Event publication

---

## Adapters

- `adapters/external`
  - Integrasi payment provider
- `adapters/repo`
  - Penyimpanan data pembayaran

---

## Messaging

- Mem-publish event hasil pembayaran (notifikasi email) -> ke service **_Notification_** dengan menggunakan **_RabbitMQ_**
- Publish dilakukan sebagai efek samping dari transaksi pembayaran

---

## Design Intent

- Logic pembayaran fokus pada transaksi menggunakan referensi pengguna (misalnya menggunakan userId); Layanan payment tidak perlu melakukan authentikasi user ulang.
- Side effect notifikasi bersifat asynchronous
- Tidak ada koordinasi lintas service selain melalui gRPC dan RabbitMQ
