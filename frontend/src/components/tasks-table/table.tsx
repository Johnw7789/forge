import {
  Link,
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
  Pagination,
} from "@nextui-org/react";
import React from "react";

import { TaskColumn, Task } from "./data";
import { RenderCell } from "./render-cell";

export const TableWrapper = (props: {TableColumns: TaskColumn[], Tasks: Task[]}) => {
  const [page, setPage] = React.useState(1);
  const rowsPerPage = 10;

  const pages = Math.ceil(props.Tasks?.length / rowsPerPage);

  const items = React.useMemo(() => {
    const start = (page - 1) * rowsPerPage;
    const end = start + rowsPerPage;

    return props.Tasks?.slice(start, end);
  }, [page, props.Tasks]);

  return (
    <div className=" w-full flex flex-col gap-4">
      <Table aria-label="Example table with custom cells"
              bottomContent={
                <div className="flex w-full justify-center">
                  <Pagination
                    isCompact
                    showControls
                    showShadow
                    color="primary"
                    page={page}
                    total={pages}
                    onChange={(page) => setPage(page)}
                  />
                </div>}
      >
        
        <TableHeader columns={props.TableColumns}>
          {(column: TaskColumn) => (
            <TableColumn
              key={column.uid}
              hideHeader={column.uid === "actions"}
              align={column.uid === "actions" ? "center" : "start"}
            >
              {column.name}
            </TableColumn>
          )}
        </TableHeader>
        <TableBody >
          {items?.map((item) => (
            <TableRow key={item.id}>
              {/* <TableCell style={{width: "20%"}}>
                {RenderCell({ task: item, columnKey: "name" })}
              </TableCell> */}
              <TableCell style={{width: "25%"}}>
                {RenderCell({ task: item, columnKey: "email" })}
              </TableCell>
              <TableCell style={{width: "15%"}}>
                {RenderCell({ task: item, columnKey: "password" })}
              </TableCell>
              <TableCell style={{width: "25%"}}>
                {RenderCell({ task: item, columnKey: "proxy" })}
              </TableCell>
              <TableCell style={{width: "25%"}}>
                {RenderCell({ task: item, columnKey: "status" })}
              </TableCell>
              <TableCell style={{width: "10%"}}>
                {RenderCell({ task: item, columnKey: "actions" })}
              </TableCell>
              {/* {(columnKey) => (
                <TableCell>
                  {RenderCell({ task: item, columnKey: columnKey })}
                </TableCell>
              )} */}
            </TableRow>
          ))
        }
        </TableBody>
      </Table>
    </div>
  );
};
