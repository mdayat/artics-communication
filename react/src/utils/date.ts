import { format } from "date-fns";

const formatDate = (dateString: string) => {
  try {
    return format(new Date(dateString), "MMM d, yyyy 'at' h:mm a");
  } catch (error) {
    console.error(new Error("invalid date string format", { cause: error }));
    return "Invalid date";
  }
};

export { formatDate };
