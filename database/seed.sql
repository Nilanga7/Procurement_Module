-- Password for all users: password123

INSERT INTO users (name, email, password_hash, role) VALUES
('Admin User', 'admin@erp.com',
'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
'admin'),

('Procurement Officer', 'officer@erp.com',
'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
'procurement_officer'),

('Procurement Manager', 'manager@erp.com',
'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
'procurement_manager'),

('Supplier', 'supplier@supplier.com',
'$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy',
'supplier');

INSERT INTO suppliers (user_id, company_name, contact_email)
SELECT id, 'ABC Suppliers Ltd', 'supplier@supplier.com'
FROM users
WHERE email = 'supplier@supplier.com';