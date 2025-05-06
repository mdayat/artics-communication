import { MeetingRoomTable } from "@/components/MeetingRoomTable";
import { ReservationTable } from "@/components/ReservationTable";
import { useAuthContext } from "@/contexts/AuthProvider";

function Home() {
  const { user } = useAuthContext();

  return (
    <div className="space-y-4 max-w-5xl mx-auto">
      {user?.role === "user" ? <MeetingRoomTable /> : <ReservationTable />}
    </div>
  );
}

export { Home };
