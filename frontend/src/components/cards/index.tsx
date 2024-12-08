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
import { AddCardModal } from "./add-card";
import { cardColumns, CardColumn, Card } from "./data";
import { cardsState } from "../state/cards/atoms";
import { useRecoilState } from "recoil";
import { RenderCell } from "./render-cell";
import { ImportCardsModal} from "./import";

import {
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
  Pagination
} from "@nextui-org/react";

export const Cards = () => {
  const [cards, setCards] = useRecoilState(cardsState)

  const [filterValue, setFilterValue] = React.useState("");
  const hasSearchFilter = Boolean(filterValue);

  const [page, setPage] = React.useState(1);
  const rowsPerPage = 10;

  const pages = Math.ceil(cards.length / rowsPerPage);

  const filteredItems = React.useMemo(() => {
    let filteredCards = [...cards];

    if (hasSearchFilter) {
      filteredCards = filteredCards.filter((card: any) =>
        card.profileName.toLowerCase().includes(filterValue.toLowerCase()),
      );
    }

    return filteredCards;
  }, [cards, filterValue]);

  const items = React.useMemo(() => {
    const start = (page - 1) * rowsPerPage;
    const end = start + rowsPerPage;

    return filteredItems.slice(start, end);
  }, [page, filteredItems]);

  const onClear = React.useCallback(()=>{
    setFilterValue("")
    setPage(1)
  },[])

  const onSearchChange = React.useCallback((value) => {
    if (value) {
      setFilterValue(value);
      setPage(1);
    } else {
      setFilterValue("");
    }
  }, []);

  return (
    <div className="my-14 lg:px-6 max-w-[95rem] mx-auto w-full flex flex-col gap-4">
      <h3 className="text-xl font-semibold">Cards</h3>
      <div className="flex justify-between flex-wrap gap-4 items-center">
        <div className="flex items-center gap-3 flex-wrap md:flex-nowrap">
          <Input
            classNames={{
              input: "w-full",
              mainWrapper: "w-full",
            }}
            placeholder="Search cards"
            value={filterValue}
            onClear={() => onClear()}
            onValueChange={onSearchChange}
          />
          {/* <SettingsIcon />
          <TrashIcon /> */}
        </div>
        {/* <AddCardModal /> */}
        <div className="flex flex-row gap-3.5 flex-wrap">
          <AddCardModal />
          <ImportCardsModal />
        </div>
      </div>
      <div className="max-w-[95rem] mx-auto w-full">
        {/* <TableWrapper TableColumns={accountColumns} Accounts={accounts} /> */}

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
        <TableHeader columns={cardColumns}>
          {/* {(column: CardColumn) => (
            <TableColumn
              key={column.uid}
              // hideHeader={column.uid === "actions"}
              // align={column.uid === "actions" ? "center" : "start"} allign at the end if actions
              align={column.uid === "actions" ? "end" : "start"}
            >
              {column.name}
            </TableColumn>
          )} */}
          <TableColumn style={{width: "20%"}} key="profileName" align="start">Profile Name</TableColumn>
          <TableColumn style={{width: "20%"}} key="name" align="start">Name</TableColumn>
          <TableColumn style={{width: "30%"}} key="number" align="start">Number</TableColumn>
          <TableColumn style={{width: "15%"}} key="expiration" align="start">Expiration</TableColumn>
          <TableColumn style={{width: "15%"}} key="actions" align="end">Actions</TableColumn>
        </TableHeader>
        <TableBody >
          {items?.map((item: Card) => (
            <TableRow> 
              {/* {(columnKey) => (
                <TableCell>
                  {RenderCell({ card: item, columnKey: columnKey })}
                </TableCell>
              )} */}
              <TableCell style={{width: "20%"}}>
                <RenderCell card={item} columnKey="profileName" />
              </TableCell>
              <TableCell style={{width: "20%"}}>
                <RenderCell card={item} columnKey="name" />
              </TableCell>
              <TableCell style={{width: "30%"}}>
                <RenderCell card={item} columnKey="number" />
              </TableCell>
              <TableCell style={{width: "20%"}}>
                <RenderCell card={item} columnKey="expiration" />
              </TableCell>
              <TableCell style={{width: "10%"}}>
                <RenderCell card={item} columnKey="actions" />
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
      </div>
    </div>
  );
};
