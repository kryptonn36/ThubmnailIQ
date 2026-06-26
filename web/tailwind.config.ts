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
        },
      },
      backgroundImage: {
        "brand-gradient": "linear-gradient(135deg, #6d3df0 0%, #a855f7 100%)",
      },
      boxShadow: {
        glow: "0 0 40px -10px rgba(124, 92, 255, 0.45)",
      },
    },
  },
  plugins: [],
};

export default config;
