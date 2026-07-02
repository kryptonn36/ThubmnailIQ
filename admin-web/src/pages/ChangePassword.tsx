import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Link, useNavigate } from "react-router-dom";
import toast from "react-hot-toast";
import { z } from "zod";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { FormField, inputClasses } from "@/components/forms/FormField";
import { useChangePassword } from "@/hooks/useProfile";
import { extractErrorMessage } from "@/lib/axios";

const changePasswordSchema = z
  .object({
    currentPassword: z.string().min(1, "Current password is required"),
    newPassword: z.string().min(8, "New password must be at least 8 characters"),
    confirmPassword: z.string().min(1, "Please confirm your new password"),
  })
  .refine((data) => data.newPassword === data.confirmPassword, {
    message: "Passwords don't match",
    path: ["confirmPassword"],
  });

type ChangePasswordForm = z.infer<typeof changePasswordSchema>;

export function ChangePasswordPage() {
  const navigate = useNavigate();
  const changePassword = useChangePassword();

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors, isSubmitting },
  } = useForm<ChangePasswordForm>({ resolver: zodResolver(changePasswordSchema) });

  async function onSubmit(values: ChangePasswordForm) {
    try {
      await changePassword.mutateAsync({
        currentPassword: values.currentPassword,
        newPassword: values.newPassword,
      });
      toast.success("Password changed");
      reset();
      navigate("/profile");
    } catch (err) {
      toast.error(extractErrorMessage(err));
    }
  }

  return (
    <div className="max-w-md space-y-6">
      <div>
        <h1 className="text-2xl font-bold text-white">Change Password</h1>
        <p className="mt-1 text-sm text-gray-400">Update your admin account password.</p>
      </div>

      <Card>
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <FormField label="Current password" error={errors.currentPassword?.message}>
            <input type="password" className={inputClasses} {...register("currentPassword")} />
          </FormField>
          <FormField label="New password" error={errors.newPassword?.message}>
            <input type="password" className={inputClasses} {...register("newPassword")} />
          </FormField>
          <FormField label="Confirm new password" error={errors.confirmPassword?.message}>
            <input type="password" className={inputClasses} {...register("confirmPassword")} />
          </FormField>

          <Button type="submit" loading={isSubmitting}>
            Update password
          </Button>
        </form>
      </Card>

      <Link to="/profile" className="inline-block text-sm text-brand-500 hover:underline">
        ← Back to profile
      </Link>
    </div>
  );
}
