import type { MeetingRoomWithTimeSlotsResponse } from "@/dtos/meetingRoom";
import { TableCell, TableRow } from "../ui/Table";
import { Button } from "../ui/Button";
import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "../ui/Dialog";
import {
  AlertDialog,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "../ui/AlertDialog";
import { toast } from "sonner";
import { axiosInstance } from "@/lib/axios";
import { formatDate } from "@/utils/date";
import type {
  CreateReservationRequest,
  ReservationResponse,
} from "@/dtos/reservation";
import type { AxiosResponse } from "axios";
import { Loader2 } from "lucide-react";

interface RowProps {
  meetingRoom: MeetingRoomWithTimeSlotsResponse;
}

function Row({ meetingRoom }: RowProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);

  const handleReserve = async (timeSlotId: string) => {
    setIsLoading(true);
    try {
      const res = await axiosInstance.post<
        ReservationResponse,
        AxiosResponse<ReservationResponse>,
        CreateReservationRequest
      >("/users/me/reservations", {
        meeting_room_id: meetingRoom.id,
        time_slot_id: timeSlotId,
      });

      if (res.status === 201) {
        toast.success("Reservation successfully created", { richColors: true });
        setIsOpen(false);
      } else if (res.status === 409) {
        toast.error("Sorry, this time slot already reserved by someone else", {
          richColors: true,
        });
      } else {
        throw new Error(`unknown status code: ${res.status}`);
      }
    } catch (error) {
      console.error(new Error("failed to reserve", { cause: error }));
      toast.error("Reservation failed, please try again", {
        richColors: true,
      });
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <>
      <TableRow>
        <TableCell>{meetingRoom.name}</TableCell>
        <TableCell>{formatDate(meetingRoom.created_at)}</TableCell>
        <TableCell className="text-center">
          <Button onClick={() => setIsOpen(true)} className="cursor-pointer">
            View Detail
          </Button>
        </TableCell>
      </TableRow>

      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogContent className="sm:max-w-md">
          <DialogHeader>
            <DialogTitle>{meetingRoom.name} - Time Slots</DialogTitle>
            <DialogDescription>
              Available time slots for reservation
            </DialogDescription>
          </DialogHeader>

          <div className="space-y-4 max-h-[60vh] overflow-y-auto">
            {meetingRoom.time_slots.map((timeSlot, index) => (
              <div
                key={timeSlot.id}
                className="border rounded-lg p-4 space-y-3"
              >
                <div className="flex justify-between items-center">
                  <div className="font-medium">Time Slot {index + 1}</div>
                  <AlertDialog>
                    <AlertDialogTrigger>
                      <Button size="sm" className="cursor-pointer">
                        Reserve
                      </Button>
                    </AlertDialogTrigger>

                    <AlertDialogContent>
                      <AlertDialogHeader>
                        <AlertDialogTitle>Are you sure?</AlertDialogTitle>
                        <AlertDialogDescription>
                          You will reserve for{" "}
                          <span className="font-semibold italic">
                            {formatDate(timeSlot.start_date)}
                          </span>{" "}
                          until{" "}
                          <span className="font-semibold italic">
                            {formatDate(timeSlot.end_date)}
                          </span>
                          .
                        </AlertDialogDescription>
                      </AlertDialogHeader>

                      <AlertDialogFooter>
                        <AlertDialogCancel className="cursor-pointer">
                          Cancel
                        </AlertDialogCancel>

                        <Button
                          disabled={isLoading}
                          onClick={() => handleReserve(timeSlot.id)}
                          className="cursor-pointer"
                        >
                          {isLoading ? (
                            <>
                              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                              Reserving
                            </>
                          ) : (
                            "Continue"
                          )}
                        </Button>
                      </AlertDialogFooter>
                    </AlertDialogContent>
                  </AlertDialog>
                </div>

                <div className="grid grid-cols-1 gap-1 text-sm">
                  <div className="grid grid-cols-3">
                    <span className="text-muted-foreground">Start Date:</span>
                    <span className="col-span-2">
                      {formatDate(timeSlot.start_date)}
                    </span>
                  </div>

                  <div className="grid grid-cols-3">
                    <span className="text-muted-foreground">End Date:</span>
                    <span className="col-span-2">
                      {formatDate(timeSlot.end_date)}
                    </span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </DialogContent>
      </Dialog>
    </>
  );
}

export { Row };
