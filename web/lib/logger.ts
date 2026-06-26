type LogLevel = "debug" | "info" | "warn" | "error";

const isProduction = process.env.NODE_ENV === "production";

function log(level: LogLevel, message: string, meta?: Record<string, unknown>): void {
  if (level === "debug" && isProduction) return;

  const timestamp = new Date().toISOString();
  const prefix = `[${timestamp}] ${level.toUpperCase()}`;
  const fn = level === "error" ? console.error : level === "warn" ? console.warn : console.log;

  if (meta) {
    fn(prefix, message, meta);
  } else {
    fn(prefix, message);
  }
}

export const logger = {
  debug: (message: string, meta?: Record<string, unknown>) => log("debug", message, meta),
  info: (message: string, meta?: Record<string, unknown>) => log("info", message, meta),
  warn: (message: string, meta?: Record<string, unknown>) => log("warn", message, meta),
  error: (message: string, meta?: Record<string, unknown>) => log("error", message, meta),
};
