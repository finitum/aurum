module.exports =  {
    parser:  "@typescript-eslint/parser",  // Specifies the ESLint parser
    plugins: ["jest"],
    extends:  [
        "eslint:recommended",
        "plugin:@typescript-eslint/eslint-recommended",
        "plugin:@typescript-eslint/recommended",
        "plugin:jest/recommended"
    ],
    parserOptions:  {
        ecmaVersion:  2018,  // Allows for the parsing of modern ECMAScript features
        sourceType:  "module",  // Allows for the use of imports
    },
    rules:  {
        // Place to specify ESLint rules. Can be used to overwrite rules specified from the extended configs
        "@typescript-eslint/explicit-function-return-type": "error",
        "semi": ["error", "always"],
        "quotes": ["error", "double"],
        "@typescript-eslint/interface-name-prefix": ["error", "always"],
        "@typescript-eslint/ban-ts-ignore": "off"
    },
    env: {
        "browser": true
    }
};
