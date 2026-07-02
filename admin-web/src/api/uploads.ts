import { api } from "@/lib/axios";
import type { PaginatedResponse, UploadSummary } from "@/types";

export interface ListUploadsParams {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
  include_deleted?: boolean;
}

export function listUploads(params: ListUploadsParams) {
  return api.get<PaginatedResponse<UploadSummary>>("/admin/uploads", { params }).then((r) => r.data);
}

export function getUpload(id: string) {
  return api.get<UploadSummary>(`/admin/uploads/${id}`).then((r) => r.data);
}

export function deleteUpload(id: string) {
  return api.delete(`/admin/uploads/${id}`).then((r) => r.data);
}

export function restoreUpload(id: string) {
  return api.post(`/admin/uploads/${id}/restore`).then((r) => r.data);
}
