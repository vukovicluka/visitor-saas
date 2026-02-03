(function () {
  "use strict";

  const domain = document.getElementById("domain");
  let period = "7d";

  document.querySelectorAll(".periods button").forEach(function (btn) {
    btn.addEventListener("click", function () {
      document.querySelector(".periods .active").classList.remove("active");
      btn.classList.add("active");
      period = btn.dataset.period;
      refresh();
    });
  });

  domain.addEventListener("change", refresh);

  function refresh() {
    const d = domain.value;
    const q = "?domain=" + encodeURIComponent(d) + "&period=" + period;

    fetch("/api/stats/summary" + q)
      .then(function (r) {
        return r.json();
      })
      .then(function (data) {
        document.getElementById("total-views").textContent = data.total_views;
        document.getElementById("unique-visitors").textContent =
          data.unique_visitors;
        renderChart(data.views_per_day || []);
      });

    fetch("/api/stats/pages" + q)
      .then(function (r) {
        return r.json();
      })
      .then(function (data) {
        renderTable("pages-table", data || []);
      });

    fetch("/api/stats/referrers" + q)
      .then(function (r) {
        return r.json();
      })
      .then(function (data) {
        renderTable("referrers-table", data || []);
      });
  }

  function renderChart(days) {
    var chart = document.getElementById("chart");
    chart.innerHTML = "";
    if (!days.length) return;

    var max = Math.max.apply(
      null,
      days.map(function (d) {
        return d.views;
      }),
    );
    if (max === 0) max = 1;

    days.forEach(function (d) {
      var pct = (d.views / max) * 100;
      const format = Intl.DateTimeFormat("en");
      var label = d.date.slice(5); // MM-DD

      var group = document.createElement("div");
      group.className = "bar-group";

      var bar = document.createElement("div");
      bar.className = "bar";
      bar.style.height = pct + "%";

      var tooltip = document.createElement("span");
      tooltip.className = "bar-tooltip";
      tooltip.textContent = d.views + " views, " + d.visitors + " visitors";
      bar.appendChild(tooltip);

      var lbl = document.createElement("span");
      lbl.className = "bar-label";
      lbl.textContent = label;

      group.appendChild(bar);
      group.appendChild(lbl);
      chart.appendChild(group);
    });
  }

  function renderTable(id, rows) {
    var tbody = document.querySelector("#" + id + " tbody");
    tbody.innerHTML = "";
    rows.forEach(function (row) {
      var tr = document.createElement("tr");
      var values = Object.values(row);
      values.forEach(function (v) {
        var td = document.createElement("td");
        td.textContent = v;
        tr.appendChild(td);
      });
      tbody.appendChild(tr);
    });
  }

  refresh();
})();
