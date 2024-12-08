import { User, Tooltip, Chip } from "@nextui-org/react";
import React from "react";
import { DeleteIcon } from "../icons/table/delete-icon";
import { Task } from "./data";

import { StopTask } from "@/wailsjs/go/main/BackgroundController";

interface Props {
  task: Task;
  columnKey: string | React.Key;
}

export const RenderCell = ({ task, columnKey }: Props) => {
  // @ts-ignore
  const cellValue = task[columnKey];
  switch (columnKey) {
    case "name":
      return (
        <User
          avatarProps={{
            src: "https://cdn.icon-icons.com/icons2/2699/PNG/512/amazon_tile_logo_icon_170594.png",
            size: "md",
            isBordered: false,
          }}
          name={task.name}
        >
        </User>
      );

    case "email":
      return (
        <div>
          <div>
            <span>{task.email}</span>
          </div>
        </div>
      );
      case "password":
        return (
          <div>
            <div>
              <span>{task.password}</span>
            </div>
          </div>
        );
    case "proxy":
        return (
          <div>
            <div>
              <span>{task.proxy ? task.proxy : "Localhost"}</span>
            </div>
          </div>
        );
    
    case "status":
      return (
        // <Chip
        //   size="sm"
        //   variant="flat"
        //   color={
        //     cellValue === "active"
        //       ? "success"
        //       : cellValue === "On Hold"
        //       ? "danger"
        //       : "warning"
        //   }
        // >
        //   <span className="capitalize text-xs">{cellValue}</span>
        // </Chip>
        <div>
        <div>
          {/* <span >{task.status}</span> */}
          {task.status === "Account Created" || task.status === "2FA Completed" ? (
            <Chip color="success" variant="faded">{task?.status?.length > 80 ? task?.status?.substring(0, 80) + "..." : task.status}</Chip>
          ) : task.status.includes("Error") ? (
            <Chip color="warning" variant="faded">{task?.status?.length > 80 ? task?.status?.substring(0, 80) + "..." : task.status}</Chip>
          ) : (
              <Chip variant="faded">{task?.status?.length > 80 ? task?.status?.substring(0, 80) + "..." : task.status}</Chip>
          )}
        </div>
      </div>
      );

    case "actions":
      return (
        <div className="flex justify-end gap-4 ">
          <div>
            <Tooltip
              content="Delete task"
              color="danger"
              onClick={() => console.log("Delete task", task.id)}
            >
              <button onClick={() => {StopTask(task.id)}}>
                <DeleteIcon size={20} fill="#FF0080" />
              </button>
            </Tooltip>
          </div>
        </div>
      );
    default:
      return cellValue;
  }
};
