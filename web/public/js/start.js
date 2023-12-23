const initialScreen = document.querySelector('.initial-screen');
const btnStart = document.querySelector('.initial-screen__btn');

btnStart.addEventListener('click', () => {
    initialScreen.classList.add('initial-screen--close');
});