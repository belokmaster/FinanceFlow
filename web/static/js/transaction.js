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

    const today = new Date();
    const formattedDate = today.toISOString().split('T')[0];
    document.getElementById('transactionDate').value = formattedDate;

    // automatic point 0.00 to amount
    const createForm = document.querySelector('form[action="/submit_transaction"]');
    if (createForm) {
        createForm.onsubmit = function (e) {
            const type = document.getElementById('transactionType').value;
            const amount = document.getElementById('transactionAmount').value.trim();

            if (type === '0' || type === '1') { // Доход или Расход
                const account = document.getElementById('transactionAccount').value;
                const category = document.getElementById('transactionCategory').value;

                if (!account) {
                    alert('Выберите счет!');
                    e.preventDefault();
                    return false;
                }

                if (!category) {
                    alert('Выберите категорию!');
                    e.preventDefault();
                    return false;
                }

            } else if (type === '2') { // Перевод
                const fromAccount = document.getElementById('fromAccount').value;
                const toAccount = document.getElementById('toAccount').value;

                if (!fromAccount) {
                    alert('Выберите счет списания!');
                    e.preventDefault();
                    return false;
                }

                if (!toAccount) {
                    alert('Выберите счет зачисления!');
                    e.preventDefault();
                    return false;
                }

                if (fromAccount === toAccount) {
                    alert('Счета списания и зачисления не могут быть одинаковыми!');
                    e.preventDefault();
                    return false;
                }
            }

            if (amount === '') {
                document.getElementById('transactionAmount').value = '0.00';
            }
        };
    }

    openModal('createTransactionModal');
}

function updateSubcategories() {
    const categorySelect = document.getElementById('transactionCategory');
    const selectedCategoryId = categorySelect.value;
    const allSubOptions = document.querySelectorAll('.select-subcategory-option');

    allSubOptions.forEach(option => {
        option.style.display = 'none';
    });

    if (selectedCategoryId) {
        const relevantSubOptions = document.querySelectorAll(`.select-subcategory-option[data-parent-id="${selectedCategoryId}"]`);
        relevantSubOptions.forEach(option => {
            option.style.display = 'flex';
        });
    }

    document.getElementById('transactionSubCategory').value = '';
    document.getElementById('selectedSubcategoryIcon').innerHTML = '';
    document.getElementById('selectedSubcategoryIcon').style.backgroundColor = '';
    document.getElementById('selectedSubcategoryName').textContent = 'Выберите подкатегорию';

    const allSubcategoryOptions = document.querySelectorAll('.select-subcategory-option');
    allSubcategoryOptions.forEach(option => {
        option.classList.remove('selected');
    });
}

function updateEditSubcategories() {
    const categorySelect = document.getElementById('editTransactionCategory');
    const selectedCategoryId = categorySelect.value;
    const allSubOptions = document.querySelectorAll('#editSubcategoryOptions .select-subcategory-option');

    allSubOptions.forEach(option => {
        option.style.display = 'none';
    });

    if (selectedCategoryId) {
        const relevantSubOptions = document.querySelectorAll(`#editSubcategoryOptions .select-subcategory-option[data-parent-id="${selectedCategoryId}"]`);
        relevantSubOptions.forEach(option => {
            option.style.display = 'flex';
        });
    }

    document.getElementById('editTransactionSubCategory').value = '';
    document.getElementById('editSelectedSubcategoryIcon').innerHTML = '';
    document.getElementById('editSelectedSubcategoryIcon').style.backgroundColor = '';
    document.getElementById('editSelectedSubcategoryName').textContent = 'Выберите подкатегорию';

    const allSubcategoryOptions = document.querySelectorAll('#editSubcategoryOptions .select-subcategory-option');
    allSubcategoryOptions.forEach(option => {
        option.classList.remove('selected');
    });
}

