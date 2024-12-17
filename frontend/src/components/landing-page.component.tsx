import styled from "styled-components";

export const PageContainer = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 20px;
  gap: 20px;
  background-size: cover;
  background-position: center;
  color: white; /* Ensure text is readable */
`;

export const AddressColumn = styled.div`
  padding: 20px;
  border: 1px solid #ddd;
  background-color: rgba(255, 255, 255, 0.8); /* White with slight transparency for readability */
  border-radius: 8px;

  @media (max-width: 768px) {
    width: 100%; /* Take full width in column layout */
  }
`;

export const Address = styled.p`
  font-size: 16px;
  line-height: 1.5;
  color: green;
  border: 1px solid #ddd;
`;

export const Overlay = styled.div`
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10;
`;

export const GuestPicker = styled.div`
  background: white;
  padding: 20px;
  border-radius: 8px;
  color: black;
`;

export const Button = styled.button<{ disabled?: boolean }>`
  width: 40px;
  height: 40px;
  font-size: 20px;
  font-weight: bold;
  color: ${(props) => (props.disabled ? "#a9a9a9" : "#000")};
  background-color: ${(props) => (props.disabled ? "#f0f0f0" : "#fff")};
  border: 2px solid ${(props) => (props.disabled ? "#d3d3d3" : "#ccc")};
  border-radius: 50%;
  cursor: ${(props) => (props.disabled ? "not-allowed" : "pointer")};
  display: flex;
  justify-content: center;
  align-items: center;

  &:hover {
    background-color: ${(props) => (props.disabled ? "#f0f0f0" : "#e6e6e6")};
  }

  &:focus {
    outline: none;
  }
`;

export const ReservationContainer = styled.div`
  padding: 20px;
  display: flex;
  justify-content: center;
  align-items: center;
  border: 1px solid #ddd;
  background-color: rgba(255, 255, 255, 0.8); /* White with slight transparency for readability */
  border-radius: 8px;

  @media (max-width: 768px) {
    width: 100%; /* Take full width in column layout */
  }
`;

export const ReservationBar = styled.div`
  display: flex;
  background: white;
  border: 1px solid #ddd;
  border-radius: 12px;
  box-shadow: 0px 2px 4px rgba(0, 0, 0, 0.1);
  padding: 8px;
  align-items: center;
  width: 100%;
  max-width: 800px;
`;

export const ReservationField = styled.div`
  flex: 1;
  padding: 10px;
  cursor: pointer;
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  font-size: 16px;
`;

export const FieldLabel = styled.div`
  font-size: 14px;
  color: #757575;
  padding: 5px;
`;

interface FieldValueProps {
  children?: React.ReactNode;
  selected: boolean;
}

export const FieldValue: React.FC<FieldValueProps> = styled.div`
  font-size: 12px;
  color: ${(props) => (props.selected ? "black" : "#757575")};
  padding: 5px;
  font-weight: 500;
`;

export const Divider = styled.div`
  width: 1px;
  height: 40px;
  background: #ddd;
  margin: 0 10px;
`;

export const ReserveButton = styled.button`
  background-color: #ff5a5f;
  color: white;
  font-size: 14px;
  font-weight: 600;
  border: none;
  border-radius: 8px;
  padding: 10px 20px;
  cursor: pointer;
  &:hover {
    background-color: #e14850;
  }
`;