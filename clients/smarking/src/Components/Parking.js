import React, { useState, useEffect } from "react";
import ListGroup from "react-bootstrap/ListGroup";
import Form from "react-bootstrap/Form";
import Button from "react-bootstrap/Button";
import api from "../APIEndpoints.js";
import Countdown from "react-countdown";

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
  const [cars, setCars] = useState();
  const [parkings, setParkings] = useState();
  const [form, setForm] = useState({
    carID: null,
    date: null,
    time: null,
    notes: null,
  });

  const setField = (field, value) => {
    setForm({
      ...form,
      [field]: value,
    });
  };

  useEffect(() => {
    getCurrentParking().then((result) => setParkings(result));
    getCurrentCars().then((result) => setCars(result));
  }, []);

  const submitForm = async (e) => {
    e.preventDefault();

    let send = { carID: +form.carID, notes: form.notes };
    send.endTime = new Date(form.date + " " + form.time);
    if (send.endTime == "Invalid Date") {
      alert("must provide a valid datetime!");
      return;
    }

    console.log(send);
    const response = await fetch(api.base + api.handlers.usersparking, {
      method: "POST",
      body: JSON.stringify(send),
      headers: new Headers({
        "Content-Type": "application/json",
        Authorization: localStorage.getItem("Authorization"),
      }),
    });
    if (response.status >= 300) {
      const error = await response.text();
      console.log(error);
      return;
    }
    const parking = await response.json();
    getCurrentParking().then((result) => setParkings(result));
  };

  const deleteParking = async (id) => {
    const response = await fetch(api.base + api.handlers.parking + id, {
      method: "PATCH",
      headers: new Headers({
        Authorization: localStorage.getItem("Authorization"),
      }),
    });
    if (response.status >= 300) {
      const error = await response.text();
      console.log(error);
      return;
    }
    const text = await response.text();
    getCurrentParking().then((result) => setParkings(result));
  };

  return (
    <div>
      {!parkings || !cars ? (
        <p>loading...</p>
      ) : cars.length == 0 ? (
        <h1>Register a car to continue!</h1>
      ) : (
        <div>
          {parkings.length != 0 && (
            <div>
              <h1>Your current parkings:</h1>
              <ListGroup>
                {parkings.map((p, i) => {
                  const car = cars.find((c) => c.id == p.carID);
                  console.log(cars);
                  if (!p.isCompleted && car) {
                    return (
                      <ListGroup.Item key={i}>
                        <h2>
                          {car.licensePlate +
                            ": " +
                            car.color +
                            " " +
                            car.year +
                            " " +
                            car.make +
                            " " +
                            car.model}
                        </h2>
                        <p>Start time: {p.startTime}</p>
                        <p>Notes: {p.notes}</p>
                        <Countdown date={p.endTime} />,
                        <Button
                          variant="success"
                          onClick={() => {
                            deleteParking(p._id);
                          }}
                          style={{ fontSize: "12px" }}
                        >
                          Mark as complete
                        </Button>
                      </ListGroup.Item>
                    );
                  }
                })}
              </ListGroup>
            </div>
          )}

          <div>
            <h1>Set up a parking!</h1>
            <Form onSubmit={submitForm}>
              <Form.Group controlId="formGroupCar">
                <Form.Label>Car</Form.Label>
                <Form.Control
                  as="select"
                  onChange={(e) => setField("carID", e.target.value)}
                >
                  {cars.map((c, i) => {
                    return (
                      <option key={i} value={c.id}>
                        {c.licensePlate +
                          ": " +
                          c.color +
                          " " +
                          c.year +
                          " " +
                          c.make +
                          " " +
                          c.model}
                      </option>
                    );
                  })}
                </Form.Control>
              </Form.Group>
              <Form.Group controlId="formGroupEndTime">
                <Form.Label>Expiry Time and Date</Form.Label>
                <Form.Control
                  type="time"
                  placeholder="Enter expiry time"
                  onChange={(e) => setField("time", e.target.value)}
                />
                <Form.Control
                  type="date"
                  placeholder="Enter expiry date"
                  onChange={(e) => setField("date", e.target.value)}
                />
              </Form.Group>
              <Form.Group controlId="formGroupNotes">
                <Form.Label>Notes</Form.Label>
                <Form.Control
                  placeholder="Enter Notes"
                  onChange={(e) => setField("notes", e.target.value)}
                />
              </Form.Group>
              <br />
              <Button type="submit">Add Parking</Button>
            </Form>
          </div>
        </div>
      )}
    </div>
  );
}

export default Parking;
