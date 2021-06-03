import React from "react";
import Navbar from "react-bootstrap/Navbar";
import Nav from "react-bootstrap/Nav";
import Button from "react-bootstrap/Button";
import NavDropdown from "react-bootstrap/NavDropdown";
import SignOutButton from "./SignOutButton.js";
import { LinkContainer } from "react-router-bootstrap";
import Image from "react-bootstrap/Image";
import st from "../img/smarkingtext.png";

//Creates header
const Header = ({ user, setUser }) => {
  return (
    <div>
      <Navbar className="nav" fixed="top" expand="lg" variant="dark">
        <LinkContainer to="/">
          <Navbar.Brand>
            <Image style={{ paddingLeft: "5px" }} src={st} height={25} />
          </Navbar.Brand>
        </LinkContainer>
        <Navbar.Toggle aria-controls="basic-navbar-nav" />
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="mr-auto">
            <LinkContainer to="/about">
              <Nav.Link>About</Nav.Link>
            </LinkContainer>
          </Nav>
          <Nav>
            {!user && (
              <Button
                className="searchbtn"
                variant="outline-light"
                href="/signin"
              >
                Sign in
              </Button>
            )}
            {user && (
              <React.Fragment>
                <Nav>
                  <LinkContainer to="/parking">
                    <Nav.Link>Current Parking</Nav.Link>
                  </LinkContainer>
                </Nav>
                <NavDropdown title="Profile" id="basic-nav-dropdown">
                  <LinkContainer to="/cars">
                    <NavDropdown.Item>My Cars</NavDropdown.Item>
                  </LinkContainer>
                  <LinkContainer to="/settings">
                    <NavDropdown.Item>Settings</NavDropdown.Item>
                  </LinkContainer>
                  <NavDropdown.Divider />
                  <SignOutButton setUser={setUser} />
                </NavDropdown>
              </React.Fragment>
            )}
          </Nav>
        </Navbar.Collapse>
      </Navbar>
    </div>
  );
};

export default Header;
