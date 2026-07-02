import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx,mdx}",
    "./components/**/*.{js,ts,jsx,tsx,mdx}",
    "./hooks/**/*.{js,ts,jsx,tsx,mdx}",
    "./lib/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: ["var(--font-inter)", "system-ui", "sans-serif"],
      },
      colors: {
        brand: {
          50: "#f1f0ff",
          100: "#e4e1ff",
          200: "#cac3ff",
          300: "#a99cff",
          400: "#8a76ff",
          500: "#7c5cff",
          600: "#6d3df0",
          700: "#5b2cd1",
          800: "#4a24a8",
          900: "#3d1f87",
        },
        surface: {
          DEFAULT: "#0b0b14",
          50: "#15151f",
          100: "#191926",
          200: "#20202e",
          300: "#2a2a3a",
          400: "#393950",
          // 500/600 continue the same lightening ramp for hover/emphasis
          // borders and popover surfaces — a dark-only theme has no use
          // for a literal 700-900 (that would just be near-white).
          500: "#45455f",
          600: "#56567a",
        },
        // Semantic aliases over the raw Tailwind shades already used ad
        // hoc across the app (emerald-400/amber-400/red-400/sky-400), so
        // new components read as `text-success` / `bg-danger/15` instead
        // of guessing which exact shade was used elsewhere. Plain hex
        // values so Tailwind's `/opacity` modifier syntax works on them.
        success: "#34d399",
        warning: "#fbbf24",
        danger: "#f87171",
        info: "#38bdf8",
      },
      backgroundImage: {
        "brand-gradient": "linear-gradient(135deg, #6d3df0 0%, #a855f7 100%)",
      },
      boxShadow: {
        glow: "0 0 40px -10px rgba(124, 92, 255, 0.45)",
        card: "0 1px 2px 0 rgba(0, 0, 0, 0.3), 0 1px 3px 0 rgba(0, 0, 0, 0.2)",
        "card-hover":
          "0 8px 24px -4px rgba(0, 0, 0, 0.45), 0 0 0 1px rgba(124, 92, 255, 0.08)",
        elevated: "0 20px 50px -12px rgba(0, 0, 0, 0.6)",
      },
    },
  },
  plugins: [],
};

export default config;
