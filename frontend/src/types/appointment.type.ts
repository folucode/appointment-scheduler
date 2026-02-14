export type Appointment = {
  id: string;
  userId: string;
  title: string;
  description: string;
  contactInformation?: {
    name: string;
    email: string;
  };
  date: string;
  startTime: string;
  endTime: string;
};
