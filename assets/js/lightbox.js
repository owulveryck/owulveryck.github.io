(function () {
  var overlay, img, close, previousFocus;

  function create() {
    overlay = document.createElement("div");
    overlay.id = "lightbox-overlay";
    overlay.setAttribute("role", "dialog");
    overlay.setAttribute("aria-modal", "true");
    overlay.setAttribute("aria-label", "Image viewer");
    overlay.innerHTML =
      '<button id="lightbox-close" aria-label="Close">&times;</button>' +
      '<img id="lightbox-img" />';
    document.body.appendChild(overlay);

    img = document.getElementById("lightbox-img");
    close = document.getElementById("lightbox-close");

    overlay.addEventListener("click", function (e) {
      if (e.target === overlay || e.target === close) hide();
    });
    document.addEventListener("keydown", function (e) {
      if (!overlay.classList.contains("active")) return;
      if (e.key === "Escape") hide();
      if (e.key === "Tab") {
        e.preventDefault();
        close.focus();
      }
    });
  }

  function show(src, alt) {
    if (!overlay) create();
    previousFocus = document.activeElement;
    img.src = src;
    img.alt = alt || "";
    overlay.classList.add("active");
    document.body.style.overflow = "hidden";
    close.focus();
  }

  function hide() {
    if (!overlay) return;
    overlay.classList.remove("active");
    document.body.style.overflow = "";
    if (previousFocus) previousFocus.focus();
  }

  document.addEventListener("DOMContentLoaded", function () {
    var content = document.querySelector(".post-content");
    if (!content) return;

    content.addEventListener("click", function (e) {
      var target = e.target;
      if (target.tagName !== "IMG") return;
      if (target.closest(".no-lightbox")) return;

      e.preventDefault();
      var figure = target.closest("figure");
      var link = target.closest("a");
      var src = (link && link.href) || target.src;
      show(src, target.alt);
    });
  });
})();
