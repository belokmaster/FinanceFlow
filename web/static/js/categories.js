const categoryModal = document.getElementById('categoryModal');
const categoryColorInput = document.getElementById('editCategoryColor');
const categoryColorPreview = document.getElementById('categoryColorPreview');
const categoryColorHexValue = document.getElementById('categoryColorHexValue');
const categoryIconSelect = document.querySelector('#categoryModal .custom-icon-select');
const categorySelectedIconDisplay = document.getElementById('selectedCategoryIconDisplay');
const categoryIconOptionsContainer = document.getElementById('categoryIconOptions');
const categoryHiddenIconInput = document.getElementById('editCategoryIcon');

const createCategoryModal = document.getElementById('createCategoryModal');
const createCategoryColorInput = document.getElementById('createCategoryColor');
const createCategoryColorPreview = document.getElementById('createCategoryColorPreview');
const createCategoryColorHexValue = document.getElementById('createCategoryColorHexValue');
const createCategoryIconSelect = document.querySelector('#createCategoryModal .custom-icon-select');
const createCategorySelectedIconDisplay = document.getElementById('createCategorySelectedIconDisplay');
const createCategoryIconOptionsContainer = document.getElementById('createCategoryIconOptionsContainer');
const createCategoryHiddenIconInput = document.getElementById('createCategoryIcon');

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

    document.getElementById('editCategoryColor').value = categoryColor;
    updateColorDisplay(categoryColor,
        document.getElementById('categoryColorPreview'),
        document.getElementById('categoryColorHexValue'));

    document.getElementById('editCategoryIcon').value = categoryIconKey;

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
