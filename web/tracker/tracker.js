(function () {
  "use strict";
  const endpoint = new URL(document.currentScript.src).origin + "/api/event";

  function send() {
    const payload = JSON.stringify({
      domain: location.hostname,
      path: location.pathname,
      referrer: document.referrer,
      screen_size: window.screen.width + "x" + window.screen.height,
    });
    if (navigator.sendBeacon) {
      navigator.sendBeacon(endpoint, payload);
    } else {
      fetch(endpoint, { method: "POST", body: payload, keepalive: true });
    }
  }

  send();
})();
