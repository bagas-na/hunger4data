# Authentication Service – Guidelines

## Tujuan

Authentication service bertanggung jawab untuk:

- Registrasi user
- Login
- Password hashing
- Pembuatan JWT (sedangkan **_pengecekan_** JWT dilakukan di **api-gateway**)

Otoritas untuk melakukan autentikasi (selain api-gateway dengan JWT) ada pada service ini.

## Inbound Structure

- Menerima request melalui gRPC saja
- gRPC handler berada di `transport/grpc`

## Service Layer

- `service/auth.go` berisi workflow autentikasi
- Domain logic (registrasi, login, penerbitan token) berada di sini

## Adapter

- `adapters/crypto`
  - Password hashing
  - JWT signing
- `adapters/repo`
  - Persistence data user

Adapter sebaiknya tidak berisi business logic. Business logic sebaiknya ditempatkan di service-layer

## Messaging

- Service boleh mem-publish event (misalnya email confirmation) -> ke service **_Notification_** dengan menggunakan **_RabbitMQ_**
- Event publishing merupakan efek-samping dari operasi yang sukses

## Design Intent

- Validasi JWT dilakukan di upstream (API Gateway)
- Pengiriman email ditangani oleh Notification service
- Tidak ada dependency langsung ke domain service lain (selain gRPC + protobuf, dan RabbitMQ)
