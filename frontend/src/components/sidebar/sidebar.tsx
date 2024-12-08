import React from "react";
import { Sidebar } from "./sidebar.styles";
import { Avatar, Tooltip } from "@nextui-org/react";
import { CompaniesDropdown } from "./companies-dropdown";
import { HomeIcon } from "../icons/sidebar/home-icon";
import { PaymentsIcon } from "../icons/sidebar/payments-icon";
import { BalanceIcon } from "../icons/sidebar/balance-icon";
import { AccountsIcon } from "../icons/sidebar/accounts-icon";
import { CustomersIcon } from "../icons/sidebar/customers-icon";
import { ProductsIcon } from "../icons/sidebar/products-icon";
import { ReportsIcon } from "../icons/sidebar/reports-icon";
import { DevIcon } from "../icons/sidebar/dev-icon";
import { ViewIcon } from "../icons/sidebar/view-icon";
import { SettingsIcon } from "../icons/sidebar/settings-icon";
import { DataIcon } from "../icons/sidebar/data-icon";
import { CollapseItems } from "./collapse-items";
import { SidebarItem, SidebarFooterItem } from "./sidebar-item";
import { SidebarMenu } from "./sidebar-menu";
import { FilterIcon } from "../icons/sidebar/filter-icon";
import { useSidebarContext } from "../layout/layout-context";
import { ChangeLogIcon } from "../icons/sidebar/changelog-icon";
import { usePathname } from "next/navigation";
import {Button, Spacer, Input} from "@nextui-org/react";
import {Icon} from "@iconify/react";
import { AddressIcon } from "../icons/sidebar/address-icon";

import { EventsEmit } from "../../wailsjs/runtime/runtime";

export const SidebarWrapper = () => {
  const pathname = usePathname();
  const { collapsed, setCollapsed } = useSidebarContext();

  React.useEffect(() => {
    EventsEmit("frontend:init")
  }, []);

  return (
    <aside className="h-screen z-[202] sticky top-0">
      {collapsed ? (
        <div className={Sidebar.Overlay()} onClick={setCollapsed} />
      ) : null}
      <div
        className={Sidebar({
          collapsed: collapsed,
        })}
      >
        <div className={Sidebar.Header()}>
          <CompaniesDropdown />
        </div>
        <div className="flex flex-col justify-between h-full">
          <div className={Sidebar.Body()}>
            {/* <SidebarItem
              title="Home"
              icon={<HomeIcon />}
              isActive={pathname === "/"}
              href="/"
            /> */}
              <SidebarItem
                isActive={pathname === "/"}
                title="Tasks"
                icon={<HomeIcon />}
                href="/"
              />
            <SidebarMenu title="Main Menu">
              <SidebarItem
                isActive={pathname === "/accounts"}
                title="Accounts"
                icon={<AccountsIcon />}
                href="accounts"
              />
              <SidebarItem
                isActive={pathname === "/data"}
                title="Data"
                icon={<DataIcon />}
                href="data"
              />
            </SidebarMenu>

            <SidebarMenu title="Profiles">
              <SidebarItem
                isActive={pathname === "/addresses"}
                title="Addresses"
                icon={<AddressIcon />}
                href="addresses"
              />
              <SidebarItem
                isActive={pathname === "/cards"}
                title="Cards"
                icon={<PaymentsIcon />}
                href="cards"
              />
            </SidebarMenu>

            <SidebarMenu title="General">
              <SidebarItem
                isActive={pathname === "/settings"}
                title="Settings"
                icon={<SettingsIcon />}
                href="settings"
              />
              {/* <SidebarItem
                isActive={pathname === "/changelog"}
                title="Changelog"
                icon={<ChangeLogIcon />}
                href="changelog"
              /> */}
            </SidebarMenu>
          </div>

          <div >
          <div className="mt-auto flex flex-col">
          {/* <SidebarFooterItem
            title="Help & Information"
            icon={ <Icon 
              className="text-default-500" 
              icon="solar:info-circle-line-duotone" 
              width={24} />
            }
          /> */}

          {/* <SidebarFooterItem
            title="Log Out"
            icon={ <Icon
              className="rotate-180 text-default-500"
              icon="solar:minus-circle-line-duotone"
              width={24} />
            }   
          /> */}
          </div>
          </div>
        </div>
      </div>
    </aside>
  );
};
