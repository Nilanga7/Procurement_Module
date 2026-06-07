CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE user_role AS ENUM ('admin','procurement_officer','procurement_manager','supplier');
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE suppliers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID UNIQUE REFERENCES users(id),
    company_name VARCHAR(150) NOT NULL,
    contact_email VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE pr_status AS ENUM (
    'draft','sent_to_finance','budget_approved','budget_rejected',
    'manager_approved','manager_rejected','open_for_bidding',
    'supplier_selected','purchase_order_created','goods_received',
    'sent_to_asset_module','sent_to_finance_for_payment','completed'
);
CREATE TABLE procurement_requests (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_number VARCHAR(20) UNIQUE NOT NULL,
    title VARCHAR(200) NOT NULL,
    description TEXT,
    department VARCHAR(100),
    project_id VARCHAR(100),
    required_date DATE,
    estimated_total NUMERIC(15,2),
    currency VARCHAR(10) DEFAULT 'LKR',
    status pr_status DEFAULT 'draft',
    created_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    approval_remarks TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TYPE item_type AS ENUM ('asset','consumable');
CREATE TABLE procurement_request_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID REFERENCES procurement_requests(id) ON DELETE CASCADE,
    item_name VARCHAR(200) NOT NULL,
    quantity INTEGER NOT NULL,
    estimated_price NUMERIC(12,2),
    category VARCHAR(100),
    item_type item_type NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE supplier_bids (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    request_id UUID REFERENCES procurement_requests(id),
    supplier_id UUID REFERENCES suppliers(id),
    offered_price NUMERIC(15,2) NOT NULL,
    delivery_days INTEGER,
    notes TEXT,
    status VARCHAR(30) DEFAULT 'pending',
    response_sent BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE purchase_orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    po_number VARCHAR(20) UNIQUE NOT NULL,
    request_id UUID REFERENCES procurement_requests(id),
    supplier_id UUID REFERENCES suppliers(id),
    total_amount NUMERIC(15,2),
    delivery_date DATE,
    status VARCHAR(30) DEFAULT 'pending',
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE purchase_order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    po_id UUID REFERENCES purchase_orders(id) ON DELETE CASCADE,
    item_name VARCHAR(200),
    quantity INTEGER,
    unit_price NUMERIC(12,2)
);

CREATE TABLE goods_receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    po_id UUID REFERENCES purchase_orders(id),
    received_date DATE,
    received_by UUID REFERENCES users(id),
    remarks TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE goods_receipt_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    receipt_id UUID REFERENCES goods_receipts(id) ON DELETE CASCADE,
    item_name VARCHAR(200),
    quantity_received INTEGER
);

CREATE TABLE asset_transfer_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    po_id UUID REFERENCES purchase_orders(id),
    payload JSONB,
    sent_at TIMESTAMP DEFAULT NOW(),
    status VARCHAR(30) DEFAULT 'sent'
);