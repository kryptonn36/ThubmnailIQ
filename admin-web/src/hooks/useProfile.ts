import { useMutation, useQuery } from "@tanstack/react-query";
import * as profileApi from "@/api/profile";

export function useProfile() {
  return useQuery({ queryKey: ["profile"], queryFn: profileApi.getProfile });
}

export function useChangePassword() {
  return useMutation({
    mutationFn: ({ currentPassword, newPassword }: { currentPassword: string; newPassword: string }) =>
      profileApi.changePassword(currentPassword, newPassword),
  });
}
