import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/pages/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/components/**/*.{js,ts,jsx,tsx,mdx}",
    "./src/app/**/*.{js,ts,jsx,tsx,mdx}",
  ],
  darkMode: "class",
  theme: {
    extend: {
      colors: {
        // Custom color palette as specified
        slate: {
          navy: "#0F172A",
          900: "#0F172A",
          800: "#1E293B",
          700: "#334155",
          600: "#475569",
          500: "#64748B",
          400: "#94A3B8",
          300: "#CBD5E1",
          200: "#E2E8F0",
          100: "#F1F5F9",
          50: "#F8FAFC",
        },
        off: {
          white: "#F8FAFC",
        },
        emerald: {
          muted: "#10B981",
          500: "#10B981",
          400: "#34D399",
        },
        amber: {
          500: "#F59E0B",
          400: "#FBBF24",
        },
        cyan: {
          950: "#083344",
          500: "#06B6D4",
        },
      },
      fontFamily: {
        sans: ["Inter", "Plus Jakarta Sans", "system-ui", "sans-serif"],
      },
      borderRadius: {
        "bento-sm": "16px",
        "bento-lg": "24px",
      },
      boxShadow: {
        "bento": "0 4px 20px rgba(0, 0, 0, 0.1)",
        "bento-hover": "0 8px 30px rgba(0, 0, 0, 0.15)",
      },
      animation: {
        "pulse-slow": "pulse 3s cubic-bezier(0.4, 0, 0.6, 1) infinite",
        "float": "float 3s ease-in-out infinite",
      },
      keyframes: {
        float: {
          "0%, 100%": { transform: "translateY(0px)" },
          "50%": { transform: "translateY(-10px)" },
        },
      },
    },
  },
  plugins: [],
};

export default config;