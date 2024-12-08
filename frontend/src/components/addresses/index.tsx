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
import { AddAddressModal } from "./add-address";
import { addressColumns, AddressColumn, Address } from "./data";
import { addressesState } from "../state/addresses/atoms";
import { useRecoilState } from "recoil";
import { ImportAddressesModal } from "./import";
import { RenderCell } from "./render-cell";

import {
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
  Pagination
} from "@nextui-org/react";

export const Addresses = () => {
  const [addresses, setAddresses] = useRecoilState(addressesState)

  const [filterValue, setFilterValue] = React.useState("");
  const hasSearchFilter = Boolean(filterValue);

  const [page, setPage] = React.useState(1);
  const rowsPerPage = 10;

  const pages = Math.ceil(addresses.length / rowsPerPage);

  const filteredItems = React.useMemo(() => {
    let filteredAddresses = [...addresses];

    if (hasSearchFilter) {
      filteredAddresses = filteredAddresses.filter((address: any) =>
        address.profileName.toLowerCase().includes(filterValue.toLowerCase()),
      );
    }

    return filteredAddresses;
  }, [addresses, filterValue]);

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
      <h3 className="text-xl font-semibold">Addresses</h3>
      <div className="flex justify-between flex-wrap gap-4 items-center">
        <div className="flex items-center gap-3 flex-wrap md:flex-nowrap">
          <Input
            classNames={{
              input: "w-full",
              mainWrapper: "w-full",
            }}
            placeholder="Search addresses"
            value={filterValue}
            onClear={() => onClear()}
            onValueChange={onSearchChange}
          />
          {/* <SettingsIcon />
          <TrashIcon /> */}
        </div>
        {/* <AddAddressModal /> */}
        <div className="flex flex-row gap-3.5 flex-wrap">
          <AddAddressModal />
          <ImportAddressesModal />
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
        <TableHeader columns={addressColumns}>
          {/* {(column: AddressColumn) => (
            <TableColumn
              key={column.uid}
              // hideHeader={column.uid === "actions"}
              // align={column.uid === "actions" ? "center" : "start"}
              align={column.uid === "actions" ? "end" : "start"}
            >
              {column.name}
            </TableColumn>
          )} */}
          <TableColumn key="profileName" style={{width: "20%"}}>
            Profile
          </TableColumn>
          <TableColumn key="name" style={{width: "20%"}}>
            Name
          </TableColumn>
          <TableColumn key="street" style={{width: "30%"}}>
            Street
          </TableColumn>
          <TableColumn key="phone" style={{width: "15%"}}>
            Phone
          </TableColumn>
          <TableColumn key="actions" style={{width: "15%"}}>
            Actions
          </TableColumn>
        </TableHeader>
        <TableBody >
          {items?.map((item: Address) => (
            <TableRow> 
              {/* {(columnKey) => (
                <TableCell>
                  {RenderCell({ address: item, columnKey: columnKey })}
                </TableCell>
              )} */}
              <TableCell style={{width: "20%"}}>
                <RenderCell address={item} columnKey="profileName" />
              </TableCell>
              <TableCell style={{width: "20%"}}>
                <RenderCell address={item} columnKey="name" />
              </TableCell>
              <TableCell style={{width: "30%"}}>
                <RenderCell address={item} columnKey="street" />
              </TableCell>
              <TableCell style={{width: "20%"}}>
                <RenderCell address={item} columnKey="phone" />
              </TableCell>
              <TableCell style={{width: "10%"}}>
                <RenderCell address={item} columnKey="actions" />
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
