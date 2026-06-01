(function () {
  var overlay, img, close;

  function create() {
    overlay = document.createElement("div");
    overlay.id = "lightbox-overlay";
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
      if (e.key === "Escape") hide();
    });
  }

  function show(src, alt) {
    if (!overlay) create();
    img.src = src;
    img.alt = alt || "";
    overlay.classList.add("active");
    document.body.style.overflow = "hidden";
  }

  function hide() {
    if (!overlay) return;
    overlay.classList.remove("active");
    document.body.style.overflow = "";
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
