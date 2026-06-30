(function () {
    const topBar = document.getElementById("topBar");
    const header = document.getElementById("siteHeader");
    const productsToggle = document.getElementById("productsToggle");
    const productsMenu = document.getElementById("productsMenu");
    const navToggle = document.getElementById("navToggle");
    const navLinks = document.getElementById("navLinks");
    const copyButton = document.getElementById("copyCommand");
    const installCommand = document.getElementById("installCommand");
    let lastScrollY = window.scrollY;
    let copyResetTimer;

    function syncHeader() {
        const currentScrollY = window.scrollY;
        const scrollingDown = currentScrollY > lastScrollY;

        header.classList.toggle("is-scrolled", currentScrollY > 24);

        if (topBar) {
            const hideTopBar = scrollingDown && currentScrollY > 120;
            topBar.classList.toggle("hidden", hideTopBar);
            header.classList.toggle("is-top", hideTopBar || window.innerWidth <= 780);
        }

        lastScrollY = currentScrollY;
    }

    window.addEventListener("scroll", syncHeader, { passive: true });
    window.addEventListener("resize", syncHeader);
    syncHeader();

    if (productsToggle && productsMenu) {
        productsToggle.addEventListener("click", function (event) {
            event.stopPropagation();
            const isOpen = productsMenu.classList.toggle("show");
            productsToggle.setAttribute("aria-expanded", String(isOpen));
        });
    }

    if (navToggle && navLinks) {
        navToggle.addEventListener("click", function () {
            const isOpen = navLinks.classList.toggle("is-open");
            navToggle.classList.toggle("is-open", isOpen);
            navToggle.setAttribute("aria-expanded", String(isOpen));
        });

        navLinks.querySelectorAll("a").forEach(function (link) {
            link.addEventListener("click", function () {
                navLinks.classList.remove("is-open");
                navToggle.classList.remove("is-open");
                navToggle.setAttribute("aria-expanded", "false");
            });
        });
    }

    function setCopyFeedback(message, status) {
        if (!copyButton) {
            return;
        }

        const originalLabel = copyButton.dataset.originalLabel || copyButton.getAttribute("aria-label") || "复制启动命令";
        copyButton.dataset.originalLabel = originalLabel;
        copyButton.setAttribute("aria-label", message);
        copyButton.setAttribute("title", message);
        copyButton.classList.toggle("is-copied", status === "success");
        copyButton.classList.toggle("is-copy-error", status === "error");

        window.clearTimeout(copyResetTimer);
        copyResetTimer = window.setTimeout(function () {
            copyButton.setAttribute("aria-label", originalLabel);
            copyButton.setAttribute("title", originalLabel);
            copyButton.classList.remove("is-copied", "is-copy-error");
        }, 1400);
    }

    async function copyText(text) {
        if (navigator.clipboard && window.isSecureContext) {
            await navigator.clipboard.writeText(text);
            return;
        }

        const textarea = document.createElement("textarea");
        textarea.value = text;
        textarea.setAttribute("readonly", "");
        textarea.style.position = "fixed";
        textarea.style.top = "-9999px";
        textarea.style.left = "-9999px";
        document.body.appendChild(textarea);
        textarea.select();

        try {
            if (!document.execCommand("copy")) {
                throw new Error("copy command failed");
            }
        } finally {
            textarea.remove();
        }
    }

    if (copyButton && installCommand) {
        copyButton.addEventListener("click", async function () {
            try {
                await copyText(installCommand.textContent.trim());
                setCopyFeedback("已复制", "success");
            } catch (error) {
                setCopyFeedback("复制失败，请手动复制", "error");
            }
        });
    }

    document.addEventListener("click", function (event) {
        if (productsMenu && productsToggle && !productsMenu.contains(event.target) && !productsToggle.contains(event.target)) {
            productsMenu.classList.remove("show");
            productsToggle.setAttribute("aria-expanded", "false");
        }
    });

    document.addEventListener("keydown", function (event) {
        if (event.key === "Escape") {
            if (productsMenu && productsToggle) {
                productsMenu.classList.remove("show");
                productsToggle.setAttribute("aria-expanded", "false");
            }
            if (navLinks && navToggle) {
                navLinks.classList.remove("is-open");
                navToggle.classList.remove("is-open");
                navToggle.setAttribute("aria-expanded", "false");
            }
        }
    });
})();
