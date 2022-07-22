module.exports = {
	rootDir: '.',
	roots: ['<rootDir>/tango-web'],

	transform: {
		'^.+\\.(t|j)sx?$': ['@swc/jest'],
	},

	collectCoverage: false,
	coverageDirectory: './build/coverage',
	coveragePathIgnorePatterns: ['\\\\node_modules\\\\'],
	coverageReporters: ['json', 'html'],
};
