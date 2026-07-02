import { Link, useNavigate, useParams } from "react-router-dom";
import toast from "react-hot-toast";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { useDeleteUpload, useRestoreUpload, useUpload } from "@/hooks/useUploads";
import { extractErrorMessage } from "@/lib/axios";

function formatBytes(bytes: number | null): string {
  if (!bytes) return "—";
  const units = ["B", "KB", "MB", "GB"];
  const i = Math.min(units.length - 1, Math.floor(Math.log(bytes) / Math.log(1024)));
  return `${(bytes / 1024 ** i).toFixed(1)} ${units[i]}`;
}

export function UploadDetailPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data: upload, isLoading, error, refetch } = useUpload(id);
  const del = useDeleteUpload();
  const restore = useRestoreUpload();

  if (isLoading) return <p className="text-sm text-gray-500">Loading…</p>;
  if (error || !upload) {
    return (
      <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">
        {error ? extractErrorMessage(error) : "Upload not found."}
      </p>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-white">{upload.keyword}</h1>
          <p className="mt-1 text-sm text-gray-400">Upload ID: {upload.id}</p>
        </div>
        <Badge status={upload.deleted_at ? "deleted" : upload.status} />
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
        <Card className="lg:col-span-1">
          <img
            src={upload.thumbnail_url}
            alt={upload.keyword}
            className="aspect-video w-full rounded-lg object-cover"
          />
        </Card>
        <Card className="lg:col-span-2">
          <dl className="grid grid-cols-2 gap-4 text-sm">
            <div>
              <dt className="text-gray-500">Score</dt>
              <dd className="mt-1 font-semibold text-white">{upload.score ?? "—"}</dd>
            </div>
            <div>
              <dt className="text-gray-500">File size</dt>
              <dd className="mt-1 font-semibold text-white">{formatBytes(upload.file_size_bytes)}</dd>
            </div>
            <div>
              <dt className="text-gray-500">Created</dt>
              <dd className="mt-1 font-semibold text-white">
                {new Date(upload.created_at).toLocaleString()}
              </dd>
            </div>
            <div>
              <dt className="text-gray-500">Workspace</dt>
              <dd className="mt-1 truncate font-mono text-xs text-gray-300">{upload.workspace_id}</dd>
            </div>
            <div>
              <dt className="text-gray-500">Owner user</dt>
              <dd className="mt-1">
                <Link to={`/users/${upload.user_id}`} className="font-mono text-xs text-brand-500 hover:underline">
                  {upload.user_id}
                </Link>
              </dd>
            </div>
          </dl>

          <div className="mt-6 flex flex-wrap gap-3 border-t border-surface-300 pt-4">
            <a href={upload.thumbnail_url} target="_blank" rel="noreferrer">
              <Button variant="secondary" type="button">
                Download
              </Button>
            </a>
            {upload.deleted_at ? (
              <Button
                variant="secondary"
                loading={restore.isPending}
                onClick={() =>
                  restore.mutate(upload.id, {
                    onSuccess: () => {
                      toast.success("Upload restored");
                      void refetch();
                    },
                    onError: (err) => toast.error(extractErrorMessage(err)),
                  })
                }
              >
                Restore
              </Button>
            ) : (
              <Button
                variant="danger"
                loading={del.isPending}
                onClick={() => {
                  if (!confirm("Delete this upload?")) return;
                  del.mutate(upload.id, {
                    onSuccess: () => {
                      toast.success("Upload deleted");
                      void refetch();
                    },
                    onError: (err) => toast.error(extractErrorMessage(err)),
                  });
                }}
              >
                Delete
              </Button>
            )}
          </div>
        </Card>
      </div>

      <button
        type="button"
        onClick={() => navigate(-1)}
        className="inline-block text-sm text-brand-500 hover:underline"
      >
        ← Back
      </button>
    </div>
  );
}
