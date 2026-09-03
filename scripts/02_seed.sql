-- Seed fixed UUIDs for deterministic testing
-- Password for all seed users is: "password123" (Bcrypt cost 10)
-- Hash: $2a$10$wN9iLOnXzH8m4B9eM68S2e8uLd2G1vA7K8x9p3q1zW5e0r2t4y6u.

-- 1. Seed System Users (Admin, Manager, Staff)
INSERT INTO users (id, email, password_hash, full_name, role, is_active) VALUES
  ('11111111-1111-1111-1111-111111111111', 'admin@example.com', '$2a$10$wN9iLOnXzH8m4B9eM68S2e8uLd2G1vA7K8x9p3q1zW5e0r2t4y6u.', 'System Admin', 'ADMIN', TRUE),
  ('22222222-2222-2222-2222-222222222222', 'manager@example.com', '$2a$10$wN9iLOnXzH8m4B9eM68S2e8uLd2G1vA7K8x9p3q1zW5e0r2t4y6u.', 'Store Manager', 'MANAGER', TRUE),
  ('33333333-3333-3333-3333-333333333333', 'staff@example.com', '$2a$10$wN9iLOnXzH8m4B9eM68S2e8uLd2G1vA7K8x9p3q1zW5e0r2t4y6u.', 'Staff Member', 'STAFF', TRUE)
ON CONFLICT (id) DO NOTHING;

-- 2. Seed Bookables
INSERT INTO bookables (id, name, description, type, metadata, created_by, updated_by) VALUES
  (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 
    'Main Conference Room', 
    'Spacious room equipped with 4K projector and video conferencing.', 
    'PHYSICAL', 
    '{"capacity": 20, "floor": 2, "has_projector": true}'::jsonb, 
    '11111111-1111-1111-1111-111111111111', 
    '11111111-1111-1111-1111-111111111111'
  ),
  (
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a22', 
    '1-on-1 Security Consultation', 
    '60-minute remote architecture review session.', 
    'SERVICE', 
    '{"duration_minutes": 60, "remote": true}'::jsonb, 
    '22222222-2222-2222-2222-222222222222', 
    '22222222-2222-2222-2222-222222222222'
  )
ON CONFLICT (id) DO NOTHING;

-- 3. Seed Sample Bookings
INSERT INTO bookings (
  id, 
  bookable_id, 
  customer_name, 
  customer_email, 
  customer_phone, 
  booking_window, 
  status, 
  created_by, 
  updated_by
) VALUES
  (
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'John Doe',
    'john.doe@client.com',
    '+1234567890',
    tstzrange('2026-10-01 09:00:00+00', '2026-10-01 10:00:00+00', '[)'),
    'CONFIRMED',
    '33333333-3333-3333-3333-333333333333',
    '33333333-3333-3333-3333-333333333333'
  )
ON CONFLICT (id) DO NOTHING;

-- 4. Seed Audit Log Entry
INSERT INTO audit_logs (
  id, 
  actor_id, 
  action, 
  entity_type, 
  entity_id, 
  new_state
) VALUES
  (
    gen_random_uuid(),
    '33333333-3333-3333-3333-333333333333',
    'BOOKING_CREATED',
    'BOOKING',
    'c0eebc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    '{"customer_email": "john.doe@client.com", "status": "CONFIRMED"}'::jsonb
  );