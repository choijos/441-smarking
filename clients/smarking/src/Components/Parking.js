import React from "react";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import api from "../APIEndpoints.js";

const getCurrentParking = async () => {
  const response = await fetch(api.base + api.handlers.usersparking, {
    method: "GET",
    headers: new Headers({
      Authorization: localStorage.getItem("Authorization"),
    }),
  });
  if (response.status >= 300) {
    const error = await response.text();
    console.log(error);
    return;
  }
  const cars = await response.json();
  return cars;
};

function Parking() {
  return <div>Let's start a new parking!</div>;
}

export default Parking;
