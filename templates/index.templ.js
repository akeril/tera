// {{define "script"}}

// add websocket handlers
const ws = new WebSocket("{{.Uri}}");

ws.addEventListener("open", () => {
  console.log("Websocket connection established!");
});

ws.addEventListener("close", () => {
  console.log("Websocket connection terminated!");
  const html = document.querySelector("html");
  html.innerHTML = "<p> Connection terminated!</p>";
});

ws.addEventListener("message", async (e) => {
  const event = JSON.parse(e.data);
  reload(event);
});

function reload(event) {
  // strip url prefix
  const url = event.Name.substr(2);

  // if entrypoint changes, reload the entire page
  if (url == "{{.Entrypoint}}") {
    loadEntryPoint();
    return;
  }

  for (let tag of ["href", "src", "data"]) {
    renderUpdates(tag, url);
  }
}

// locate changed elements and update them
function renderUpdates(tag, url) {
  const elements = document.querySelectorAll(`[${tag}*="${url}"]`);
  for (let element of elements) {
    const newUrl = `${url}?t=${new Date().getTime()}`;
    element.setAttribute(tag, newUrl);
  }
}

// initial data fetch
async function loadEntryPoint() {
  const entryPoint = "{{.Entrypoint}}";
  const html = document.querySelector("html");

  let data;
  // handle html content

  if (entryPoint.endsWith(".pdf")) {
    data = `
    <object 
      data="${entryPoint}" 
      type="application/pdf" 
      style="width: 100vw; height: 100vh; position: fixed; top: 0; left: 0; border: none;"
    >
    </object>
  `;
  } else {
    const resp = await fetch("{{.Entrypoint}}");
    data = await resp.text();
  }
  html.innerHTML = data;
}

loadEntryPoint();

// {{end}}
