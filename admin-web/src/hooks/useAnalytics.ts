import { useQuery } from "@tanstack/react-query";
import { getAnalytics } from "@/api/analytics";

export function useAnalytics() {
  return useQuery({ queryKey: ["analytics"], queryFn: getAnalytics });
}
