// Core TypeScript interfaces matching the ThumbnailIQ API contract.

export interface User {
  id: string;
  email: string;
  full_name: string;
  avatar_url?: string | null;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  expires_in: number;
  user: User;
}

export interface Workspace {
  id: string;
  name: string;
  slug: string;
  plan: string;
  analyses_this_month: number;
  analyses_limit: number;
}

export type AnalysisStatus = "pending" | "processing" | "complete" | "failed";

export interface AnalysisSummary {
  id: string;
  keyword: string;
  thumbnail_url: string;
  status: AnalysisStatus;
  score: number | null;
  created_at: string;
}

export interface AnalysisCreateResponse {
  id: string;
  status: AnalysisStatus;
  keyword: string;
  thumbnail_url: string;
  created_at: string;
}

export interface DominantColor {
  hex: string;
  percentage: number;
}

export interface CVResults {
  face_count: number;
  primary_emotion: string | null;
  text_detected: boolean;
  text_strings: string[];
  text_density_pct: number;
  word_count: number;
  dominant_colors: DominantColor[];
  contrast_score: number;
  brightness_score: number;
  saturation_score: number;
  object_count: number;
  clutter_score: number;
  visual_complexity: number;
}

export interface CompetitorAvg {
  score: number;
  face_count: number;
  text_density_pct: number;
  clutter_score: number;
}

export type ImpactLevel = "high" | "medium" | "low";

export interface Suggestion {
  type: string;
  impact_level: ImpactLevel;
  headline: string;
  explanation: string;
  estimated_ctr_min: number;
  estimated_ctr_max: number;
}

export interface Competitor {
  video_id: string;
  video_title: string;
  channel_name: string;
  thumbnail_url: string;
  view_count: number;
  rank_position: number;
  score: number;
}

export interface ThumbnailVersion {
  id: string;
  version_number: number;
  thumbnail_url: string;
  score: number;
  cv_results: CVResults;
}

export interface AnalysisFull {
  id: string;
  status: AnalysisStatus;
  keyword: string;
  thumbnail_url: string;
  created_at: string;
  score: number | null;
  visibility_score: number | null;
  contrast_score: number | null;
  attention_score: number | null;
  mobile_score: number | null;
  branding_score: number | null;
  curiosity_score: number | null;
  cv_results: CVResults | null;
  competitor_avg: CompetitorAvg | null;
  suggestions: Suggestion[];
  competitor_count: number;
  rank_in_competitors: number | null;
  competitors: Competitor[];
  versions?: ThumbnailVersion[];
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  per_page: number;
  total: number;
}

export type TrackingType = "channel" | "keyword";
export type TrackingStatus = "active" | "paused" | "failed";

export interface TrackingJob {
  id: string;
  type: TrackingType;
  channel_id?: string | null;
  keyword?: string | null;
  status: TrackingStatus;
  interval_hours: number;
  last_checked_at: string | null;
  next_check_at: string | null;
}

export interface ViralThumbnail {
  id: string;
  video_title: string;
  channel_name: string;
  thumbnail_url: string;
  niche: string;
  score: number;
  view_count: number;
}

export interface BillingPlan {
  id: string;
  name: string;
  price_monthly: number;
  features: string[];
  analyses_limit: number;
  api_requests_limit: number;
}

export interface SubscribeResponse {
  plan: string;
  status: string;
}
