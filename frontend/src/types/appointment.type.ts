import type { Appointment } from "../gen/appointment_pb";

export interface FormData {
  title: string;
  description: string;
  date: string;
  startTime: string;
  endTime: string;
  contactInformation: {
    email: string;
    name: string;
  };
}

export interface AppointmentFormProps {
  formData: FormData;
  setFormData: React.Dispatch<React.SetStateAction<FormData>>;
  onSubmit: (e: React.FormEvent) => Promise<void>;
  error: string | null;
}

export interface AppointmentListProps {
  items: Appointment[];
  selectedId: string;
  onSelect: (id: string) => void;
}

export interface AppointmentDetailsProps {
  appointment?: Appointment;
  onDelete: () => Promise<void>;
}
