import { axiosInstance } from "@/lib/axios";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "../ui/Table";
import { RowSkeleton } from "./RowSkeleton";
import { Row } from "./Row";
import type { EnrichedReservationResponse } from "@/dtos/reservation";

function ReservationTable() {
  const [isLoading, setIsLoading] = useState(true);
  const [reservations, setReservations] = useState<
    EnrichedReservationResponse[]
  >([]);

  useEffect(() => {
    (async () => {
      try {
        const res = await axiosInstance.get<EnrichedReservationResponse[]>(
          "/reservations"
        );

        if (res.status === 200) {
          setReservations(res.data);
        } else if (res.status === 403) {
          toast.error("Insufficient permissions", { richColors: true });
        } else {
          throw new Error(`unknown status code: ${res.status}`);
        }
      } catch (error) {
        console.error(
          new Error("failed to get reservations", { cause: error })
        );

        toast.error("Cannot display reservations, please try again", {
          richColors: true,
        });
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  return (
    <>
      <h1 className="text-2xl font-bold">Reservation</h1>
      <p>You can search a reservation through reservation list.</p>

      <div className="grid grid-cols-1 rounded-lg border-2">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>User</TableHead>
              <TableHead>Meeting Room</TableHead>
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
    </>
  );
}

export { ReservationTable };
