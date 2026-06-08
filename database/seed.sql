-- Password for all users: password123

INSERT INTO users (name, email, password_hash, role) VALUES
('Admin User', 'admin@erp.com',
'$2a$10$8xZdmPdjZqq6tG4JEAwxAuzRNbyk0hRJpMGgc.z65NqcmB2bKhtI.',
'admin'),

('Procurement Officer', 'officer@erp.com',
'$2a$10$8xZdmPdjZqq6tG4JEAwxAuzRNbyk0hRJpMGgc.z65NqcmB2bKhtI.',
'procurement_officer'),

('Procurement Manager', 'manager@erp.com',
'$2a$10$8xZdmPdjZqq6tG4JEAwxAuzRNbyk0hRJpMGgc.z65NqcmB2bKhtI.',
'procurement_manager'),

('Supplier', 'supplier@supplier.com',
'$2a$10$8xZdmPdjZqq6tG4JEAwxAuzRNbyk0hRJpMGgc.z65NqcmB2bKhtI.',
'supplier');

INSERT INTO suppliers (user_id, company_name, contact_email)
SELECT id, 'ABC Suppliers Ltd', 'supplier@supplier.com'
FROM users
WHERE email = 'supplier@supplier.com';