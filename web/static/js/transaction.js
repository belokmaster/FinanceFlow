// creating transaction
const createTransactionModal = document.getElementById('createTransactionModal');
const createTransactionAccountSelect = document.getElementById('transactionAccount');
const createTransactionCategorySelect = document.getElementById('transactionCategory');
const createTransactionAmountInput = document.getElementById('transactionAmount');
const createTransactionCommentInput = document.getElementById('transactionDescription');
const createTransactionDateInput = document.getElementById('transactionDate');

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = "flex";
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = "none";
    }
}

function openCreateTransactionModal() {
    document.getElementById('transactionAccount').value = '';
    document.getElementById('transactionCategory').value = '';
    createTransactionAmountInput.value = '';
    document.getElementById('transactionDescription').value = '';
    document.getElementById('transactionDate').value = '';

    const typeButtons = document.querySelectorAll('.type-btn');
    typeButtons.forEach(btn => btn.classList.remove('active'));
    document.querySelector('.income-btn').classList.add('active');
    document.getElementById('transactionType').value = '0';

    const today = new Date().toISOString().split('T')[0];
    document.getElementById('transactionDate').value = today;

    // automatic point 0.00 to amount
    const createAmmountForm = document.querySelector('form[action="/submit_transaction"]');
    if (createAmmountForm) {
        createAmmountForm.onsubmit = function () {
            const amountInput = document.getElementById('transactionAmount');
            if (!amountInput.value.trim()) {
                amountInput.value = '0.00';
            }
        };
    }

    openModal('createTransactionModal');
}


function updateSubcategories() {
    const categorySelect = document.getElementById('transactionCategory');
    const subCategorySelect = document.getElementById('transactionSubCategory');
    const selectedCategoryId = categorySelect.value;

    const allSubOptions = subCategorySelect.querySelectorAll('option[data-parent]');
    allSubOptions.forEach(option => {
        option.style.display = 'none';
    });

    if (selectedCategoryId) {
        const relevantSubOptions = subCategorySelect.querySelectorAll(`option[data-parent="${selectedCategoryId}"]`);
        relevantSubOptions.forEach(option => {
            option.style.display = 'block';
        });
    }

    subCategorySelect.value = "";
}

document.addEventListener('DOMContentLoaded', function () {
    const categorySelect = document.getElementById('transactionCategory');
    if (categorySelect) {
        categorySelect.addEventListener('change', updateSubcategories);
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const typeButtons = document.querySelectorAll('.type-btn');
    const typeInput = document.getElementById('transactionType');

    typeButtons.forEach(button => {
        button.addEventListener('click', function () {
            typeButtons.forEach(btn => btn.classList.remove('active'));
            this.classList.add('active');
            typeInput.value = this.getAttribute('data-type');
        });
    });
});

// close when click non a model window
document.addEventListener('click', function (e) {
    const modals = document.getElementsByClassName("modal");
    for (let i = 0; i < modals.length; i++) {
        if (e.target == modals[i]) {
            modals[i].style.display = "none";
        }
    }
});

// close when tap escape
document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' || e.key === 'Esc') {
        const modals = document.getElementsByClassName("modal");
        for (let i = 0; i < modals.length; i++) {
            modals[i].style.display = "none";
        }
    }
});

createTransactionAmountInput.addEventListener('input', formatAmountInput);

function formatAmountInput(e) {
    let value = e.target.value;

    // regular for everything without figures and point
    value = value.replace(/[^\d.]/g, '');
    const parts = value.split('.');
    if (parts.length > 2) {
        value = parts[0] + '.' + parts.slice(1).join('');
    }

    if (parts.length === 2 && parts[1].length > 2) {
        value = parts[0] + '.' + parts[1].substring(0, 2);
    }

    e.target.value = value;
}

