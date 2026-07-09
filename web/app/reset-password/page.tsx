"use client";

import { Suspense, useEffect, useRef, useState } from "react";
import Link from "next/link";
import { useRouter, useSearchParams } from "next/navigation";
import { motion } from "framer-motion";
import { ShieldCheck } from "lucide-react";
import { useAuth } from "@/hooks/useAuth";
import { ApiError } from "@/lib/api";
import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Input } from "@/components/ui/Input";
import { useMotionVariants } from "@/lib/motion";

const RESEND_COOLDOWN_SECONDS = 30;

function ResetPasswordInner() {
  const { resetPassword, requestPasswordReset } = useAuth();
  const router = useRouter();
  const { fadeInUp } = useMotionVariants();
  const searchParams = useSearchParams();

  const [email, setEmail] = useState(searchParams.get("email") ?? "");
  const [code, setCode] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [cooldown, setCooldown] = useState(0);

  const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
  useEffect(() => {
    return () => {
      if (timerRef.current) clearInterval(timerRef.current);
    };
  }, []);

  function startCooldown() {
    setCooldown(RESEND_COOLDOWN_SECONDS);
    if (timerRef.current) clearInterval(timerRef.current);
    timerRef.current = setInterval(() => {
      setCooldown((s) => {
        if (s <= 1 && timerRef.current) {
          clearInterval(timerRef.current);
          timerRef.current = null;
        }
        return s - 1;
      });
    }, 1000);
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setNotice(null);
    setSubmitting(true);
    try {
      await resetPassword(email, code.trim(), password);
      // Reset revokes all sessions, so send them to log in with the new password.
      router.push("/login?reset=1");
    } catch (err) {
      setError(
        err instanceof ApiError
          ? "That code is invalid or expired, or the password is too short (min 8 characters)."
          : "Unable to reset your password right now. Please try again.",
      );
      setSubmitting(false);
    }
  }

  async function handleResend() {
    if (cooldown > 0 || !email) return;
    setError(null);
    setNotice(null);
    try {
      await requestPasswordReset(email);
      setNotice("If the account exists, a new reset code is on its way.");
      startCooldown();
    } catch {
      setError("Couldn't resend the code. Please try again in a moment.");
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center px-6">
      <motion.div initial="hidden" animate="visible" variants={fadeInUp} className="w-full max-w-sm">
        <div className="mb-8 text-center">
          <Link href="/" className="mb-4 inline-flex items-center gap-2">
            <span className="flex h-9 w-9 items-center justify-center rounded-lg bg-brand-gradient text-sm font-bold text-white shadow-glow">
              IQ
            </span>
            <span className="text-lg font-semibold tracking-tight text-white">ThumbnailIQ</span>
          </Link>
          <h1 className="text-2xl font-bold tracking-tight text-white">Set a new password</h1>
          <p className="mt-1 text-sm text-gray-400">
            {email ? (
              <>
                Enter the code sent to{" "}
                <span className="font-medium text-gray-200">{email}</span> and choose a new password
              </>
            ) : (
              "Enter your email, the reset code, and a new password"
            )}
          </p>
        </div>

        <Card as="form" onSubmit={handleSubmit} className="space-y-4">
          {error && <Alert variant="danger">{error}</Alert>}
          {notice && <Alert variant="info">{notice}</Alert>}

          {!searchParams.get("email") && (
            <div>
              <label htmlFor="email" className="mb-1.5 block text-sm font-medium text-gray-300">
                Email
              </label>
              <Input
                id="email"
                type="email"
                autoComplete="email"
                required
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
              />
            </div>
          )}

          <div>
            <label htmlFor="code" className="mb-1.5 block text-sm font-medium text-gray-300">
              Reset code
            </label>
            <Input
              id="code"
              type="text"
              inputMode="numeric"
              autoComplete="one-time-code"
              pattern="\d{6}"
              maxLength={6}
              required
              value={code}
              onChange={(e) => setCode(e.target.value.replace(/\D/g, ""))}
              placeholder="123456"
              className="text-center text-lg tracking-[0.5em]"
            />
          </div>

          <div>
            <label htmlFor="password" className="mb-1.5 block text-sm font-medium text-gray-300">
              New password
            </label>
            <Input
              id="password"
              type="password"
              autoComplete="new-password"
              required
              minLength={8}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="At least 8 characters"
            />
          </div>

          <Button
            type="submit"
            loading={submitting}
            icon={<ShieldCheck className="h-4 w-4" aria-hidden="true" />}
            className="w-full"
          >
            Reset Password
          </Button>

          <Button
            type="button"
            variant="ghost"
            size="sm"
            onClick={handleResend}
            disabled={cooldown > 0 || !email}
            className="w-full"
          >
            {cooldown > 0 ? `Resend code in ${cooldown}s` : "Resend code"}
          </Button>
        </Card>

        <p className="mt-6 text-center text-sm text-gray-400">
          Remembered it?{" "}
          <Link href="/login" className="font-medium text-brand-300 transition-colors hover:text-brand-200">
            Back to log in
          </Link>
        </p>
      </motion.div>
    </div>
  );
}

export default function ResetPasswordPage() {
  // useSearchParams requires a Suspense boundary in the app router.
  return (
    <Suspense>
      <ResetPasswordInner />
    </Suspense>
  );
}
