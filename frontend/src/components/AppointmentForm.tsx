import type { AppointmentFormProps } from "../types/appointment.type";

export const AppointmentForm = ({
  onSubmit,
  formData,
  setFormData,
  error,
}: AppointmentFormProps) => (
  <section className="card">
    <h3>CREATE NEW APPOINTMENT</h3>
    {error && <div className="error-banner">{error}</div>}
    <form onSubmit={onSubmit}>
      <div className="input-group">
        <label>Title</label>
        <input
          required
          type="text"
          value={formData.title}
          onChange={(e) => setFormData({ ...formData, title: e.target.value })}
        />
      </div>
      <div className="input-group">
        <label>Description</label>
        <textarea
          required
          value={formData.description}
          onChange={(e) =>
            setFormData({ ...formData, description: e.target.value })
          }
        />
      </div>
      <div className="input-group">
        <label>Name</label>
        <input
          required
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
          required
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
      <div className="input-group">
        <label>Date</label>
        <input
          required
          type="date"
          value={formData.date}
          onChange={(e) => setFormData({ ...formData, date: e.target.value })}
        />
      </div>
      <div className="grid-2">
        <div className="input-group">
          <label>Start</label>
          <input
            required
            type="time"
            value={formData.startTime}
            onChange={(e) =>
              setFormData({ ...formData, startTime: e.target.value })
            }
          />
        </div>
        <div className="input-group">
          <label>End</label>
          <input
            required
            type="time"
            value={formData.endTime}
            onChange={(e) =>
              setFormData({ ...formData, endTime: e.target.value })
            }
          />
        </div>
      </div>
      {/* ... other inputs ... */}
      <button type="submit" className="btn-primary">
        BOOK APPOINTMENT
      </button>
    </form>
  </section>
);