function toggleAccountDropdown() {
    const options = document.getElementById('accountOptions');
    const selected = document.querySelector('.select-selected-account');

    if (options && selected) {
        const isShowing = options.classList.contains('show');

        closeAllDropdowns();

        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function toggleCategoryDropdown() {
    const options = document.getElementById('categoryOptions');
    const selected = document.querySelector('.select-selected-category');

    if (options && selected) {
        const isShowing = options.classList.contains('show');

        closeAllDropdowns();

        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function selectAccountOption(optionElement) {
    const accountId = optionElement.getAttribute('data-account-id');
    const accountName = optionElement.getAttribute('data-account-name');
    const accountBalance = optionElement.getAttribute('data-account-balance');
    const accountColor = optionElement.getAttribute('data-account-color');
    const accountIconElement = optionElement.querySelector('.option-account-icon');
    const accountIconHTML = accountIconElement.innerHTML;

    document.getElementById('transactionAccount').value = accountId;

    const selectedIcon = document.getElementById('selectedAccountIcon');
    const selectedName = document.getElementById('selectedAccountName');
    const selectedBalance = document.getElementById('selectedAccountBalance');

    if (selectedIcon && selectedName && selectedBalance) {
        selectedIcon.innerHTML = accountIconHTML;
        selectedIcon.style.backgroundColor = accountColor;

        selectedIcon.style.display = 'flex';
        selectedIcon.style.alignItems = 'center';
        selectedIcon.style.justifyContent = 'center';
        selectedIcon.style.borderRadius = '10px';
        selectedIcon.style.width = '40px';
        selectedIcon.style.height = '40px';
        selectedIcon.style.fontSize = '18px';

        selectedName.textContent = accountName;
        selectedBalance.textContent = accountBalance;
    }

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('.select-account-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');
}

function selectCategoryOption(optionElement) {
    const categoryId = optionElement.getAttribute('data-category-id');
    const categoryName = optionElement.getAttribute('data-category-name');
    const categoryColor = optionElement.getAttribute('data-category-color');
    const categoryIconElement = optionElement.querySelector('.option-category-icon');
    const categoryIconHTML = categoryIconElement.innerHTML;

    document.getElementById('transactionCategory').value = categoryId;

    const selectedIcon = document.getElementById('selectedCategoryIcon');
    const selectedName = document.getElementById('selectedCategoryName');

    if (selectedIcon && selectedName) {
        selectedIcon.innerHTML = categoryIconHTML;
        selectedIcon.style.backgroundColor = categoryColor;

        selectedIcon.style.display = 'flex';
        selectedIcon.style.alignItems = 'center';
        selectedIcon.style.justifyContent = 'center';
        selectedIcon.style.borderRadius = '10px';
        selectedIcon.style.width = '40px';
        selectedIcon.style.height = '40px';
        selectedIcon.style.fontSize = '18px';

        selectedName.textContent = categoryName;
    }

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('.select-category-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');

    updateSubcategories();
}

function closeAllDropdowns() {
    const accountOptions = document.getElementById('accountOptions');
    const accountSelected = document.querySelector('.select-selected-account');
    const categoryOptions = document.getElementById('categoryOptions');
    const categorySelected = document.querySelector('.select-selected-category');

    if (accountOptions) accountOptions.classList.remove('show');
    if (accountSelected) accountSelected.classList.remove('active');
    if (categoryOptions) categoryOptions.classList.remove('show');
    if (categorySelected) categorySelected.classList.remove('active');
}

document.addEventListener('click', function (event) {
    const accountSelectContainer = document.querySelector('.custom-account-select');
    const categorySelectContainer = document.querySelector('.custom-category-select');

    if ((!accountSelectContainer.contains(event.target)) &&
        (!categorySelectContainer.contains(event.target))) {
        closeAllDropdowns();
    }
});

document.addEventListener('DOMContentLoaded', function () {
    const accountOptions = document.querySelectorAll('.select-account-option');
    accountOptions.forEach(option => {
        option.addEventListener('click', function (e) {
            e.stopPropagation();
            selectAccountOption(this);
        });
    });

    document.addEventListener('click', function (event) {
        const selectContainer = document.getElementById('accountSelect');
        if (selectContainer && !selectContainer.contains(event.target)) {
            closeAllDropdowns();
        }
    });

    const firstAccountOption = document.querySelector('.select-account-option');
    if (firstAccountOption) {
        selectAccountOption(firstAccountOption);
    }

    const categoryOptions = document.querySelectorAll('.select-category-option');
    categoryOptions.forEach(option => {
        option.addEventListener('click', function (e) {
            e.stopPropagation();
            selectCategoryOption(this);
        });
    });

    const firstCategoryOption = document.querySelector('.select-category-option');
    if (firstCategoryOption) {
        selectCategoryOption(firstCategoryOption);
    }
});

function openCreateTransactionModal() {
    document.getElementById('transactionAccount').value = '';
    document.getElementById('transactionCategory').value = '';
    createTransactionAmountInput.value = '';
    document.getElementById('transactionDescription').value = '';
    document.getElementById('transactionDate').value = '';

    const typeButtons = document.querySelectorAll('.type-btn');
    typeButtons.forEach(btn => btn.classList.remove('active'));
    document.querySelector('.income-btn').classList.add('active');
    document.getElementById('transactionType').value = '0';

    const today = new Date().toISOString().split('T')[0];
    document.getElementById('transactionDate').value = today;

    setTimeout(() => {
        const firstAccountOption = document.querySelector('.select-account-option');
        if (firstAccountOption) {
            selectAccountOption(firstAccountOption);
        }

        const firstCategoryOption = document.querySelector('.select-category-option');
        if (firstCategoryOption) {
            selectCategoryOption(firstCategoryOption);
        }
    }, 100);

    // automatic point 0.00 to amount
    const createAmmountForm = document.querySelector('form[action="/submit_transaction"]');
    if (createAmmountForm) {
        createAmmountForm.onsubmit = function () {
            const amountInput = document.getElementById('transactionAmount');
            if (!amountInput.value.trim()) {
                amountInput.value = '0.00';
            }
        };
    }

    openModal('createTransactionModal');
}