(function () {
  "use strict";
  const endpoint = new URL(document.currentScript.src).origin + "/api/event";

  function send() {
    const payload = JSON.stringify({
      domain: location.hostname,
      path: location.pathname,
      referrer: document.referrer,
    });
    if (navigator.sendBeacon) {
      navigator.sendBeacon(endpoint, payload);
    } else {
      fetch(endpoint, { method: "POST", body: payload, keepalive: true });
    }
  }

  send();
})();
