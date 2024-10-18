import type { Config } from "tailwindcss";

export default {
  content: ["./templates/*.html"],
  theme: {
    extend: {},
  },
  plugins: [require("@tailwindcss/typography"), require("daisyui")],
} satisfies Config;
