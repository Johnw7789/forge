"use client";

import React from "react";
import { AddAccountModal } from "./add-account";
import { accountColumns, AccountColumn, Account } from "./data";
import { accountsState, selectedAccountsState } from "../state/accounts/atoms";
import { useRecoilState } from "recoil";
import { RenderCell } from "./render-cell";
import {ExportAccountsModal} from "./export";

import {
  Table,
  TableBody,
  TableCell,
  TableColumn,
  TableHeader,
  TableRow,
  Pagination,
  Selection,
  Input
} from "@nextui-org/react";

export const Accounts = () => {
  const [accounts, setAccounts] = useRecoilState(accountsState)
  const [selectedAccounts, setSelectedAccounts] = useRecoilState(selectedAccountsState)
  const [selectedKeys, setSelectedKeys] = React.useState<Selection>(new Set([]));

  const [filterValue, setFilterValue] = React.useState("");
  const hasSearchFilter = Boolean(filterValue);

  const [page, setPage] = React.useState(1);
  const rowsPerPage = 10;

  const pages = Math.ceil(accounts.length / rowsPerPage);

  React.useEffect(() => {
      if (selectedKeys !== "all") {
        setSelectedAccounts(accounts.filter((account: Account) => selectedKeys.has(account.id)));
      } else {
        setSelectedAccounts(accounts);          
      }
  }, [selectedKeys])

  const filteredItems = React.useMemo(() => {
    let filteredUsers = [...accounts];

    if (hasSearchFilter) {
      filteredUsers = filteredUsers.filter((account: any) =>
        account.email.toLowerCase().includes(filterValue.toLowerCase()),
      );
    }

    return filteredUsers;
  }, [accounts, filterValue]);

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
      <h3 className="text-xl font-semibold">Accounts</h3>
      <div className="flex justify-between flex-wrap gap-4 items-center">
        <div className="flex items-center gap-3 flex-wrap md:flex-nowrap">
          <Input
            classNames={{
              input: "w-full",
              mainWrapper: "w-full",
            }}
            placeholder="Search accounts"
            value={filterValue}
            onClear={() => onClear()}
            onValueChange={onSearchChange}
          />
          {/* <SettingsIcon />
          <TrashIcon /> */}
        </div>
        {/* <AddAccountModal /> */}
        <div className="flex flex-row gap-3.5 flex-wrap">
          <AddAccountModal />
          <ExportAccountsModal/>
        </div>
      </div>
      <div className="max-w-[95rem] mx-auto w-full">
        {/* <TableWrapper TableColumns={accountColumns} Accounts={accounts} /> */}

        <div className=" w-full flex flex-col gap-4">
      <Table aria-label="Example table with custom cells"
        selectedKeys={selectedKeys}
        onSelectionChange={setSelectedKeys}
        selectionMode="multiple"
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
        <TableHeader columns={accountColumns}>
          {/* {(column: AccountColumn) => (
            <TableColumn
              key={column.uid}
              // hideHeader={column.uid === "actions"}
              align={column.uid === "actions" ? "end" : "start"}
            >
              {column.name}
            </TableColumn>
          )} */}
          <TableColumn key="name" align="start" style={{width: "20%"}}>
            Name
          </TableColumn>
          <TableColumn key="email" align="start" style={{width: "25%"}}>
            Email
          </TableColumn>
          <TableColumn key="password" align="start" style={{width: "15%"}}>
            Password
          </TableColumn>
          <TableColumn key="prime" align="start" style={{width: "5%"}}>
            Prime
          </TableColumn>
          <TableColumn key="status" align="start" style={{width: "20%"}}>
            Status
          </TableColumn>
          <TableColumn key="actions" align="end" style={{width: "15%"}}>
            Actions
          </TableColumn>
        </TableHeader>
        <TableBody >
          {items?.map((item: Account) => (
            <TableRow key={item.id}> 
              {/* {(columnKey) => (
                <TableCell>
                  {RenderCell({ account: item, columnKey: columnKey })}
                </TableCell>
              )} */}

              <TableCell style={{width: "20%"}}>
                <RenderCell account={item} columnKey="name" />
              </TableCell>
              <TableCell style={{width: "25%"}}>
                <RenderCell account={item} columnKey="email" />
              </TableCell>
              <TableCell style={{width: "15%"}}>
                <RenderCell account={item} columnKey="password" />
              </TableCell>
              <TableCell style={{width: "5%"}}>
                <RenderCell account={item} columnKey="prime" />
              </TableCell>
              <TableCell style={{width: "25%"}}>
                <RenderCell account={item} columnKey="status" />
              </TableCell>
              <TableCell style={{width: "10%"}}>
                <RenderCell account={item} columnKey="actions" />
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
