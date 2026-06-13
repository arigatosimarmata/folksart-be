/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import express from "express";
import path from "path";
import { createServer as createViteServer } from "vite";
import { IAMUser, AuditLog, KYCStatus, UserRole, UserStatus } from "./src/types";

// In-memory persistent arrays populated with elite corporate data
let users: IAMUser[] = [
  {
    id: "usr-8821",
    name: "Sarah Connor",
    username: "sconnor",
    email: "sconnor@cyberdyne.org",
    phone: "+1 (555) 019-9042",
    role: "Administrator",
    status: "Active",
    kycStatus: "Verified",
    department: "Security Operations",
    riskScore: 12,
    mfaEnabled: true,
    createdAt: "2026-01-10T08:34:21Z"
  },
  {
    id: "usr-0101",
    name: "Marcus Wright",
    username: "mwright",
    email: "mwright@cyberdyne.org",
    phone: "+1 (555) 012-3849",
    role: "Security Officer",
    status: "Active",
    kycStatus: "Verified",
    department: "Security Operations",
    riskScore: 45,
    mfaEnabled: true,
    createdAt: "2026-02-14T09:15:00Z"
  },
  {
    id: "usr-1991",
    name: "John Connor",
    username: "jconnor",
    email: "jconnor@cyberdyne.org",
    phone: "+1 (555) 018-1991",
    role: "End User",
    status: "Active",
    kycStatus: "Verified",
    department: "Engineering",
    riskScore: 15,
    mfaEnabled: true,
    createdAt: "2026-03-01T10:45:12Z"
  },
  {
    id: "usr-0800",
    name: "T-800 Guardian",
    username: "t800",
    email: "t800@cyberdyne.org",
    phone: "+1 (555) 012-0800",
    role: "End User",
    status: "Active",
    kycStatus: "Suspicious",
    department: "Logistics",
    riskScore: 68,
    mfaEnabled: false,
    createdAt: "2026-03-12T00:00:00Z"
  },
  {
    id: "usr-1000",
    name: "T-1000 Mimetic",
    username: "t1000",
    email: "t1000@cyberdyne.org",
    phone: "+1 (555) 011-1000",
    role: "End User",
    status: "Banned",
    kycStatus: "Failed",
    department: "Finance",
    riskScore: 98,
    mfaEnabled: false,
    createdAt: "2026-04-18T14:22:10Z"
  },
  {
    id: "usr-7762",
    name: "Miles Dyson",
    username: "mdyson",
    email: "mdyson@cyberdyne.org",
    phone: "+1 (555) 014-9981",
    role: "Administrator",
    status: "Active",
    kycStatus: "Verified",
    department: "Research & Development",
    riskScore: 28,
    mfaEnabled: true,
    createdAt: "2026-01-05T11:20:34Z"
  },
  {
    id: "usr-2931",
    name: "Dr. Peter Silberman",
    username: "psilberman",
    email: "psilberman@cyberdyne.org",
    phone: "+1 (555) 015-8821",
    role: "End User",
    status: "Active",
    kycStatus: "Pending",
    department: "Legal",
    riskScore: 52,
    mfaEnabled: false,
    createdAt: "2026-05-20T16:40:02Z"
  },
  {
    id: "usr-4411",
    name: "Kate Brewster",
    username: "kbrewster",
    email: "kbrewster@cyberdyne.org",
    phone: "+1 (555) 016-3023",
    role: "Security Officer",
    status: "Active",
    kycStatus: "Verified",
    department: "Engineering",
    riskScore: 31,
    mfaEnabled: true,
    createdAt: "2026-05-22T13:10:55Z"
  }
];

let auditLogs: AuditLog[] = [
  {
    id: "log-001",
    timestamp: "2026-06-12T10:00:00Z",
    actor: "sarah_connor",
    action: "System Initialization",
    target: "IAM Console Enterprise",
    severity: "Low"
  },
  {
    id: "log-002",
    timestamp: "2026-06-12T14:22:15Z",
    actor: "system_daemon",
    action: "Automated Threat Score Recalculation",
    target: "T-1000 Mimetic",
    severity: "High"
  },
  {
    id: "log-003",
    timestamp: "2026-06-12T14:23:40Z",
    actor: "sarah_connor",
    action: "Identity State Interdiction (Status: Banned)",
    target: "T-1000 Mimetic",
    severity: "Critical"
  },
  {
    id: "log-004",
    timestamp: "2026-06-12T18:05:12Z",
    actor: "marcus_wright",
    action: "KYC Suspicious Flag Appended",
    target: "T-800 Guardian",
    severity: "Medium"
  }
];

// Helper to log administrative actions
function addAuditEntry(actor: string, action: string, target: string, severity: "Critical" | "High" | "Medium" | "Low") {
  const newLog: AuditLog = {
    id: `log-${Math.floor(1000 + Math.random() * 9000)}`,
    timestamp: new Date().toISOString(),
    actor: actor || "anonymous_operator",
    action,
    target,
    severity
  };
  auditLogs.unshift(newLog); // Prepend to show most recent first
}

