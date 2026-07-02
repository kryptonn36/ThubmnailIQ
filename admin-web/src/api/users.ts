import { api } from "@/lib/axios";
import type { PaginatedResponse, UploadSummary, UserDetail, UserSummary } from "@/types";

export interface ListUsersParams {
  page?: number;
  per_page?: number;
  search?: string;
  status?: string;
}

export function listUsers(params: ListUsersParams) {
  return api.get<PaginatedResponse<UserSummary>>("/admin/users", { params }).then((r) => r.data);
}

export function getUser(id: string) {
  return api.get<UserDetail>(`/admin/users/${id}`).then((r) => r.data);
}

export function suspendUser(id: string) {
  return api.patch(`/admin/users/${id}/suspend`).then((r) => r.data);
}

export function activateUser(id: string) {
  return api.patch(`/admin/users/${id}/activate`).then((r) => r.data);
}

export function deleteUser(id: string) {
  return api.delete(`/admin/users/${id}`).then((r) => r.data);
}

export function resetUserPassword(id: string) {
  return api
    .post<{ temporary_password: string }>(`/admin/users/${id}/reset-password`)
    .then((r) => r.data);
}

export function changeUserRole(id: string, workspaceId: string, role: string) {
  return api
    .patch(`/admin/users/${id}/role`, { workspace_id: workspaceId, role })
    .then((r) => r.data);
}

export function listUserUploads(id: string, params: { page?: number; per_page?: number }) {
  return api
    .get<PaginatedResponse<UploadSummary>>(`/admin/users/${id}/uploads`, { params })
    .then((r) => r.data);
}
