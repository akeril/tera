// add websocket handlers
const teraWebSocket = new WebSocket("{{.Uri}}");

teraWebSocket.addEventListener("open", () => {
  console.log("Websocket connection established!");
});

teraWebSocket.addEventListener("close", () => {
  console.log("Websocket connection terminated!");
  const html = document.querySelector("html");
  html.innerHTML = "<p> Connection terminated!</p>";
});

teraWebSocket.addEventListener("message", async (e) => {
  const event = JSON.parse(e.data);
  reload(event);
});

function reload(event) {
  // strip url prefix
  const url = event.Name.substr(2);

  // if entrypoint changes, reload the entire page
  if (url == "{{.Entrypoint}}") {
    location.reload();
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
    element.setAttribute(tag, stamp(url));
  }
}

// add timestamp to url
function stamp(url) {
  return `${url}?t=${new Date().getTime()}`;
}
