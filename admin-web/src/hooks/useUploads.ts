import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as uploadsApi from "@/api/uploads";
import type { ListUploadsParams } from "@/api/uploads";

export function useUploads(params: ListUploadsParams) {
  return useQuery({
    queryKey: ["uploads", params],
    queryFn: () => uploadsApi.listUploads(params),
  });
}

export function useUpload(id: string | undefined) {
  return useQuery({
    queryKey: ["uploads", id],
    queryFn: () => uploadsApi.getUpload(id as string),
    enabled: Boolean(id),
  });
}

function useInvalidateUploads() {
  const queryClient = useQueryClient();
  return (id?: string) => {
    void queryClient.invalidateQueries({ queryKey: ["uploads"] });
    void queryClient.invalidateQueries({ queryKey: ["dashboard"] });
    if (id) void queryClient.invalidateQueries({ queryKey: ["uploads", id] });
  };
}

export function useDeleteUpload() {
  const invalidate = useInvalidateUploads();
  return useMutation({
    mutationFn: uploadsApi.deleteUpload,
    onSuccess: (_data, id) => invalidate(id),
  });
}

export function useRestoreUpload() {
  const invalidate = useInvalidateUploads();
  return useMutation({
    mutationFn: uploadsApi.restoreUpload,
    onSuccess: (_data, id) => invalidate(id),
  });
}
