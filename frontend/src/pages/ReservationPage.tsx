import React from "react";
import { useLocation } from "react-router-dom";
import { formatDate } from "../util";

const NIGHTLY_RATE = 500_000;

export const ReservationPage: React.FC = () => {
  const { state } = useLocation();
  const { startDate, endDate, guests } = state || {};
  if (!startDate || !endDate || guests === undefined) {
    return <div>Invalid reservation details. Please go back and try again.</div>;
  }

  const calculateTotalCost = (startDate: Date, endDate: Date) => {
    const days =
      (endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24);
    let cost = days > 0 ? days * NIGHTLY_RATE : 0;

    return cost.toLocaleString();
  };

  return (
    <div style={{ padding: "20px" }}>
      <h1>Payment</h1>
      <p>
        <strong>Start Date:</strong> {formatDate(startDate)}
      </p>
      <p>
        <strong>End Date:</strong> {formatDate(endDate)}
      </p>
      <p>
        <strong>Number of Guests:</strong> {guests}
      </p>
      <p>
        <strong>Total cost:</strong> {calculateTotalCost(startDate, endDate)} VND
      </p>

    </div>
  );
};