# Subscription Service – Guidelines

## Tujuan

Subscription service mengelola:

- Subscription user ke cause tertentu
- Periodic data ingestion
- Publishing update subscription

---

## Inbound Structure

- gRPC for subscription management
- Internal scheduler for periodic jobs

---

## Service Layer

- `service/subscription.go`
  - Subscription lifecycle
- `service/scheduler.go`
  - Fetch cadence and orchestration

---

## Adapters

- `adapters/external`
  - Humanitarian API client
- `adapters/repo`
  - Subscription persistence
- `adapters/rabbitmq`
  - Publishing update events

---

## Messaging

- Publishes subscription update events
- Does not consume messages

---

## Design Intent

- Polling logic is isolated here
- Notifications are delegated downstream
- No payment or authentication logic
