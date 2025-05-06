import { Skeleton } from "../ui/Skeleton";
import { TableCell, TableRow } from "../ui/Table";

function RowSkeleton() {
  return (
    <TableRow>
      <TableCell>
        <Skeleton className="h-4 w-28" />
      </TableCell>

      <TableCell>
        <Skeleton className="h-4 w-28" />
      </TableCell>

      <TableCell>
        <Skeleton className="h-4 w-48" />
      </TableCell>

      <TableCell>
        <Skeleton className="h-4 w-28" />
      </TableCell>

      <TableCell>
        <Skeleton className="h-4 w-48" />
      </TableCell>

      <TableCell className="flex justify-center items-center">
        <Skeleton className="h-9 w-28" />
      </TableCell>
    </TableRow>
  );
}

export { RowSkeleton };
