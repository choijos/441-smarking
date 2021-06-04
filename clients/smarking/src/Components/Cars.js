import React, { useState, useEffect } from "react";
import api from "../APIEndpoints.js";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import { Alert } from "react-bootstrap";

const getCurrentCars = async () => {
  const response = await fetch(api.base + api.handlers.cars, {
    method: "GET",
    headers: new Headers({
      Authorization: localStorage.getItem("Authorization"),
    }),
  });
  if (response.status >= 300) {
    const error = await response.text();
    alert(error);
    return;
  }
  const cars = await response.json();
  return cars;
};

const Cars = () => {
  const [cars, setCars] = useState();
  const [form, setForm] = useState({
    licensePlate: null,
    make: null,
    model: null,
    year: null,
    color: null,
  });

  useEffect(() => {
    getCurrentCars().then((result) => setCars(result));
  }, []);

  const setField = (field, value) => {
    setForm({
      ...form,
      [field]: value,
    });
  };

  const deleteCar = async (id) => {
    const response = await fetch(api.base + api.handlers.carspec + id, {
      method: "DELETE",
      headers: new Headers({
        Authorization: localStorage.getItem("Authorization"),
      }),
    });
    if (response.status >= 300) {
      const error = await response.text();
      alert(error);
      return;
    }
    const text = await response.text();
    getCurrentCars().then((result) => setCars(result));
  };

  const submitForm = async (e) => {
    e.preventDefault();

    const response = await fetch(api.base + api.handlers.cars, {
      method: "POST",
      body: JSON.stringify(form),
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
    const car = await response.json();
    getCurrentCars().then((result) => setCars(result));
  };

  return !cars ? (
    <p>Loading...</p>
  ) : (
    <div>
      <h1>Car List</h1>
      {cars.length == 0 ? (
        <p>Your car list is empty! Add some below!</p>
      ) : (
        <div>
          {cars.map((c, i) => {
            return (
              <p key={i}>
                {c.licensePlate +
                  ": " +
                  c.color +
                  " " +
                  c.year +
                  " " +
                  c.make +
                  " " +
                  c.model}{" "}
                <Button
                  variant="danger"
                  onClick={() => deleteCar(c.id)}
                  style={{ fontSize: "12px" }}
                >
                  Delete
                </Button>
              </p>
            );
          })}
        </div>
      )}
      <Form onSubmit={submitForm}>
        <Form.Group controlId="formGroupLicensePlate">
          <Form.Label>License Plate</Form.Label>
          <Form.Control
            placeholder="Enter License Plate"
            onChange={(e) => setField("licensePlate", e.target.value)}
          />
        </Form.Group>
        <Form.Group controlId="formGroupColor">
          <Form.Label>Color</Form.Label>
          <Form.Control
            placeholder="Enter Car Color"
            onChange={(e) => setField("color", e.target.value)}
          />
        </Form.Group>
        <Form.Group controlId="formGroupYear">
          <Form.Label>Year</Form.Label>
          <Form.Control
            placeholder="Enter Manufacturing Year"
            onChange={(e) => setField("year", e.target.value)}
          />
        </Form.Group>
        <Form.Group controlId="formGroupMake">
          <Form.Label>Make</Form.Label>
          <Form.Control
            placeholder="Enter Car Make"
            onChange={(e) => setField("make", e.target.value)}
          />
        </Form.Group>
        <Form.Group controlId="formGroupModel">
          <Form.Label>Model</Form.Label>
          <Form.Control
            placeholder="Enter Car Model"
            onChange={(e) => setField("model", e.target.value)}
          />
        </Form.Group>
        <br />
        <Button type="submit">Add Car</Button>
      </Form>
    </div>
  );
};

export default Cars;
