const form = document.querySelector("form#select-rom")

form.addEventListener("submit", (e) => {
  e.preventDefault()

  const data = new FormData(form)

  const rom = data.get("rom")

  document.rom = rom

  form.remove()

  const go = new Go()

  WebAssembly
    .instantiateStreaming(fetch(`chip8.wasm`), go.importObject)
    .then(result => {
      go.run(result.instance)
      const canvas = document.getElementsByTagName("canvas")[0]

      canvas.addEventListener("keydown", (e) => {
        console.log(e.key)
          if (e.key == "Escape") {
            location.reload()
          }
        })
    })
})
