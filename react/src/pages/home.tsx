import { Row } from "@/components/MeetingRoomTable/Row";
import { RowSkeleton } from "@/components/MeetingRoomTable/RowSkeleton";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table";
import type { MeetingRoomWithTimeSlotsResponse } from "@/dtos/meetingRoom";
import { axiosInstance } from "@/lib/axios";
import { useEffect, useState } from "react";
import { toast } from "sonner";

function Home() {
  const [isLoading, setIsLoading] = useState(true);
  const [meetingRooms, setMeetingRooms] = useState<
    MeetingRoomWithTimeSlotsResponse[]
  >([]);

  useEffect(() => {
    (async () => {
      try {
        const res = await axiosInstance.get<MeetingRoomWithTimeSlotsResponse[]>(
          "/meeting-rooms/available"
        );

        if (res.status === 200) {
          setMeetingRooms(res.data);
        } else {
          throw new Error(`unknown status code: ${res.status}`);
        }
      } catch (error) {
        console.error(
          new Error("failed to get available meeting rooms", { cause: error })
        );

        toast.error(
          "Cannot display available meeting rooms, please try again",
          { richColors: true }
        );
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  return (
    <div className="space-y-4 max-w-5xl mx-auto">
      <h1 className="text-2xl font-bold">Available Meeting Rooms</h1>
      <p>You can search the available meeting rooms and make a reservation.</p>

      <div className="grid grid-cols-1 rounded-lg border-2">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>Name</TableHead>
              <TableHead>Created At</TableHead>
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
              meetingRooms.map((item) => (
                <Row key={item.id} meetingRoom={item} />
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {meetingRooms.length === 0 ? (
        <div className="h-64 rounded-lg border-2 border-dashed flex items-center justify-center mt-6">
          <p className="text-muted-foreground">
            There are no available meeting rooms
          </p>
        </div>
      ) : (
        <></>
      )}
    </div>
  );
}

export { Home };
