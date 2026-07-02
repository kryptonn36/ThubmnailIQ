import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { z } from "zod";
import { Button } from "@/components/ui/Button";
import { FormField, inputClasses } from "@/components/forms/FormField";
import { useAuth } from "@/hooks/useAuth";
import { extractErrorMessage } from "@/lib/axios";

const loginSchema = z.object({
  email: z.string().min(1, "Email is required").email("Enter a valid email"),
  password: z.string().min(1, "Password is required"),
});

type LoginForm = z.infer<typeof loginSchema>;

export function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const [serverError, setServerError] = useState<string | null>(null);

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginForm>({ resolver: zodResolver(loginSchema) });

  async function onSubmit(values: LoginForm) {
    setServerError(null);
    try {
      await login(values.email, values.password);
      navigate("/", { replace: true });
    } catch (err) {
      setServerError(extractErrorMessage(err));
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-[#0e0e14] px-6">
      <div className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <span className="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-lg bg-brand-600 text-sm font-bold text-white">
            IQ
          </span>
          <h1 className="text-xl font-bold text-white">Admin Panel</h1>
          <p className="mt-1 text-sm text-gray-400">Sign in to manage ThumbnailIQ</p>
        </div>

        <form
          onSubmit={handleSubmit(onSubmit)}
          className="space-y-4 rounded-2xl border border-surface-300 bg-surface-100 p-6"
        >
          {serverError && (
            <p className="rounded-lg bg-red-500/10 px-3 py-2 text-sm text-red-400">{serverError}</p>
          )}

          <FormField label="Email" error={errors.email?.message}>
            <input
              type="email"
              className={inputClasses}
              placeholder="admin@thumbnailiq.local"
              {...register("email")}
            />
          </FormField>

          <FormField label="Password" error={errors.password?.message}>
            <input type="password" className={inputClasses} {...register("password")} />
          </FormField>

          <Button type="submit" loading={isSubmitting} className="w-full">
            Sign In
          </Button>
        </form>
      </div>
    </div>
  );
}
