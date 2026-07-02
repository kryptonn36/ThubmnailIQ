export interface Admin {
  id: string;
  email: string;
  full_name: string;
  role: string;
  is_active: boolean;
  last_login_at: string | null;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  admin: Admin;
}

export interface UserSummary {
  id: string;
  email: string;
  full_name: string;
  status: "active" | "suspended";
  plan: string;
  analyses_count: number;
  storage_used_bytes: number;
  created_at: string;
  deleted_at: string | null;
}

export interface UserDetail extends UserSummary {
  workspace_id: string | null;
}

export interface UploadSummary {
  id: string;
  workspace_id: string;
  user_id: string;
  keyword: string;
  thumbnail_url: string;
  status: string;
  score: number | null;
  file_size_bytes: number | null;
  created_at: string;
  deleted_at: string | null;
}

export interface SystemHealth {
  database: boolean;
  redis: boolean;
  cv_service: boolean;
}

export interface DashboardStats {
  total_users: number;
  active_users: number;
  total_uploads: number;
  storage_used_bytes: number;
  daily_uploads: number;
  monthly_uploads: number;
  recent_users: UserSummary[];
  recent_uploads: UploadSummary[];
  system_health: SystemHealth;
}

export interface TrendPoint {
  date: string;
  count: number;
}

export interface APIUsageStats {
  total_requests_this_month: number;
  active_keys: number;
}

export interface Analytics {
  daily_user_signups: TrendPoint[];
  monthly_user_signups: TrendPoint[];
  daily_upload_trend: TrendPoint[];
  storage_used_bytes: number;
  top_active_users: UserSummary[];
  file_type_breakdown: Record<string, number>;
  api_usage: APIUsageStats;
}

export interface Settings {
  max_upload_size_bytes: number;
  allowed_extensions: string[];
  feature_flags: Record<string, boolean>;
  storage_provider: string;
  email_provider: string;
  email_from_address: string;
  updated_at: string;
}

export interface AuditLog {
  id: string;
  admin_id: string;
  action: string;
  target_type: string;
  target_id: string;
  metadata: Record<string, unknown> | null;
  created_at: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
}

export interface ApiErrorBody {
  error?: string;
  message?: string;
}
