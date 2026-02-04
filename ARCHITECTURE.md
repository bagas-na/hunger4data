# CODE ARCHITECTURE

## 1. Introduction

Repository adalah implementasi arsitektur microservices dengan bahasa pemrograman Go untuk sistem backend dari API service. Teknologi gRPC digunakan untuk komunikasi antar microservices secara synchronous, sedangkan komunikasi asynchronous menggunkana teknologi message broker RabbitMQ.

## 2. Design Principles

Beberapa prinsip utama dalam repository ini:

- **_API Gateway as Single Entry Point_**: Semua trafik dari klien masuk melalui satu gerbang utama yang menangani routing, protokol translasi (REST ke gRPC), dan validasi awal.

- **_Separation of Concerns_**: Setiap layanan memiliki satu domain tanggung jawab yang spesifik (misalnya: Authentication, Payment, Notification).

- **_Hybrid Communication_**:
  - _Synchronous (gRPC)_: Digunakan untuk request yang membutuhkan respon langsung (seperti proses login subscription).
  - _Asynchronous (RabbitMQ)_: Digunakan untuk _background task_ yang tidak memerlukan respon secara langsung (seperti pengiriman email).

Setiap service memiliki ruang lingkup fungsi masing-masing:

Client tidak pernah berkomunikasi langsung dengan domain service.
Tanggung jawab API Gateway:

- Mengelola autentikasi via JWT
- Meneruskan request ke backend via gRPC
- Tidak melakukan penyimpanan data
  |

## 3. System Overview

### 3.1 Request Flowchart

