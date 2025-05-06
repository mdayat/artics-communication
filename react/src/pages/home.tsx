import { MeetingRoomTable } from "@/components/MeetingRoomTable";

function Home() {
  return (
    <div className="space-y-4">
      <h1 className="text-2xl font-bold">Available Meeting Rooms</h1>
      <p>You can search the available meeting rooms and make a reservation.</p>
      <MeetingRoomTable />
    </div>
  );
}

export { Home };
