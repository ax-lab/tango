const path = require('path');

const OUTPUT = path.resolve(__dirname, 'build');

const DIR = `./tango-web`;
const APP = `${DIR}/index.ts`;
const PORT = 29899;

const PreactRefreshPlugin = require('@prefresh/webpack');

const config = {
	mode: 'development',
	plugins: [new PreactRefreshPlugin()],
	module: {
		rules: [
			{
				test: /\.tsx?$/,
				use: {
					loader: 'swc-loader',
					options: {
						jsc: {
							target: 'es2015',
							parser: {
								syntax: 'typescript',
								tsx: true,
							},
							transform: {
								react: {
									runtime: 'automatic',
									pragma: 'h',
									pragmaFrag: 'Fragment',
								},
							},
						},
					},
				},
				exclude: /node_modules/,
			},
			{
				test: /\.less$/,
				use: ['style-loader', 'css-loader', 'less-loader'],
			},
		],
	},
	resolve: {
		extensions: ['.tsx', '.ts', '.js'],
		alias: {
			react: 'preact/compat',
			'react-dom/test-utils': 'preact/test-utils',
			'react-dom': 'preact/compat', // Must be below test-utils
			'react/jsx-runtime': 'preact/jsx-runtime',
		},
	},
};

module.exports = (env, args) => {
	const { server = false } = env || {};
	const app = {
		entry: server ? ['webpack-dev-server/client', APP] : APP,
		devtool: args.mode == 'production' ? undefined : 'inline-source-map',
		output: {
			filename: 'index.js',
			path: OUTPUT,
			publicPath: '/',
		},
		devServer: {
			host: '0.0.0.0',
			hot: true,
			port: PORT,
			static: {
				directory: `${DIR}/public`,
				serveIndex: true,
				watch: true,
			},
			proxy: {
				'/api': {
					target: 'http://localhost:29801',
				},
			},
		},
	};

	return Object.assign({}, config, app);
};
