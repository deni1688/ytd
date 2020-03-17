async function main() {
	const form = document.getElementById('yt-form');
	const loader = document.getElementById('loader');

	form.addEventListener('submit', e => {
		e.preventDefault();

		loader.classList.add('d-flex');

		const urlInput = document.getElementById('yt-url').value;
		const filename = document.getElementById('yt-filename').value;

		if (urlInput && filename) {
			getMP3(urlInput, filename).then(() => {
				form.reset();
				loader.classList.remove('d-flex');
			});
		} else {
			alert('Please make sure all fields are filled out!');
		}
	});
}

function getMP3(url, filename) {
	return fetch(`/api/v1/convert?url=${url}`)
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
