import { useEffect, useState } from "react";
import {
  Table,
  TableBody,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/Table";
import { Card, CardContent } from "@/components/ui/Card";
import type { MeetingRoomWithTimeSlotsResponse } from "@/dtos/meetingRoom";
import { toast } from "sonner";
import { axiosInstance } from "@/lib/axios";
import { MeetingRoomTableRow } from "./Row";
import { MeetingRoomTableRowSkeleton } from "./RowSkeleton";

function MeetingRoomTable() {
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
    <Card className="overflow-hidden py-0">
      <CardContent className="p-0">
        <div className="overflow-x-auto">
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
                  <MeetingRoomTableRowSkeleton />
                  <MeetingRoomTableRowSkeleton />
                  <MeetingRoomTableRowSkeleton />
                </>
              ) : (
                meetingRooms.map((item) => (
                  <MeetingRoomTableRow key={item.id} meetingRoom={item} />
                ))
              )}
            </TableBody>
          </Table>
        </div>
      </CardContent>
    </Card>
  );
}

export { MeetingRoomTable };
