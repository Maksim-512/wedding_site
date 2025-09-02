document.addEventListener('DOMContentLoaded', function() {
    // Мобильное меню
    const mobileMenuBtn = document.querySelector('.mobile-menu-btn');
    const nav = document.querySelector('nav');

    if (mobileMenuBtn) {
        mobileMenuBtn.addEventListener('click', function() {
            nav.classList.toggle('active');
        });
    }

    // Плавная прокрутка
    document.querySelectorAll('nav a').forEach(anchor => {
        anchor.addEventListener('click', function(e) {
            e.preventDefault();
            const targetId = this.getAttribute('href');
            const targetElement = document.querySelector(targetId);

            if (targetElement) {
                if (nav.classList.contains('active')) nav.classList.remove('active');
                targetElement.scrollIntoView({ behavior: 'smooth' });
            }
        });
    });

    // Форма RSVP
    const rsvpForm = document.getElementById('rsvp-form');
    const companionField = document.querySelector('.companion');

    if (rsvpForm && companionField) {
        companionField.classList.remove('visible');

        // Показать/скрыть поле спутника
        const attendanceRadios = document.querySelectorAll('input[name="attendance"]');
        attendanceRadios.forEach(radio => {
            radio.addEventListener('change', function() {
                if (this.value === 'yes') {
                    companionField.classList.add('visible');
                } else {
                    companionField.classList.remove('visible');
                    document.getElementById('companion').value = '';
                }
            });
        });

        // Отправка формы
        rsvpForm.addEventListener('submit', async function(e) {
            e.preventDefault();

            const data = {
                name: document.getElementById('name').value,
                attendance: document.querySelector('input[name="attendance"]:checked')?.value,
                companion: document.getElementById('companion').value
            };

            try {
                const response = await fetch("/", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(data)
                });

                const result = await response.json();

                if (response.ok && result.status === "success") {
                    alert(result.message);
                    rsvpForm.reset();
                    companionField.classList.remove('visible');
                } else {
                    alert(result.message || "Ошибка при отправке формы");
                }
            } catch (error) {
                console.error("Ошибка:", error);
                alert("Произошла ошибка при отправке формы");
            }
        });
    }

    // Таймер
    function updateCountdown() {
        const weddingDate = new Date('2025-10-25T16:00:00');
        const now = new Date();
        const diff = weddingDate - now;

        const days = Math.floor(diff / (1000 * 60 * 60 * 24));
        const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
        const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
        const seconds = Math.floor((diff % (1000 * 60)) / 1000);

        document.getElementById('days').textContent = days;
        document.getElementById('hours').textContent = hours.toString().padStart(2, '0');
        document.getElementById('minutes').textContent = minutes.toString().padStart(2, '0');
        document.getElementById('seconds').textContent = seconds.toString().padStart(2, '0');
    }

    updateCountdown();
    setInterval(updateCountdown, 1000);
});
