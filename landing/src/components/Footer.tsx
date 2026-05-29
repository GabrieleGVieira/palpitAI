export default function Footer() {
  return (
    <footer className="site-footer">
      <div className="section-shell footer-content">
        <div>
          <p>PalpitAI</p>
          <p>App de bolão com IA para futebol. Beta Android, iPhone em breve.</p>
        </div>
        <nav className="footer-links" aria-label="Links legais">
          <a href="/privacy">Política de Privacidade</a>
          <a href="/terms">Termos de Uso</a>
        </nav>
      </div>
    </footer>
  );
}
