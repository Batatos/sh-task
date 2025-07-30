-- Simple Skyhawk Security Database Schema
-- PostgreSQL 13+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ========================================
-- SIMPLE EVENTS TABLE
-- ========================================

-- Security Events Table (simplified)
CREATE TABLE security_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_id VARCHAR(255) UNIQUE NOT NULL,
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL CHECK (severity IN ('low', 'medium', 'high', 'critical')),
    source VARCHAR(255) NOT NULL,
    description TEXT,
    event_data JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- ========================================
-- BASIC INDEXES
-- ========================================

-- Indexes for common queries
CREATE INDEX idx_security_events_event_type ON security_events(event_type);
CREATE INDEX idx_security_events_severity ON security_events(severity);
CREATE INDEX idx_security_events_created_at ON security_events(created_at);
CREATE INDEX idx_security_events_event_data ON security_events USING GIN (event_data);

-- ========================================
-- TRIGGER FOR UPDATED_AT
-- ========================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create trigger for updated_at
CREATE TRIGGER update_security_events_updated_at 
    BEFORE UPDATE ON security_events 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- ========================================
-- SAMPLE DATA
-- ========================================

-- Insert sample events
INSERT INTO security_events (event_id, event_type, severity, source, description, event_data) VALUES
('event-20240115103015-123456789', 'login', 'high', 'web-application', 'Multiple failed login attempts', '{"ip": "192.168.1.100", "user": "admin", "attempts": 5}'),
('event-20240115103016-123456790', 'data_access', 'medium', 'database', 'Unusual data access pattern', '{"table": "users", "rows_accessed": 1000, "user": "analyst"}'),
('event-20240115103017-123456791', 'file_access', 'low', 'file-system', 'File access outside business hours', '{"file": "/etc/passwd", "user": "developer", "time": "02:30"}'); 