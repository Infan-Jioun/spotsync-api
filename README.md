# SpotSync API 🚗

Smart Parking & EV Charging Reservation System



## Tech Stack
- Go (Golang) + Echo Framework
- GORM + PostgreSQL (NeonDB)
- JWT Authentication
- bcrypt Password Hashing

## Architecture
Handler → Service → Repository → Database

## Setup Locally
1. Clone the repo
2. Create `.env` file:
   DB_URL=your_neondb_url
   JWT_SECRET=your_secret
   PORT=8080
3. Run: `go run main.go`

## API Endpoints

### Auth (Public)
- POST /api/v1/auth/register
- POST /api/v1/auth/login

### Zones (Public GET, Admin POST)
- GET  /api/v1/zones
- GET  /api/v1/zones/:id
- POST /api/v1/zones (Admin only)

### Reservations (Authenticated)
- POST   /api/v1/reservations
- GET    /api/v1/reservations/my-reservations
- DELETE /api/v1/reservations/:id
- GET    /api/v1/reservations (Admin only)
