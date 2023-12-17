const go = new Go()

WebAssembly
  .instantiateStreaming(fetch("chip8.wasm"), go.importObject)
  .then(result => go.run(result.instance))
