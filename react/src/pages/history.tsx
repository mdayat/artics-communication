import { Row } from "@/components/ReservationHistoryTable/Row";
import { RowSkeleton } from "@/components/ReservationHistoryTable/RowSkeleton";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table";
import type { UserReservationResponse } from "@/dtos/user";
import { axiosInstance } from "@/lib/axios";
import { useEffect, useState } from "react";
import { toast } from "sonner";

function History() {
  const [isLoading, setIsLoading] = useState(true);
  const [reservations, setReservations] = useState<UserReservationResponse[]>(
    []
  );

  useEffect(() => {
    (async () => {
      try {
        const res = await axiosInstance.get<UserReservationResponse[]>(
          "/users/me/reservations"
        );

        if (res.status === 200) {
          setReservations(res.data);
        } else {
          throw new Error(`unknown status code: ${res.status}`);
        }
      } catch (error) {
        console.error(
          new Error("failed to get reservation history", { cause: error })
        );

        toast.error("Cannot display reservation history, please try again", {
          richColors: true,
        });
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  return (
    <div className="space-y-4 max-w-5xl mx-auto">
      <h1 className="text-2xl font-bold">Reservation History</h1>
      <p>You can search your reservation history.</p>

      <div className="grid grid-cols-1 rounded-lg border-2">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Meeting Room Name</TableHead>
              <TableHead>Start Date - End Date</TableHead>
              <TableHead>Cancellation Status</TableHead>
              <TableHead>Reserved At</TableHead>
              <TableHead className="text-center">Action</TableHead>
            </TableRow>
          </TableHeader>

          <TableBody>
            {isLoading ? (
              <>
                <RowSkeleton />
                <RowSkeleton />
                <RowSkeleton />
              </>
            ) : (
              reservations.map((item) => (
                <Row
                  key={item.id}
                  reservation={item}
                  setReservations={setReservations}
                />
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {reservations.length === 0 ? (
        <div className="h-64 rounded-lg border-2 border-dashed flex items-center justify-center mt-6">
          <p className="text-muted-foreground">
            You don't have reservation history
          </p>
        </div>
      ) : (
        <></>
      )}
    </div>
  );
}

export { History };
