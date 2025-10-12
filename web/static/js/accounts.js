// for edit account
const editModal = document.getElementById('accountModal');
const editColorInput = document.getElementById('editAccountColor');
const editColorPreview = document.getElementById('colorPreview');
const editColorHexValue = document.getElementById('colorHexValue');
const editIconSelect = editModal.querySelector('.custom-icon-select');
const editSelectedIconDisplay = document.getElementById('selectedIconDisplay');
const editIconOptionsContainer = document.getElementById('iconOptions');
const editHiddenIconInput = document.getElementById('editAccountIcon');
const editAccountBalanceInput = document.getElementById('editAccountBalance');

// creating account
const createModal = document.getElementById('createAccountModal');
const createColorInput = document.getElementById('createAccountColor');
const createColorPreview = document.getElementById('createColorPreview');
const createColorHexValue = document.getElementById('createColorHexValue');
const createIconSelect = createModal.querySelector('.custom-icon-select');
const createSelectedIconDisplay = document.getElementById('createSelectedIconDisplay');
const createIconOptionsContainer = document.getElementById('createIconOptions');
const createHiddenIconInput = document.getElementById('createAccountIcon');
const createAccountBalanceInput = document.getElementById('createAccountBalance');

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = "flex"; // make it visible
    }
}

function closeModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.style.display = "none";
    }
}

// delete it later when account will done
function showForm(formId) {
    const formsToHide = ['transferForm', 'transactionForm', 'createCategoryForm', 'deleteCategoryForm', 'createSubCategoryForm', 'deleteSubCategoryForm', 'deleteAccountForm'];
    formsToHide.forEach(id => {
        const formElement = document.getElementById(id);
        if (formElement) {
            formElement.style.display = 'none';
        }
    });
    const formToShow = document.getElementById(formId);
    if (formToShow) {
        formToShow.style.display = 'block';
    }
}

function updateColorDisplay(newColor, previewEl, hexEl) {
    if (!newColor || !newColor.startsWith('#')) {
        newColor = '#000000';
    }
    previewEl.style.backgroundColor = newColor;
    hexEl.textContent = newColor.toUpperCase();
}

function openCreateModal() {
    // zero values
    document.getElementById('createAccountName').value = '';
    createAccountBalanceInput.value = '0.00';
    document.getElementById('createAccountCurrency').value = '0';
    createColorInput.value = '#4cd67a';
    createHiddenIconInput.value = '';

    createSelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = '';
    createSelectedIconDisplay.querySelector('.selected-icon-key').textContent = 'Выберите иконку';

    updateColorDisplay('#4cd67a', createColorPreview, createColorHexValue);

    openModal('createAccountModal');
}

function openAccountModal(accountId, accountName, accountColor, accountIconKey, accountBalance, accountCurrency) {
    // fil
    document.getElementById('editAccountId').value = accountId;
    document.getElementById('editAccountName').value = accountName;
    document.getElementById('editAccountCurrency').value = accountCurrency;

    let formattedBalance = parseFloat(accountBalance);
    editAccountBalanceInput.value = !isNaN(formattedBalance) ? formattedBalance.toFixed(2) : '0.00';

    editColorInput.value = accountColor;
    updateColorDisplay(accountColor, editColorPreview, editColorHexValue);
    editHiddenIconInput.value = accountIconKey;

    const targetOption = editIconOptionsContainer.querySelector(`.select-icon-option[data-key="${accountIconKey}"]`);
    if (targetOption) {
        const iconSvgHTML = targetOption.querySelector('.option-icon-svg').innerHTML;
        editSelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = iconSvgHTML;
        editSelectedIconDisplay.querySelector('.selected-icon-key').textContent = accountIconKey;
    } else {
        editSelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = '';
        editSelectedIconDisplay.querySelector('.selected-icon-key').textContent = 'Выберите иконку';
    }

    openModal('accountModal');
}

function deleteAccount() {
    const accountId = document.getElementById('editAccountId').value;
    if (confirm('Вы уверены, что хотите удалить этот счет? Это действие необратимо.')) {
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/delete_account';

        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'ID';
        input.value = accountId;

        form.appendChild(input);
        document.body.appendChild(form);
        form.submit();
    }
}

function formatBalanceInput(e) {
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

createColorInput.addEventListener('input', () => {
    updateColorDisplay(createColorInput.value, createColorPreview, createColorHexValue);
});

createSelectedIconDisplay.addEventListener('click', (e) => {
    e.stopPropagation();
    createIconOptionsContainer.classList.toggle('show');
    createIconSelect.classList.toggle('active');
});

createIconOptionsContainer.addEventListener('click', (e) => {
    const option = e.target.closest('.select-icon-option');
    if (option) {
        const iconKey = option.dataset.key;
        const iconSvgHTML = option.querySelector('.option-icon-svg').innerHTML;
        const iconText = option.querySelector('.option-icon-key').textContent;

        createHiddenIconInput.value = iconKey;
        createSelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = iconSvgHTML;
        createSelectedIconDisplay.querySelector('.selected-icon-key').textContent = iconText;

        createIconOptionsContainer.classList.remove('show');
        createIconSelect.classList.remove('active');
    }
});

// edit accont block
editColorInput.addEventListener('input', () => {
    updateColorDisplay(editColorInput.value, editColorPreview, editColorHexValue);
});

editSelectedIconDisplay.addEventListener('click', (e) => {
    e.stopPropagation();
    editIconOptionsContainer.classList.toggle('show');
    editIconSelect.classList.toggle('active');
});

editIconOptionsContainer.addEventListener('click', (e) => {
    const option = e.target.closest('.select-icon-option');
    if (option) {
        const iconKey = option.dataset.key;
        const iconSvgHTML = option.querySelector('.option-icon-svg').innerHTML;

        editHiddenIconInput.value = iconKey;
        editSelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = iconSvgHTML;
        editSelectedIconDisplay.querySelector('.selected-icon-key').textContent = iconKey;

        editIconOptionsContainer.classList.remove('show');
        editIconSelect.classList.remove('active');
    }
});

createAccountBalanceInput.addEventListener('input', formatBalanceInput);
editAccountBalanceInput.addEventListener('input', formatBalanceInput);

// close when click non a model window
document.addEventListener('click', function (e) {
    const modals = document.getElementsByClassName("modal");
    for (let i = 0; i < modals.length; i++) {
        if (e.target == modals[i]) {
            modals[i].style.display = "none";
        }
    }

    // closed that windows
    if (editIconSelect && !editIconSelect.contains(e.target)) {
        editIconOptionsContainer.classList.remove('show');
        editIconSelect.classList.remove('active');
    }

    if (createIconSelect && !createIconSelect.contains(e.target)) {
        createIconOptionsContainer.classList.remove('show');
        createIconSelect.classList.remove('active');
    }
});

// close when tap escape
document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' || e.key === 'Esc') {
        const modals = document.getElementsByClassName("modal");
        for (let i = 0; i < modals.length; i++) {
            modals[i].style.display = "none";
        }

        if (editIconSelect) {
            editIconOptionsContainer.classList.remove('show');
            editIconSelect.classList.remove('active');
        }

        if (createIconSelect) {
            createIconOptionsContainer.classList.remove('show');
            createIconSelect.classList.remove('active');
        }
    }
});