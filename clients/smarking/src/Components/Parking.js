import React, { useState, useEffect } from "react";
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
  const parking = await response.json();
  return parking;
};

const getCurrentCars = async () => {
  const response = await fetch(api.base + api.handlers.cars, {
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
  const [parkings, setParkings] = useState();
  const [form, setForm] = useState({
    carID: null,
    endDate: null,
  });

  useEffect(() => {
    getCurrentParking().then((result) => setParkings(result));
  }, []);
  return (
    <div>
      {!parkings ? (
        <p>loading...</p>
      ) : parkings.length == 0 ? (
        <h1>Set up a parking!</h1>
      ) : (
        <div>
          <h1>Your current parkings:</h1>
          {parkings.map}
        </div>
      )}
    </div>
  );
}

export default Parking;
