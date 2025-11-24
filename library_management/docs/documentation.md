# Concurrent Book Reservation System

## Concurrency Design

- Goroutines: Each reservation request runs in its own goroutine
- Channels: `ReservationChan` safely queues incoming requests
- sync.Mutex: Protects shared state (books & members) from race conditions
- Auto-cancellation: Reserved books are freed after 5 seconds if not borrowed
- Safe borrowing: Only the reserving member can borrow during reservation window

## Key Features

- Only one member can reserve a book at a time
- Race-condition free
- High throughput via worker pool pattern
- Real-time feedback in simulation
- Clean separation: concurrency logic isolated in `concurrency/` package

## Run Simulation

```bash
go run .`