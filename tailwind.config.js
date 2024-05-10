/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx}", "./pages/*.html", "/pages/*.html"],
  theme: {
    extend: {
      container: {
        maxWidth: {
          "2xl": "[1368]px", // Set the custom max-width for screens wider than 1536px
        },
      },
      colors: {
        secondary: "var(--clr-jet)",
        secondaryDeep: 'var(--clr-eerieBlack)',
        primary: "var(--clr-neonBlue)",
        primaryMiddle: "var(--clr-savoyBlue)",
        primaryDeep: "var(--clr-palatinateBlue)",
        raisinBlack: "var(--clr-raisinBlack)",
        graniteGray: "var(--clr-graniteGray)",
        ghostWhite: "var(--clr-ghostWhite)",
        white: "var(--clr-white)",
        // Add more custom colors using the CSS variables
      },
      boxShadow: {
        custom: 'rgba(100, 100, 111, 0.2) 0px 7px 29px 0px',
        discord: '0 8px 15px rgba(0,0,0,.2)',
      },
    },
  },
  plugins: [require("daisyui")],
};
