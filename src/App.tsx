/**
 * @license
 * SPDX-License-Identifier: Apache-2.0
 */

import React, { useState, useEffect } from "react";
import { motion, AnimatePresence } from "motion/react";
import {
  Shield,
  ShieldCheck,
  ShieldAlert,
  Users,
  Search,
  FileSpreadsheet,
  Plus,
  Trash2,
  RefreshCw,
  AlertTriangle,
  UserCheck,
  Building2,
  KeyRound,
  FileText,
  X,
  Gauge,
  Sliders,
  Send,
  HelpCircle,
  ChevronDown,
  ChevronUp,
  History,
  Key,
  ShieldHalf,
  Clock
} from "lucide-react";
import { IAMUser, AuditLog, UserRole, UserStatus, KYCStatus, LogSeverity } from "./types";

export default function App() {
  // 1. RBAC Simulated Context Clearance Level
  const [activeRole, setActiveRole] = useState<UserRole>("Administrator");
  const [operatorId, setOperatorId] = useState<string>("sarah_connor");

  // Simulated JWT secret visualization
  const simulatedToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c3IiOiI${operatorId.substring(0, 5)}Iiwicm9sIjoi${activeRole.replace(" ", "")}IiwiaWF0IjoxNzkyODg5NjAwfQ.sig_hash_enterprise`;

  // 2. Data Directories State
  const [users, setUsers] = useState<IAMUser[]>([]);
  const [auditLogs, setAuditLogs] = useState<AuditLog[]>([]);
  const [stats, setStats] = useState({
    totalUsers: 0,
    criticalRisks: 0,
    kycComplianceRate: 0,
    mfaAdoptionRate: 0
  });

  // 3. User Directories Filtering & Querying Parameters
  const [search, setSearch] = useState("");
  const [selectedRole, setSelectedRole] = useState("");
  const [selectedStatus, setSelectedStatus] = useState("");
  const [selectedDept, setSelectedDept] = useState("");
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [alertMessage, setAlertMessage] = useState<{ type: "success" | "error"; text: string } | null>(null);

  // 4. Modals and Triggers
  const [isEnrollModeOpen, setIsEnrollModeOpen] = useState(false);
  const [isEditModeOpen, setIsEditModeOpen] = useState(false);
  const [selectedUser, setSelectedUser] = useState<IAMUser | null>(null);
  const [expandedUserId, setExpandedUserId] = useState<string | null>(null);
  const [loadingIntelligence, setLoadingIntelligence] = useState<string | null>(null);

  // Enrollment fields
  const [enrollName, setEnrollName] = useState("");
  const [enrollEmail, setEnrollEmail] = useState("");
  const [enrollRole, setEnrollRole] = useState<UserRole>("End User");
  const [enrollDept, setEnrollDept] = useState("Engineering");

  // Editorial attributes fields
  const [patchStatus, setPatchStatus] = useState<UserStatus>("Active");
  const [patchKYC, setPatchKYC] = useState<KYCStatus>("Pending");
  const [patchRisk, setPatchRisk] = useState<number>(30);
  const [patchMFA, setPatchMFA] = useState<boolean>(false);

  // 5. Audit Log Operations filter & custom manual log form
  const [filteredSeverity, setFilteredSeverity] = useState("");
  const [manualAction, setManualAction] = useState("");
  const [manualTarget, setManualTarget] = useState("");
  const [manualSeverity, setManualSeverity] = useState<LogSeverity>("Low");

  const departments = [
    "Security Operations",
    "Engineering",
    "Human Resources",
    "Research & Development",
    "Legal",
    "Finance",
    "Logistics"
  ];

  // Auto-adjust operator identity based on chosen simulated role
  const handleRoleChange = (role: UserRole) => {
    setActiveRole(role);
    if (role === "Administrator") {
      setOperatorId("sarah_connor");
    } else if (role === "Security Officer") {
      setOperatorId("marcus_wright");
    } else {
      setOperatorId("john_connor_read_only");
    }
  };

  // 6. Async Fetch Users from full-stack API
  const fetchDirectories = async () => {
    setIsRefreshing(true);
    try {
      const params = new URLSearchParams();
      if (search) params.append("search", search);
      if (selectedRole) params.append("role", selectedRole);
      if (selectedStatus) params.append("status", selectedStatus);
      if (selectedDept) params.append("department", selectedDept);
      params.append("page", String(currentPage));
      params.append("limit", "8");

      const response = await fetch(`/api/v1/users?${params.toString()}`);
      if (!response.ok) throw new Error("API retrieval interdiction");
      const result = await response.json();
      
      setUsers(result.users);
      setTotalPages(result.pagination.pages);

      // Trigger automatic recalculations for stats banner
      evaluateGlobalStats(result.users);
    } catch (err: any) {
      triggerNotification("error", "Failed to retrieve corporate directory: " + err.message);
    } finally {
      setIsRefreshing(false);
    }
  };

  // 7. Async Fetch Logs
  const fetchAuditLogs = async () => {
    try {
      const params = new URLSearchParams();
      if (filteredSeverity) params.append("severity", filteredSeverity);
      params.append("limit", "20");

      const response = await fetch(`/api/v1/audit-logs?${params.toString()}`);
      if (!response.ok) throw new Error("Failed to query records database");
      const records = await response.json();
      setAuditLogs(records);
    } catch (err: any) {
      console.error(err);
    }
  };

  // Helper to re-calc summary metrics across the loaded list
  const evaluateGlobalStats = (userList: IAMUser[]) => {
    if (!userList.length) return;
    const total = userList.length;
    const highRiskCount = userList.filter(u => u.riskScore > 75).length;
    const compliantCount = userList.filter(u => u.kycStatus === "Verified").length;
    const mfaCount = userList.filter(u => u.mfaEnabled).length;

    setStats({
      totalUsers: total,
      criticalRisks: highRiskCount,
      kycComplianceRate: Math.round((compliantCount / total) * 100),
      mfaAdoptionRate: Math.round((mfaCount / total) * 100)
    });
  };

  // Run effects
  useEffect(() => {
    fetchDirectories();
  }, [search, selectedRole, selectedStatus, selectedDept, currentPage]);

  useEffect(() => {
    fetchAuditLogs();
  }, [filteredSeverity]);

  // Notifications handler triggers
  const triggerNotification = (type: "success" | "error", text: string) => {
    setAlertMessage({ type, text });
    setTimeout(() => setAlertMessage(null), 5000);
  };

  // 8. Enroll Subject Principal (POST /api/v1/users)
  const handleEnrollPrincipal = async (e: React.FormEvent) => {
    e.preventDefault();
    if (activeRole === "End User") {
      triggerNotification("error", "Access Denied: Read-Only Privilege Levels cannot enroll identities");
      return;
    }

    if (!enrollName || !enrollEmail) {
      triggerNotification("error", "Please provide name and corporate email");
      return;
    }

    try {
      const response = await fetch("/api/v1/users", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          name: enrollName,
          email: enrollEmail,
          role: enrollRole,
          department: enrollDept,
          operator: operatorId
        })
      });

      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.error || "Enrollment rejected");
      }

      triggerNotification("success", `Security credentials successfully generated for ${enrollName}`);
      setIsEnrollModeOpen(false);
      // Reset variables
      setEnrollName("");
      setEnrollEmail("");
      setEnrollRole("End User");
      setEnrollDept("Engineering");

      // Refetch core datasets
      fetchDirectories();
      fetchAuditLogs();
    } catch (err: any) {
      triggerNotification("error", err.message);
    }
  };

  // 9. Patch Principal attributes (PATCH /api/v1/users/:id)
  const handleOpenEditPrincipal = (user: IAMUser) => {
    if (activeRole === "End User") {
      triggerNotification("error", "Access Denied: Viewer clearance is read-only");
      return;
    }
    setSelectedUser(user);
    setPatchStatus(user.status);
    setPatchKYC(user.kycStatus);
    setPatchRisk(user.riskScore);
    setPatchMFA(user.mfaEnabled);
    setIsEditModeOpen(true);
  };

  const handlePatchPrincipal = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedUser) return;

    try {
      const response = await fetch(`/api/v1/users/${selectedUser.id}`, {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          status: patchStatus,
          kycStatus: patchKYC,
          riskScore: patchRisk,
          mfaEnabled: patchMFA,
          operator: operatorId
        })
      });

      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.error || "Altering operation rejected");
      }

      triggerNotification("success", `Governance state patched for ${selectedUser.name}`);
      setIsEditModeOpen(false);
      fetchDirectories();
      fetchAuditLogs();
    } catch (err: any) {
      triggerNotification("error", err.message);
    }
  };

  // 10. Decommission/Offboard Principal (DELETE /api/v1/users/:id)
  const handleDecommissionPrincipal = async (id: string, name: string) => {
    if (activeRole !== "Administrator") {
      triggerNotification("error", "Access Denied: Only Administrator role levels can decommission subject accounts!");
      return;
    }

    if (!window.confirm(`Are you absolutely sure you want to permanently decommission and offboard ${name}? This action is immediate and write-locked.`)) {
      return;
    }

    try {
      const response = await fetch(`/api/v1/users/${id}?operator=${operatorId}`, {
        method: "DELETE"
      });

      if (!response.ok) {
        const errData = await response.json();
        throw new Error(errData.error || "Offboarding failure");
      }

      triggerNotification("success", `Identity directory permanently decommissioned for ${name}`);
      fetchDirectories();
      fetchAuditLogs();
    } catch (err: any) {
      triggerNotification("error", err.message);
    }
  };

  const toggleRowExpansion = async (userId: string) => {
    if (expandedUserId === userId) {
      setExpandedUserId(null);
      return;
    }

    setExpandedUserId(userId);
    
    // Check if we already have intelligence for this user
    const user = users.find(u => u.id === userId);
    if (user && user.intelligence) return;

    setLoadingIntelligence(userId);
    
    // Simulate fetching user intelligence (Audit Trail, Login Attempts, Permissions)
    // In a real scenario, this would be an API call like /api/v1/users/:id/intelligence
    setTimeout(() => {
      const mockIntelligence = {
        auditTrail: [
          { id: `log-det-1-${userId}`, timestamp: new Date().toISOString(), actor: "System", action: "Access Token Refresh", target: userId, severity: "Low" as const },
          { id: `log-det-2-${userId}`, timestamp: new Date(Date.now() - 3600000).toISOString(), actor: operatorId, action: "Attribute Inspection", target: userId, severity: "Low" as const },
        ],
        loginAttempts: [
          { id: `auth-1-${userId}`, timestamp: new Date().toISOString(), ipAddress: "182.16.2.45", device: "Chrome / macOS", location: "Singapore", status: "Success" as const },
          { id: `auth-2-${userId}`, timestamp: new Date(Date.now() - 86400000).toISOString(), ipAddress: "182.16.2.45", device: "Chrome / macOS", location: "Singapore", status: "Success" as const },
          { id: `auth-3-${userId}`, timestamp: new Date(Date.now() - 172800000).toISOString(), ipAddress: "103.45.1.12", device: "Safari / iOS", location: "Unknown", status: "Failed" as const },
        ],
        permissions: rolePermissions[users.find(u => u.id === userId)?.role || "End User"]
      };

      setUsers(prev => prev.map(u => u.id === userId ? { ...u, intelligence: mockIntelligence } : u));
      setLoadingIntelligence(null);
    }, 800);
  };

  const rolePermissions: Record<string, string[]> = {
    "Administrator": ["identity.write", "identity.read", "identity.delete", "audit.read", "system.config"],
    "Security Officer": ["identity.read", "identity.patch", "audit.read", "risk.manage"],
    "End User": ["identity.read.self", "profile.edit"]
  };

  // 11. Write custom manual system log entry (POST /api/v1/audit-logs)
  const handleRecordManualLog = async (e: React.FormEvent) => {
    e.preventDefault();
    if (activeRole === "End User") {
      triggerNotification("error", "Access Denied: Viewer role cannot update logs");
      return;
    }

    if (!manualAction || !manualTarget) {
      triggerNotification("error", "Please provide action and target description fields");
      return;
    }

    try {
      const response = await fetch("/api/v1/audit-logs", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          actor: operatorId,
          action: manualAction,
          target: manualTarget,
          severity: manualSeverity
        })
      });

      if (!response.ok) throw new Error("Logger rejection");

      triggerNotification("success", "Governance audit statement successfully generated");
      setManualAction("");
      setManualTarget("");
      setManualSeverity("Low");
      fetchAuditLogs();
    } catch (err: any) {
      triggerNotification("error", "Logging registration block: " + err.message);
    }
  };

  // 12. Trigger Secure CSV Excel Export
  const handleTriggerExport = () => {
    const params = new URLSearchParams();
    if (selectedDept) params.append("department", selectedDept);
    if (selectedRole) params.append("role", selectedRole);
    if (selectedStatus) params.append("status", selectedStatus);

    // This browser trigger actually grabs the download from Express '/api/v1/export/csv'
    window.open(`/api/v1/export/csv?${params.toString()}`);
    triggerNotification("success", "Compiled Identity Directory spreadsheet downloaded successfully.");
    setTimeout(() => {
      fetchAuditLogs();
    }, 1500);
  };

  // Compute real-time compliance metrics breakdown
  const kycVerifiedCount = users.filter((u) => u.kycStatus === "Verified").length;
  const kycPendingCount = users.filter((u) => u.kycStatus === "Pending").length;
  const kycFailedCount = users.filter((u) => u.kycStatus === "Failed").length;
  const kycSuspiciousCount = users.filter((u) => u.kycStatus === "Suspicious").length;

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 font-sans flex flex-col selection:bg-teal-500 selection:text-slate-900" id="main_frame">
      
      {/* 1. Header & Identity Access Switcher */}
      <header className="border-b border-slate-900 bg-slate-950/80 backdrop-blur-md sticky top-0 z-40 px-6 py-4" id="header_container">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-start md:items-center justify-between gap-4">
          
          <div className="flex items-center gap-3">
            <div className="bg-gradient-to-tr from-teal-500 to-indigo-600 p-2.5 rounded-xl shadow-lg shadow-teal-500/10">
              <Shield className="w-6 h-6 text-white" />
            </div>
            <div>
              <h1 className="text-xl font-bold tracking-tight text-white flex items-center gap-2">
                FolksArt <span className="text-teal-400 font-normal">IAM Governance Console</span>
              </h1>
              <p className="text-xs text-slate-400">Enterprise Subject Principal Directories & Audit Ledger</p>
            </div>
          </div>

          {/* Role simulation selector */}
          <div className="bg-slate-900 p-1.5 rounded-xl flex items-center gap-1 border border-slate-800 w-full md:w-auto">
            <span className="text-xs font-mono text-slate-500 px-2 select-none uppercase tracking-wider hidden lg:inline">Clearance Level:</span>
            
            {(["Administrator", "Security Officer", "End User"] as UserRole[]).map((role) => (
              <button
                key={role}
                onClick={() => handleRoleChange(role)}
                className={`text-xs px-3 py-1.5 rounded-lg font-medium transition-all cursor-pointer ${
                  activeRole === role
                    ? "bg-gradient-to-r from-teal-500 to-teal-600 text-slate-950 font-semibold shadow-md shadow-teal-500/10"
                    : "text-slate-400 hover:text-white hover:bg-slate-800"
                }`}
              >
                {role}
              </button>
            ))}
          </div>

        </div>
      </header>

      {/* Security Context & Active Simulated Clearance Info */}
      <section className="bg-slate-900/40 border-b border-slate-900 px-6 py-3 text-xs" id="security_banner">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row md:items-center justify-between gap-3 text-slate-400">
          <div className="flex items-center gap-5 flex-wrap">
            <div className="flex items-center gap-1.5 text-teal-400">
              <KeyRound className="w-3.5 h-3.5" />
              <span className="font-mono">User Principal ID:</span>
              <strong className="font-semibold">{operatorId}</strong>
            </div>
            <div className="flex items-center gap-1.5">
              <span>Security Clearance:</span>
              <span className={`px-2 py-0.5 rounded text-[10px] uppercase tracking-wider font-mono ${
                activeRole === "Administrator" 
                  ? "bg-red-500/10 text-red-400 border border-red-500/20" 
                  : activeRole === "Security Officer" 
                    ? "bg-amber-500/10 text-amber-400 border border-amber-500/20" 
                    : "bg-slate-500/10 text-slate-400 border border-slate-500/20"
              }`}>
                {activeRole}
              </span>
            </div>
          </div>

          <div className="flex items-center gap-2 max-w-full overflow-hidden">
            <span className="text-slate-500 whitespace-nowrap font-mono">Simulated JWT:</span>
            <span className="font-mono bg-slate-950 px-2 py-0.5 rounded text-[10px] text-slate-500 truncate max-w-xs md:max-w-md border border-slate-900" title={simulatedToken}>
              {simulatedToken}
            </span>
          </div>
        </div>
      </section>

      {/* Main Corporate Space */}
      <main className="flex-1 max-w-7xl w-full mx-auto p-6 flex flex-col gap-8">

        {/* Global Notifications system */}
        <AnimatePresence>
          {alertMessage && (
            <motion.div
              initial={{ opacity: 0, y: -20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
              className={`p-4 rounded-xl flex items-center justify-between border shadow-lg ${
                alertMessage.type === "success"
                  ? "bg-teal-950/80 text-teal-200 border-teal-500/20 shadow-teal-500/10"
                  : "bg-red-950/80 text-red-200 border-red-500/20 shadow-red-500/10"
              }`}
            >
              <div className="flex items-center gap-3">
                {alertMessage.type === "success" ? (
                  <ShieldCheck className="w-5 h-5 text-teal-400" />
                ) : (
                  <AlertTriangle className="w-5 h-5 text-red-400" />
                )}
                <span className="text-sm font-medium">{alertMessage.text}</span>
              </div>
              <button onClick={() => setAlertMessage(null)} className="text-slate-400 hover:text-white">
                <X className="w-4 h-4" />
              </button>
            </motion.div>
          )}
        </AnimatePresence>

        {/* 2. Executive Stat Blocks */}
        <section className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4" id="stats_panel">
          
          <div className="bg-slate-900/60 border border-slate-900 p-5 rounded-2xl flex items-center justify-between">
            <div className="flex flex-col gap-1">
              <span className="text-xs font-medium text-slate-400 tracking-wide">Subject Principals</span>
              <span className="text-2xl font-bold tracking-tight text-white">{stats.totalUsers}</span>
              <span className="text-[10px] text-slate-500">Total active directory size</span>
            </div>
            <div className="p-3 bg-slate-950 rounded-xl text-indigo-400 border border-slate-800">
              <Users className="w-5 h-5" />
            </div>
          </div>

          <div className="bg-slate-900/60 border border-slate-900 p-5 rounded-2xl flex items-center justify-between">
            <div className="flex flex-col gap-1">
              <span className="text-xs font-medium text-slate-400 tracking-wide">Threat Vectors</span>
              <span className={`text-2xl font-bold tracking-tight ${stats.criticalRisks > 0 ? "text-red-400" : "text-white"}`}>
                {stats.criticalRisks}
              </span>
              <span className="text-[10px] text-slate-500">Risk profiles &gt; 75% score</span>
            </div>
            <div className={`p-3 bg-slate-950 rounded-xl border border-slate-800 ${stats.criticalRisks > 0 ? "text-red-400" : "text-slate-400"}`}>
              <ShieldAlert className="w-5 h-5" />
            </div>
          </div>

          <div className="bg-slate-900/60 border border-slate-900 p-5 rounded-2xl flex items-center justify-between">
            <div className="flex flex-col gap-1">
              <span className="text-xs font-medium text-slate-400 tracking-wide">KYC Verification Rate</span>
              <span className="text-2xl font-bold tracking-tight text-teal-400">{stats.kycComplianceRate}%</span>
              <span className="text-[10px] text-slate-500">Identity compliance verified</span>
            </div>
            <div className="p-3 bg-slate-950 rounded-xl text-teal-400 border border-slate-800">
              <UserCheck className="w-5 h-5" />
            </div>
          </div>

          <div className="bg-slate-900/60 border border-slate-900 p-5 rounded-2xl flex items-center justify-between">
            <div className="flex flex-col gap-1">
              <span className="text-xs font-medium text-slate-400 tracking-wide">MFA Armed Adoption</span>
              <span className="text-2xl font-bold tracking-tight text-white">{stats.mfaAdoptionRate}%</span>
              <span className="text-[10px] text-slate-500">Two-factor protection coverage</span>
            </div>
            <div className="p-3 bg-slate-950 rounded-xl text-indigo-400 border border-slate-800">
              <ShieldCheck className="w-5 h-5" />
            </div>
          </div>

        </section>

        {/* 3. Splitting Directory List and Audit Sidebar Logs */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          
          {/* Identity directory panel (2/3 width) */}
          <div className="lg:col-span-2 flex flex-col gap-5">
            
            <div className="bg-slate-900/40 border border-slate-900 p-5 rounded-2xl flex flex-col gap-5">
              
              {/* Directory actions top toolbar */}
              <div className="flex flex-col sm:flex-row items-stretch sm:items-center justify-between gap-4">
                <div>
                  <h2 className="text-base font-bold text-white flex items-center gap-2">
                    Corporate Directory
                  </h2>
                  <p className="text-xs text-slate-400">Manage identity configurations, access levels, and risk calculations</p>
                </div>

                <div className="flex items-center gap-2">
                  <button
                    onClick={fetchDirectories}
                    disabled={isRefreshing}
                    className="p-2 text-slate-400 hover:text-white bg-slate-950 hover:bg-slate-900 rounded-xl border border-slate-800 transition-all cursor-pointer flex items-center justify-center disabled:opacity-50"
                    title="Refresh listing"
                  >
                    <RefreshCw className={`w-4 h-4 ${isRefreshing ? "animate-spin" : ""}`} />
                  </button>

                  <button
                    onClick={handleTriggerExport}
                    className="flex items-center gap-1.5 px-3 py-2 bg-slate-950 hover:bg-slate-900 text-xs text-slate-300 font-medium rounded-xl border border-slate-800 transition-all cursor-pointer"
                    title="Export as CSV"
                  >
                    <FileSpreadsheet className="w-4 h-4 text-emerald-400" />
                    <span>Export</span>
                  </button>

                  <button
                    onClick={() => {
                      if (activeRole === "End User") {
                        triggerNotification("error", "Access Denied: Read-only accounts cannot enroll users");
                        return;
                      }
                      setIsEnrollModeOpen(true);
                    }}
                    className={`flex items-center gap-1.5 px-4 py-2 bg-gradient-to-r from-teal-500 to-teal-600 text-xs text-slate-950 font-bold rounded-xl transition-all shadow-md shadow-teal-500/10 cursor-pointer ${
                      activeRole === "End User" ? "opacity-50 cursor-not-allowed" : "hover:brightness-110"
                    }`}
                  >
                    <Plus className="w-4 h-4" />
                    <span>Enroll Principal</span>
                  </button>
                </div>
              </div>

              {/* Filtering block */}
              <div className="grid grid-cols-1 sm:grid-cols-4 gap-3 bg-slate-950 p-4 rounded-xl border border-slate-900">
                
                <div className="relative sm:col-span-1">
                  <span className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none text-slate-500">
                    <Search className="w-3.5 h-3.5" />
                  </span>
                  <input
                    type="text"
                    value={search}
                    onChange={(e) => {
                      setSearch(e.target.value);
                      setCurrentPage(1);
                    }}
                    placeholder="Search identities..."
                    className="w-full bg-slate-900 text-xs text-slate-100 pl-9 pr-3 py-2 rounded-lg border border-slate-800 focus:outline-none focus:border-teal-500 transition-colors placeholder:text-slate-500"
                  />
                </div>

                <div>
                  <select
                    value={selectedRole}
                    onChange={(e) => {
                      setSelectedRole(e.target.value);
                      setCurrentPage(1);
                    }}
                    className="w-full bg-slate-900 text-xs text-slate-300 px-3 py-2 rounded-lg border border-slate-800 focus:outline-none focus:border-teal-500 transition-colors cursor-pointer"
                  >
                    <option value="">All Roles</option>
                    <option value="Administrator">Administrator</option>
                    <option value="Security Officer">Security Officer</option>
                    <option value="End User">End User</option>
                  </select>
                </div>

                <div>
                  <select
                    value={selectedStatus}
                    onChange={(e) => {
                      setSelectedStatus(e.target.value);
                      setCurrentPage(1);
                    }}
                    className="w-full bg-slate-900 text-xs text-slate-300 px-3 py-2 rounded-lg border border-slate-800 focus:outline-none focus:border-teal-500 transition-colors cursor-pointer"
                  >
                    <option value="">All Statuses</option>
                    <option value="Active">Active</option>
                    <option value="Deactivated">Deactivated</option>
                    <option value="Banned">Banned</option>
                  </select>
                </div>

                <div>
                  <select
                    value={selectedDept}
                    onChange={(e) => {
                      setSelectedDept(e.target.value);
                      setCurrentPage(1);
                    }}
                    className="w-full bg-slate-900 text-xs text-slate-300 px-3 py-2 rounded-lg border border-slate-800 focus:outline-none focus:border-teal-500 transition-colors cursor-pointer"
                  >
                    <option value="">All Departments</option>
                    {departments.map(d => (
                      <option key={d} value={d}>{d}</option>
                    ))}
                  </select>
                </div>

              </div>

              {/* Grid or List Table of Identities */}
              <div className="overflow-x-auto rounded-xl border border-slate-900 bg-slate-950">
                <table className="w-full text-left text-xs border-collapse">
                  <thead>
                    <tr className="bg-slate-900/60 border-b border-slate-900 text-slate-400 font-mono">
                      <th className="p-4 w-10"></th>
                      <th className="p-4 font-medium uppercase tracking-wider">Identity Details</th>
                      <th className="p-4 font-medium uppercase tracking-wider">Role & Dept</th>
                      <th className="p-4 font-medium uppercase tracking-wider">Account Status</th>
                      <th className="p-4 font-medium uppercase tracking-wider">KYC Compliance</th>
                      <th className="p-4 font-medium uppercase tracking-wider text-center">Threat Risk</th>
                      <th className="p-4 font-medium uppercase tracking-wider text-right">Actions</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.length === 0 ? (
                      <tr>
                        <td colSpan={6} className="p-8 text-center text-slate-500">
                          <HelpCircle className="w-8 h-8 mx-auto text-slate-600 mb-2" />
                          No subject principals discovered matching the criteria query.
                        </td>
                      </tr>
                    ) : (
                      users.flatMap((user) => {
                        const scoreColor = 
                          user.riskScore > 75 
                            ? "bg-red-500 text-red-500" 
                            : user.riskScore > 40 
                              ? "bg-amber-500 text-amber-500" 
                              : "bg-teal-500 text-teal-400";

                        const isSimulatedActionDisabled = activeRole === "End User";
                        const isExpanded = expandedUserId === user.id;

                        return [
                          <tr key={`${user.id}-main`} className={`border-b border-slate-900/80 hover:bg-slate-900/20 transition-all ${isExpanded ? "bg-slate-900/10" : ""}`}>
                            
                            {/* Expand toggle */}
                            <td className="p-4">
                              <button 
                                onClick={() => toggleRowExpansion(user.id)}
                                className="p-1 text-slate-500 hover:text-teal-400 transition-colors"
                              >
                                {isExpanded ? <ChevronUp className="w-4 h-4" /> : <ChevronDown className="w-4 h-4" />}
                              </button>
                            </td>

                            {/* Personal identifiers details */}
                            <td className="p-4">
                              <div className="flex items-center gap-3">
                                <div className="w-9 h-9 rounded-full bg-slate-900 border border-slate-800 flex items-center justify-center font-bold text-teal-400">
                                  {user.name.split(" ").map(w => w[0]).join("")}
                                </div>
                                <div className="flex flex-col gap-0.5">
                                  <div className="font-semibold text-white flex items-center gap-1.5">
                                    <span>{user.name}</span>
                                    {user.mfaEnabled && (
                                      <ShieldCheck className="w-3.5 h-3.5 text-teal-400" title="MFA Armed Profile" />
                                    )}
                                  </div>
                                  <div className="text-[10px] text-slate-400 font-mono">@{user.username} • {user.id}</div>
                                  <div className="text-[10px] text-slate-500">{user.email}</div>
                                </div>
                              </div>
                            </td>

                            {/* Role representation & Department */}
                            <td className="p-4 vertical-align-middle">
                              <div className="flex flex-col gap-1">
                                <span className="font-semibold text-white">{user.role}</span>
                                <span className="text-[10px] text-slate-400 flex items-center gap-1">
                                  <Building2 className="w-3 h-3 text-slate-500" />
                                  {user.department}
                                </span>
                              </div>
                            </td>

                            {/* Status tags */}
                            <td className="p-4">
                              <span className={`inline-flex items-center gap-1 px-2.5 py-0.5 text-[10px] font-mono font-medium rounded-full ${
                                user.status === "Active"
                                  ? "bg-teal-500/10 text-teal-400 border border-teal-500/20"
                                  : user.status === "Banned"
                                    ? "bg-red-500/10 text-red-400 border border-red-500/20"
                                    : "bg-slate-500/10 text-slate-400 border border-slate-500/20"
                              }`}>
                                <span className={`w-1 h-1 rounded-full ${
                                  user.status === "Active"
                                    ? "bg-teal-400"
                                    : user.status === "Banned"
                                      ? "bg-red-400 animate-pulse"
                                      : "bg-slate-400"
                                }`} />
                                {user.status}
                              </span>
                            </td>

                            {/* KYC statuses info */}
                            <td className="p-4">
                              <span className={`inline-flex px-2 py-0.5 text-[10px] font-mono font-semibold rounded ${
                                user.kycStatus === "Verified"
                                  ? "bg-emerald-500/10 text-emerald-400 border border-emerald-500/20"
                                  : user.kycStatus === "Suspicious"
                                    ? "bg-amber-500/10 text-amber-400 border border-amber-500/20 animate-pulse"
                                    : user.kycStatus === "Failed"
                                      ? "bg-red-500/10 text-red-400 border border-red-500/20"
                                      : "bg-slate-500/10 text-slate-400 border border-slate-500/10"
                              }`}>
                                {user.kycStatus}
                              </span>
                            </td>

                            {/* Threat index slider indicator */}
                            <td className="p-4 text-center">
                              <div className="flex flex-col items-center gap-1 mx-auto max-w-[80px]">
                                <span className={`font-mono text-xs font-semibold ${
                                  user.riskScore > 75 ? "text-red-400" : user.riskScore > 40 ? "text-amber-400" : "text-teal-400"
                                }`}>
                                  {user.riskScore}%
                                </span>
                                <div className="w-full bg-slate-900 rounded-full h-1 overflow-hidden border border-slate-800">
                                  <div className={`h-full rounded-full ${scoreColor.split(" ")[0]}`} style={{ width: `${user.riskScore}%` }}></div>
                                </div>
                              </div>
                            </td>

                            {/* User offboarding controls */}
                            <td className="p-4 text-right">
                              <div className="flex items-center justify-end gap-1.5">
                                <button
                                  onClick={() => handleOpenEditPrincipal(user)}
                                  disabled={isSimulatedActionDisabled}
                                  className={`p-1.5 text-slate-400 hover:text-white bg-slate-900 border border-slate-800 rounded-lg transition-all cursor-pointer ${
                                    isSimulatedActionDisabled ? "opacity-30 cursor-not-allowed" : "hover:bg-slate-800"
                                  }`}
                                  title={isSimulatedActionDisabled ? "Viewer Privileges are read-only" : "Alter account governance metrics"}
                                >
                                  <Sliders className="w-3.5 h-3.5" />
                                </button>
                                
                                <button
                                  onClick={() => handleDecommissionPrincipal(user.id, user.name)}
                                  disabled={activeRole !== "Administrator"}
                                  className={`p-1.5 text-slate-400 hover:text-red-400 bg-slate-900 border border-slate-800 rounded-lg transition-all cursor-pointer ${
                                    activeRole !== "Administrator" ? "opacity-30 cursor-not-allowed" : "hover:bg-red-950/20 hover:border-red-500/30"
                                  }`}
                                  title={activeRole !== "Administrator" ? "Only Administrators can offboard profiles" : "Decommission this identity"}
                                >
                                  <Trash2 className="w-3.5 h-3.5" />
                                </button>
                              </div>
                            </td>

                          </tr>,
                          isExpanded && (
                            <tr key={`${user.id}-expansion`} className="bg-slate-900/30">
                              <td colSpan={7} className="p-0 border-b border-slate-800">
                                <motion.div
                                  initial={{ height: 0, opacity: 0 }}
                                  animate={{ height: "auto", opacity: 1 }}
                                  exit={{ height: 0, opacity: 0 }}
                                  className="overflow-hidden"
                                >
                                  <div className="p-6 grid grid-cols-1 md:grid-cols-3 gap-6">
                                    
                                    {/* Audit Trail Column */}
                                    <div className="flex flex-col gap-3">
                                      <h4 className="text-xs font-bold text-slate-300 flex items-center gap-2 uppercase tracking-wider">
                                        <History className="w-3.5 h-3.5 text-teal-400" />
                                        Subject Audit Trail
                                      </h4>
                                      <div className="bg-slate-950/50 border border-slate-800 rounded-xl p-3 flex flex-col gap-2 min-h-[160px]">
                                        {loadingIntelligence === user.id ? (
                                          <div className="flex-1 flex flex-col items-center justify-center gap-2 text-slate-600 font-mono text-[10px]">
                                            <RefreshCw className="w-4 h-4 animate-spin" />
                                            <span>Querying Logs...</span>
                                          </div>
                                        ) : user.intelligence?.auditTrail.length === 0 ? (
                                          <div className="flex-1 flex items-center justify-center text-slate-600 italic text-[10px]">No recent audit logs found.</div>
                                        ) : (
                                          user.intelligence?.auditTrail.map(log => (
                                            <div key={log.id} className="flex flex-col gap-0.5 border-l border-slate-800 pl-3 pb-2 last:pb-0">
                                              <div className="flex items-center justify-between">
                                                <span className="text-white font-medium">{log.action}</span>
                                                <span className="text-[9px] text-slate-500">{new Date(log.timestamp).toLocaleTimeString()}</span>
                                              </div>
                                              <span className="text-[10px] text-slate-400">Actor: {log.actor}</span>
                                            </div>
                                          ))
                                        )}
                                      </div>
                                    </div>

                                    {/* Login Attempts Column */}
                                    <div className="flex flex-col gap-3">
                                      <h4 className="text-xs font-bold text-slate-300 flex items-center gap-2 uppercase tracking-wider">
                                        <Clock className="w-3.5 h-3.5 text-teal-400" />
                                        Recent Access Logic
                                      </h4>
                                      <div className="bg-slate-950/50 border border-slate-800 rounded-xl p-3 flex flex-col gap-2 min-h-[160px]">
                                        {loadingIntelligence === user.id ? (
                                          <div className="flex-1 flex flex-col items-center justify-center gap-2 text-slate-600 font-mono text-[10px]">
                                            <RefreshCw className="w-4 h-4 animate-spin" />
                                            <span>Analyzing sessions...</span>
                                          </div>
                                        ) : user.intelligence?.loginAttempts.length === 0 ? (
                                          <div className="flex-1 flex items-center justify-center text-slate-600 italic text-[10px]">No login records detected.</div>
                                        ) : (
                                          user.intelligence?.loginAttempts.map(attempt => (
                                            <div key={attempt.id} className="flex items-center justify-between gap-3 p-2 bg-slate-900/50 rounded-lg border border-slate-800/50">
                                              <div className="flex flex-col">
                                                <span className="text-[10px] font-mono text-slate-300">{attempt.ipAddress}</span>
                                                <span className="text-[9px] text-slate-500">{attempt.device} • {attempt.location}</span>
                                              </div>
                                              <span className={`text-[9px] font-bold uppercase rounded px-1.5 py-0.5 ${
                                                attempt.status === "Success" ? "text-teal-400 bg-teal-400/5" : "text-red-400 bg-red-400/5"
                                              }`}>
                                                {attempt.status}
                                              </span>
                                            </div>
                                          ))
                                        )}
                                      </div>
                                    </div>

                                    {/* Permissions Column */}
                                    <div className="flex flex-col gap-3">
                                      <h4 className="text-xs font-bold text-slate-300 flex items-center gap-2 uppercase tracking-wider">
                                        <ShieldHalf className="w-3.5 h-3.5 text-teal-400" />
                                        Active Permission Sets
                                      </h4>
                                      <div className="bg-slate-950/50 border border-slate-800 rounded-xl p-4 flex flex-wrap gap-2 content-start min-h-[160px]">
                                        {loadingIntelligence === user.id ? (
                                          <div className="flex-1 flex flex-col items-center justify-center gap-2 text-slate-600 font-mono text-[10px]">
                                            <RefreshCw className="w-4 h-4 animate-spin" />
                                            <span>Building R-MAP...</span>
                                          </div>
                                        ) : (
                                          user.intelligence?.permissions.map(perm => (
                                            <span key={perm} className="flex items-center gap-1.5 px-2 py-1 bg-teal-500/5 border border-teal-500/10 rounded text-[10px] text-teal-400 font-mono">
                                              <Key className="w-2.5 h-2.5" />
                                              {perm}
                                            </span>
                                          ))
                                        )}
                                      </div>
                                    </div>

                                  </div>
                                </motion.div>
                              </td>
                            </tr>
                          )
                        ];
                      })
                    )}
                  </tbody>
                </table>
              </div>

              {/* Simple pagination block controls */}
              {totalPages > 1 && (
                <div className="flex items-center justify-between border-t border-slate-900 pt-4" id="pagination_bar">
                  <span className="text-xs text-slate-400 font-mono">
                    Page <b className="text-white">{currentPage}</b> of {totalPages}
                  </span>
                  
                  <div className="flex items-center gap-1.5">
                    <button
                      disabled={currentPage === 1}
                      onClick={() => setCurrentPage(prev => Math.max(1, prev - 1))}
                      className="px-3 py-1.5 text-xs text-slate-400 bg-slate-950 border border-slate-900 rounded-lg hover:text-white disabled:opacity-30 disabled:cursor-not-allowed cursor-pointer"
                    >
                      Prev
                    </button>
                    
                    <button
                      disabled={currentPage === totalPages}
                      onClick={() => setCurrentPage(prev => Math.min(totalPages, prev + 1))}
                      className="px-3 py-1.5 text-xs text-slate-400 bg-slate-950 border border-slate-900 rounded-lg hover:text-white disabled:opacity-30 disabled:cursor-not-allowed cursor-pointer"
                    >
                      Next
                    </button>
                  </div>
                </div>
              )}

            </div>

          </div>

          {/* Audit Logging Trails Sidebar (1/3 width) */}
          <div className="flex flex-col gap-6 lg:col-span-1">

            {/* Compliance Summary Circular Progress Widget */}
            <div className="bg-slate-900/40 border border-slate-900 p-5 rounded-2xl flex flex-col gap-5" id="compliance_summary_card">
              <div>
                <h2 className="text-base font-bold text-white flex items-center gap-2">
                  Compliance Summary
                </h2>
                <p className="text-[11px] text-slate-400">KYC verification lifecycle visualization across active entries</p>
              </div>

              <div className="flex flex-col sm:flex-row items-center justify-around gap-6 py-2 bg-slate-950 p-4 rounded-xl border border-slate-900">
                {/* Visual Circular Progress Ring */}
                <div className="relative flex items-center justify-center w-28 h-28">
                  <svg className="w-full h-full transform -rotate-90">
                    {/* Background circle */}
                    <circle
                      cx="56"
                      cy="56"
                      r="46"
                      className="text-slate-800 stroke-current"
                      strokeWidth="8"
                      fill="none"
                    />
                    {/* Foreground compliance ring */}
                    <motion.circle
                      cx="56"
                      cy="56"
                      r="46"
                      className="text-teal-500 stroke-current"
                      strokeWidth="8"
                      strokeDasharray={2 * Math.PI * 46}
                      initial={{ strokeDashoffset: 2 * Math.PI * 46 }}
                      animate={{ strokeDashoffset: (2 * Math.PI * 46) - (stats.kycComplianceRate / 100) * (2 * Math.PI * 46) }}
                      transition={{ duration: 1.2, ease: "easeOut" }}
                      fill="none"
                      strokeLinecap="round"
                    />
                  </svg>
                  <div className="absolute flex flex-col items-center justify-center">
                    <span className="text-2xl font-bold text-teal-400 tracking-tight">{stats.kycComplianceRate}%</span>
                    <span className="text-[9px] font-mono uppercase tracking-widest text-slate-500">Verified</span>
                  </div>
                </div>

                {/* Legend list details */}
                <div className="flex flex-col gap-2 w-full sm:w-auto">
                  <div className="flex items-center justify-between gap-5 text-xs">
                    <div className="flex items-center gap-1.5 text-slate-400">
                      <span className="w-2.5 h-2.5 rounded-full bg-teal-400 inline-block" />
                      <span>Verified</span>
                    </div>
                    <span className="font-mono text-white font-semibold">{kycVerifiedCount}</span>
                  </div>

                  <div className="flex items-center justify-between gap-5 text-xs">
                    <div className="flex items-center gap-1.5 text-slate-400">
                      <span className="w-2.5 h-2.5 rounded-full bg-amber-500 inline-block animate-pulse" />
                      <span>Pending</span>
                    </div>
                    <span className="font-mono text-white font-semibold">{kycPendingCount}</span>
                  </div>

                  <div className="flex items-center justify-between gap-5 text-xs">
                    <div className="flex items-center gap-1.5 text-slate-400">
                      <span className="w-2.5 h-2.5 rounded-full bg-yellow-500 inline-block" />
                      <span>Suspicious</span>
                    </div>
                    <span className="font-mono text-white font-semibold">{kycSuspiciousCount}</span>
                  </div>

                  <div className="flex items-center justify-between gap-5 text-xs">
                    <div className="flex items-center gap-1.5 text-slate-400">
                      <span className="w-2.5 h-2.5 rounded-full bg-red-500 inline-block" />
                      <span>Failed</span>
                    </div>
                    <span className="font-mono text-white font-semibold">{kycFailedCount}</span>
                  </div>
                </div>
              </div>

              {/* Status Indicator text context */}
              <div className="text-[11px] text-slate-400 bg-slate-950/50 p-3 rounded-lg border border-slate-900 flex items-center gap-2">
                <ShieldCheck className="w-4 h-4 text-teal-400 shrink-0" />
                <span>
                  {stats.kycComplianceRate >= 80
                    ? "Healthy security profile. The organization boasts a solid compliance rate above recommended standards."
                    : stats.kycComplianceRate >= 50
                      ? "Moderate compliance coverage. Address pending and suspicious profiles immediately."
                      : "Action Required! High volume of unverified subject principals in current directory."}
                </span>
              </div>
            </div>

            <div className="bg-slate-900/40 border border-slate-900 p-5 rounded-2xl flex flex-col gap-4">
              
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-base font-bold text-white flex items-center gap-2">
                    Action Trail Ledger
                  </h2>
                  <p className="text-[11px] text-slate-400">High-fidelity immutable security logging stream</p>
                </div>
                
                <select
                  value={filteredSeverity}
                  onChange={(e) => setFilteredSeverity(e.target.value)}
                  className="bg-slate-950 text-[10px] font-mono text-slate-300 px-2 py-1 rounded border border-slate-800 focus:outline-none cursor-pointer"
                >
                  <option value="">All</option>
                  <option value="Critical">Critical</option>
                  <option value="High">High</option>
                  <option value="Medium">Medium</option>
                  <option value="Low">Low</option>
                </select>
              </div>

              {/* Real system logs list container details */}
              <div className="flex flex-col gap-2 max-h-[300px] overflow-y-auto pr-1 bg-slate-950 p-3 rounded-xl border border-slate-900 font-mono text-[11px]">
                {auditLogs.length === 0 ? (
                  <span className="text-slate-600 text-center py-6">No historical actions logged in database.</span>
                ) : (
                  auditLogs.map((log) => {
                    const severityClass = 
                      log.severity === "Critical" 
                        ? "text-red-400 bg-red-400/5 border border-red-500/20" 
                        : log.severity === "High" 
                          ? "text-amber-400 bg-amber-400/5 border border-amber-500/20" 
                          : log.severity === "Medium" 
                            ? "text-yellow-300 bg-yellow-400/5 border border-yellow-500/10" 
                            : "text-slate-400 bg-slate-900/40 border border-slate-800";

                    return (
                      <div key={log.id} className="p-2.5 rounded-lg bg-slate-900/60 border border-slate-900 flex flex-col gap-1 hover:border-slate-800 transition-all">
                        <div className="flex items-center justify-between gap-2">
                          <span className="text-[10px] text-slate-500 font-mono">
                            {new Date(log.timestamp).toLocaleTimeString()}
                          </span>
                          <span className={`px-1.5 py-0.2 rounded text-[9px] uppercase tracking-wider font-semibold ${severityClass}`}>
                            {log.severity}
                          </span>
                        </div>
                        <div className="text-slate-200">{log.action}</div>
                        <div className="flex items-center justify-between text-[10px] text-slate-400">
                          <span>Target: <b className="text-slate-300 font-normal">{log.target}</b></span>
                          <span className="text-slate-500">@{log.actor}</span>
                        </div>
                      </div>
                    );
                  })
                )}
              </div>

              {/* Log Registration Form (RBAC locked for End Users) */}
              <div className="bg-slate-950 p-4 rounded-xl border border-slate-900 flex flex-col gap-3">
                <span className="text-xs font-semibold text-slate-300 flex items-center gap-1.5">
                  <FileText className="w-3.5 h-3.5 text-teal-400" />
                  Record Governance Log
                </span>
                
                <form onSubmit={handleRecordManualLog} className="flex flex-col gap-2.5">
                  <input
                    type="text"
                    value={manualAction}
                    onChange={(e) => setManualAction(e.target.value)}
                    disabled={activeRole === "End User"}
                    placeholder="Action description (e.g. Identity Review)"
                    className="w-full bg-slate-900 border border-slate-800 px-3 py-1.5 rounded text-[11px] focus:outline-none focus:border-teal-500 text-white placeholder:text-slate-500 disabled:opacity-40"
                  />
                  <input
                    type="text"
                    value={manualTarget}
                    onChange={(e) => setManualTarget(e.target.value)}
                    disabled={activeRole === "End User"}
                    placeholder="Target scope (e.g. John Connor)"
                    className="w-full bg-slate-900 border border-slate-800 px-3 py-1.5 rounded text-[11px] focus:outline-none focus:border-teal-500 text-white placeholder:text-slate-500 disabled:opacity-40"
                  />
                  
                  <div className="flex items-center gap-2">
                    <select
                      value={manualSeverity}
                      onChange={(e) => setManualSeverity(e.target.value as LogSeverity)}
                      disabled={activeRole === "End User"}
                      className="flex-1 bg-slate-900 border border-slate-800 px-2 py-1.5 rounded text-[11px] text-slate-300 focus:outline-none cursor-pointer disabled:opacity-40"
                    >
                      <option value="Low">Low Severity</option>
                      <option value="Medium">Medium Severity</option>
                      <option value="High">High Severity</option>
                      <option value="Critical">Critical Severity</option>
                    </select>

                    <button
                      type="submit"
                      disabled={activeRole === "End User"}
                      className={`px-4 py-1.5 rounded text-slate-950 font-bold text-[11px] bg-gradient-to-r from-teal-500 to-teal-600 transition-all flex items-center gap-1 cursor-pointer disabled:opacity-40`}
                    >
                      <Send className="w-3 h-3" />
                      <span>Emit</span>
                    </button>
                  </div>
                </form>
              </div>

            </div>

          </div>

        </div>

      </main>

      {/* FOOTER */}
      <footer className="border-t border-slate-900 mt-12 bg-slate-950/20 py-6 px-6 text-center text-xs text-slate-500" id="footer">
        <div className="max-w-7xl mx-auto flex flex-col md:flex-row items-center justify-between gap-4">
          <span>FolksArt Web Console Enterprise Edition. All Rights Reserved.</span>
          <span className="font-mono text-slate-600">Secure Vault Session: {simulatedToken.slice(0, 15)}...</span>
        </div>
      </footer>

      {/* 4. ENROLLMENT DOCK DIALOGUE - MODAL */}
      <AnimatePresence>
        {isEnrollModeOpen && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            {/* Backdrop opacity layer */}
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsEnrollModeOpen(false)}
              className="absolute inset-0 bg-slate-950/70 backdrop-blur-sm cursor-pointer"
            />

            {/* Modal Body container */}
            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 15 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 15 }}
              className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative z-10 shadow-2xl"
              id="enrollment_modal"
            >
              
              <div className="flex items-center justify-between border-b border-slate-800 pb-4 mb-4">
                <div className="flex items-center gap-2">
                  <Users className="w-5 h-5 text-teal-400" />
                  <h3 className="text-sm font-bold text-white">Enroll Corporate Identity</h3>
                </div>
                <button
                  onClick={() => setIsEnrollModeOpen(false)}
                  className="text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition-colors"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>

              <form onSubmit={handleEnrollPrincipal} className="flex flex-col gap-4">
                
                <div className="flex flex-col gap-1.5">
                  <label className="text-[11px] font-mono text-slate-400 uppercase tracking-wider">Full Name</label>
                  <input
                    type="text"
                    required
                    value={enrollName}
                    onChange={(e) => setEnrollName(e.target.value)}
                    placeholder="E.g. Sarah Connor"
                    className="bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs focus:outline-none focus:border-teal-500 text-white whitespace-pre-line"
                  />
                </div>

                <div className="flex flex-col gap-1.5">
                  <label className="text-[11px] font-mono text-slate-400 uppercase tracking-wider">Corporate Email</label>
                  <input
                    type="email"
                    required
                    value={enrollEmail}
                    onChange={(e) => setEnrollEmail(e.target.value)}
                    placeholder="sconnor@cyberdyne.org"
                    className="bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs focus:outline-none focus:border-teal-500 text-white"
                  />
                </div>

                <div className="grid grid-cols-2 gap-3">
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[11px] font-mono text-slate-400 uppercase tracking-wider">Role Designation</label>
                    <select
                      value={enrollRole}
                      onChange={(e) => setEnrollRole(e.target.value as UserRole)}
                      className="bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-slate-200 focus:outline-none cursor-pointer"
                    >
                      <option value="End User">End User</option>
                      <option value="Security Officer">Security Officer</option>
                      <option value="Administrator">Administrator</option>
                    </select>
                  </div>

                  <div className="flex flex-col gap-1.5">
                    <label className="text-[11px] font-mono text-slate-400 uppercase tracking-wider">Department</label>
                    <select
                      value={enrollDept}
                      onChange={(e) => setEnrollDept(e.target.value)}
                      className="bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-slate-200 focus:outline-none cursor-pointer"
                    >
                      {departments.map((d) => (
                        <option key={d} value={d}>
                          {d}
                        </option>
                      ))}
                    </select>
                  </div>
                </div>

                <div className="bg-slate-950 p-3 rounded-lg border border-slate-800/60 mt-2 flex items-start gap-2.5">
                  <ShieldCheck className="w-4 h-4 text-teal-400 mt-0.5 shrink-0" />
                  <div className="flex flex-col gap-0.5">
                    <span className="text-[11px] font-semibold text-slate-200">Baseline Registration Criteria</span>
                    <span className="text-[10px] text-slate-400">The account will be initially activated in "Active" state. KYC compliance starts as "Pending" pending evaluation.</span>
                  </div>
                </div>

                <button
                  type="submit"
                  className="w-full bg-gradient-to-r from-teal-500 to-teal-600 hover:brightness-110 text-slate-950 font-bold py-2.5 rounded-lg text-xs mt-3 flex items-center justify-center gap-1 opacity-100 transition-all cursor-pointer"
                >
                  <UserCheck className="w-4 h-4" />
                  <span>Verify and Enroll User</span>
                </button>

              </form>

            </motion.div>
          </div>
        )}
      </AnimatePresence>

      {/* 5. PATCH GOVERNANCE ATTRIBUTES DIALOGUE - MODAL */}
      <AnimatePresence>
        {isEditModeOpen && selectedUser && (
          <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
            <motion.div
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              exit={{ opacity: 0 }}
              onClick={() => setIsEditModeOpen(false)}
              className="absolute inset-0 bg-slate-950/70 backdrop-blur-sm cursor-pointer"
            />

            <motion.div
              initial={{ opacity: 0, scale: 0.95, y: 15 }}
              animate={{ opacity: 1, scale: 1, y: 0 }}
              exit={{ opacity: 0, scale: 0.95, y: 15 }}
              className="bg-slate-900 border border-slate-800 rounded-2xl max-w-md w-full p-6 relative z-10 shadow-2xl"
              id="patch_modal"
            >
              
              <div className="flex items-center justify-between border-b border-slate-800 pb-4 mb-4">
                <div className="flex items-center gap-2">
                  <Sliders className="w-5 h-5 text-teal-400" />
                  <h3 className="text-sm font-bold text-white">Alter Governance Attributes</h3>
                </div>
                <button
                  onClick={() => setIsEditModeOpen(false)}
                  className="text-slate-400 hover:text-white p-1 rounded-lg hover:bg-slate-800 transition-colors"
                >
                  <X className="w-4 h-4" />
                </button>
              </div>

              <div className="mb-4 bg-slate-950 p-3 rounded-lg border border-slate-900 flex items-center gap-3">
                <div className="w-8 h-8 rounded-full bg-slate-800 flex items-center justify-center font-bold text-teal-400 text-xs">
                  {selectedUser.name.split(" ").map(w => w[0]).join("")}
                </div>
                <div>
                  <h4 className="text-xs font-bold text-white">{selectedUser.name}</h4>
                  <p className="text-[10px] text-slate-500 font-mono">{selectedUser.id} • {selectedUser.department}</p>
                </div>
              </div>

              <form onSubmit={handlePatchPrincipal} className="flex flex-col gap-4">
                
                <div className="grid grid-cols-2 gap-3">
                  <div className="flex flex-col gap-1.5">
                    <label className="text-[11px] font-mono text-slate-400 uppercase tracking-wider">Account Status</label>
                    <select
                      value={patchStatus}
                      onChange={(e) => setPatchStatus(e.target.value as UserStatus)}
                      className="bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-slate-200 focus:outline-none cursor-pointer"
                    >
                      <option value="Active">Active</option>
                      <option value="Deactivated">Deactivated</option>
                      <option value="Banned">Banned</option>
                    </select>
                  </div>

                  <div className="flex flex-col gap-1.5">
                    <label className="text-[11px] font-mono text-slate-400 uppercase tracking-wider">Compliance Status</label>
                    <select
                      value={patchKYC}
                      onChange={(e) => setPatchKYC(e.target.value as KYCStatus)}
                      className="bg-slate-950 border border-slate-800 rounded-lg px-3 py-2 text-xs text-slate-200 focus:outline-none cursor-pointer"
                    >
                      <option value="Pending">Pending</option>
                      <option value="Verified">Verified</option>
                      <option value="Failed">Failed</option>
                      <option value="Suspicious">Suspicious</option>
                    </select>
                  </div>
                </div>

                <div className="flex flex-col gap-1.5">
                  <div className="flex items-center justify-between text-[11px] font-mono text-slate-400 uppercase tracking-wider">
                    <span>Threat Vector Risk Score</span>
                    <span className={`font-mono text-xs font-bold ${
                      patchRisk > 75 ? "text-red-400" : patchRisk > 40 ? "text-amber-400" : "text-teal-400"
                    }`}>{patchRisk}%</span>
                  </div>
                  
                  <div className="flex items-center gap-3">
                    <input
                      type="range"
                      min="0"
                      max="100"
                      value={patchRisk}
                      onChange={(e) => setPatchRisk(Number(e.target.value))}
                      className="flex-1 accent-teal-500 h-1 bg-slate-950 rounded-lg cursor-pointer"
                    />
                  </div>
                </div>

                {/* MFA check toggle */}
                <div className="flex items-center justify-between bg-slate-950 p-3 rounded-lg border border-slate-800/60 mt-1">
                  <div className="flex items-start gap-2">
                    <Shield className="w-4 h-4 text-indigo-400 mt-0.5 shrink-0" />
                    <div className="flex flex-col">
                      <span className="text-[11px] font-semibold text-slate-200">Two-Factor Authentication</span>
                      <span className="text-[10px] text-slate-500">Require MFA for session enrollment</span>
                    </div>
                  </div>
                  
                  <input
                    type="checkbox"
                    checked={patchMFA}
                    onChange={(e) => setPatchMFA(e.target.checked)}
                    className="w-4 h-4 text-teal-600 bg-slate-900 border-slate-800 rounded focus:ring-teal-500 accent-teal-500 cursor-pointer"
                  />
                </div>

                <button
                  type="submit"
                  className="w-full bg-gradient-to-r from-teal-500 to-teal-600 hover:brightness-110 text-slate-950 font-bold py-2.5 rounded-lg text-xs mt-3 flex items-center justify-center gap-1 opacity-100 transition-all cursor-pointer"
                >
                  <Gauge className="w-4 h-4" />
                  <span>Update Principal State</span>
                </button>

              </form>

            </motion.div>
          </div>
        )}
      </AnimatePresence>

    </div>
  );
}
