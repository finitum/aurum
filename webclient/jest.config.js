module.exports = {
  // Enabel typescript support
  preset: 'ts-jest',

  // testEnvironment: 'node',
  testEnvironment: "jest-environment-jsdom",

  // Automatically clear mock calls and instances between every test
  clearMocks: true,

  // The directory where Jest should output its coverage files
  coverageDirectory: "coverage",

  coverageReporters: ["lcov", "text", "text-summary"],

  setupFiles: ["jest-localstorage-mock"]
};