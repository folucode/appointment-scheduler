import { timestampFromDate, type Timestamp } from "@bufbuild/protobuf/wkt";

export const convertToProtobufTime = (
  date: string,
  startTime: string,
  endTime: string,
): {
  date: Timestamp;
  start: Timestamp;
  end: Timestamp;
} => {
  const start = new Date(`${date} ${startTime}`);
  const end = new Date(`${date} ${endTime}`);

  return {
    date: timestampFromDate(new Date(date)),
    start: timestampFromDate(start),
    end: timestampFromDate(end),
  };
};
