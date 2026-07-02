import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect } from "react";
import { useForm } from "react-hook-form";
import toast from "react-hot-toast";
import { z } from "zod";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { FormField, inputClasses } from "@/components/forms/FormField";
import { useSettings, useUpdateSettings } from "@/hooks/useSettings";
import { extractErrorMessage } from "@/lib/axios";

const settingsSchema = z.object({
  max_upload_size_bytes: z.coerce.number().int().positive("Must be greater than zero"),
  allowed_extensions: z.string().min(1, "List at least one extension"),
  storage_provider: z.string().min(1, "Required"),
  email_provider: z.string(),
  email_from_address: z.string(),
  feature_flags_json: z.string().refine(
    (value) => {
      try {
        const parsed = JSON.parse(value) as unknown;
        return (
          typeof parsed === "object" &&
          parsed !== null &&
          Object.values(parsed).every((v) => typeof v === "boolean")
        );
      } catch {
        return false;
      }
    },
    { message: "Must be a JSON object of booleans, e.g. {\"viral_db\": true}" },
  ),
});

// z.coerce.number() means the raw form value (a string from <input>) and the
// parsed output (a number) are different types — RHF's third generic lets
// the resolver accept the raw input shape while onSubmit receives the
// coerced output shape.
type SettingsFormInput = z.input<typeof settingsSchema>;
type SettingsFormOutput = z.output<typeof settingsSchema>;

export function SettingsPage() {
  const { data, isLoading, error } = useSettings();
  const update = useUpdateSettings();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<SettingsFormInput, unknown, SettingsFormOutput>({ resolver: zodResolver(settingsSchema) });

  useEffect(() => {
    if (data) {
      reset({
        max_upload_size_bytes: data.max_upload_size_bytes,
        allowed_extensions: data.allowed_extensions.join(", "),
        storage_provider: data.storage_provider,
        email_provider: data.email_provider,
        email_from_address: data.email_from_address,
        feature_flags_json: JSON.stringify(data.feature_flags, null, 2),
      });
    }
  }, [data, reset]);

  async function onSubmit(values: SettingsFormOutput) {
    try {
      await update.mutateAsync({
        max_upload_size_bytes: values.max_upload_size_bytes,
        allowed_extensions: values.allowed_extensions
          .split(",")
          .map((s) => s.trim().toLowerCase())
          .filter(Boolean),
        storage_provider: values.storage_provider,
        email_provider: values.email_provider,
        email_from_address: values.email_from_address,
        feature_flags: JSON.parse(values.feature_flags_json),
        updated_at: data?.updated_at ?? new Date().toISOString(),
      });
      toast.success("Settings saved");
    } catch (err) {
      toast.error(extractErrorMessage(err));
    }
  }

  if (isLoading) return <p className="text-sm text-gray-500">Loading settings…</p>;
  if (error || !data) {
    return (
      <p className="rounded-lg bg-red-500/10 px-4 py-3 text-sm text-red-400">
        {error ? extractErrorMessage(error) : "Failed to load settings."}
      </p>
    );
  }

  return (
    <div className="max-w-2xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Settings</h1>
        <p className="mt-1 text-sm text-gray-400">App-wide configuration.</p>
      </div>

      <Card>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <FormField label="Max upload size (bytes)" error={errors.max_upload_size_bytes?.message}>
            <input
              type="number"
              className={inputClasses}
              {...register("max_upload_size_bytes")}
            />
          </FormField>

          <FormField label="Allowed extensions (comma-separated)" error={errors.allowed_extensions?.message}>
            <input type="text" className={inputClasses} placeholder="jpg, jpeg, png" {...register("allowed_extensions")} />
          </FormField>

          <FormField label="Storage provider" error={errors.storage_provider?.message}>
            <input type="text" className={inputClasses} {...register("storage_provider")} />
          </FormField>

          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <FormField label="Email provider" error={errors.email_provider?.message}>
              <input type="text" className={inputClasses} placeholder="resend" {...register("email_provider")} />
            </FormField>
            <FormField label="Email from address" error={errors.email_from_address?.message}>
              <input
                type="email"
                className={inputClasses}
                placeholder="noreply@thumbnailiq.com"
                {...register("email_from_address")}
              />
            </FormField>
          </div>

          <FormField label="Feature flags (JSON)" error={errors.feature_flags_json?.message}>
            <textarea
              rows={5}
              className={`${inputClasses} font-mono text-xs`}
              {...register("feature_flags_json")}
            />
          </FormField>

          <Button type="submit" loading={isSubmitting}>
            Save settings
          </Button>
        </form>
      </Card>
    </div>
  );
}
