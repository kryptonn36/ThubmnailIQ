import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as usersApi from "@/api/users";
import type { ListUsersParams } from "@/api/users";

export function useUsers(params: ListUsersParams) {
  return useQuery({
    queryKey: ["users", params],
    queryFn: () => usersApi.listUsers(params),
  });
}

export function useUser(id: string | undefined) {
  return useQuery({
    queryKey: ["users", id],
    queryFn: () => usersApi.getUser(id as string),
    enabled: Boolean(id),
  });
}

export function useUserUploads(id: string | undefined, params: { page?: number; per_page?: number }) {
  return useQuery({
    queryKey: ["users", id, "uploads", params],
    queryFn: () => usersApi.listUserUploads(id as string, params),
    enabled: Boolean(id),
  });
}

function useInvalidateUsers() {
  const queryClient = useQueryClient();
  return (id?: string) => {
    void queryClient.invalidateQueries({ queryKey: ["users"] });
    void queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    if (id) void queryClient.invalidateQueries({ queryKey: ["users", id] });
  };
}

export function useSuspendUser() {
  const invalidate = useInvalidateUsers();
  return useMutation({
    mutationFn: usersApi.suspendUser,
    onSuccess: (_data, id) => invalidate(id),
  });
}

export function useActivateUser() {
  const invalidate = useInvalidateUsers();
  return useMutation({
    mutationFn: usersApi.activateUser,
    onSuccess: (_data, id) => invalidate(id),
  });
}

export function useDeleteUser() {
  const invalidate = useInvalidateUsers();
  return useMutation({
    mutationFn: usersApi.deleteUser,
    onSuccess: (_data, id) => invalidate(id),
  });
}

export function useResetUserPassword() {
  return useMutation({ mutationFn: usersApi.resetUserPassword });
}

export function useChangeUserRole() {
  const invalidate = useInvalidateUsers();
  return useMutation({
    mutationFn: ({ id, workspaceId, role }: { id: string; workspaceId: string; role: string }) =>
      usersApi.changeUserRole(id, workspaceId, role),
    onSuccess: (_data, { id }) => invalidate(id),
  });
}
