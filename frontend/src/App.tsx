import { useEffect, useState } from "react";
import "./App.css";

import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient, type Transport } from "@connectrpc/connect";
import { AppointmentService } from "./gen/appointment_pb";
import type { Appointment } from "./gen/appointment_pb";
import { convertToProtobufTime } from "./utils/date.util";
import { timestampDate } from "@bufbuild/protobuf/wkt";

const transport: Transport = createConnectTransport({
  baseUrl: "http://localhost:8080",
});

const appointmentClient = createClient(AppointmentService, transport);

const App = () => {
  const [appointments, setAppointments] = useState<Appointment[]>([]);
  const [selectedAppointmentId, setSelectedAppointmentId] =
    useState<string>("");

  const [formData, setFormData] = useState({
    title: "",
    description: "",
    contactInformation: {
      email: "",
      name: "",
    },
    date: "",
    startTime: "",
    endTime: "",
  });

  const [loading, setLoading] = useState(false);

  const selectedAppointment = appointments.find(
    (a) => a.id === selectedAppointmentId,
  );

  const handleCreateAppointment = async (e: React.FormEvent) => {
    e.preventDefault();

    try {
      const { start, end, date } = convertToProtobufTime(
        formData.date,
        formData.startTime,
        formData.endTime,
      );

      const response = await appointmentClient.createAppointment({
        title: formData.title,
        description: formData.description,
        contactInformation: {
          name: formData.contactInformation.name,
          email: formData.contactInformation.email,
        },
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
    } catch (err) {
      console.error("Failed to book appointment:", err);
    }
  };

  const handleFetchAppointments = async () => {
    setLoading(true);

    const userId = localStorage.getItem("userId");

    if (!userId) return;
    try {
      const response = await appointmentClient.getUserAppointments({
        userId,
      });

      setAppointments(response.appointments);
    } catch (err) {
      console.error("Failed to book appointment:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    handleFetchAppointments();
  }, []);

  return (
    <div className="dashboard-container">
      <div className="grid-layout">
        <section className="card">
          <h3>CREATE NEW APPOINTMENT</h3>

          <form onSubmit={handleCreateAppointment}>
            <div className="input-group">
              <label>Title</label>

              <input
                type="text"
                value={formData.title}
                onChange={(e) =>
                  setFormData({ ...formData, title: e.target.value })
                }
              />
            </div>

            <div className="input-group">
              <label>Description</label>
              <textarea
                value={formData.description}
                onChange={(e) =>
                  setFormData({ ...formData, description: e.target.value })
                }
              />
            </div>

            <div className="input-group">
              <label>Date</label>
              <input
                type="text"
                placeholder="December 17, 1995"
                value={formData.date}
                onChange={(e) =>
                  setFormData({ ...formData, date: e.target.value })
                }
              />
            </div>

            <div className="input-group">
              <label>Start Time</label>
              <input
                type="text"
                placeholder="12:10"
                value={formData.startTime}
                onChange={(e) =>
                  setFormData({ ...formData, startTime: e.target.value })
                }
              />
            </div>

            <div className="input-group">
              <label>End Time</label>
              <input
                type="text"
                placeholder="13:10"
                value={formData.endTime}
                onChange={(e) =>
                  setFormData({ ...formData, endTime: e.target.value })
                }
              />
            </div>

            <div className="input-group">
              <label>Name</label>
              <input
                type="text"
                value={formData.contactInformation.name}
                onChange={(e) =>
                  setFormData({
                    ...formData,

                    contactInformation: {
                      ...formData.contactInformation,

                      name: e.target.value,
                    },
                  })
                }
              />
            </div>

            <div className="input-group">
              <label>Email</label>
              <input
                type="text"
                value={formData.contactInformation.email}
                onChange={(e) =>
                  setFormData({
                    ...formData,

                    contactInformation: {
                      ...formData.contactInformation,

                      email: e.target.value,
                    },
                  })
                }
              />
            </div>

            <button type="submit" className="btn-primary">
              BOOK APPOINTMENT
            </button>
          </form>
        </section>

        <section className="card">
          <h3>VIEW APPOINTMENTS</h3>
          <div className="schedule-list">
            {!loading
              ? appointments.map((item) => (
                  <div
                    key={item.id}
                    className={`schedule-item ${selectedAppointmentId === item.id ? "selected" : ""}`}
                    onClick={() => setSelectedAppointmentId(item.id)}
                  >
                    <div className="schedule-info">
                      <span className="title-text">{item.title}</span>
                      <span className="time-text">
                        {item.startTime &&
                          timestampDate(item.startTime).toLocaleTimeString([], {
                            hour: "2-digit",
                            minute: "2-digit",
                          })}
                      </span>
                    </div>
                    <div className="user-tag">
                      {item.contactInformation?.name}
                    </div>
                  </div>
                ))
              : "↻"}
          </div>
        </section>

        <section className="card">
          <h3>APPOINTMENT DETAILS</h3>
          {selectedAppointment ? (
            <>
              <span className="title-text">{selectedAppointment.title}</span>
              <p>{selectedAppointment.description}</p>
              <span className="time-text">
                {selectedAppointment.startTime &&
                  timestampDate(selectedAppointment.startTime).toLocaleString()}
              </span>
              <div className="action-buttons">
                <button className="btn-action cancel">CANCEL</button>
              </div>
            </>
          ) : (
            <span>Select an appointment to view details</span>
          )}
        </section>
      </div>
    </div>
  );
};

export default App;
