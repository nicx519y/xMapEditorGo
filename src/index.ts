declare var Go;
declare var WebAssembly;
const go = new Go()
const mainBox = document.getElementById('main-box');
const canvas = document.getElementById('render-canvas') as HTMLCanvasElement;
const ctx = canvas.getContext('2d');

(function main() {
	WebAssembly.instantiateStreaming(fetch('./assembly/engine.wasm'),go.importObject)
	.then( res => go.run(res.instance) )
	
	resizeCanvas();
	window.addEventListener('resize', resizeCanvas);
})();

function resizeCanvas() {
	let w = mainBox.clientWidth,
		h = mainBox.clientHeight;
	canvas.setAttribute('width', w + 'px');
	canvas.setAttribute('height', h + 'px');
	canvas.style.width = w + 'px';
	canvas.style.height = h + 'px';
}

window['isReady'] = function(callback) {
	callback(mainBox.clientWidth, mainBox.clientHeight)
}

var n = 0;

window['printer'] = function(arr, x, y, width, height) {
	
	if(width <= 0 || height <= 0) return;
	n ++
	let carr = new Uint32Array(arr);
	let narr = new Uint8ClampedArray(carr.buffer);
	
	let imageData = new ImageData(width, height);
	imageData.data.set(narr)

	ctx.putImageData(imageData, x, y);
	console.log('Render complete: x: ', x, ' y: ', y, ' total: ', n);
	carr = null;
	narr = null;

}