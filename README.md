# E-Commerce Backend API

Backend API สำหรับระบบ E-Commerce พัฒนาโดยใช้ Go (Gin + GORM) และ PostgreSQL  
รองรับระบบสั่งซื้อสินค้า, ชำระเงิน, รีวิวสินค้า และการติดตามสถานะออเดอร์

---

## Tech Stack

- Go (Golang)
- Gin (HTTP framework)
- GORM (ORM)
- PostgreSQL
- Docker & Docker Compose
- Swagger (API Documentation)

---

## Features

### Authentication
- Register / Login (JWT)
- Role-based access control (user / admin)

### Product
- Create / Update / Delete product
- Search / Filter / Pagination
- Stock management

### Cart
- Add to cart
- Update quantity
- Remove item

### Order
- Checkout (Transaction-safe)
- Order items snapshot
- Stock deduction
- Cancel order + restore stock

### Payment
- Mock payment system
- Payment status (pending / success / failed)

### Review
- User can review product (1 ครั้ง / user)
- Average rating calculation
- Constraint validation

### Order Status Tracking
- Status flow (pending → paid → shipped → completed)
- Order status history (audit log)

---

## Database Design

- users
- products
- carts
- cart_items
- orders
- order_items
- payments
- reviews
- order_status_histories

ใช้ **Database Migration (golang-migrate)** แทน AutoMigrate

---

## Getting Started

### 1. Clone project

```bash
git clone https://github.com/your-username/ecommerce-api.git
cd ecommerce-api
