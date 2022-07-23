// This file provides the top-level application configuration.

// Local override for `APP_CONFIG` values. Not source-controlled, provides
// support for local development config.
const LOCAL_CONFIG_FILE = './app.config.local.js';

const APP_CONFIG = {
	// public listening port for the application server
	port: 29899,
	// address to bind the listener to
	host: 'localhost',
};

// The local override, if provided, is expected to export an object
// with `APP_CONFIG` fields.
const LOCAL_CONFIG =
	(() => {
		const fs = require('fs');
		if (fs.existsSync(LOCAL_CONFIG_FILE)) {
			return require(LOCAL_CONFIG_FILE);
		}
	})() || {};

// Merge the default and local configuration.
module.exports = {
	...APP_CONFIG,
	...LOCAL_CONFIG,
};
