const analyzeState = {
    mode: "day",
    view: "report",
    customStart: "",
    customEnd: "",
    flowChart: null,
    budgetChart: null,
};

const monthNamesRu = [
    "январь",
    "февраль",
    "март",
    "апрель",
    "май",
    "июнь",
    "июль",
    "август",
    "сентябрь",
    "октябрь",
    "ноябрь",
    "декабрь",
];

function parseDateOnly(value) {
    const [y, m, d] = value.split("-").map(Number);
    return new Date(y, m - 1, d);
}

function formatAmount(value, currency = "₽") {
    const abs = Math.abs(value);
    const formatted = abs.toLocaleString("ru-RU", {
        minimumFractionDigits: 2,
        maximumFractionDigits: 2,
    });
    const sign = value < 0 ? "-" : "";
    return `${sign}${formatted} ${currency}`;
}

function formatMonthShort(date) {
    const short = monthNamesRu[date.getMonth()].slice(0, 3);
    return short.charAt(0).toUpperCase() + short.slice(1);
}

function toISODate(date) {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const d = String(date.getDate()).padStart(2, "0");
    return `${y}-${m}-${d}`;
}

function startOfWeek(date) {
    const result = new Date(date);
    const day = (result.getDay() + 6) % 7;
    result.setDate(result.getDate() - day);
    result.setHours(0, 0, 0, 0);
    return result;
}

function endOfWeek(date) {
    const result = startOfWeek(date);
    result.setDate(result.getDate() + 6);
    return result;
}

function periodShiftDays(mode) {
    if (mode === "day") return 1;
    if (mode === "week") return 7;
    return 0;
}

function getCurrentPeriodBounds() {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    if (analyzeState.mode === "custom") {
        const startValue = analyzeState.customStart;
        const endValue = analyzeState.customEnd;
        if (!startValue || !endValue) {
            return { start: today, end: today };
        }

        const start = parseDateOnly(startValue);
        const end = parseDateOnly(endValue);
        return start <= end ? { start, end } : { start: end, end: start };
    }

    if (analyzeState.mode === "day") {
        return { start: today, end: today };
    }

    if (analyzeState.mode === "week") {
        return { start: startOfWeek(today), end: endOfWeek(today) };
    }

    if (analyzeState.mode === "month") {
        const start = new Date(today.getFullYear(), today.getMonth(), 1);
        const end = new Date(today.getFullYear(), today.getMonth() + 1, 0);
        return { start, end };
    }

    const start = new Date(today.getFullYear(), 0, 1);
    const end = new Date(today.getFullYear(), 11, 31);
    return { start, end };
}

function getPreviousPeriodBounds(current) {
    if (analyzeState.mode === "month") {
        const start = new Date(current.start.getFullYear(), current.start.getMonth() - 1, 1);
        const end = new Date(current.start.getFullYear(), current.start.getMonth(), 0);
        return { start, end };
    }

    if (analyzeState.mode === "year") {
        const year = current.start.getFullYear() - 1;
        return {
            start: new Date(year, 0, 1),
            end: new Date(year, 11, 31),
        };
    }

    const days = analyzeState.mode === "custom"
        ? Math.max(1, Math.round((current.end - current.start) / 86400000) + 1)
        : periodShiftDays(analyzeState.mode);

    const end = new Date(current.start);
    end.setDate(end.getDate() - 1);
    const start = new Date(end);
    start.setDate(start.getDate() - (days - 1));

    return { start, end };
}

function inRange(dateValue, period) {
    const date = parseDateOnly(dateValue);
    return date >= period.start && date <= period.end;
}

function getPeriodLabel(period) {
    if (analyzeState.mode === "month") {
        return `${monthNamesRu[period.start.getMonth()]} ${period.start.getFullYear()}`;
    }

    if (analyzeState.mode === "year") {
        return `${period.start.getFullYear()} год`;
    }

    if (period.start.getTime() === period.end.getTime()) {
        return period.start.toLocaleDateString("ru-RU");
    }

    return `${period.start.toLocaleDateString("ru-RU")} - ${period.end.toLocaleDateString("ru-RU")}`;
}

function collectTotals(transactions, period) {
    const totals = {
        income: 0,
        expense: 0,
        budget: 0,
        byCategoryIncome: new Map(),
        byCategoryExpense: new Map(),
    };

    transactions.forEach((tx) => {
        if (!inRange(tx.date, period)) {
            return;
        }

        const category = tx.categoryName || "Без категории";
        if (tx.type === 0) {
            totals.income += tx.amount;
            totals.byCategoryIncome.set(category, (totals.byCategoryIncome.get(category) || 0) + tx.amount);
        } else if (tx.type === 1) {
            totals.expense += tx.amount;
            totals.byCategoryExpense.set(category, (totals.byCategoryExpense.get(category) || 0) + tx.amount);
        }
    });

    totals.budget = totals.income - totals.expense;
    return totals;
}

