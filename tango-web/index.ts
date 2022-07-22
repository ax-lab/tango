import 'preact/debug';

import { render } from 'preact';
import App from './App';

import './css/main.less';

function init() {
	const loading = document.getElementById('loading');
	loading && loading.parentElement?.removeChild(loading);
	render(
		App(),
		document.getElementById('root') ??
			(() => {
				document.body.innerHTML = 'Error rendering application';
				throw new Error('root rendering element not found');
			})(),
	);
}

init();
