import { TableCell, TableRow } from "../ui/Table";
import { Button } from "../ui/Button";
import { useState, type Dispatch, type SetStateAction } from "react";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "../ui/AlertDialog";
import { toast } from "sonner";
import { axiosInstance } from "@/lib/axios";
import { formatDate } from "@/utils/date";
import type { ReservationResponse } from "@/dtos/reservation";
import type { AxiosResponse } from "axios";
import type { UserReservationResponse } from "@/dtos/user";
import { Badge } from "../ui/Badge";

interface RowProps {
  reservation: UserReservationResponse;
  setReservations: Dispatch<SetStateAction<UserReservationResponse[]>>;
}

function Row({ reservation, setReservations }: RowProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);

  const handleCancel = async () => {
    setIsLoading(true);
    try {
      const res = await axiosInstance.patch<
        ReservationResponse,
        AxiosResponse<ReservationResponse>
      >(`/users/me/reservations/${reservation.id}`);

      if (res.status === 200) {
        toast.success("Reservation successfully canceled", {
          richColors: true,
        });

        setReservations((reservations) =>
          reservations.map((item) => {
            if (item.id !== reservation.id) {
              return item;
            }

            return {
              ...item,
              canceled: true,
              canceled_at: res.data.canceled_at,
            };
          })
        );

        setIsOpen(false);
      } else if (res.status === 404) {
        toast.error("Reservation not found", { richColors: true });
      } else {
        throw new Error(`unknown status code: ${res.status}`);
      }
    } catch (error) {
      console.error(
        new Error("failed to cancel reservation", { cause: error })
      );

      toast.error("Cancel Reservation failed, please try again", {
        richColors: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <TableRow>
      <TableCell>{reservation.meeting_room.name}</TableCell>
      <TableCell>
        {formatDate(reservation.time_slot.start_date)} -{" "}
        {formatDate(reservation.time_slot.end_date)}
      </TableCell>

      <TableCell>
        <Badge className={reservation.canceled ? "bg-red-600" : "bg-green-600"}>
          {reservation.canceled ? "Canceled" : "Not Canceled"}
        </Badge>
      </TableCell>

      <TableCell>{formatDate(reservation.reserved_at)}</TableCell>

      <TableCell className="text-center">
        <Button onClick={() => setIsOpen(true)} className="cursor-pointer">
          Cancel
        </Button>

        <AlertDialog open={isOpen} onOpenChange={setIsOpen}>
          <AlertDialogContent>
            <AlertDialogHeader>
              <AlertDialogTitle>Are you sure?</AlertDialogTitle>
              <AlertDialogDescription>
                You will cancel your reservation on{" "}
                <span className="font-semibold">
                  {reservation.meeting_room.name}
                </span>{" "}
                Meeting Room
              </AlertDialogDescription>
            </AlertDialogHeader>

            <AlertDialogFooter>
              <AlertDialogCancel className="cursor-pointer">
                Cancel
              </AlertDialogCancel>

              <Button
                disabled={isLoading}
                onClick={handleCancel}
                className="cursor-pointer"
              >
                {isLoading ? "Cancelling..." : "Continue"}
              </Button>
            </AlertDialogFooter>
          </AlertDialogContent>
        </AlertDialog>
      </TableCell>
    </TableRow>
  );
}

export { Row };
