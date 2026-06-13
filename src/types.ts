/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

export type UserRole = "Administrator" | "Security Officer" | "End User";

export type UserStatus = "Active" | "Banned" | "Deactivated";

export type KYCStatus = "Pending" | "Verified" | "Failed" | "Suspicious";

export interface IAMUser {
  id: string;
  name: string;
  username: string;
  email: string;
  phone: string;
  role: UserRole;
  status: UserStatus;
  kycStatus: KYCStatus;
  department: string;
  riskScore: number; // 0 to 100
  mfaEnabled: boolean;
  createdAt: string;
}

export type LogSeverity = "Critical" | "High" | "Medium" | "Low";

export interface AuditLog {
  id: string;
  timestamp: string;
  actor: string;
  action: string;
  target: string;
  severity: LogSeverity;
}

export interface GovernanceStats {
  totalUsers: number;
  criticalRisks: number;
  kycComplianceRate: number;
  mfaAdoptionRate: number;
}
