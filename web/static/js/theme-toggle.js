(function () {
    const STORAGE_KEY = "financeflow-theme";
    const DARK_THEME = "dark";
    const LIGHT_THEME = "light";

    function applyTheme(theme) {
        const isDark = theme === DARK_THEME;
        document.body.classList.toggle("theme-dark", isDark);

        const toggleButton = document.getElementById("themeToggleBtn");
        if (toggleButton) {
            toggleButton.setAttribute("aria-pressed", String(isDark));
            toggleButton.setAttribute("title", isDark ? "Переключить на светлую тему" : "Переключить на темную тему");
            toggleButton.innerHTML = isDark
                ? '<i class="fa-solid fa-sun"></i><span>Светлая</span>'
                : '<i class="fa-solid fa-moon"></i><span>Темная</span>';
        }
    }

    function getInitialTheme() {
        const savedTheme = localStorage.getItem(STORAGE_KEY);
        if (savedTheme === DARK_THEME || savedTheme === LIGHT_THEME) {
            return savedTheme;
        }

        const prefersDark = window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches;
        return prefersDark ? DARK_THEME : LIGHT_THEME;
    }

    document.addEventListener("DOMContentLoaded", function () {
        applyTheme(getInitialTheme());

        const toggleButton = document.getElementById("themeToggleBtn");
        if (!toggleButton) {
            return;
        }

        toggleButton.addEventListener("click", function () {
            const nextTheme = document.body.classList.contains("theme-dark") ? LIGHT_THEME : DARK_THEME;
            localStorage.setItem(STORAGE_KEY, nextTheme);
            applyTheme(nextTheme);
        });
    });
})();