function updateSummary(currentTotals, currency) {
    const incomeEl = document.getElementById("summaryIncome");
    const expenseEl = document.getElementById("summaryExpense");
    const budgetEl = document.getElementById("summaryBudget");

    incomeEl.textContent = formatAmount(currentTotals.income, currency);
    expenseEl.textContent = formatAmount(-currentTotals.expense, currency);
    budgetEl.textContent = formatAmount(currentTotals.budget, currency);

    budgetEl.classList.toggle("negative", currentTotals.budget < 0);
}

function categoryMetaMap(categories) {
    const map = new Map();
    categories.forEach((category) => {
        map.set(category.name, category);
    });
    return map;
}

function renderCategoryTable(currentTotals, previousTotals, categories) {
    const tableEl = document.getElementById("analyzeTable");
    const categoryMeta = categoryMetaMap(categories);

    const allIncomeKeys = new Set([
        ...currentTotals.byCategoryIncome.keys(),
        ...previousTotals.byCategoryIncome.keys(),
    ]);

    const allExpenseKeys = new Set([
        ...currentTotals.byCategoryExpense.keys(),
        ...previousTotals.byCategoryExpense.keys(),
    ]);

    function renderRows(keys, currentMap, previousMap, sign) {
        return [...keys]
            .sort((a, b) => {
                const bValue = currentMap.get(b) || 0;
                const aValue = currentMap.get(a) || 0;
                return bValue - aValue;
            })
            .map((key) => {
                const meta = categoryMeta.get(key) || {
                    name: key,
                    color: "#9aa3af",
                    iconHtml: '<i class="fa-solid fa-layer-group"></i>',
                };
                const current = currentMap.get(key) || 0;
                const previous = previousMap.get(key) || 0;
                const currentValue = sign === "negative" ? -Math.abs(current) : Math.abs(current);
                const previousValue = sign === "negative" ? -Math.abs(previous) : Math.abs(previous);

                return `
                    <div class="analyze-row">
                        <div class="analyze-row-name">
                            <span class="analyze-icon" style="background-color: ${meta.color || "#9aa3af"};">${meta.iconHtml}</span>
                            <span>${meta.name}</span>
                        </div>
                        <div class="analyze-row-value">${formatAmount(currentValue)}</div>
                        <div class="analyze-row-value prev">${formatAmount(previousValue)}</div>
                    </div>
                `;
            })
            .join("");
    }

    const incomeRows = renderRows(
        allIncomeKeys,
        currentTotals.byCategoryIncome,
        previousTotals.byCategoryIncome,
        "positive"
    );

    const expenseRows = renderRows(
        allExpenseKeys,
        currentTotals.byCategoryExpense,
        previousTotals.byCategoryExpense,
        "negative"
    );

    tableEl.innerHTML = `
        <div class="analyze-section">
            <div class="analyze-section-title">
                <span>Общий доход</span>
                <span class="section-amount">${formatAmount(currentTotals.income)}</span>
                <span class="section-amount prev">${formatAmount(previousTotals.income)}</span>
            </div>
            ${incomeRows || '<div class="analyze-empty">Нет данных за период</div>'}
        </div>
        <div class="analyze-section">
            <div class="analyze-section-title expense">
                <span>Общие расходы</span>
                <span class="section-amount">${formatAmount(-currentTotals.expense)}</span>
                <span class="section-amount prev">${formatAmount(-previousTotals.expense)}</span>
            </div>
            ${expenseRows || '<div class="analyze-empty">Нет данных за период</div>'}
        </div>
    `;
}

function chartLabels(period) {
    if (analyzeState.mode === "year") {
        return ["Янв", "Фев", "Мар", "Апр", "Май", "Июн", "Июл", "Авг", "Сен", "Окт", "Ноя", "Дек"];
    }

    const labels = [];
    const cursor = new Date(period.start);
    while (cursor <= period.end) {
        labels.push(cursor.toLocaleDateString("ru-RU", { day: "2-digit", month: "2-digit" }));
        cursor.setDate(cursor.getDate() + 1);
    }
    return labels;
}

function aggregateForCharts(transactions, period) {
    const labels = chartLabels(period);
    const income = labels.map(() => 0);
    const expense = labels.map(() => 0);
    const budget = labels.map(() => 0);

    const labelToIndex = new Map();
    labels.forEach((label, i) => labelToIndex.set(label, i));

    transactions.forEach((tx) => {
        if (!inRange(tx.date, period)) {
            return;
        }

        const date = parseDateOnly(tx.date);
        const key = analyzeState.mode === "year"
            ? formatMonthShort(date)
            : date.toLocaleDateString("ru-RU", { day: "2-digit", month: "2-digit" });

        const index = labelToIndex.get(key);
        if (index === undefined) {
            return;
        }

        if (tx.type === 0) {
            income[index] += tx.amount;
        } else if (tx.type === 1) {
            expense[index] += tx.amount;
        }
    });

    let running = 0;
    for (let i = 0; i < labels.length; i += 1) {
        running += income[i] - expense[i];
        budget[i] = running;
    }

    return { labels, income, expense, budget };
}

