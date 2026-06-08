-- Run this ONCE against your existing database to fix the missing supplier profile.
-- This inserts the supplier profile row that seed.sql failed to create
-- (because it had the wrong email 'abc@supplier.com' instead of 'supplier@supplier.com').

INSERT INTO suppliers (user_id, company_name, contact_email)
SELECT id, 'ABC Suppliers Ltd', 'supplier@supplier.com'
FROM users
WHERE email = 'supplier@supplier.com'
ON CONFLICT DO NOTHING;   -- safe to re-run: no-ops if the row already exists