async function startServer() {
  const app = express();
  const PORT = 3000;

  // JSON Body Parser Middleware
  app.use(express.json());

  // API v1 Endpoints
  
  // 1. GET /api/v1/users - List users with query filtering & pagination
  app.get("/api/v1/users", (req, res) => {
    try {
      const { search, role, status, department, page = "1", limit = "10" } = req.query;

      let filtered = [...users];

      // Contextual search
      if (search) {
        const query = String(search).toLowerCase();
        filtered = filtered.filter(u => 
          u.name.toLowerCase().includes(query) || 
          u.username.toLowerCase().includes(query) || 
          u.email.toLowerCase().includes(query) || 
          u.id.toLowerCase().includes(query)
        );
      }

      // Exact Filters
      if (role) {
        filtered = filtered.filter(u => u.role === role);
      }
      if (status) {
        filtered = filtered.filter(u => u.status === status);
      }
      if (department) {
        filtered = filtered.filter(u => u.department === department);
      }

      // Simple Pagination Calculations
      const pageNum = parseInt(String(page)) || 1;
      const limitNum = parseInt(String(limit)) || 10;
      const totalCount = filtered.length;
      const startIndex = (pageNum - 1) * limitNum;
      const paginatedData = filtered.slice(startIndex, startIndex + limitNum);

      res.status(200).json({
        users: paginatedData,
        pagination: {
          total: totalCount,
          page: pageNum,
          limit: limitNum,
          pages: Math.ceil(totalCount / limitNum)
        }
      });
    } catch (error: any) {
      res.status(500).json({ error: "Failed to query identities", details: error.message });
    }
  });

  // 2. POST /api/v1/users - Enroll a new Subject Principal Identity
  app.post("/api/v1/users", (req, res) => {
    try {
      const { name, email, role, department, operator } = req.body;

      if (!name || !email || !role || !department) {
        return res.status(400).json({ error: "Missing required identity components (name, email, role, department)" });
      }

      const cleanUsername = name.toLowerCase().replace(/[^a-z0-9]/g, "").slice(0, 15);
      const randomId = `usr-${Math.floor(1000 + Math.random() * 9000)}`;

      const newUser: IAMUser = {
        id: randomId,
        name,
        username: cleanUsername,
        email,
        phone: `+1 (555) 01${Math.floor(10 + Math.random() * 89)}-${Math.floor(1000 + Math.random() * 8999)}`,
        role: role as UserRole,
        status: "Active",
        kycStatus: "Pending",
        department,
        riskScore: Math.floor(10 + Math.random() * 30), // Initialized with non-critical risk
        mfaEnabled: false,
        createdAt: new Date().toISOString()
      };

      users.unshift(newUser);

      // Audit Governance Action
      addAuditEntry(
        operator || "operator",
        `Enrolled Subject Principal (Role: ${role})`,
        name,
        "Medium"
      );

      res.status(201).json(newUser);
    } catch (error: any) {
      res.status(500).json({ error: "Failed to enroll identity", details: error.message });
    }
  });

  // 3. PATCH /api/v1/users/:id - Update identity attributes or compliance statuses
  app.patch("/api/v1/users/:id", (req, res) => {
    try {
      const { id } = req.params;
      const { status, kycStatus, riskScore, mfaEnabled, operator } = req.body;

      const userIndex = users.findIndex(u => u.id === id);
      if (userIndex === -1) {
        return res.status(404).json({ error: `Identity with Reference ID ${id} not found` });
      }

      const existingUser = users[userIndex];
      const updatedUser = { ...existingUser };

      const changes: string[] = [];

      if (status !== undefined) {
        updatedUser.status = status as UserStatus;
        changes.push(`Status altered from '${existingUser.status}' to '${status}'`);
      }
      if (kycStatus !== undefined) {
        updatedUser.kycStatus = kycStatus as KYCStatus;
        changes.push(`KYC compliance status updated to '${kycStatus}'`);
      }
      if (riskScore !== undefined) {
        const safeScore = Math.max(0, Math.min(100, Number(riskScore)));
        updatedUser.riskScore = safeScore;
        changes.push(`Risk vector recalculated from ${existingUser.riskScore}% to ${safeScore}%`);
      }
      if (mfaEnabled !== undefined) {
        updatedUser.mfaEnabled = !!mfaEnabled;
        changes.push(`MFA security protection ${mfaEnabled ? "Enabled" : "Disabled"}`);
      }

      users[userIndex] = updatedUser;

      // Log Severity Determination
      let logSeverity: "Critical" | "High" | "Medium" | "Low" = "Low";
      if (status === "Banned" || (riskScore && riskScore > 75)) {
        logSeverity = "Critical";
      } else if (changes.some(c => c.includes("KYC") || c.includes("Risk"))) {
        logSeverity = "High";
      } else if (changes.length > 0) {
        logSeverity = "Medium";
      }

      if (changes.length > 0) {
        addAuditEntry(
          operator || "operator",
          `Attributes Patch: [${changes.join(" | ")}]`,
          existingUser.name,
          logSeverity
        );
      }

      res.status(200).json(updatedUser);
    } catch (error: any) {
      res.status(500).json({ error: "Failed to update identity vectors", details: error.message });
    }
  });

  // 4. DELETE /api/v1/users/:id - Offboard/Decommission an identity
  app.delete("/api/v1/users/:id", (req, res) => {
    try {
      const { id } = req.params;
      const { operator } = req.query;

      const user = users.find(u => u.id === id);
      if (!user) {
        return res.status(404).json({ error: `Identity with Reference ID ${id} not found` });
      }

      // Filter user out
      users = users.filter(u => u.id !== id);

      // Audit Governance Action
      addAuditEntry(
        String(operator || "operator"),
        "Permanent Governance Offboarding (Account Decommissioned)",
        user.name,
        "Critical"
      );

      res.status(200).json({ success: true, message: `Identity ${user.name} offboarded successfully.` });
    } catch (error: any) {
      res.status(500).json({ error: "Decommissioning failure", details: error.message });
    }
  });

  // 5. GET /api/v1/audit-logs - Retrieve high-fidelity governance trails
  app.get("/api/v1/audit-logs", (req, res) => {
    try {
      const { limit = "50", severity } = req.query;

      let filtered = [...auditLogs];

      if (severity) {
        filtered = filtered.filter(log => log.severity.toLowerCase() === String(severity).toLowerCase());
      }

      const limitNum = parseInt(String(limit)) || 50;
      res.status(200).json(filtered.slice(0, limitNum));
    } catch (error: any) {
      res.status(500).json({ error: "Failed to query historical log trail", details: error.message });
    }
  });

  // 6. POST /api/v1/audit-logs - Record a manual system/human action
  app.post("/api/v1/audit-logs", (req, res) => {
    try {
      const { actor, action, target, severity } = req.body;

      if (!action || !target || !severity) {
        return res.status(400).json({ error: "Missing required audit specifications" });
      }

      addAuditEntry(actor, action, target, severity);
      res.status(201).json(auditLogs[0]);
    } catch (error: any) {
      res.status(500).json({ error: "Logging interdiction", details: error.message });
    }
  });

  // 7. GET /api/v1/export/csv - Stream real IAM identity dataset as CSV
  app.get("/api/v1/export/csv", (req, res) => {
    try {
      const { department, role, status } = req.query;

      let filtered = [...users];
      if (department) filtered = filtered.filter(u => u.department === department);
      if (role) filtered = filtered.filter(u => u.role === role);
      if (status) filtered = filtered.filter(u => u.status === status);

      // Construct legal CSV
      const headers = ["Reference ID", "Full Name", "Username", "Email Address", "Phone Number", "Assigned Role", "Account Status", "KYC Compliance Status", "Department", "Threat Risk Score", "MFA Armed Status", "Enrollment Date"];
      
      const rows = filtered.map(u => [
        u.id,
        `"${u.name.replace(/"/g, '""')}"`,
        u.username,
        u.email,
        u.phone,
        u.role,
        u.status,
        u.kycStatus,
        `"${u.department.replace(/"/g, '""')}"`,
        `${u.riskScore}%`,
        u.mfaEnabled ? "Armed" : "Disarmed",
        u.createdAt
      ]);

      const csvContent = [headers.join(","), ...rows.map(r => r.join(","))].join("\n");

      // Log system export activity
      addAuditEntry(
        "security_officer",
        `Secure CSV Export Generation (${filtered.length} Identity Records Captured)`,
        "Identity Directory DB Export",
        "Low"
      );

      res.setHeader("Content-Type", "text/csv");
      res.setHeader("Content-Disposition", "attachment; filename=iam_governance_report.csv");
      res.status(200).send(csvContent);
    } catch (error: any) {
      res.status(500).json({ error: "Failed to stream compiled CSV data", details: error.message });
    }
  });

  // Vite development vs production compiler modes
  if (process.env.NODE_ENV !== "production") {
    const vite = await createViteServer({
      server: { middlewareMode: true },
      appType: "spa"
    });
    app.use(vite.middlewares);
  } else {
    const distPath = path.join(process.cwd(), "dist");
    app.use(express.static(distPath));
    app.get("*", (req, res) => {
      res.sendFile(path.join(distPath, "index.html"));
    });
  }

  app.listen(PORT, "0.0.0.0", () => {
    console.log(`[IAM-CORE] Node.js Service listening persistently on Port ${PORT}`);
  });
}

startServer();
