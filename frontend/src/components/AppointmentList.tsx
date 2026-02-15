import { timestampDate } from "@bufbuild/protobuf/wkt";
import type { AppointmentListProps } from "../types/appointment.type";

export const AppointmentList = ({
  items,
  selectedId,
  onSelect,
}: AppointmentListProps) => {
  if (items.length === 0)
    return <div className="empty-state">No appointments found</div>;

  return (
    <div className="schedule-list">
      {items.map((item) => (
        <div
          key={item.id}
          className={`schedule-item ${selectedId === item.id ? "selected" : ""}`}
          onClick={() => onSelect(item.id)}
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
          <div className="user-tag">{item.contactInformation?.name}</div>
        </div>
      ))}
    </div>
  );
};
