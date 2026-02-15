import React, { useEffect, useState } from "react";
import "./App.css";

import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient, type Transport } from "@connectrpc/connect";
import { AppointmentService } from "./gen/appointment_pb";
import type { Appointment } from "./gen/appointment_pb";
import { convertToProtobufTime } from "./utils/date.util";
import { AppointmentForm } from "./components/AppointmentForm";
import { AppointmentDetails } from "./components/AppointmentDetails";
import { AppointmentList } from "./components/AppointmentList";

const transport: Transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

const appointmentClient = createClient(AppointmentService, transport);

const App = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [selectedAppointmentId, setSelectedAppointmentId] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [formData, setFormData] = useState({
    title: "",
    description: "",
    date: "",
    startTime: "",
    endTime: "",
    contactInformation: { email: "", name: "" },
  });

  const isFormValid = () => {
    const { title, date, startTime, endTime, contactInformation, description } =
      formData;
    if (
      !title ||
      !date ||
      !startTime ||
      !endTime ||
      !contactInformation.email ||
      !contactInformation.name ||
      !description
    ) {
      setError("Please fill in all required fields.");
      return false;
    }
    if (startTime >= endTime) {
      setError("End time must be after start time.");
      return false;
    }
    return true;
  };

  const handleFetchAppointments = async () => {
    const userId = localStorage.getItem("userId");
    if (!userId) {
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);
    try {
      const response = await appointmentClient.getUserAppointments({ userId });
      setAppointments(response.appointments);
    } catch (err: any) {
      setError(err.message || "Failed to load appointments");
    } finally {
      setLoading(false);
    }
  };

  const handleCreateAppointment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isFormValid()) return;

    setError(null);
    try {
      const { start, end, date } = convertToProtobufTime(
        formData.date,
        formData.startTime,
        formData.endTime,
      );

      const response = await appointmentClient.createAppointment({
        ...formData,
        startTime: start,
        endTime: end,
        date,
        userId: localStorage.getItem("userId") ?? "",
      });

      localStorage.setItem("userId", response.userId);

      setAppointments((prev) => [...prev, response]);

      setFormData({
        title: "",
        description: "",
        contactInformation: { email: "", name: "" },
        date: "",
        startTime: "",
        endTime: "",
      });
    } catch (err: any) {
      setError(err.message || "Failed to create appointment");
    }
  };

  const handleDeleteAppointment = async () => {
    try {
      await appointmentClient.deleteAppointment({
        id: selectedAppointmentId,
      });

      await handleFetchAppointments();
    } catch (err: any) {
      setError(err.message || "Failed to delete appointment");
    }
  };

  useEffect(() => {
    handleFetchAppointments();
  }, []);

  return (
    <div className="dashboard-container">
      <div className="grid-layout">
        <AppointmentForm
          formData={formData}
          setFormData={setFormData}
          onSubmit={handleCreateAppointment}
          error={error}
        />

        <section className="card">
          <h3>YOUR SCHEDULE</h3>
          {loading ? (
            <div className="loader">Loading...</div>
          ) : (
            <AppointmentList
              items={appointments}
              selectedId={selectedAppointmentId}
              onSelect={setSelectedAppointmentId}
            />
          )}
        </section>

        <AppointmentDetails
          appointment={appointments.find((a) => a.id === selectedAppointmentId)}
          onDelete={handleDeleteAppointment}
        />
      </div>
    </div>
  );
};

export default App;
