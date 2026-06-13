-- 
-- FolksArt IAM Governance Suite: Relational Schema Mapping (MySQL RDBMS)
-- 

CREATE DATABASE IF NOT EXISTS folksart_iam;
USE folksart_iam;

-- 1. Subject Principals Table (IAM Users)
CREATE TABLE IF NOT EXISTS iam_users (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL,
    phone VARCHAR(30),
    role VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL,
    kyc_status VARCHAR(50) NOT NULL,
    department VARCHAR(100) NOT NULL,
    risk_score INT DEFAULT 0,
    mfa_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_role (role),
    INDEX idx_user_status (status),
    INDEX idx_user_department (department)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 2. Audit Governance Trail Table
CREATE TABLE IF NOT EXISTS audit_logs (
    id VARCHAR(50) PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    actor VARCHAR(50) NOT NULL,
    action TEXT NOT NULL,
    target VARCHAR(100) NOT NULL,
    severity VARCHAR(30) NOT NULL,
    INDEX idx_log_severity (severity),
    INDEX idx_log_timestamp (timestamp)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 3. Mock Corporate Seed Dataset
INSERT INTO iam_users (id, name, username, email, phone, role, status, kyc_status, department, risk_score, mfa_enabled, created_at) VALUES
('usr-8821', 'Sarah Connor', 'sconnor', 'sconnor@cyberdyne.org', '+1 (555) 019-9042', 'Administrator', 'Active', 'Verified', 'Security Operations', 12, TRUE, '2026-01-10 08:34:21'),
('usr-0101', 'Marcus Wright', 'mwright', 'mwright@cyberdyne.org', '+1 (555) 012-3849', 'Security Officer', 'Active', 'Verified', 'Security Operations', 45, TRUE, '2026-02-14 09:15:00'),
('usr-1991', 'John Connor', 'jconnor', 'jconnor@cyberdyne.org', '+1 (555) 018-1991', 'End User', 'Active', 'Verified', 'Engineering', 15, TRUE, '2026-03-01 10:45:12'),
('usr-0800', 'T-800 Guardian', 't800', 't800@cyberdyne.org', '+1 (555) 012-0800', 'End User', 'Active', 'Suspicious', 'Logistics', 68, FALSE, '2026-03-12 00:00:00'),
('usr-1000', 'T-1000 Mimetic', 't1000', 't1000@cyberdyne.org', '+1 (555) 011-1000', 'End User', 'Banned', 'Failed', 'Finance', 98, FALSE, '2026-04-18 14:22:10'),
('usr-7762', 'Miles Dyson', 'mdyson', 'mdyson@cyberdyne.org', '+1 (555) 014-9981', 'Administrator', 'Active', 'Verified', 'Research & Development', 28, TRUE, '2026-01-05 11:20:34'),
('usr-2931', 'Dr. Peter Silberman', 'psilberman', 'psilberman@cyberdyne.org', '+1 (555) 015-8821', 'End User', 'Active', 'Pending', 'Legal', 52, FALSE, '2026-05-20 16:40:02'),
('usr-4411', 'Kate Brewster', 'kbrewster', 'kbrewster@cyberdyne.org', '+1 (555) 016-3023', 'Security Officer', 'Active', 'Verified', 'Engineering', 31, TRUE, '2026-05-22 13:10:55');

INSERT INTO audit_logs (id, timestamp, actor, action, target, severity) VALUES
('log-001', '2026-06-12 10:00:00', 'sarah_connor', 'System Initialization', 'IAM Console Enterprise', 'Low'),
('log-002', '2026-06-12 14:22:15', 'system_daemon', 'Automated Threat Score Recalculation', 'T-1000 Mimetic', 'High'),
('log-003', '2026-06-12 14:23:40', 'sarah_connor', 'Identity State Interdiction (Status: Banned)', 'T-1000 Mimetic', 'Critical'),
('log-004', '2026-06-12 18:05:12', 'marcus_wright', 'KYC Suspicious Flag Appended', 'T-800 Guardian', 'Medium');