function destroyCharts() {
    if (analyzeState.flowChart) {
        analyzeState.flowChart.destroy();
        analyzeState.flowChart = null;
    }
    if (analyzeState.budgetChart) {
        analyzeState.budgetChart.destroy();
        analyzeState.budgetChart = null;
    }
}

function buildCharts(transactions, period) {
    const chartData = aggregateForCharts(transactions, period);
    const flowCanvas = document.getElementById("flowChart");
    const budgetCanvas = document.getElementById("budgetChart");

    destroyCharts();

    analyzeState.flowChart = new Chart(flowCanvas, {
        type: "bar",
        data: {
            labels: chartData.labels,
            datasets: [
                {
                    label: "Доход",
                    data: chartData.income,
                    backgroundColor: "rgba(76, 214, 122, 0.8)",
                    borderRadius: 6,
                },
                {
                    label: "Расход",
                    data: chartData.expense,
                    backgroundColor: "rgba(223, 68, 68, 0.78)",
                    borderRadius: 6,
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: { legend: { position: "top" } },
        },
    });

    analyzeState.budgetChart = new Chart(budgetCanvas, {
        type: "line",
        data: {
            labels: chartData.labels,
            datasets: [
                {
                    label: "Баланс",
                    data: chartData.budget,
                    borderColor: "#1f3c88",
                    backgroundColor: "rgba(31, 60, 136, 0.14)",
                    fill: true,
                    tension: 0.3,
                },
            ],
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
            plugins: { legend: { position: "top" } },
        },
    });
}

function updateDashboard() {
    const payload = window.analyzePayload || { transactions: [], categories: [] };
    const transactions = payload.transactions || [];
    const categories = payload.categories || [];
    const currency = transactions[0]?.currency || "₽";

    const current = getCurrentPeriodBounds();
    const previous = getPreviousPeriodBounds(current);

    const currentTotals = collectTotals(transactions, current);
    const previousTotals = collectTotals(transactions, previous);

    document.getElementById("currentPeriodLabel").textContent = getPeriodLabel(current);
    document.getElementById("previousPeriodLabel").textContent = getPeriodLabel(previous);

    updateSummary(currentTotals, currency);
    renderCategoryTable(currentTotals, previousTotals, categories);

    const chartsPanel = document.getElementById("chartsPanel");
    if (analyzeState.view === "graph") {
        buildCharts(transactions, current);
    }
}

function setAnalyzeView(view) {
    analyzeState.view = view;

    const reportPanel = document.getElementById("reportPanel");
    const chartsPanel = document.getElementById("chartsPanel");
    const reportBtn = document.getElementById("reportViewBtn");
    const graphBtn = document.getElementById("graphViewBtn");

    const isReport = view === "report";

    reportPanel.hidden = !isReport;
    chartsPanel.hidden = isReport;

    reportBtn.classList.toggle("active", isReport);
    graphBtn.classList.toggle("active", !isReport);
    reportBtn.setAttribute("aria-selected", String(isReport));
    graphBtn.setAttribute("aria-selected", String(!isReport));

    if (!isReport) {
        updateDashboard();
    }
}

function initAnalyzeControls() {
    const chips = document.querySelectorAll(".period-chip");
    const customPanel = document.getElementById("customPeriodPanel");
    const customStart = document.getElementById("customStartDate");
    const customEnd = document.getElementById("customEndDate");

    const today = new Date();
    const isoToday = toISODate(today);
    customStart.value = isoToday;
    customEnd.value = isoToday;
    analyzeState.customStart = isoToday;
    analyzeState.customEnd = isoToday;

    chips.forEach((chip) => {
        chip.addEventListener("click", () => {
            chips.forEach((btn) => btn.classList.remove("active"));
            chip.classList.add("active");
            analyzeState.mode = chip.dataset.mode;

            customPanel.hidden = analyzeState.mode !== "custom";
            updateDashboard();
        });
    });

    customStart.addEventListener("change", (event) => {
        analyzeState.customStart = event.target.value;
        if (analyzeState.mode === "custom") {
            updateDashboard();
        }
    });

    customEnd.addEventListener("change", (event) => {
        analyzeState.customEnd = event.target.value;
        if (analyzeState.mode === "custom") {
            updateDashboard();
        }
    });

    const reportBtn = document.getElementById("reportViewBtn");
    const graphBtn = document.getElementById("graphViewBtn");

    reportBtn.addEventListener("click", () => {
        setAnalyzeView("report");
    });

    graphBtn.addEventListener("click", () => {
        setAnalyzeView("graph");
    });

    setAnalyzeView("report");
}

document.addEventListener("DOMContentLoaded", () => {
    initAnalyzeControls();
    updateDashboard();
});
