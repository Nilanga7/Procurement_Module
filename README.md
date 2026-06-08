# ERP Procurement Module

> \*\*Technology Stack: \*\* Vue.js · Go (Gin) · PostgreSQL · Docker

\---

## Project Description

The **Procurement Module** is part of a larger ERP system built collaboratively by multiple groups. It manages the full procurement lifecycle inside an organisation — from a department raising a purchase request through finance approval, supplier bidding, purchase order creation, and final goods receipt.



### Key Features

* **Role-based access control** — four roles: `admin`, `procurement\_officer`, `procurement\_manager`, `supplier`
* **JWT authentication** — all protected routes require a Bearer token
* **Full procurement workflow** — Draft → Finance Check → Manager Approval → Open for Bidding → Supplier Selected → Purchase Order → Goods Receipt → Completed
* **Cross-module integration** — pushes budget checks to the Finance Module and registers assets with the Asset Module after goods are received
* **Transaction history** — filterable, paginated audit view across all procurement activity

\---

## Quick Start

### Prerequisites

* [Docker Desktop](https://www.docker.com/products/docker-desktop/) installed and running
* Ports `3000`, `8080`, `5432` free on your machine

### 1\. Clone / unzip the project

```bash
unzip "Procurement Module.zip"
cd "Procurement Module"
```

### 2\. Configure environment variables

Copy the example env file and adjust if needed:

```bash
cp .env.example .env   # or edit .env directly
```

The default `.env` values work out of the box for local Docker. Only change `FINANCE\_MODULE\_URL` / `ASSET\_MODULE\_URL` when integrating with the other teams' modules.

### 3\. Install frontend dependencies

> **Note:** `node_modules` is not included in the shipped package. You must install dependencies once before building.

```bash
cd frontend
npm install
cd ..
```

### 4\. Start the system

```bash
docker-compose up --build
```

Docker will:

1. Start **PostgreSQL** and run `init.sql` + `seed.sql` automatically
2. Build and start the **Go backend** (waits for DB health check)
3. Build and start the **Vue frontend** via nginx

### 5\. Open the application

|Service|URL|
|-|-|
|Frontend (Vue)|http://localhost:3000|
|Backend API|http://localhost:8080|
|PostgreSQL|localhost:5432|

### 6\. Seed credentials (from `seed.sql`)

|Role|Email|Password|
|-|-|-|
|Admin|`admin@erp.com`|`admin123`|
|Procurement Officer|`officer@erp.com`|`officer123`|
|Procurement Manager|`manager@erp.com`|`manager123`|
|Supplier|`supplier@erp.com`|`supplier123`|

\---

## Port Details

|Container|Internal Port|Host Port|Purpose|
|-|-|-|-|
|`frontend`|80|**3000**|Vue.js SPA served by nginx|
|`backend`|8080|**8080**|Go/Gin REST API|
|`postgres`|5432|**5432**|PostgreSQL 16|

\---

## API Endpoints

All API calls are prefixed with `http://localhost:8080`.  
Protected routes require the header: `Authorization: Bearer <token>`

### Authentication (Public)

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/auth/login`|Login — returns JWT token|
|POST|`/api/auth/register`|Register new user (officer / manager / supplier)|

**Login body:**

```json
{ "email": "officer@erp.com", "password": "officer123" }
```

**Register body:**

```json
{
  "name": "Jane Smith",
  "email": "jane@example.com",
  "password": "password123",
  "role": "procurement\_officer"
}
```

\---

### Procurement Requests 🔒

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/procurement/requests`|Create a new procurement request|
|GET|`/api/procurement/requests`|List all requests|
|GET|`/api/procurement/requests/:id`|Get a single request|
|PUT|`/api/procurement/requests/:id`|Update a request|
|DELETE|`/api/procurement/requests/:id`|Delete a request|

**Create request body:**

```json
{
  "title": "Office Laptops Q3",
  "description": "10 laptops for the IT department",
  "department": "IT",
  "required\_date": "2026-08-01",
  "estimated\_total": 1500000.00,
  "items": \[
    {
      "item\_name": "Dell Laptop 15",
      "quantity": 10,
      "estimated\_price": 150000.00,
      "category": "Electronics",
      "item\_type": "asset"
    }
  ]
}
```

**Request status flow:**

```
draft → sent\_to\_finance → budget\_approved / budget\_rejected
      → manager\_approved / manager\_rejected
      → open\_for\_bidding → supplier\_selected
      → purchase\_order\_created → goods\_received
      → sent\_to\_asset\_module → sent\_to\_finance\_for\_payment → completed
```

\---

### Finance Integration 🔒

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/procurement/requests/:id/send-to-finance`|Send request to Finance Module for budget check|
|POST|`/api/procurement/finance-response`|Receive budget approval/rejection from Finance Module|

**Finance response body (called by Finance Module):**

```json
{
  "request\_id": "PR0001",
  "budget\_status": "APPROVED",
  "remarks": "Within Q3 budget allocation"
}
```

\---

### Approval 🔒 `procurement\_manager only`

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/procurement/requests/:id/approve`|Approve a budget-approved request (opens for bidding)|
|POST|`/api/procurement/requests/:id/reject`|Reject a request|

\---

### Supplier Bids 🔒

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/supplier/bids`|Submit a bid (supplier role)|
|GET|`/api/supplier/bids`|View my submitted bids (supplier role)|
|GET|`/api/procurement/requests/:id/bids`|View all bids for a request|
|POST|`/api/procurement/bids/:id/accept`|Accept a bid|
|POST|`/api/procurement/bids/:id/reject`|Reject a bid|

\---

### Purchase Orders 🔒

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/procurement/purchase-orders`|Create PO from accepted bid|
|GET|`/api/procurement/purchase-orders`|List all purchase orders|
|GET|`/api/procurement/purchase-orders/:id`|Get a single PO|

\---

### Goods Receipts 🔒

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/procurement/goods-receipts`|Record goods received against a PO|
|GET|`/api/procurement/goods-receipts`|List all goods receipts|

On creation, this endpoint automatically notifies both the **Asset Module** (for asset-type items) and the **Finance Module** (for payment processing).

\---

### Asset Integration 🔒

|Method|Endpoint|Description|
|-|-|-|
|POST|`/api/procurement/assets/send`|Manually push asset items to Asset Module|
|POST|`/api/procurement/asset-response`|Receive registration confirmation from Asset Module|

\---

### Transactions 🔒

|Method|Endpoint|Description|
|-|-|-|
|GET|`/api/procurement/transactions`|Paginated, filterable transaction history|

**Query parameters:**

|Parameter|Type|Example|Description|
|-|-|-|-|
|`status`|string|`completed`|Filter by request status|
|`department`|string|`IT`|Filter by department (partial match)|
|`supplier`|string|`ABC`|Filter by supplier name|
|`request\_id`|string|`PR0001`|Filter by request number|
|`from`|date|`2026-01-01`|Date range start|
|`to`|date|`2026-12-31`|Date range end|
|`page`|int|`1`|Page number (default: 1)|
|`limit`|int|`10`|Results per page (default: 10, max: 100)|

\---

## Integration Details

This module integrates with two other ERP modules over HTTP using REST + JSON.

### Finance Module

|Direction|Trigger|Endpoint Called|
|-|-|-|
|Outbound (Procurement → Finance)|Officer clicks "Send to Finance"|`POST {FINANCE\_MODULE\_URL}/api/finance/budget-check`|
|Inbound  (Finance → Procurement)|Finance approves/rejects|`POST /api/procurement/finance-response`|
|Outbound (Procurement → Finance)|Goods received|`POST {FINANCE\_MODULE\_URL}/api/finance/payment-request` (async, triggered by goods receipt)|

**Outbound payload to Finance:**

```json
{
  "request\_id": "PR0001",
  "estimated\_total": 1500000.00,
  "currency": "LKR"
}
```

### Asset Module

|Direction|Trigger|Endpoint Called|
|-|-|-|
|Outbound (Procurement → Asset)|Goods received \& items are type `asset`|`POST {ASSET\_MODULE\_URL}/api/assets/register`|
|Inbound  (Asset → Procurement)|Asset registered confirmation|`POST /api/procurement/asset-response`|

**Outbound payload to Asset Module:**

```json
{
  "po\_id": "<uuid>",
  "items": \[
    { "item\_name": "Dell Laptop 15", "quantity": 10, "unit\_price": 150000.00 }
  ]
}
```

### Environment variables for integration

Set these in `.env` before running:

```env
FINANCE\_MODULE\_URL=http://finance-backend:8081
ASSET\_MODULE\_URL=http://asset-backend:8082
```

When running in a shared Docker network with the other teams, replace the hostnames with the actual Docker service names from their `docker-compose.yml`.

\---

## Database Schema

9 tables in PostgreSQL:

|Table|Description|
|-|-|
|`users`|All system users with hashed passwords and roles|
|`suppliers`|Supplier profiles linked to supplier users|
|`procurement\_requests`|Core request records with status tracking|
|`procurement\_request\_items`|Line items for each request|
|`supplier\_bids`|Bids submitted by suppliers per request|
|`purchase\_orders`|POs generated from accepted bids|
|`purchase\_order\_items`|Line items copied from request items|
|`goods\_receipts`|Delivery confirmations against POs|
|`goods\_receipt\_items`|Individual item receipt quantities|

\---

## Project Structure

```
Procurement Module/
├── backend/
│   ├── cmd/main.go              # Entry point, route registration
│   ├── internal/
│   │   ├── auth/auth.go         # JWT generation \& validation
│   │   ├── db/db.go             # PostgreSQL connection pool
│   │   ├── handlers/            # Route handlers (one file per domain)
│   │   ├── middleware/          # CORS, JWT auth, role guard
│   │   └── models/              # Shared struct definitions
│   ├── Dockerfile
│   └── go.mod
├── frontend/
│   ├── src/
│   │   ├── views/               # Vue page components
│   │   ├── stores/auth.js       # Pinia auth store
│   │   ├── api/index.js         # Axios API client
│   │   └── router/index.js      # Vue Router
│   ├── nginx.conf
│   └── Dockerfile
├── database/
│   ├── init.sql                 # Schema creation
│   ├── seed.sql                 # Sample data \& demo users
│   └── fix\_supplier\_link.sql    # Patch script
├── docker-compose.yml
└── .env
```

\---

## Stopping the System

```bash
docker-compose down          # Stop containers
docker-compose down -v       # Stop and remove DB volume (full reset)
```

