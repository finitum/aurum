/* eslint-disable */
const { colors } = require("tailwindcss/defaultTheme");

module.exports = {
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true
  },
  purge: ["./src/**/*.html", "./src/**/*.vue"],
  theme: {
    extend: {
      colors: {
        primary: "#8E0103",
        primarylight: "#B70104",
        primarydark: "#7A0103",

        secondary: "#ADB6C4",
        secondarylight: "#B8BFCC",
        secondarydark: "#8995A9",

        cancel: colors.red[500],
        accept: colors.green[500],
        info: colors.blue[500]
      }
    }
  },
  variants: {},
  plugins: []
};
