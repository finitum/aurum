module.exports = {
  preset: "@vue/cli-plugin-unit-jest/presets/typescript",
  transform: {
    "^.+\\.vue$": "vue-jest"
  },

  // preset: "ts-jest"
  testMatch: ["**/*.test.ts"]
};
