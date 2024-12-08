import { Navbar, NavbarContent, Chip } from "@nextui-org/react";
import React from "react";

import { UserDropdown } from "./user-dropdown";
import { WindowMinimise } from "../../wailsjs/runtime/runtime"
import { Quit } from "../../wailsjs/runtime/runtime"

interface Props {
  children: React.ReactNode;
}

export const NavbarWrapper = ({ children }: Props) => {
  return (
    <div className="relative flex flex-col flex-1 ">
      <div style={{"--wails-draggable":"drag"} as React.CSSProperties} className="flex justify-end mr-4 mb-0" >
        <div className="mr-3 cursor-pointer" onClick={() => WindowMinimise()} id="minimize">–</div>
        <div className="cursor-pointer" onClick={() => Quit()} id="close" >✕</div>
      </div>
      <div style={{"--wails-draggable":"drag"} as React.CSSProperties}>
      <Navbar 
        isBordered
        className="w-full h-14"
        classNames={{
          wrapper: "w-full max-w-full",
        }}
      >
        <NavbarContent className="w-full max-md:hidden">
          {/* <Chip color="success" variant="dot">Connected</Chip> */}
          <Chip
            className="mb-5 select-none cursor-default"
            startContent={<span className="w-2 h-2 ml-1 rounded-full bg-success"></span>}
            variant="faded"
            color="success"
          >
            Connected
          </Chip>

          <Chip
            className="mb-5 select-none cursor-default"
            // startContent={<span className="w-2 h-2 rounded-full"></span>}
            variant="faded"
          >
            v0.9
          </Chip>
        </NavbarContent>
        <NavbarContent
          justify="end"
          className="w-fit data-[justify=end]:flex-grow-0"
        >
          <NavbarContent>
            <UserDropdown />
          </NavbarContent>
        </NavbarContent>
      </Navbar>
      </div>
      {children}
    </div>
  );
};
