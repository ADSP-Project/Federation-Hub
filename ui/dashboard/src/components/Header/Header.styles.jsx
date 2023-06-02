import styled from 'styled-components';

export const HeaderContainer = styled.header`
  background-color: #808080;
  padding: 1em;
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const Navigation = styled.nav`
  & > a {
    margin-right: 1em;
  }
`;

export const LoginButton = styled.button`
  background-color: #007bff;
  color: white;
  border: none;
  padding: 0.5em 1em;
  border-radius: 0.25em;
  cursor: pointer;
`;

export const Link = styled.a`
  color: white;
`