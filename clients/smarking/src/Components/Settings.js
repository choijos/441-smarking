import React, { useState } from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import api from "../APIEndpoints.js";

const Settings = ({ user, setUser }) => {
  const [form, setForm] = useState({
    firstName: user ? user.firstName : null,
    lastName: user ? user.lastName : null,
    phoneNumber: user ? user.phoneNumber : null,
  });
  const [errors, setErrors] = useState({});

  if (!user) {
    return <p>Loading...</p>;
  }
  const setField = (field, value) => {
    setForm({
      ...form,
      [field]: value,
    });
  };

  const submitForm = async (e) => {
    e.preventDefault();

    let { firstName, lastName, phoneNumber } = form;

    if (phoneNumber.charAt(0) != "+") {
      this.setError("Phone number must be valid with country code (e.g. +12061234567)");
      return;

    }

    phoneNumber = phoneNumber.replace(/\D/g, "");
    if (phoneNumber.length < 10 || phoneNumber.length > 15) {
      this.setError("Phone number must be valid!");
      return;
    }

    const sendData = {
      firstName: firstName,
      lastName: lastName,
      phoneNumber: "+" + phoneNumber,
    };

    const response = await fetch(api.base + api.handlers.myprofile, {
      method: "PATCH",
      body: JSON.stringify(sendData),
      headers: new Headers({
        "Content-Type": "application/json",
        Authorization: localStorage.getItem("Authorization"),
      }),
    });
    if (response.status >= 300) {
      const error = await response.text();
      alert(error);
      return;
    }
    const newUser = await response.json();
    setUser(newUser);
  };

  return (
    <div>
      <h1>Settings</h1>
      <Form onSubmit={submitForm}>
        <Form.Group controlId="formGroupFirstName">
          <Form.Label>First Name</Form.Label>
          <Form.Control
            placeholder="Enter First Name"
            defaultValue={user.firstName}
            onChange={(e) => setField("firstName", e.target.value)}
          />
        </Form.Group>
        <Form.Group controlId="formGroupLastName">
          <Form.Label>Last Name</Form.Label>
          <Form.Control
            placeholder="Enter Last Name"
            defaultValue={user.lastName}
            onChange={(e) => setField("lastName", e.target.value)}
          />
        </Form.Group>
        <Form.Group controlId="formGroupPhoneNumber">
          <Form.Label>Phone Number</Form.Label>
          <Form.Control
            placeholder="Enter Phone Number"
            defaultValue={user.phoneNumber}
            onChange={(e) => setField("phoneNumber", e.target.value)}
          />
        </Form.Group>
        <br />
        <Button type="submit">Change Settings</Button>
      </Form>
    </div>
  );
};

export default Settings;