document.addEventListener('DOMContentLoaded', function () {
    const createCategorySelect = document.getElementById('transactionCategory');
    if (createCategorySelect) {
        createCategorySelect.addEventListener('change', updateSubcategories);
    }

    const createSubcategoryOptions = document.querySelectorAll('#subcategoryOptions .select-subcategory-option');
    createSubcategoryOptions.forEach(option => {
        option.addEventListener('click', function (e) {
            e.stopPropagation();
            selectSubcategoryOption(this);
        });
    });

    const editCategorySelect = document.getElementById('editTransactionCategory');
    if (editCategorySelect) {
        editCategorySelect.addEventListener('change', updateEditSubcategories);
    }

    const editSubcategoryOptions = document.querySelectorAll('#editSubcategoryOptions .select-subcategory-option');
    editSubcategoryOptions.forEach(option => {
        option.addEventListener('click', function (e) {
            e.stopPropagation();
            selectEditSubcategoryOption(this);
        });
    });

    const createTypeButtons = document.querySelectorAll('#createTransactionModal .type-btn');
    const createTypeInput = document.getElementById('transactionType');

    createTypeButtons.forEach(button => {
        button.addEventListener('click', function () {
            createTypeButtons.forEach(btn => btn.classList.remove('active'));
            this.classList.add('active');
            createTypeInput.value = this.getAttribute('data-type');
        });
    });

    const editTypeButtons = document.querySelectorAll('#editTransactionModal .type-btn');
    const editTypeInput = document.getElementById('editTransactionType');

    editTypeButtons.forEach(button => {
        button.addEventListener('click', function () {
            editTypeButtons.forEach(btn => btn.classList.remove('active'));
            this.classList.add('active');
            editTypeInput.value = this.getAttribute('data-type');
        });
    });

    const amountInputs = document.querySelectorAll('input[type="text"][name="Amount"]');
    amountInputs.forEach(input => {
        input.addEventListener('input', formatAmountInput);
    });

    const accountOptions = document.querySelectorAll('.select-account-option');
    accountOptions.forEach(option => {
        option.addEventListener('click', function (e) {
            e.stopPropagation();
            const container = this.closest('.select-account-options');
            let prefix = 'income';

            if (container && container.id === 'fromAccountOptions') {
                prefix = 'from';
            } else if (container && container.id === 'toAccountOptions') {
                prefix = 'to';
            }

            selectAccountOption(this, prefix);
        });
    });

    const categoryOptions = document.querySelectorAll('.select-category-option');
    categoryOptions.forEach(option => {
        option.addEventListener('click', function (e) {
            e.stopPropagation();
            if (this.closest('#editCategoryOptions')) {
                selectEditCategoryOption(this);
            } else {
                selectCategoryOption(this);
            }
        });
    });

    const transactionCards = document.querySelectorAll('.transaction-card');
    transactionCards.forEach(card => {
        card.addEventListener('click', function (e) {
            e.stopPropagation();

            const transactionId = this.getAttribute('data-transaction-id');
            const transactionType = parseInt(this.getAttribute('data-transaction-type'));
            const amount = this.getAttribute('data-amount');
            const accountId = this.getAttribute('data-account-id');
            const accountName = this.getAttribute('data-account-name');
            const accountColor = this.getAttribute('data-account-color');
            const categoryId = this.getAttribute('data-category-id');
            const categoryName = this.getAttribute('data-category-name');
            const categoryColor = this.getAttribute('data-category-color');
            const subCategoryId = this.getAttribute('data-subcategory-id') || '';
            const description = this.getAttribute('data-description') || '';
            const date = this.getAttribute('data-date');

            console.log('Opening edit modal for transaction:', transactionId);

            openEditTransactionModal(
                transactionId,
                transactionType,
                amount,
                accountId,
                accountName,
                accountColor,
                categoryId,
                categoryName,
                categoryColor,
                subCategoryId,
                description,
                date
            );
        });

        card.style.cursor = 'pointer';
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

function formatAmountInput(e) {
    let value = e.target.value;

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

function toggleAccountDropdown(prefix = 'income') {
    let options, selected;

    if (prefix === 'from') {
        options = document.querySelector('#fromAccount + .select-selected-account').nextElementSibling;
        selected = document.querySelector('#fromAccount + .select-selected-account');
    } else if (prefix === 'to') {
        options = document.querySelector('#toAccount + .select-selected-account').nextElementSibling;
        selected = document.querySelector('#toAccount + .select-selected-account');
    } else {
        options = document.getElementById('accountOptions');
        selected = document.querySelector('#createTransactionModal .select-selected-account');
    }

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
    const selected = document.querySelector('#createTransactionModal .select-selected-category');

    if (options && selected) {
        const isShowing = options.classList.contains('show');
        closeAllDropdowns();
        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function toggleSubcategoryDropdown() {
    const options = document.getElementById('subcategoryOptions');
    const selected = document.querySelector('#createTransactionModal .select-selected-subcategory');

    if (options && selected) {
        const isShowing = options.classList.contains('show');
        closeAllDropdowns();
        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function toggleEditAccountDropdown() {
    const options = document.getElementById('editAccountOptions');
    const selected = document.querySelector('#editTransactionModal .select-selected-account');

    if (options && selected) {
        const isShowing = options.classList.contains('show');
        closeAllDropdowns();
        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function toggleEditCategoryDropdown() {
    const options = document.getElementById('editCategoryOptions');
    const selected = document.querySelector('#editTransactionModal .select-selected-category');

    if (options && selected) {
        const isShowing = options.classList.contains('show');
        closeAllDropdowns();
        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function toggleEditSubcategoryDropdown() {
    const options = document.getElementById('editSubcategoryOptions');
    const selected = document.querySelector('#editTransactionModal .select-selected-subcategory');

    if (options && selected) {
        const isShowing = options.classList.contains('show');
        closeAllDropdowns();
        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function selectAccountOption(optionElement, prefix = 'income') {
    const accountId = optionElement.getAttribute('data-account-id');
    const accountName = optionElement.getAttribute('data-account-name');
    const accountBalance = optionElement.getAttribute('data-account-balance');
    const accountColor = optionElement.getAttribute('data-account-color');
    const accountIconElement = optionElement.querySelector('.option-account-icon');
    const accountIconHTML = accountIconElement.innerHTML;

    let accountInput, selectedIcon, selectedName, selectedBalance;

    if (prefix === 'from') {
        accountInput = document.getElementById('fromAccount');
        selectedIcon = document.getElementById('fromSelectedAccountIcon');
        selectedName = document.getElementById('fromSelectedAccountName');
        selectedBalance = document.getElementById('fromSelectedAccountBalance');
        optionsContainer = document.getElementById('fromAccountOptions');
    } else if (prefix === 'to') {
        accountInput = document.getElementById('toAccount');
        selectedIcon = document.getElementById('toSelectedAccountIcon');
        selectedName = document.getElementById('toSelectedAccountName');
        selectedBalance = document.getElementById('toSelectedAccountBalance');
        optionsContainer = document.getElementById('toAccountOptions');
    } else {
        accountInput = document.getElementById('transactionAccount');
        selectedIcon = document.getElementById('selectedAccountIcon');
        selectedName = document.getElementById('selectedAccountName');
        selectedBalance = document.getElementById('selectedAccountBalance');
        optionsContainer = document.getElementById('accountOptions');
    }

    if (accountInput) {
        accountInput.value = accountId.trim();
        console.log(`Set ${prefix} account input to:`, accountId.trim());
    }

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

    if (optionsContainer) {
        const allOptions = optionsContainer.querySelectorAll('.select-account-option');
        allOptions.forEach(option => {
            option.classList.remove('selected');
        });
        optionElement.classList.add('selected');
    }
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

    const allOptions = document.querySelectorAll('#categoryOptions .select-category-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');

    updateSubcategories();
}

function selectSubcategoryOption(optionElement) {
    const subcategoryId = optionElement.getAttribute('data-subcategory-id');
    const subcategoryName = optionElement.getAttribute('data-subcategory-name');
    const subcategoryColor = optionElement.getAttribute('data-subcategory-color');
    const subcategoryIconElement = optionElement.querySelector('.option-subcategory-icon');
    const subcategoryIconHTML = subcategoryIconElement.innerHTML;

    document.getElementById('transactionSubCategory').value = subcategoryId;

    const selectedIcon = document.getElementById('selectedSubcategoryIcon');
    const selectedName = document.getElementById('selectedSubcategoryName');

    if (selectedIcon && selectedName) {
        selectedIcon.innerHTML = subcategoryIconHTML;
        selectedIcon.style.backgroundColor = subcategoryColor;

        selectedIcon.style.display = 'flex';
        selectedIcon.style.alignItems = 'center';
        selectedIcon.style.justifyContent = 'center';
        selectedIcon.style.borderRadius = '10px';
        selectedIcon.style.width = '40px';
        selectedIcon.style.height = '40px';
        selectedIcon.style.fontSize = '18px';

        selectedName.textContent = subcategoryName;
    }

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('#subcategoryOptions .select-subcategory-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');
}

function selectEditAccountOption(optionElement) {
    const accountId = optionElement.getAttribute('data-account-id');
    const accountName = optionElement.getAttribute('data-account-name');
    const accountBalance = optionElement.getAttribute('data-account-balance');
    const accountColor = optionElement.getAttribute('data-account-color');
    const accountIconHTML = optionElement.querySelector('.option-account-icon').innerHTML;

    document.getElementById('editTransactionAccount').value = accountId;

    const selectedIcon = document.getElementById('editSelectedAccountIcon');
    const selectedName = document.getElementById('editSelectedAccountName');
    const selectedBalance = document.getElementById('editSelectedAccountBalance');

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

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('#editAccountOptions .select-account-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');
}

function selectEditCategoryOption(optionElement) {
    const categoryId = optionElement.getAttribute('data-category-id');
    const categoryName = optionElement.getAttribute('data-category-name');
    const categoryColor = optionElement.getAttribute('data-category-color');
    const categoryIconHTML = optionElement.querySelector('.option-category-icon').innerHTML;

    document.getElementById('editTransactionCategory').value = categoryId;

    const selectedIcon = document.getElementById('editSelectedCategoryIcon');
    const selectedName = document.getElementById('editSelectedCategoryName');

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

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('#editCategoryOptions .select-category-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');

    updateEditSubcategories();
}

function selectEditSubcategoryOption(optionElement) {
    const subcategoryId = optionElement.getAttribute('data-subcategory-id');
    const subcategoryName = optionElement.getAttribute('data-subcategory-name');
    const subcategoryColor = optionElement.getAttribute('data-subcategory-color');
    const subcategoryIconHTML = optionElement.querySelector('.option-subcategory-icon').innerHTML;

    document.getElementById('editTransactionSubCategory').value = subcategoryId;

    const selectedIcon = document.getElementById('editSelectedSubcategoryIcon');
    const selectedName = document.getElementById('editSelectedSubcategoryName');

    selectedIcon.innerHTML = subcategoryIconHTML;
    selectedIcon.style.backgroundColor = subcategoryColor;

    selectedIcon.style.display = 'flex';
    selectedIcon.style.alignItems = 'center';
    selectedIcon.style.justifyContent = 'center';
    selectedIcon.style.borderRadius = '10px';
    selectedIcon.style.width = '40px';
    selectedIcon.style.height = '40px';
    selectedIcon.style.fontSize = '18px';

    selectedName.textContent = subcategoryName;

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('#editSubcategoryOptions .select-subcategory-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');
}

function closeAllDropdowns() {
    const accountOptions = document.getElementById('accountOptions');
    const accountSelected = document.querySelector('#createTransactionModal .select-selected-account');
    const categoryOptions = document.getElementById('categoryOptions');
    const categorySelected = document.querySelector('#createTransactionModal .select-selected-category');
    const subcategoryOptions = document.getElementById('subcategoryOptions');
    const subcategorySelected = document.querySelector('#createTransactionModal .select-selected-subcategory');

    const editAccountOptions = document.getElementById('editAccountOptions');
    const editAccountSelected = document.querySelector('#editTransactionModal .select-selected-account');
    const editCategoryOptions = document.getElementById('editCategoryOptions');
    const editCategorySelected = document.querySelector('#editTransactionModal .select-selected-category');
    const editSubcategoryOptions = document.getElementById('editSubcategoryOptions');
    const editSubcategorySelected = document.querySelector('#editTransactionModal .select-selected-subcategory');

    const fromAccountOptions = document.getElementById('fromAccountOptions');
    const fromAccountSelected = document.querySelector('#createTransactionModal #fromAccount + .select-selected-account');
    const toAccountOptions = document.getElementById('toAccountOptions');
    const toAccountSelected = document.querySelector('#createTransactionModal #toAccount + .select-selected-account');

    const editFromAccountOptions = document.getElementById('editFromAccountOptions');
    const editFromAccountSelected = document.querySelector('#editTransferModal #editFromAccount + .select-selected-account');
    const editToAccountOptions = document.getElementById('editToAccountOptions');
    const editToAccountSelected = document.querySelector('#editTransferModal #editToAccount + .select-selected-account');

    if (accountOptions) accountOptions.classList.remove('show');
    if (accountSelected) accountSelected.classList.remove('active');
    if (categoryOptions) categoryOptions.classList.remove('show');
    if (categorySelected) categorySelected.classList.remove('active');
    if (subcategoryOptions) subcategoryOptions.classList.remove('show');
    if (subcategorySelected) subcategorySelected.classList.remove('active');

    if (editAccountOptions) editAccountOptions.classList.remove('show');
    if (editAccountSelected) editAccountSelected.classList.remove('active');
    if (editCategoryOptions) editCategoryOptions.classList.remove('show');
    if (editCategorySelected) editCategorySelected.classList.remove('active');
    if (editSubcategoryOptions) editSubcategoryOptions.classList.remove('show');
    if (editSubcategorySelected) editSubcategorySelected.classList.remove('active');

    if (fromAccountOptions) fromAccountOptions.classList.remove('show');
    if (fromAccountSelected) fromAccountSelected.classList.remove('active');
    if (toAccountOptions) toAccountOptions.classList.remove('show');
    if (toAccountSelected) toAccountSelected.classList.remove('active');

    if (editFromAccountOptions) editFromAccountOptions.classList.remove('show');
    if (editFromAccountSelected) editFromAccountSelected.classList.remove('active');
    if (editToAccountOptions) editToAccountOptions.classList.remove('show');
    if (editToAccountSelected) editToAccountSelected.classList.remove('active');
}

document.addEventListener('click', function (event) {
    const accountSelectContainers = document.querySelectorAll('.custom-account-select');
    const categorySelectContainers = document.querySelectorAll('.custom-category-select');
    const subcategorySelectContainers = document.querySelectorAll('.custom-subcategory-select');

    let isClickInside = false;

    accountSelectContainers.forEach(container => {
        if (container.contains(event.target)) {
            isClickInside = true;
        }
    });

    categorySelectContainers.forEach(container => {
        if (container.contains(event.target)) {
            isClickInside = true;
        }
    });

    subcategorySelectContainers.forEach(container => {
        if (container.contains(event.target)) {
            isClickInside = true;
        }
    });

    if (!isClickInside) {
        closeAllDropdowns();
    }
});

function openEditTransactionModal(id, type, amount, accountId, accountName, accountColor, categoryId, categoryName, categoryColor, subCategoryId, description, date) {
    document.getElementById('editTransactionId').value = id;
    document.getElementById('editTransactionAmount').value = parseFloat(amount).toFixed(2);
    document.getElementById('editTransactionDescription').value = description || '';

    let transactionDate;
    if (date) {
        const originalDate = new Date(date);
        if (!isNaN(originalDate.getTime())) {
            transactionDate = originalDate;
        } else {
            transactionDate = new Date();
        }
    } else {
        transactionDate = new Date();
    }

    const year = transactionDate.getFullYear();
    const month = String(transactionDate.getMonth() + 1).padStart(2, '0');
    const day = String(transactionDate.getDate()).padStart(2, '0');
    document.getElementById('editTransactionDate').value = `${year}-${month}-${day}`;

    const typeButtons = document.querySelectorAll('#editTransactionModal .type-btn');
    typeButtons.forEach(btn => btn.classList.remove('active'));

    if (type === 0 || type === '0') {
        document.querySelector('#editTransactionModal .income-btn').classList.add('active');
        document.getElementById('editTransactionType').value = '0';
    } else {
        document.querySelector('#editTransactionModal .expense-btn').classList.add('active');
        document.getElementById('editTransactionType').value = '1';
    }

    document.getElementById('editSelectedAccountIcon').innerHTML = '';
    document.getElementById('editSelectedAccountIcon').style.backgroundColor = '';
    document.getElementById('editSelectedAccountName').textContent = 'Выберите счет';
    document.getElementById('editSelectedAccountBalance').textContent = '';

    document.getElementById('editSelectedCategoryIcon').innerHTML = '';
    document.getElementById('editSelectedCategoryIcon').style.backgroundColor = '';
    document.getElementById('editSelectedCategoryName').textContent = 'Выберите категорию';

    document.getElementById('editSelectedSubcategoryIcon').innerHTML = '';
    document.getElementById('editSelectedSubcategoryIcon').style.backgroundColor = '';
    document.getElementById('editSelectedSubcategoryName').textContent = 'Выберите подкатегорию';

    const accountOption = document.querySelector(`#editAccountOptions .select-account-option[data-account-id="${accountId}"]`);
    if (accountOption) selectEditAccountOption(accountOption);

    const categoryOption = document.querySelector(`#editCategoryOptions .select-category-option[data-category-id="${categoryId}"]`);
    if (categoryOption) selectEditCategoryOption(categoryOption);

    if (subCategoryId && subCategoryId !== 'null' && subCategoryId !== '') {
        const subCategoryOption = document.querySelector(`#editSubcategoryOptions .select-subcategory-option[data-subcategory-id="${subCategoryId}"]`);
        if (subCategoryOption) selectEditSubcategoryOption(subCategoryOption);
    }

    openModal('editTransactionModal');
}

function deleteTransaction() {
    const transactionId = document.getElementById('editTransactionId').value;
    if (confirm('Вы уверены, что хотите удалить эту транзакцию? Это действие необратимо.')) {
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/delete_transaction';

        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'ID';
        input.value = transactionId;

        form.appendChild(input);
        document.body.appendChild(form);
        form.submit();
    }
}

function toggleTransactionGroup(header) {
    const group = header.parentElement;
    const button = header.querySelector('.toggle-transactions-btn');
    const icon = header.querySelector('i');

    group.classList.toggle('collapsed');
    button.classList.toggle('collapsed');

    if (group.classList.contains('collapsed')) {
        icon.classList.remove('fa-chevron-down');
        icon.classList.add('fa-chevron-right');
    } else {
        icon.classList.remove('fa-chevron-right');
        icon.classList.add('fa-chevron-down');
    }
}

function handleTypeChange(type) {
    const transactionTypeInput = document.getElementById('transactionType');
    const typeButtons = document.querySelectorAll('#createTransactionModal .type-buttons .type-btn');
    const incomeExpenseFields = document.getElementById('incomeExpenseFields');
    const transferFields = document.getElementById('transferFields');
    const form = document.getElementById('transactionForm');
    const formActionInput = document.getElementById('formAction');

    const incomeExpenseInputs = incomeExpenseFields.querySelectorAll('input, button, select');
    const transferInputs = transferFields.querySelectorAll('input, button, select');

    transactionTypeInput.value = type;

    typeButtons.forEach(button => {
        button.classList.remove('active');
        if (button.getAttribute('data-type') === type) {
            button.classList.add('active');
        }
    });

    if (type === '0' || type === '1') {
        incomeExpenseFields.style.display = 'block';
        transferFields.style.display = 'none';

        incomeExpenseInputs.forEach(input => input.disabled = false);
        transferInputs.forEach(input => input.disabled = true);

        if (form) form.action = '/submit_transaction';
        if (formActionInput) formActionInput.value = '/submit_transaction';
    } else if (type === '2') {
        incomeExpenseFields.style.display = 'none';
        transferFields.style.display = 'block';

        incomeExpenseInputs.forEach(input => input.disabled = true);
        transferInputs.forEach(input => input.disabled = false);

        if (form) form.action = '/transfer';
        if (formActionInput) formActionInput.value = '/transfer';
    }
}

function handleEditTypeChange(type) {
    const transactionTypeInput = document.getElementById('editTransactionType');
    const typeButtons = document.querySelectorAll('#editTransactionModal .type-buttons .type-btn');

    transactionTypeInput.value = type;

    typeButtons.forEach(button => {
        button.classList.remove('active');
        if (button.getAttribute('data-type') === type) {
            button.classList.add('active');
        }
    });
}

document.addEventListener('DOMContentLoaded', () => {
    handleTypeChange(document.getElementById('transactionType').value || '0');
});

window.onclick = function (event) {
    if (!event.target.closest('.custom-account-select')) {
        const dropdowns = document.querySelectorAll('.select-account-options');
        dropdowns.forEach(function (dropdown) {
            if (dropdown.classList.contains('show')) {
                dropdown.classList.remove('show');
            }
        });
    }
}

document.addEventListener('DOMContentLoaded', function () {
    const createTypeButtons = document.querySelectorAll('#createTransactionModal .type-btn');
    createTypeButtons.forEach(button => {
        button.addEventListener('click', function () {
            const type = this.getAttribute('data-type');
            handleTypeChange(type);
        });
    });

    const editTypeButtons = document.querySelectorAll('#editTransactionModal .type-btn');
    editTypeButtons.forEach(button => {
        button.addEventListener('click', function () {
            const type = this.getAttribute('data-type');
            handleEditTypeChange(type);
        });
    });
});

function setupFormFields() {
    const form = document.getElementById('transactionForm');
    const type = document.getElementById('transactionType').value;

    const existingDynamicFields = form.querySelectorAll('[data-dynamic-field]');
    existingDynamicFields.forEach(field => field.remove());

    if (type === '0' || type === '1') {
        const accountId = document.getElementById('transactionAccount').value;
        const categoryId = document.getElementById('transactionCategory').value;
        const subCategoryId = document.getElementById('transactionSubCategory').value;

        if (accountId) {
            const accountField = document.createElement('input');
            accountField.type = 'hidden';
            accountField.name = 'AccountID';
            accountField.value = accountId;
            accountField.setAttribute('data-dynamic-field', 'true');
            form.appendChild(accountField);
        }

        if (categoryId) {
            const categoryField = document.createElement('input');
            categoryField.type = 'hidden';
            categoryField.name = 'CategoryID';
            categoryField.value = categoryId;
            categoryField.setAttribute('data-dynamic-field', 'true');
            form.appendChild(categoryField);
        }

        if (subCategoryId) {
            const subCategoryField = document.createElement('input');
            subCategoryField.type = 'hidden';
            subCategoryField.name = 'SubCategoryID';
            subCategoryField.value = subCategoryId;
            subCategoryField.setAttribute('data-dynamic-field', 'true');
            form.appendChild(subCategoryField);
        }

    } else if (type === '2') {
        const fromAccountId = document.getElementById('fromAccount').value;
        const toAccountId = document.getElementById('toAccount').value;

        if (fromAccountId) {
            const fromAccountField = document.createElement('input');
            fromAccountField.type = 'hidden';
            fromAccountField.name = 'AccountID';
            fromAccountField.value = fromAccountId;
            fromAccountField.setAttribute('data-dynamic-field', 'true');
            form.appendChild(fromAccountField);
        }

        if (toAccountId) {
            const toAccountField = document.createElement('input');
            toAccountField.type = 'hidden';
            toAccountField.name = 'TransferAccountID';
            toAccountField.value = toAccountId;
            toAccountField.setAttribute('data-dynamic-field', 'true');
            form.appendChild(toAccountField);
        }
    }
}

document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('transactionForm');
    if (form) {
        form.onsubmit = function (e) {
            setupFormFields();

            const type = document.getElementById('transactionType').value;
            const amount = document.getElementById('transactionAmount').value.trim();

            if (!amount || isNaN(parseFloat(amount))) {
                alert('Введите корректную сумму!');
                e.preventDefault();
                return false;
            }

            if (type === '0' || type === '1') {
                const account = document.getElementById('transactionAccount').value;
                const category = document.getElementById('transactionCategory').value;

                if (!account || account.trim() === '') {
                    alert('Выберите счет!');
                    e.preventDefault();
                    return false;
                }

                if (!category || category.trim() === '') {
                    alert('Выберите категорию!');
                    e.preventDefault();
                    return false;
                }

            } else if (type === '2') {
                const fromAccount = document.getElementById('fromAccount').value;
                const toAccount = document.getElementById('toAccount').value;

                if (!fromAccount || fromAccount.trim() === '') {
                    alert('Выберите счет списания!');
                    e.preventDefault();
                    return false;
                }

                if (!toAccount || toAccount.trim() === '') {
                    alert('Выберите счет зачисления!');
                    e.preventDefault();
                    return false;
                }

                if (fromAccount.trim() === toAccount.trim()) {
                    alert('Счета списания и зачисления не могут быть одинаковыми!');
                    e.preventDefault();
                    return false;
                }
            }

            if (amount === '') {
                document.getElementById('transactionAmount').value = '0.00';
            }

            return true;
        };
    }
});

function openEditTransferModal(id, amount, accountId, accountName, accountColor,
    transferAccountId, transferAccountName, transferAccountColor,
    description, date) {

    document.getElementById('editTransferId').value = id;
    document.getElementById('editTransferAmount').value = parseFloat(amount).toFixed(2);
    document.getElementById('editTransferDescription').value = description || '';

    let transferDate;
    if (date) {
        const originalDate = new Date(date);
        if (!isNaN(originalDate.getTime())) {
            transferDate = originalDate;
        } else {
            transferDate = new Date();
        }
    } else {
        transferDate = new Date();
    }

    const year = transferDate.getFullYear();
    const month = String(transferDate.getMonth() + 1).padStart(2, '0');
    const day = String(transferDate.getDate()).padStart(2, '0');
    document.getElementById('editTransferDate').value = `${year}-${month}-${day}`;

    document.getElementById('editFromSelectedAccountIcon').innerHTML = '';
    document.getElementById('editFromSelectedAccountIcon').style.backgroundColor = '';
    document.getElementById('editFromSelectedAccountName').textContent = 'Выберите счет';
    document.getElementById('editFromSelectedAccountBalance').textContent = '';

    document.getElementById('editToSelectedAccountIcon').innerHTML = '';
    document.getElementById('editToSelectedAccountIcon').style.backgroundColor = '';
    document.getElementById('editToSelectedAccountName').textContent = 'Выберите счет для перевода';
    document.getElementById('editToSelectedAccountBalance').textContent = '';

    const fromAccountOption = document.querySelector(`#editFromAccountOptions .select-account-option[data-account-id="${accountId}"]`);
    if (fromAccountOption) selectEditFromAccountOption(fromAccountOption);

    const toAccountOption = document.querySelector(`#editToAccountOptions .select-account-option[data-account-id="${transferAccountId}"]`);
    if (toAccountOption) selectEditToAccountOption(toAccountOption);

    openModal('editTransferModal');
}

function toggleEditFromAccountDropdown() {
    const options = document.getElementById('editFromAccountOptions');
    const selected = document.querySelector('#editTransferModal #editFromAccount + .select-selected-account');

    if (options && selected) {
        const isShowing = options.classList.contains('show');
        closeAllDropdowns();
        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function toggleEditToAccountDropdown() {
    const options = document.getElementById('editToAccountOptions');
    const selected = document.querySelector('#editTransferModal #editToAccount + .select-selected-account');

    if (options && selected) {
        const isShowing = options.classList.contains('show');
        closeAllDropdowns();
        if (!isShowing) {
            options.classList.add('show');
            selected.classList.add('active');
        }
    }
}

function selectEditFromAccountOption(optionElement) {
    const accountId = optionElement.getAttribute('data-account-id');
    const accountName = optionElement.getAttribute('data-account-name');
    const accountBalance = optionElement.getAttribute('data-account-balance');
    const accountColor = optionElement.getAttribute('data-account-color');
    const accountIconHTML = optionElement.querySelector('.option-account-icon').innerHTML;

    document.getElementById('editFromAccount').value = accountId;

    const selectedIcon = document.getElementById('editFromSelectedAccountIcon');
    const selectedName = document.getElementById('editFromSelectedAccountName');
    const selectedBalance = document.getElementById('editFromSelectedAccountBalance');

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

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('#editFromAccountOptions .select-account-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');
}

function selectEditToAccountOption(optionElement) {
    const accountId = optionElement.getAttribute('data-account-id');
    const accountName = optionElement.getAttribute('data-account-name');
    const accountBalance = optionElement.getAttribute('data-account-balance');
    const accountColor = optionElement.getAttribute('data-account-color');
    const accountIconHTML = optionElement.querySelector('.option-account-icon').innerHTML;

    document.getElementById('editToAccount').value = accountId;

    const selectedIcon = document.getElementById('editToSelectedAccountIcon');
    const selectedName = document.getElementById('editToSelectedAccountName');
    const selectedBalance = document.getElementById('editToSelectedAccountBalance');

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

    closeAllDropdowns();

    const allOptions = document.querySelectorAll('#editToAccountOptions .select-account-option');
    allOptions.forEach(option => {
        option.classList.remove('selected');
    });
    optionElement.classList.add('selected');
}

function deleteTransfer() {
    const transferId = document.getElementById('editTransferId').value;
    if (confirm('Вы уверены, что хотите удалить этот перевод? Это действие необратимо.')) {
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/delete_transfer';

        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'ID';
        input.value = transferId;

        form.appendChild(input);
        document.body.appendChild(form);
        form.submit();
    }
}

const transferCards = document.querySelectorAll('.transfer-card');
transferCards.forEach(card => {
    card.addEventListener('click', function (e) {
        e.stopPropagation();

        const transferId = this.getAttribute('data-transfer-id');
        const amount = this.getAttribute('data-amount');
        const accountId = this.getAttribute('data-account-id');
        const accountName = this.getAttribute('data-account-name');
        const accountColor = this.getAttribute('data-account-color');
        const transferAccountId = this.getAttribute('data-transfer-account-id');
        const transferAccountName = this.getAttribute('data-transfer-account-name');
        const transferAccountColor = this.getAttribute('data-transfer-account-color');
        const description = this.getAttribute('data-description') || '';
        const date = this.getAttribute('data-date');

        console.log('Opening edit modal for transfer:', transferId);

        openEditTransferModal(
            transferId,
            amount,
            accountId,
            accountName,
            accountColor,
            transferAccountId,
            transferAccountName,
            transferAccountColor,
            description,
            date
        );
    });

    card.style.cursor = 'pointer';
});

function deleteCategory(categoryId) {
    if (confirm('Внимание! При удалении категории будут также удалены все связанные подкатегории и транзакции. Вы уверены, что хотите продолжить?')) {
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/delete_category';

        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'ID';
        input.value = categoryId;

        form.appendChild(input);
        document.body.appendChild(form);
        form.submit();
    }
}

function deleteCategory(categoryId) {
    if (confirm('Внимание! При удалении счета будут также удалены все связанные с ним транзакции. Вы уверены, что хотите продолжить?')) {
        const form = document.createElement('form');
        form.method = 'POST';
        form.action = '/delete_account';

        const input = document.createElement('input');
        input.type = 'hidden';
        input.name = 'ID';
        input.value = categoryId;

        form.appendChild(input);
        document.body.appendChild(form);
        form.submit();
    }
}