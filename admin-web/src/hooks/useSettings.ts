import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import * as settingsApi from "@/api/settings";
import type { Settings } from "@/types";

export function useSettings() {
  return useQuery({ queryKey: ["settings"], queryFn: settingsApi.getSettings });
}

export function useUpdateSettings() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (settings: Settings) => settingsApi.updateSettings(settings),
    onSuccess: (updated) => {
      queryClient.setQueryData(["settings"], updated);
    },
  });
}
