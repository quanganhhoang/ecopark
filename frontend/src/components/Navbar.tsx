import React from 'react';
import styled from 'styled-components';

const NavbarContainer = styled.nav`
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 20px;
  background-color: #f8f9fa;
  border-bottom: 1px solid #e9ecef;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

  @media (max-width: 768px) {
    flex-direction: column;
    padding: 15px;
    gap: 10px;
  }
`;

const LogoLink = styled.a`
  display: flex;
  align-items: center;
  text-decoration: none;
`;

const LogoImage = styled.img`
  height: 50px;
  margin-right: 10px;
`;

const LogoText = styled.span`
  font-size: 20px;
  font-weight: bold;
  color: #343a40;
`;

const NavLinks = styled.div`
  display: flex;
  gap: 20px;

  @media (max-width: 768px) {
    gap: 10px;
  }
`;

const NavLink = styled.a`
  text-decoration: none;
  font-size: 16px;
  color: #343a40;
  font-weight: 500;
  transition: color 0.3s;

  &:hover {
    color: #007bff;
  }
`;

const Navbar: React.FC = () => {
  return (
    <NavbarContainer>
      <LogoLink href="/">
        <LogoImage src="/logo.jpeg" alt="Ecopark Logo" />
        <LogoText>Ecopark</LogoText>
      </LogoLink>

      <NavLinks>
        <NavLink href="/reservation">Reservations</NavLink>
        <NavLink href="/activities">Activities</NavLink>
      </NavLinks>
    </NavbarContainer>
  );
};

export default Navbar;