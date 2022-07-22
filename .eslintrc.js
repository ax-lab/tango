module.exports = {
	root: true,
	parser: '@typescript-eslint/parser',
	parserOptions: {
		project: 'tsconfig.json',
		tsconfigRootDir: __dirname,
	},
	env: {
		es6: true,
		node: true,
		browser: true,
		jest: true,
	},
	ignorePatterns: ['node_modules', 'build', 'coverage', 'dist', '*.config.js', '.*.js'],
	plugins: ['@typescript-eslint', 'eslint-comments', 'jest', 'only-warn'],
	extends: [
		'eslint:recommended',
		'plugin:eslint-comments/recommended',
		'plugin:@typescript-eslint/recommended',
		'plugin:@typescript-eslint/recommended-requiring-type-checking',
		'plugin:jest/recommended',
	],
	globals: {
		console: true,
	},
	rules: {
		// allow tabs for indent and spaces for alignment
		'no-mixed-spaces-and-tabs': 'off',

		// allow using any when we don't care about the type
		'@typescript-eslint/no-explicit-any': 'off',
	},
};
