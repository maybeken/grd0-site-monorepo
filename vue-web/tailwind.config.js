/** @type {import('tailwindcss').Config} */
const colors = require('tailwindcss/colors');

export default {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    colors: {
      ...colors,
      background: colors.neutral['950'],
      shade: colors.neutral['900'],
      foreground: colors.neutral['100'],
      secondary: colors.neutral['400'],
      accent: colors.neutral['600'],
    },
    extend: {},
  },
  plugins: [],
}

