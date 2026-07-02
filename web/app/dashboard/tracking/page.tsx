"use client";

import { useEffect, useState } from "react";
import { motion } from "framer-motion";
import { Radar } from "lucide-react";
import { api, ApiError } from "@/lib/api";
import { Alert } from "@/components/ui/Alert";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { EmptyState } from "@/components/ui/EmptyState";
import { Input } from "@/components/ui/Input";
import { PageHeader } from "@/components/ui/PageHeader";
import { Skeleton } from "@/components/ui/Skeleton";
import { useMotionVariants } from "@/lib/motion";
import type { TrackingJob, TrackingType } from "@/types";

export default function TrackingPage() {
  const { fadeInUp, staggerContainer } = useMotionVariants();
  const [jobs, setJobs] = useState<TrackingJob[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [type, setType] = useState<TrackingType>("keyword");
  const [channelId, setChannelId] = useState("");
  const [keyword, setKeyword] = useState("");
  const [intervalHours, setIntervalHours] = useState(24);
  const [submitting, setSubmitting] = useState(false);

  async function loadJobs() {
    setLoading(true);
    try {
      const data = await api.get<TrackingJob[]>("/tracking");
      setJobs(data);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to load tracking jobs.");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    loadJobs();
  }, []);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      await api.post("/tracking", {
        type,
        channel_id: type === "channel" ? channelId.trim() : undefined,
        keyword: type === "keyword" ? keyword.trim() : undefined,
        interval_hours: intervalHours,
      });
      setChannelId("");
      setKeyword("");
      await loadJobs();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "Failed to create tracking job.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="space-y-8">
      <PageHeader
        title="Tracking"
        description="Track channels or keywords to monitor competitor thumbnail changes over time."
      />

      <Card as="form" onSubmit={handleSubmit} className="space-y-4">
        <h2 className="text-base font-semibold text-white">Add Tracking Job</h2>

        <div className="flex gap-2">
          {(["keyword", "channel"] as TrackingType[]).map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => setType(t)}
              className={`rounded-full px-3 py-1.5 text-sm font-medium capitalize transition-colors duration-150 ${
                type === t ? "bg-brand-600 text-white" : "bg-surface-200 text-gray-400 hover:text-gray-200"
              }`}
            >
              {t}
            </button>
          ))}
        </div>

        {type === "keyword" ? (
          <div>
            <label htmlFor="keyword" className="mb-1.5 block text-sm font-medium text-gray-300">
              Keyword
            </label>
            <Input
              id="keyword"
              type="text"
              required
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              placeholder="e.g. crypto trading tips"
            />
          </div>
        ) : (
          <div>
            <label htmlFor="channelId" className="mb-1.5 block text-sm font-medium text-gray-300">
              Channel ID
            </label>
            <Input
              id="channelId"
              type="text"
              required
              value={channelId}
              onChange={(e) => setChannelId(e.target.value)}
              placeholder="e.g. UC1234567890"
            />
          </div>
        )}

        <div>
          <label htmlFor="interval" className="mb-1.5 block text-sm font-medium text-gray-300">
            Check interval (hours)
          </label>
          <Input
            id="interval"
            type="number"
            min={1}
            max={168}
            value={intervalHours}
            onChange={(e) => setIntervalHours(Number(e.target.value))}
          />
        </div>

        {error && <Alert variant="danger">{error}</Alert>}

        <Button type="submit" loading={submitting}>
          Add Tracking Job
        </Button>
      </Card>

      <div>
        <h2 className="mb-4 text-lg font-semibold text-white">Tracked Items</h2>
        {loading ? (
          <div className="space-y-3">
            {[1, 2].map((i) => (
              <Skeleton key={i} className="h-16 rounded-xl" />
            ))}
          </div>
        ) : jobs.length === 0 ? (
          <EmptyState
            icon={<Radar className="h-5 w-5" aria-hidden="true" />}
            title="You're not tracking anything yet."
          />
        ) : (
          <motion.div initial="hidden" animate="visible" variants={staggerContainer} className="space-y-3">
            {jobs.map((job) => (
              <motion.div
                key={job.id}
                variants={fadeInUp}
                className="flex flex-wrap items-center justify-between gap-3 rounded-xl border border-surface-300 bg-surface-100 p-4"
              >
                <div>
                  <p className="font-medium text-white">
                    {job.type === "keyword" ? job.keyword : job.channel_id}
                  </p>
                  <p className="text-xs capitalize text-gray-500">
                    {job.type} &middot; every {job.interval_hours}h
                  </p>
                </div>
                <div className="text-right text-xs text-gray-500">
                  <Badge variant={job.status === "active" ? "success" : "neutral"}>{job.status}</Badge>
                  <p className="mt-1">
                    Last checked: {job.last_checked_at ? new Date(job.last_checked_at).toLocaleString() : "Never"}
                  </p>
                </div>
              </motion.div>
            ))}
          </motion.div>
        )}
      </div>
    </div>
  );
}
