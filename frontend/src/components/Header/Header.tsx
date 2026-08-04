import { Link, useNavigate } from "react-router-dom";
import { Container, Logo, LogoutBtn } from "../index.ts";
import { useAppSelector } from "../../store/hooks.ts";

function Header() {
  const authStatus = useAppSelector((state) => state.auth.status);
  const navigate = useNavigate();

  const navItems = [
    { name: "Home", slug: "/", active: authStatus },
    { name: "Login", slug: "/login", active: !authStatus },
    { name: "Signup", slug: "/register", active: !authStatus },
  ];

  return (
    <header className="border-b border-gray-800 bg-gray-900">
      <Container>
        <nav className="flex py-3">
          <div className="mr-4 flex items-center">
            <Link to="/">
              <Logo width="100px" />
            </Link>
          </div>
          <ul className="flex ml-auto items-center gap-2">
            {navItems.map(
              (item) =>
                item.active && (
                  <li key={item.name}>
                    <button
                      onClick={() => navigate(item.slug)}
                      className="inline-block px-4 py-2 rounded-full duration-200 hover:bg-gray-800 text-gray-200"
                    >
                      {item.name}
                    </button>
                  </li>
                ),
            )}
            {authStatus && (
              <li>
                <LogoutBtn />
              </li>
            )}
          </ul>
        </nav>
      </Container>
    </header>
  );
}

export default Header;
