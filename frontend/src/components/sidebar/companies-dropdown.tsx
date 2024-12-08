"use client";
import {
  Dropdown,
  DropdownItem,
  DropdownMenu,
  DropdownSection,
  DropdownTrigger,
} from "@nextui-org/react";
import React, { useState } from "react";
import { AcmeIcon } from "../icons/acme-icon";
import { BottomIcon } from "../icons/sidebar/bottom-icon";
import {Avatar, Button, Spacer, useDisclosure} from "@nextui-org/react";
import type {SVGProps} from "react";
// import logo from "./appicon.png";

type IconSvgProps = SVGProps<SVGSVGElement> & {
  size?: number;
};

interface Company {
  name: string;
  location: string;
  logo: React.ReactNode;
}

let img = process.env.PUBLIC_URL + "/appicon.png"
console.log(img);

const AcmeLogo = (props: IconSvgProps) => (
  <svg fill="none" height="36" viewBox="0 0 32 32" width="36" {...props}>
    <path
      clipRule="evenodd"
      d="M17.6482 10.1305L15.8785 7.02583L7.02979 22.5499H10.5278L17.6482 10.1305ZM19.8798 14.0457L18.11 17.1983L19.394 19.4511H16.8453L15.1056 22.5499H24.7272L19.8798 14.0457Z"
      fill="currentColor"
      fillRule="evenodd"
    />
  </svg>
);
import { User, Tooltip, Chip } from "@nextui-org/react";


export const CompaniesDropdown = () => {
  const [company, setCompany] = useState<Company>({
    name: "Acme Co.",
    location: "Palo Alto, CA",
    logo: <AcmeIcon />,
  });
  return (
    // <Dropdown
    //   classNames={{
    //     base: "w-full min-w-[260px]",
    //   }}
    // >
    //   <DropdownTrigger className="cursor-pointer">
    //     <div className="flex items-center gap-2">
    //       {company.logo}
    //       <div className="flex flex-col gap-4">
    //         <h3 className="text-xl font-medium m-0 text-default-900 -mb-4 whitespace-nowrap">
    //           {company.name}
    //         </h3>
    //         <span className="text-xs font-medium text-default-500">
    //           {company.location}
    //         </span>
    //       </div>
    //       <BottomIcon />
    //     </div>
    //   </DropdownTrigger>
    //   <DropdownMenu
    //     onAction={(e) => {
    //       if (e === "1") {
    //         setCompany({
    //           name: "Facebook",
    //           location: "San Fransico, CA",
    //           logo: <AcmeIcon />,
    //         });
    //       }
    //       if (e === "2") {
    //         setCompany({
    //           name: "Instagram",
    //           location: "Austin, Tx",
    //           logo: <AcmeLogo />,
    //         });
    //       }
    //       if (e === "3") {
    //         setCompany({
    //           name: "Twitter",
    //           location: "Brooklyn, NY",
    //           logo: <AcmeIcon />,
    //         });
    //       }
    //       if (e === "4") {
    //         setCompany({
    //           name: "Acme Co.",
    //           location: "Palo Alto, CA",
    //           logo: <AcmeIcon />,
    //         });
    //       }
    //     }}
    //     aria-label="Avatar Actions"
    //   >
    //     <DropdownSection title="Companies">
    //       <DropdownItem
    //         key="1"
    //         startContent={<AcmeIcon />}
    //         description="San Fransico, CA"
    //         classNames={{
    //           base: "py-4",
    //           title: "text-base font-semibold",
    //         }}
    //       >
    //         Facebook
    //       </DropdownItem>
    //       <DropdownItem
    //         key="2"
    //         startContent={<AcmeLogo />}
    //         description="Austin, Tx"
    //         classNames={{
    //           base: "py-4",
    //           title: "text-base font-semibold",
    //         }}
    //       >
    //         Instagram
    //       </DropdownItem>
    //       <DropdownItem
    //         key="3"
    //         startContent={<AcmeIcon />}
    //         description="Brooklyn, NY"
    //         classNames={{
    //           base: "py-4",
    //           title: "text-base font-semibold",
    //         }}
    //       >
    //         Twitter
    //       </DropdownItem>
    //       <DropdownItem
    //         key="4"
    //         startContent={<AcmeIcon />}
    //         description="Palo Alto, CA"
    //         classNames={{
    //           base: "py-4",
    //           title: "text-base font-semibold",
    //         }}
    //       >
    //         Acme Co.
    //       </DropdownItem>
    //     </DropdownSection>
    //   </DropdownMenu>
    // </Dropdown>
    <div>
    <div className="flex items-center gap-2">
        <div >
          {/* <AcmeLogo className="text-background" /> */}
          <Avatar size="md" src="./appicon.png" />
        </div>
        <span className="text-szmall font-bold uppercase text-foreground">The Forge</span>
    </div>
    </div>

  );
};
