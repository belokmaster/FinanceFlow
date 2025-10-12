const createCategoryModal = document.getElementById('createCategoryModal');
const createCategoryColorInput = document.getElementById('createCategoryColor');
const createCategoryColorPreview = document.getElementById('createCategoryColorPreview');
const createCategoryColorHexValue = document.getElementById('createCategoryColorHexValue');
const createCategoryIconSelect = document.querySelector('#createCategoryModal .custom-icon-select');
const createCategorySelectedIconDisplay = document.getElementById('createCategorySelectedIconDisplay');
const createCategoryIconOptionsContainer = document.getElementById('createCategoryIconOptionsContainer');
const createCategoryHiddenIconInput = document.getElementById('createCategoryIcon');

const editCategoryModal = document.getElementById('categoryModal');
const editCategoryColorInput = document.getElementById('editCategoryColor');
const editCategoryColorPreview = document.getElementById('categoryColorPreview');
const editCategoryColorHexValue = document.getElementById('categoryColorHexValue');
const editCategoryIconSelect = document.querySelector('#categoryModal .custom-icon-select');
const editCategorySelectedIconDisplay = document.getElementById('categorySelectedIconDisplay');
const editCategoryIconOptionsContainer = document.getElementById('categoryIconOptions');
const editCategoryHiddenIconInput = document.getElementById('editCategoryIcon');

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

function updateColorDisplay(newColor, previewEl, hexEl) {
    if (!newColor || !newColor.startsWith('#')) {
        newColor = '#000000';
    }
    previewEl.style.backgroundColor = newColor;
    hexEl.textContent = newColor.toUpperCase();
}

function openCreateCategoryModal() {
    document.getElementById('createCategoryName').value = '';
    document.getElementById('createCategoryColor').value = '#4cd67a';
    document.getElementById('createCategoryIcon').value = '';

    createCategorySelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = '';
    createCategorySelectedIconDisplay.querySelector('.selected-icon-key').textContent = 'Выберите иконку';

    updateColorDisplay('#4cd67a', createCategoryColorPreview, createCategoryColorHexValue);

    openModal('createCategoryModal');
}

function openCategoryModal(categoryId, categoryName, categoryColor, categoryIconKey) {
    document.getElementById('editCategoryId').value = categoryId;
    document.getElementById('editCategoryName').value = categoryName;

    editCategoryColorInput.value = categoryColor;
    updateColorDisplay(categoryColor, editCategoryColorPreview, editCategoryColorHexValue);
    editCategoryHiddenIconInput.value = categoryIconKey;

    const targetOption = editCategoryIconOptionsContainer.querySelector(`.select-icon-option[data-key="${categoryIconKey}"]`);
    if (targetOption) {
        const iconSvgHTML = targetOption.querySelector('.option-icon-svg').innerHTML;
        editCategorySelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = iconSvgHTML;
        editCategorySelectedIconDisplay.querySelector('.selected-icon-key').textContent = categoryIconKey;
    } else {
        editCategorySelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = '';
        editCategorySelectedIconDisplay.querySelector('.selected-icon-key').textContent = 'Выберите иконку';
    }

    openModal('categoryModal');
}

createCategoryColorInput.addEventListener('input', () => {
    updateColorDisplay(createCategoryColorInput.value, createCategoryColorPreview, createCategoryColorHexValue);
});

createCategorySelectedIconDisplay.addEventListener('click', (e) => {
    e.stopPropagation();
    createCategoryIconOptionsContainer.classList.toggle('show');
    createCategoryIconSelect.classList.toggle('active');
});

createCategoryIconOptionsContainer.addEventListener('click', (e) => {
    const option = e.target.closest('.select-icon-option');
    if (option) {
        const iconKey = option.dataset.key;
        const iconSvgHTML = option.querySelector('.option-icon-svg').innerHTML;
        const iconText = option.querySelector('.option-icon-key').textContent;

        createCategoryHiddenIconInput.value = iconKey;
        createCategorySelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = iconSvgHTML;
        createCategorySelectedIconDisplay.querySelector('.selected-icon-key').textContent = iconText;

        createCategoryIconOptionsContainer.classList.remove('show');
        createCategoryIconSelect.classList.remove('active');
    }
});

editCategoryColorInput.addEventListener('input', () => {
    updateColorDisplay(editCategoryColorInput.value, editCategoryColorPreview, editCategoryColorHexValue);
});

editCategorySelectedIconDisplay.addEventListener('click', (e) => {
    e.stopPropagation();
    editCategoryIconOptionsContainer.classList.toggle('show');
    editCategoryIconSelect.classList.toggle('active');
});

editCategoryIconOptionsContainer.addEventListener('click', (e) => {
    const option = e.target.closest('.select-icon-option');
    if (option) {
        const iconKey = option.dataset.key;
        const iconSvgHTML = option.querySelector('.option-icon-svg').innerHTML;
        const iconText = option.querySelector('.option-icon-key').textContent;

        editCategoryHiddenIconInput.value = iconKey;
        editCategorySelectedIconDisplay.querySelector('.selected-icon-svg').innerHTML = iconSvgHTML;
        editCategorySelectedIconDisplay.querySelector('.selected-icon-key').textContent = iconText;

        editCategoryIconOptionsContainer.classList.remove('show');
        editCategoryIconSelect.classList.remove('active');
    }
});

document.addEventListener('click', function (e) {
    if (e.target.classList.contains('modal')) {
        closeModal(e.target.id);
    }
    if (createCategoryIconSelect && !createCategoryIconSelect.contains(e.target)) {
        createCategoryIconOptionsContainer.classList.remove('show');
        createCategoryIconSelect.classList.remove('active');
    }
    if (editCategoryIconSelect && !editCategoryIconSelect.contains(e.target)) {
        editCategoryIconOptionsContainer.classList.remove('show');
        editCategoryIconSelect.classList.remove('active');
    }
});

document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape' || e.key === 'Esc') {
        closeModal('createCategoryModal');
        closeModal('categoryModal');
    }
});
