CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "btree_gist";

-- Custom Enums
CREATE TYPE bookable_type AS ENUM ('SERVICE', 'PHYSICAL');
CREATE TYPE booking_status AS ENUM ('PENDING', 'CONFIRMED', 'CANCELLED', 'COMPLETED');
CREATE TYPE user_role AS ENUM ('ADMIN', 'MANAGER', 'STAFF');

-- 1. System Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'STAFF',
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 2. Bookables Table
CREATE TABLE bookables (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    type bookable_type NOT NULL,
    metadata JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 3. Bookings Table (Client info remains purely informational)
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    bookable_id UUID NOT NULL REFERENCES bookables(id) ON DELETE CASCADE,
    
    -- Client info (Non-system actors)
    customer_name VARCHAR(255) NOT NULL,
    customer_email VARCHAR(255) NOT NULL,
    customer_phone VARCHAR(50),
    
    booking_window TSTZRANGE NOT NULL,
    status booking_status NOT NULL DEFAULT 'CONFIRMED',
    
    -- Action performed by System User (null if self-registered guest booking via public endpoint)
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    updated_by UUID REFERENCES users(id) ON DELETE SET NULL,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    -- DB Engine Level Double-Booking Exclusion Constraint
    CONSTRAINT no_overlapping_bookings EXCLUDE USING gist (
        bookable_id WITH =,
        booking_window WITH &&
    ) WHERE (status != 'CANCELLED')
);

-- 4. Audit Logs Table (Immutable Event Ledger)
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES users(id) ON DELETE SET NULL, -- Who performed the action
    action VARCHAR(100) NOT NULL,                           -- e.g., 'BOOKABLE_CREATED', 'BOOKING_CANCELLED'
    entity_type VARCHAR(50) NOT NULL,                        -- 'BOOKABLE' or 'BOOKING'
    entity_id UUID NOT NULL,                                -- ID of the target resource
    old_state JSONB,                                        -- State prior to update (NULL for INSERT)
    new_state JSONB,                                        -- State after update (NULL for DELETE)
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for Fast Querying & Audit Lookups
CREATE INDEX idx_bookables_metadata ON bookables USING gin (metadata);
CREATE INDEX idx_audit_logs_entity ON audit_logs (entity_type, entity_id);
CREATE INDEX idx_audit_logs_actor ON audit_logs (actor_id);