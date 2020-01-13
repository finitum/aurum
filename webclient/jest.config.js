module.exports = {
  // Enabel typescript support
  preset: "ts-jest",

  // testEnvironment: 'node',
  testEnvironment: "jest-environment-jsdom",

  // Automatically clear mock calls and instances between every test
  clearMocks: true,

  // The directory where Jest should output its coverage files
  coverageDirectory: "coverage",

  coverageReporters: ["text", "text-summary"],

  setupFiles: ["jest-localstorage-mock", "./jest.setup.js"],

  globals: {
    "ts-jest": {
      diagnostics: {
        // Disable strict null check for tests (tests would fail anyway if result is null)
        ignoreCodes: [2531]
      }
    }
  }
};