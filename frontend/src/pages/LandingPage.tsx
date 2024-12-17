import React, { useEffect, useState } from 'react';
import { useNavigate } from "react-router-dom";
// import axios from 'axios';
import DatePicker from "react-datepicker";
import "react-datepicker/dist/react-datepicker.css";

import { GridItem } from '../components/grid-item.component';
import { GridContainer } from '../components/grid-container.component';
import { Image } from '../components/image.component';
import {
  Button,
  ReserveButton,
  GuestPicker,
  Overlay,

  ReservationContainer,
  ReservationBar,
  ReservationField,
  FieldLabel,
  FieldValue,
  Divider,

  PageContainer,
  AddressColumn,
  Address
} from '../components/landing-page.component';
import { formatDate } from '../util';

interface Picture {
  id: string;
  url: string;
  alt: string;
}

// const UNSPLASH_KEY = import.meta.env.VITE_UNSPLASH_ACCESS_KEY;
// const UNSPLASH_API = 'https://api.unsplash.com/photos/random';
const MIN_NUM_GUESTS = 1;
const MAX_NUM_GUESTS = 4;

// Fetch random pictures from Unsplash
async function getRandomPictures(): Promise<Picture[]> {
  // const params = {
  //   count: 9, // Number of random photos to fetch
  //   client_id: UNSPLASH_KEY,
  // };

  // try {
  //   const response = await axios.get(UNSPLASH_API, { params });
  //   return response.data.map((pic: any) => ({
  //     id: pic.id,
  //     url: pic.urls.small,
  //     alt: pic.alt_description || 'Random image from Unsplash',
  //   }));
  // } catch (error) {
  //   console.error('Error fetching pictures from Unsplash:', error);
  //   return [];
  // }

  return [];
}

const LandingPage: React.FC = () => {
  const [pictures, setPictures] = useState<Picture[]>([]);
  // const [loading, setLoading] = useState(true);

  const [startDate, setStartDate] = useState<Date | null>(null);
  const [endDate, setEndDate] = useState<Date | null>(null);

  const [showDatePicker, setShowDatePicker] = useState<"checkin" | "checkout" | null>(null);
  const [guests, setGuests] = useState(1);
  const [showGuestPicker, setShowGuestPicker] = useState(false);

  const navigate = useNavigate();

  const unavailableDates = [
    new Date(2024, 10, 28), // November 28, 2024
    new Date(2024, 10, 29), // November 29, 2024
    new Date(2024, 10, 30), // November 30, 2024
  ];

  const handleReserve = () => {
    if (!startDate || !endDate) {
      alert("Please select both start and end dates.");
      return;
    }

    if (endDate <= startDate) {
      alert("End date must be after the start date.");
      return;
    }

    navigate("/reservation", {
      state: {
        startDate,
        endDate,
        guests,
      },
    });
  };

  const incrementGuests = () => {
    if (guests < MAX_NUM_GUESTS) {
      setGuests(guests + 1);
    }
  };

  const decrementGuests = () => {
    if (guests > MIN_NUM_GUESTS) {
      setGuests(guests - 1);
    }
  };

  const nightlyRate = 500_000;

  // Calculate total cost
  const calculateTotalCost = () => {
    let cost = 0;
    if (startDate && endDate) {
      const days =
        (endDate.getTime() - startDate.getTime()) / (1000 * 60 * 60 * 24);
      cost = days > 0 ? days * nightlyRate : 0;
    }
    return cost.toLocaleString();
  };

  useEffect(() => {
    async function fetchPictures() {
      try {
      const fetchedPictures = await getRandomPictures();
      setPictures(fetchedPictures);
      } catch (error) {
        console.error('Error in fetchPictures:', error);
      } finally {
        // setLoading(false);
      }
    }

    fetchPictures();
  }, []);

  // if (loading) {
  //   return <div>Loading...</div>;
  // }

  // if (pictures.length === 0) {
  //   return <div>No pictures available. Check API or network settings.</div>;
  // }

  return (
    <div style={{ textAlign: 'center', padding: '20px' }}>
      <h1>Welcome to our Ecopark Homestay</h1>
      <PageContainer>
        <ReservationContainer>
          <ReservationBar>
            <ReservationField onClick={() => setShowDatePicker("checkin")}>
              <FieldLabel>Check-in</FieldLabel>
              <FieldValue selected={!!startDate}>
                {startDate ?
                  formatDate(startDate)
                  : "Add dates"
                }
              </FieldValue>
            </ReservationField>
            <Divider />
            <ReservationField onClick={() => setShowDatePicker("checkout")}>
              <FieldLabel>Check-out</FieldLabel>
              <FieldValue selected={!!endDate}>
                {endDate ?
                  formatDate(endDate)
                  : "Add dates"}
              </FieldValue>
            </ReservationField>
            <Divider />
            <ReservationField onClick={() => setShowGuestPicker(true)}>
              <FieldLabel>Who</FieldLabel>
              <FieldValue selected={!!guests}>
                {guests ? `${guests} guests` : "Add guests"}
              </FieldValue>
            </ReservationField>
            <ReserveButton onClick={handleReserve}>Reserve</ReserveButton>
          </ReservationBar>
        </ReservationContainer>
        {showDatePicker && (
            <Overlay onClick={() => setShowDatePicker(null)}>
              <div onClick={(e) => e.stopPropagation()}>
                <DatePicker
                  selected={showDatePicker === "checkin" ? startDate : endDate}
                  onChange={(date) => {
                    if (showDatePicker === "checkin") {
                      setStartDate(date);
                    } else {
                      setEndDate(date);
                    }
                    setShowDatePicker(null);
                  }}
                  inline
                  excludeDates={unavailableDates}
                />
              </div>
            </Overlay>
          )}

          {showGuestPicker && (
            <Overlay onClick={() => setShowGuestPicker(false)}>
              <GuestPicker onClick={(e) => e.stopPropagation()}>
                <h3>Select Guests</h3>
                <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
                  <Button onClick={decrementGuests} disabled={guests <= MIN_NUM_GUESTS}>-</Button>
                    <span>{guests}</span>
                  <Button onClick={incrementGuests} disabled={guests >= MAX_NUM_GUESTS}>+</Button>
                </div>
              </GuestPicker>
            </Overlay>
          )}
        <AddressColumn>
          <h2 style={{color: 'green'}}>Our Location</h2>
          <Address>
            Toà Landmark 2,<br />
            Khu đô thị Ecopark,<br />
            Văn Giang, Hưng Yên<br />
            Hà Nội
          </Address>
        </AddressColumn>
      </PageContainer>
      <GridContainer>
        {pictures.map(pic => (
          <GridItem key={pic.id}>
            <Image src={pic.url} alt={pic.alt} />
          </GridItem>
        ))}
      </GridContainer>
    </div>
  );
}

export default LandingPage;