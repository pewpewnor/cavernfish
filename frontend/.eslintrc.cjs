module.exports = {
  root: true,
  parser: '@typescript-eslint/parser',
  parserOptions: {
    ecmaVersion: 2020,
    sourceType: 'module',
  },
  plugins: ['@typescript-eslint', 'svelte3'],
  overrides: [
    {
      files: ['*.svelte'],
      processor: 'svelte3/svelte3',
    },
  ],
  extends: ['plugin:@typescript-eslint/recommended'],
  settings: {
    'svelte3/typescript': () => require('typescript'),
  },
  rules: {
    '@typescript-eslint/no-explicit-any': 'off',
    '@typescript-eslint/no-unused-vars': ['warn', { argsIgnorePattern: '^_' }],
    '@typescript-eslint/ban-ts-comment': 'off',
  },
  ignorePatterns: ['node_modules/', 'dist/', 'wailsjs/'],
}
