// Restaurer la position de scroll sans animation
  const savedScrollY = localStorage.getItem("scrollY");
  if (savedScrollY !== null) {
    // Attendre que le layout soit prêt
    requestAnimationFrame(() => {
      window.scrollTo(0, parseInt(savedScrollY));
    });
  }


// Sauvegarder la position de scroll à chaque scroll
window.addEventListener("scroll", () => {
  localStorage.setItem("scrollY", window.scrollY);
});