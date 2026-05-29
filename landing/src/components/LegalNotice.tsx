export default function LegalNotice({ className = '' }: { className?: string }) {
  return (
    <p className={`legal-notice ${className}`.trim()}>
      Ao continuar, você concorda com os{' '}
      <a href="/terms">Termos de Uso</a> e a{' '}
      <a href="/privacy">Política de Privacidade</a>.
    </p>
  );
}
