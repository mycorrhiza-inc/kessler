import type { Config } from "tailwindcss";

export default {
  content: ["../**/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
} satisfies Config;
