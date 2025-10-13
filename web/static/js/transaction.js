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
    document.getElementById('transactionAmount').value = '';
    document.getElementById('transactionDescription').value = '';
    document.getElementById('transactionDate').value = '';

    const typeButtons = document.querySelectorAll('.type-btn');
    typeButtons.forEach(btn => btn.classList.remove('active'));
    document.querySelector('.income-btn').classList.add('active');
    document.getElementById('transactionType').value = '0';

    const today = new Date().toISOString().split('T')[0];
    document.getElementById('transactionDate').value = today;

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