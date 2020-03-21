const loader = document.getElementById('loader');

async function main() {
	ensureAuth();

	const form = document.getElementById('yt-form');
	const out = document.getElementById('results');

	form.addEventListener('submit', e => {
		e.preventDefault();

		const query = document.getElementById('yt-search').value;

		if (query) {
			search(query).then(d => {
				render(out, d);
				form.reset();
			});
		}
	});
}

function ensureAuth() {
	const auth = localStorage.getItem('auth');

	if (!auth) {
		const auth = prompt('Please enter your password to move forward!');
		localStorage.setItem('auth', auth);
	}
}

function render(outlet, d) {
	const items = d.items;

	outlet.innerHTML = '';

	for (const item of items) {
		outlet.innerHTML += `
			<div class="search-item">
				<img src="${item.snippet.thumbnails.default.url}"/>
				<div class="vid-details">
					<h5>${item.snippet.title}</h5>
					<button class="btn btn-primary btn-sm download" data-yt-url="https://youtu.be/${item.id.videoId}" data-title="${item.snippet.title}">Download</button>
				</div>
			</div>`;
	}

	const downloadBtns = document.getElementsByClassName('download');

	Array.from(downloadBtns).forEach(element => {
		element.addEventListener('click', e => {
			loader.classList.add('d-flex');

			getMP3(e.target.getAttribute('data-yt-url'), e.target.getAttribute('data-title')).then(() => {
				loader.remove('d-flex');
				Array.from(downloadBtns).forEach(element => element.removeEventListener('click', () => {}));
				outlet.innerHTML = '';
			});
		});
	});
}

async function search(query) {
	return fetch(`/api/v1/search?query=${query}`)
		.then(r => {
			if (r.status != 200) {
				return r.text().then(alert);
			}

			return r.json();
		})
		.then(d => d);
}

function getMP3(url, filename) {
	return fetch(`/api/v1/convert?url=${url}&key=${localStorage.getItem('auth')}`)
		.then(r => {
			if (r.status != 200) {
				throw r;
			}

			return r.blob();
		})
		.then(d => download(d, filename))
		.catch(e => e.text())
		.then(e => e && alert(e));
}

function download(data, filename) {
	const a = document.createElement('a');

	a.href = URL.createObjectURL(data);
	a.setAttribute('download', filename + '.mp3');
	a.click();

	return false;
}

window.onload = main;
