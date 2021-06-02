import React from "react";
import Row from "react-bootstrap/Row";
import Col from "react-bootstrap/Col";
import Image from "react-bootstrap/Image";
import Button from "react-bootstrap/Button";
import { LinkContainer } from "react-router-bootstrap";
import { Link } from "react-router-dom";

import caf from "../img/car1.jpg";
import eat from "../img/parking.jpg";

function About() {
  return (
    <div>
      <h1
        style={{
          textAlign: "center",
          fontFamily: "Raleway",
          fontWeight: 900,
          marginTop: "20px",
        }}
      >
        Never lose your car again <br />
        with <span style={{ color: "#009999" }}>Smarking</span>
      </h1>
      <Row style={{ marginTop: "30px" }}>
        <Col md={7}>
          <Image
            rounded
            fluid
            src={caf}
            style={{
              paddingLeft: "5vw",
              paddingBottom: "10px",
            }}
          />
        </Col>
        <Col md={5}>
          <h2
            style={{
              textAlign: "center",
              fontFamily: "Raleway",
              fontWeight: 600,
              marginLeft: "3vw",
              marginRight: "3vw",
            }}
          >
            Store notes about where you parked your car
          </h2>
          <p
            style={{
              textAlign: "left",
              fontFamily: "Roboto",
              fontSize: "15px",
              color: "gray",
              marginTop: "30px",
              marginLeft: "7vw",
              marginRight: "7vw",
            }}
          >
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat. Duis aute irure dolor in
            reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
            culpa qui officia deserunt mollit anim id est laborum.
          </p>
          <Link to="/signup">
            <p
              style={{
                textAlign: "Center",
                fontFamily: "Roboto",
                fontSize: "15px",
                color: "#009999",
                fontWeight: 700,
              }}
            >
              Start parking
            </p>
          </Link>
        </Col>
      </Row>
      <Row
        className="flex-column-reverse flex-md-row"
        style={{ marginTop: "40px" }}
      >
        <Col md={5}>
          <h2
            style={{
              textAlign: "center",
              fontFamily: "Raleway",
              fontWeight: 600,
              marginLeft: "3vw",
              marginRight: "3vw",
            }}
          >
            Get notified when your payment will run out
          </h2>
          <p
            style={{
              textAlign: "left",
              fontFamily: "Roboto",
              fontSize: "15px",
              color: "gray",
              marginTop: "30px",
              marginLeft: "7vw",
              marginRight: "7vw",
            }}
          >
            Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
            eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim
            ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
            aliquip ex ea commodo consequat. Duis aute irure dolor in
            reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
            pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
            culpa qui officia deserunt mollit anim id est laborum.
          </p>
          <Link to="/signup">
            <p
              style={{
                textAlign: "Center",
                fontFamily: "Roboto",
                fontSize: "15px",
                color: "#009999",
                fontWeight: 700,
              }}
            >
              Create Profile
            </p>
          </Link>
        </Col>
        <Col md={7}>
          <Image
            rounded
            fluid
            src={eat}
            style={{
              paddingRight: "5vw",
              paddingBottom: "10px",
            }}
          />
        </Col>
      </Row>
      <p
        style={{
          textAlign: "Center",
          fontFamily: "Roboto",
          fontSize: "18px",
          color: "black",
          fontWeight: 700,
          marginTop: "40px",
        }}
      >
        Take charge of your new parking experiences with{" "}
        <span style={{ color: "#009999" }}>Smarking</span>.<br /> The
        revolutionary parking focused site.
      </p>
      <Row className="d-flex justify-content-center"></Row>
    </div>
  );
}

export default About;
