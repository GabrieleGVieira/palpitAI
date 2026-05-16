const eslintConfigPrettier = require('eslint-config-prettier');
const expoConfig = require('eslint-config-expo/flat');

module.exports = [
  ...expoConfig,
  eslintConfigPrettier,
  {
    ignores: ['.expo/*', 'dist/*', 'node_modules/*', 'web-build/*'],
  },
];
