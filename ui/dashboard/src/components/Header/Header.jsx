import { HeaderContainer, Navigation, LoginButton, Link } from "./Header.styles"

const Header = () => (
  <HeaderContainer>
    <Navigation>
      <Link href="/">Home</Link>
      <Link href="/about">About</Link>
      <Link href="/shops">Shops</Link>
    </Navigation>
    <LoginButton>Login</LoginButton>
  </HeaderContainer>
);

export default Header;
