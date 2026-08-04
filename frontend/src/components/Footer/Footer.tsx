import { Container } from "../index.ts";

function Footer() {
  return (
    <footer className="border-t border-gray-800 bg-gray-900">
      <Container>
        <div className="py-4 text-center text-sm text-gray-400">
          Sravas — a video sharing platform
        </div>
      </Container>
    </footer>
  );
}

export default Footer;
