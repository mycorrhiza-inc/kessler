const svgToDataUri = require("mini-svg-data-uri");
const {
  default: flattenColorPalette,
} = require("tailwindcss/lib/util/flattenColorPalette");
import type { Config } from "tailwindcss";

const config = {
  darkMode: ["class"],
  content: [
    "./pages/**/*.{ts,tsx}",
    "./components/**/*.{ts,tsx}",
    "./app/**/*.{ts,tsx}",
    "./src/**/*.{ts,tsx}",
  ],
  prefix: "",
  theme: {
    extend: {
      animation: {
        aurora: "aurora 60s linear infinite",
      },
      keyframes: {
        aurora: {
          from: {
            backgroundPosition: "50% 50%, 50% 50%",
          },
          to: {
            backgroundPosition: "350% 50%, 350% 50%",
          },
        },
      },
    },
  },
  plugins: [
    function ({ matchUtilities, theme }: any) {
      matchUtilities(
        {
          "bg-dot-thick": (value: any) => ({
            backgroundImage: `url("${svgToDataUri(
              `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 32 32" width="16" height="16" fill="none"><circle fill="${value}" id="pattern-circle" cx="10" cy="10" r="2.5"></circle></svg>`,
            )}")`,
          }),
        },
        {
          values: flattenColorPalette(theme("backgroundColor")),
          type: "color",
        },
      );
    },
    addVariablesForColors,
    require("@tailwindcss/typography"),
    require("tailwindcss-animate"),
    require("daisyui"),
  ],
  daisyui: {
    themes: [
      {
        light: {
          ...require("daisyui/src/theming/themes")["light"],
          //"base-100": "#FFFFFF",
          "base-content": "#000000",
        },
        dark: {
          ...require("daisyui/src/theming/themes")["dim"],
          // "base-100": "#000000",
          "base-content": "#FFFFFF",
          "success-content": "#FFFFFF",
        },
        bumblebee: {
          ...require("daisyui/src/theming/themes")["bumblebee"],
          //"base-100": "#FFFFFF",
          "base-content": "#000000",
        },
        cmyk: {
          ...require("daisyui/src/theming/themes")["cmyk"],
          //"base-100": "#FFFFFF",
          "base-content": "#000000",
        },
        emerald: {
          ...require("daisyui/src/theming/themes")["emerald"],
          //"base-100": "#FFFFFF",
          "base-content": "#000000",
        },
      },
      "black",
      "forest",
      "corporate",
      "sunset",
      "acid",
    ],
  },
} satisfies Config;

export default config;

// This plugin adds each Tailwind color as a global CSS variable, e.g. var(--gray-200).
function addVariablesForColors({ addBase, theme }: any) {
  let allColors = flattenColorPalette(theme("colors"));
  let newVars = Object.fromEntries(
    Object.entries(allColors).map(([key, val]) => [`--${key}`, val]),
  );

  addBase({
    ":root": newVars,
  });
}
