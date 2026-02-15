import { timestampDate } from "@bufbuild/protobuf/wkt";
import type { AppointmentDetailsProps } from "../types/appointment.type";

export const AppointmentDetails = ({
  appointment,
  onDelete,
}: AppointmentDetailsProps) => {
  if (!appointment) {
    return (
      <section className="card empty">
        <h3>APPOINTMENT DETAILS</h3>
        <span className="muted-text">
          Select an appointment to view details
        </span>
      </section>
    );
  }

  return (
    <section className="card">
      <h3>APPOINTMENT DETAILS</h3>
      <div className="details-content">
        <h4 className="detail-title">{appointment.title}</h4>
        <p className="detail-desc">
          {appointment.description || "No description provided."}
        </p>

        <div className="detail-meta">
          <strong>Time:</strong>{" "}
          {appointment.startTime &&
            timestampDate(appointment.startTime).toLocaleString()}
        </div>

        <div className="detail-meta">
          <strong>Contact:</strong> {appointment.contactInformation?.name} (
          {appointment.contactInformation?.email})
        </div>

        <div className="action-buttons">
          <button className="btn-action cancel" onClick={onDelete}>
            CANCEL APPOINTMENT
          </button>
        </div>
      </div>
    </section>
  );
};
