"use client";

import { useEffect, useState } from "react";
import { api, ApiError } from "@/lib/api";
import type { TrackingJob, TrackingType } from "@/types";

export default function TrackingPage() {
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
      <div>
        <h1 className="text-2xl font-bold text-white">Tracking</h1>
        <p className="mt-1 text-sm text-gray-400">
          Track channels or keywords to monitor competitor thumbnail changes over time.
        </p>
      </div>

      <form
        onSubmit={handleSubmit}
        className="space-y-4 rounded-2xl border border-surface-300 bg-surface-100 p-6"
      >
        <h2 className="text-base font-semibold text-white">Add Tracking Job</h2>

        <div className="flex gap-2">
          {(["keyword", "channel"] as TrackingType[]).map((t) => (
            <button
              key={t}
              type="button"
              onClick={() => setType(t)}
              className={`rounded-full px-3 py-1.5 text-sm font-medium capitalize transition ${
                type === t
                  ? "bg-brand-600 text-white"
                  : "bg-surface-200 text-gray-400 hover:text-gray-200"
              }`}
            >
              {t}
            </button>
          ))}
        </div>

        {type === "keyword" ? (
          <div>
            <label htmlFor="keyword" className="mb-1 block text-sm font-medium text-gray-300">
              Keyword
            </label>
            <input
              id="keyword"
              type="text"
              required
              value={keyword}
              onChange={(e) => setKeyword(e.target.value)}
              placeholder="e.g. crypto trading tips"
              className="w-full rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
            />
          </div>
        ) : (
          <div>
            <label htmlFor="channelId" className="mb-1 block text-sm font-medium text-gray-300">
              Channel ID
            </label>
            <input
              id="channelId"
              type="text"
              required
              value={channelId}
              onChange={(e) => setChannelId(e.target.value)}
              placeholder="e.g. UC1234567890"
              className="w-full rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
            />
          </div>
        )}

        <div>
          <label htmlFor="interval" className="mb-1 block text-sm font-medium text-gray-300">
            Check interval (hours)
          </label>
          <input
            id="interval"
            type="number"
            min={1}
            max={168}
            value={intervalHours}
            onChange={(e) => setIntervalHours(Number(e.target.value))}
            className="w-full rounded-lg border border-surface-300 bg-surface-200 px-3 py-2 text-sm text-white outline-none focus:border-brand-500"
          />
        </div>

        {error && <p className="text-sm text-red-400">{error}</p>}

        <button
          type="submit"
          disabled={submitting}
          className="rounded-lg bg-brand-gradient px-4 py-2.5 text-sm font-semibold text-white shadow-glow transition hover:opacity-90 disabled:opacity-60"
        >
          {submitting ? "Adding..." : "Add Tracking Job"}
        </button>
      </form>

      <div>
        <h2 className="mb-4 text-lg font-semibold text-white">Tracked Items</h2>
        {loading ? (
          <p className="text-sm text-gray-500">Loading...</p>
        ) : jobs.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-surface-300 p-10 text-center text-sm text-gray-500">
            You&apos;re not tracking anything yet.
          </div>
        ) : (
          <div className="space-y-3">
            {jobs.map((job) => (
              <div
                key={job.id}
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
                  <p className="capitalize">{job.status}</p>
                  <p>
                    Last checked:{" "}
                    {job.last_checked_at
                      ? new Date(job.last_checked_at).toLocaleString()
                      : "Never"}
                  </p>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