[![](https://mermaid.ink/img/pako:eNqlWPtv4kYQ_lesle7ESSQ8kwBqT00ITahCwvFo2hqEFnsDq_Pr1nYSGvjfb3bttdfEpkePH0LGM9_sPL4Z27whwzUJ6qAny30x1pgF2t1o5swcDT4fPmi_Fn0Si65FiRNoH7XLYV-7wQF5wZsfhU99wnT-R6vEfubybD9crhj21tztzaOuOJ9HBvxz87i4nUyGesCw43suCyrrIPB-WbLKZ-aGAWGnK1dIa-yYlhCz6MHjm01NUL1gRqTxH48T7Rlb1MQBdZ1dBjD-s6tDvM_UIJVVFM87pzejYVfHJvYggMqKeYbwisNgLU8wRK4pkDimzFtU4-Tk83bUG0-2MsNIFwtcHQWfXB48CoyIWoP4t3Gs-wbUUUz4UYkBGEu_PP6jKHAJqUE-1BAF08ZRfX7UQdrp6eRWz_elFDhjvqjphusEmDqEKTb8I9SiFSk7kmbEfPCzvUtgapfVvjGyoj50FdhuuSvq5GFHveGD7L5fYcRzfRq4bCMchODVX8ClYMVI0emDLwl7GF4uaWB_E2AvXFrUX2cyFcyRwvXVQtSwxIsI0qf39IrbK8hwx1OAyRvFWW3TmkXGiSiYISuj6GLSbIFOlR5jLtOa1dp2n0OqaVojRcdFoYwTyANuh1H2GrExtRbQ9SfKbEGSrSzbUaQd4o3NF9f_Zevw8m99z0ceS8HsAEm59j84-p4kHKQy1IuiyDc8hkwS857CESQ65xB_Obr31yQFk1f4crAlHLwCEelenPsE5mUtybqmFI6dw21A1xvMPBnCnWqT9BAuzw9zvcsIbOvYfpvUPTKUkuBgXN1UU7AZVVWcdqoAQSogtiJIOgZSklMAcg4omYG4EwvoQWjF6RxJ_3G49A1GvZ_a2OPplZ7nKG8QwPbAIHDt0YPAQVCaBXyPk2nwlXiKIXf98SSBQE2DAtPube86dW2siRkWhnJw86thHZog7qh4gtahDU8lOB-WjrpfMOv5uO7o4b50DVt1o3WZy-8IE2p8JexTwZTeTge6vsYePZXhuGw1n2fGmDOjlKHG3u2I__N-TGPEkkAQU8eX0jbhRwSRkhgXtaWpWl7Jjm6uNi55qpXTC4mqgeaeyrmX9cuvJCbpiEtJjjjIcmBlE1LHnHaKVy4WuEx0yXJQmbYIPWhQXL9oQxyxI0aCRNrgy7ELYnR5ddWfwIG_vWk-MIV0NI-5hl_WLLwkVkeboYigEBLaZW73cBhPNPGQLsF8TZRYniYRhfL-YdL__dglee8G9OlnH2vFyXqeq7w1GcVZvChlHsqizMw64PzQzhv1CKk-ODhKTEX2B5ea6uDATotPHsCbWuLJt-MXNf44tx_u_kNBVMJSpobZZ4Pe4LJ_p_MjtCFzn6mZFE55IJC1Uwgh7vSqVt6f0-xVrTrA4sohMI8mox3E724i2JmDymjFqIk6AQtJGUHToBYgojeOmiF4EbLJDPFpMTH7OkMzZwcYDzv_uK4tYfCiu1qjzhO2fJCieb-mGNiUmkANCOu6oROgTq3eblSFF9R5Q6-oc1Jrthqn1Va10W5dNM6arfZ5GW3gerPZOG3X6rVWs9Gsts_rFxe7MvpXHN04rZ-Ddb1dbbTqjdZZ_ayMiMmJMYh-ThC_Kuy-AwBh3IM?type=png)](https://mermaid.live/edit#pako:eNqlWPtv4kYQ_lesle7ESSQ8kwBqT00ITahCwvFo2hqEFnsDq_Pr1nYSGvjfb3bttdfEpkePH0LGM9_sPL4Z27whwzUJ6qAny30x1pgF2t1o5swcDT4fPmi_Fn0Si65FiRNoH7XLYV-7wQF5wZsfhU99wnT-R6vEfubybD9crhj21tztzaOuOJ9HBvxz87i4nUyGesCw43suCyrrIPB-WbLKZ-aGAWGnK1dIa-yYlhCz6MHjm01NUL1gRqTxH48T7Rlb1MQBdZ1dBjD-s6tDvM_UIJVVFM87pzejYVfHJvYggMqKeYbwisNgLU8wRK4pkDimzFtU4-Tk83bUG0-2MsNIFwtcHQWfXB48CoyIWoP4t3Gs-wbUUUz4UYkBGEu_PP6jKHAJqUE-1BAF08ZRfX7UQdrp6eRWz_elFDhjvqjphusEmDqEKTb8I9SiFSk7kmbEfPCzvUtgapfVvjGyoj50FdhuuSvq5GFHveGD7L5fYcRzfRq4bCMchODVX8ClYMVI0emDLwl7GF4uaWB_E2AvXFrUX2cyFcyRwvXVQtSwxIsI0qf39IrbK8hwx1OAyRvFWW3TmkXGiSiYISuj6GLSbIFOlR5jLtOa1dp2n0OqaVojRcdFoYwTyANuh1H2GrExtRbQ9SfKbEGSrSzbUaQd4o3NF9f_Zevw8m99z0ceS8HsAEm59j84-p4kHKQy1IuiyDc8hkwS857CESQ65xB_Obr31yQFk1f4crAlHLwCEelenPsE5mUtybqmFI6dw21A1xvMPBnCnWqT9BAuzw9zvcsIbOvYfpvUPTKUkuBgXN1UU7AZVVWcdqoAQSogtiJIOgZSklMAcg4omYG4EwvoQWjF6RxJ_3G49A1GvZ_a2OPplZ7nKG8QwPbAIHDt0YPAQVCaBXyPk2nwlXiKIXf98SSBQE2DAtPube86dW2siRkWhnJw86thHZog7qh4gtahDU8lOB-WjrpfMOv5uO7o4b50DVt1o3WZy-8IE2p8JexTwZTeTge6vsYePZXhuGw1n2fGmDOjlKHG3u2I__N-TGPEkkAQU8eX0jbhRwSRkhgXtaWpWl7Jjm6uNi55qpXTC4mqgeaeyrmX9cuvJCbpiEtJjjjIcmBlE1LHnHaKVy4WuEx0yXJQmbYIPWhQXL9oQxyxI0aCRNrgy7ELYnR5ddWfwIG_vWk-MIV0NI-5hl_WLLwkVkeboYigEBLaZW73cBhPNPGQLsF8TZRYniYRhfL-YdL__dglee8G9OlnH2vFyXqeq7w1GcVZvChlHsqizMw64PzQzhv1CKk-ODhKTEX2B5ea6uDATotPHsCbWuLJt-MXNf44tx_u_kNBVMJSpobZZ4Pe4LJ_p_MjtCFzn6mZFE55IJC1Uwgh7vSqVt6f0-xVrTrA4sohMI8mox3E724i2JmDymjFqIk6AQtJGUHToBYgojeOmiF4EbLJDPFpMTH7OkMzZwcYDzv_uK4tYfCiu1qjzhO2fJCieb-mGNiUmkANCOu6oROgTq3eblSFF9R5Q6-oc1Jrthqn1Va10W5dNM6arfZ5GW3gerPZOG3X6rVWs9Gsts_rFxe7MvpXHN04rZ-Ddb1dbbTqjdZZ_ayMiMmJMYh-ThC_Kuy-AwBh3IM)

| Service       | Tanggung Jawab Utama                                 |
| ------------- | ---------------------------------------------------- |
| API Gateway   | HTTP entry point, JWT validation, routing            |
| Authenticator | Autentikasi identitas, login, registrasi, JWT issuer |
| Payment       | Pemrosesan pembayaran melalui payment processor      |
| Subscription  | Subscription user dan core data management           |
| Notification  | Pengiriman dan logging notifikasi                    |

### 3.2 Struktur Repository

```cli
├── apps
│   ├── <nama microservice>
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── cmd
│   │   │   └── main.go
│   │   └── internal
│   │       ├── transport         -> handler inbound (http, gRPC, MQ)
│   │       │   ├── http          -> hanya ada di api-gateway
│   │       │   │   ├── handler.go
│   │       │   │   ├── middleware.go
│   │       │   │   └── router.go
│   │       │   ├── grpc          -> menerima arahan dari api-gateway
│   │       │   │   ├── handler.go
│   │       │   │   └── server.go
│   │       │   └── rabbitMQ      -> hanya ada di notification
│   │       │       └── consumer.go
│   │       ├── service
│   │       │   └── <nama-service>.go -> berisi business logic, orkestrasi
│   │       │                            dari inbound ke outbound
│   │       ├── adapter           -> outbound (DB, API, MQ, SMTP)
│   │       │   └── grpc
│   │       │       ├── <...>.go
│   │       │       └── client.go
|   │     (*dst)                  -> fungsi lain seperti config, auth, dll.
│   │       └── config
│   │           └── config.go
├── proto
│   ├── <nama microservice>
│   │   └── <nama microservice>.proto -> kontrak protobuf
|
└── diagram                       -> dokumentasi arsitektur
    └── flowchart_end_to_end.md
```

## 4. Internal Service Structure

### 4.1 Arah Dependency

> transport -> service -> adapters

- Arah dependency selalu satu arah
- Adapters tidak pernah memanggil service atau transport

\*Khusus untuk cron job:

> Scheduler -> service -> adapters

### 4.2 Peran Tiap Layer

**Transport**

- Menerima input dari eksternal
- Melakukan parsing dan validasi ringan
- Meneruskan ke service layer

**Service**

- Mengandung business logic dan workflow
- Menentukan urutan operasi
- Menghasilkan event (misal "user berhasil melakukan donasi")

**Adapters**

- Menghubungkan ke sistem eksternal
- Tidak memiliki business logic
- Testing dilakukan dengan menggunakan mock

## 5. Autentikasi & Authorization

- JWT:
  - Dibuat oleh Authentication service
  - Diverifikasi di API Gateway
- Microservices:
  - Tidak memvalidasi JWT
  - Mengasumsikan request telah diautentikasi

Authorization detail (role/permission) dapat ditambahkan di API Gateway jika diperlukan.

## 6. Messaging & Eventing

### 6.1 Event Philosophy

- Event bersifat fakta yang telah terjadi
- Event bersifat immutable
- Event tidak digunakan untuk request/response

### 6.2 RabbitMQ

Digunakan untuk:

- Email confirmation
- Payment result
- Subscription updates
- Notification service adalah consumer utama

## 7. External Integrations

- Payment provider: Xendit
- Humanitarian Data API: (hapi.humdata.org)
- SMTP provider
