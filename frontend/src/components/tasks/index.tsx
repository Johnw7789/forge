"use client";
import { Button, Input } from "@nextui-org/react";
import Link from "next/link";
import React from "react";
import { DotsIcon } from "@/components/icons/accounts/dots-icon";
import { ExportIcon } from "@/components/icons/accounts/export-icon";
import { InfoIcon } from "@/components/icons/accounts/info-icon";
import { TrashIcon } from "@/components/icons/accounts/trash-icon";
import { HouseIcon } from "@/components/icons/breadcrumb/house-icon";
import { UsersIcon } from "@/components/icons/breadcrumb/users-icon";
import { SettingsIcon } from "@/components/icons/sidebar/settings-icon";
import { TableWrapper } from "@/components/tasks-table/table";

import { taskColumns } from "../tasks-table/data";
import { tasksState } from "../state/tasks/atoms";
import { useRecoilState } from "recoil";
import { CreateTask } from "@/wailsjs/go/main/BackgroundController";


export const Tasks = () => {
  const [tasks, setTasks] = useRecoilState(tasksState)

  const createTask = () => {
    CreateTask();
  }    

  return (
    <div className="my-14 lg:px-6 max-w-[95rem] mx-auto w-full flex flex-col gap-4">
      <h3 className="text-xl font-semibold">Tasks</h3>
      <div className="flex justify-between flex-wrap gap-4 items-center">
        <div className="flex items-center gap-3 flex-wrap md:flex-nowrap">
        <div>
          <Button onClick={createTask} size="md" color="primary" >
            Start Task
          </Button>
        </div>

          {/* <SettingsIcon />
          <TrashIcon /> */}
        </div>
      </div>
      <div className="max-w-[95rem] mx-auto w-full">
        <TableWrapper TableColumns={taskColumns} Tasks={tasks} />
      </div>
    </div>
  );
};