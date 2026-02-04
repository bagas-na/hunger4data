# API Gateway – Guideline

## Tujuan

API Gateway adalah **satu-satunya _entry point_ eksternal** sistem.
Bertugas untuk:

- Menerima HTTP request dari client
- Melakukan autentikasi (validasi JWT)
- Meneruskan request ke backend service melalui gRPC

API Gateway fokus untuk _orchestration_ (mengarahkan request ke backend lain secara tepat)

## Inbound Structure

- All inbound HTTP handling lives in `internal/transport/http`
- `router.go` berguna untuk menyatukan handler ke satu tempat (RegisterRoutes)
- `xxx_handler.go` parse HTTP, cek kelengkapan json (kalau ada) -> service calls
- `middleware.go` melakukan pengecekan jwt, atau fungsi umum lainnya (misal logger, jika perlu)

### JWT Handling

- Validasi JWT dilakukan di `middleware.go`
- Handler dan service dapat mengasumsikan request sudah terautentikasi

## Service Layer

- `service/gateway.go` mengkoordinasikan pemanggilan ke microservices lain
- Service layer boleh:
  - Memanggil gRPC adapter
  - Menggabungkan atau meneruskan response
- Service layer sebaiknya menghindari:
  - Business / domain logic (karena akan diimplementasikan di microservices lainnya)
  - Urusan penyimpanan data (akan diimplementasikan oleh microservices)

## Adapter

- `adapter/grpc` berisi seluruh client untuk komunikasi outbound gRPC (gRPC ke luar, yang menuju microservices lainnya)
- Adapter mengkonversi _internal calls_ menjadi _external protocol_
- Adapter sebisa mungkin dibuat ringan, atau dengan kata lain, jika beban pekerjaan bisa dialihkan ke microservices, sebaiknya dilakukan di microservices saja

## Design Intent

- Tidak melakukan pemrosesan data -> secukupnya saja agar bisa meneruskan pesan yang valid ke layanan (microservices) berikutnya
- Tidak memiliki _business logic_
- Tidak ada akses langsung ke database
- Tidak ada pemanggilan external API secara langsung (external API dipanggil oleh microservices lainnya)
