import React, { useState } from "react";
import PropTypes from "prop-types";
import NavDropdown from "react-bootstrap/NavDropdown";
import api from "../APIEndpoints";
import Errors from "../Errors/Errors.js";
import { Redirect } from "react-router-dom";

const SignOutButton = ({ setUser }) => {
  const [error, setError] = useState("");
  const [redir, setRedir] = false;

  return (
    <div className="sign-out">
      {redir && <Redirect to="/" />}
      <NavDropdown.Item
        onClick={async (e) => {
          e.preventDefault();
          const response = await fetch(api.base + api.handlers.mysession, {
            method: "DELETE",
            headers: new Headers({
              Authorization: localStorage.getItem("Authorization"),
              "Content-Type": "application/json",
            }),
          });
          if (response.status >= 300) {
            const error = await response.text();
            setError(error);
            return;
          }
          localStorage.removeItem("Authorization");
          setError("");
          setUser(null);
          setRedir(true);
        }}
      >
        Sign out
      </NavDropdown.Item>
      {error && (
        <div>
          <Errors error={error} setError={setError} />
        </div>
      )}
    </div>
  );
};

SignOutButton.propTypes = {
  setUser: PropTypes.func.isRequired,
};

export default SignOutButton;
