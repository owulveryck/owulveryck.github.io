(function () {
  var footnoteElements = [];

  window.addEventListener("beforeprint", function () {
    var content = document.querySelector(".post-content");
    if (!content) return;

    var links = content.querySelectorAll('a[href^="http"]');
    if (!links.length) return;

    var urlMap = {};
    var counter = 0;

    links.forEach(function (link) {
      var href = link.href;
      if (!urlMap[href]) {
        counter++;
        urlMap[href] = counter;
      }
      var sup = document.createElement("sup");
      sup.className = "print-fn";
      sup.textContent = " [" + urlMap[href] + "]";
      link.parentNode.insertBefore(sup, link.nextSibling);
      footnoteElements.push(sup);
    });

    var section = document.createElement("section");
    section.className = "print-footnotes";
    var title = document.createElement("h3");
    title.textContent = "Liens";
    section.appendChild(title);

    var ol = document.createElement("ol");
    Object.keys(urlMap).forEach(function (url) {
      var li = document.createElement("li");
      li.setAttribute("value", urlMap[url]);
      li.textContent = url;
      ol.appendChild(li);
    });
    section.appendChild(ol);

    content.appendChild(section);
    footnoteElements.push(section);
  });

  window.addEventListener("afterprint", function () {
    footnoteElements.forEach(function (el) {
      if (el.parentNode) el.parentNode.removeChild(el);
    });
    footnoteElements = [];
  });
})();
