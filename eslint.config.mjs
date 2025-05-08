//  @ts-check

import { tanstackConfig } from '@tanstack/eslint-config';
import globals from 'globals';
import pluginVitest from 'eslint-plugin-vitest';
import pluginUnusedImports from 'eslint-plugin-unused-imports';
import pluginReact from 'eslint-plugin-react';
import pluginReactHooks from 'eslint-plugin-react-hooks';
import eslintConfigPrettier from 'eslint-config-prettier';
import eslint from '@eslint/js';
import tseslint from 'typescript-eslint';

export default tseslint.config(
  eslint.configs.recommended,
  tseslint.configs.recommended,
  tanstackConfig,
  {
    ignores: ['**/*.gen.ts'],
    files: ['**/*.tsx', '**/*.ts'],
    languageOptions: {
      ecmaVersion: 2022,
      sourceType: 'module',
      globals: {
        ...globals.browser,
        ...globals.nodeBuiltin,
        ...globals.node,
      },
    },
    plugins: {
      vitest: pluginVitest,
      react: pluginReact,
      'react-hooks': pluginReactHooks,
      'unused-imports': pluginUnusedImports,
    },
    rules: {
      '@typescript-eslint/no-unnecessary-condition': 'off',
      '@typescript-eslint/no-explicit-any': 'off',
      '@typescript-eslint/naming-convention': [
        'error',
        {
          selector: ['import', 'variable'],
          format: ['camelCase', 'PascalCase', 'UPPER_CASE'],
          leadingUnderscore: 'allow',
        },
      ],
      '@typescript-eslint/no-unused-vars': 'off',
      'unused-imports/no-unused-imports': 'error',
      'unused-imports/no-unused-vars': [
        'error',
        {
          vars: 'all',
          varsIgnorePattern: '^_',
          args: 'after-used',
          argsIgnorePattern: '^_',
        },
      ],
      curly: 'error',
      eqeqeq: 'error',
      'no-throw-literal': 'warn',
      semi: 'error',
      ...pluginReactHooks.configs.recommended.rules,
    },
  },
  eslintConfigPrettier,
);
