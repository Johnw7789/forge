import NextLink from "next/link";
import React from "react";
import { useSidebarContext } from "../layout/layout-context";
import clsx from "clsx";
import {Button, Spacer, Input} from "@nextui-org/react";
// import "./sidebar-item.css";
import { MouseEvent } from 'react';

interface ItemProps {
  title: string;
  icon: React.ReactNode;
  isActive?: boolean;
  href?: string;
}

// we dont always want an onclick
interface Props {
  title: string;
  icon: React.ReactNode;
  onClick?: () => void;
}

export const SidebarItem = ({ icon, title, isActive, href = "" }: ItemProps) => {
  const { collapsed, setCollapsed } = useSidebarContext();

  const handleDragStart = (event: MouseEvent<HTMLAnchorElement>) => {
    event.preventDefault();
  };

  const handleClick = () => {
    if (window.innerWidth < 768) {
      setCollapsed();
    }
  };
  return (
    <NextLink
      href={href}
      className="text-default-900 active:bg-none max-w-full"
      onDragStart={handleDragStart}
    >
      <div
        className={clsx(
          isActive
            ? "bg-primary-500 [&_svg_path]:fill-white"
            : "hover:bg-default-100",
          "flex gap-2 w-full min-h-[44px] h-full items-center px-3.5 rounded-xl cursor-pointer transition-all duration-150 active:scale-[0.98]"
        )}
        onClick={handleClick}
      >
        {icon}
        <span className="text-default-900 ">{title}</span>
      </div>
    </NextLink>
  );
};

export const SidebarFooterItem = ({ icon, title, onClick}: Props) => {
  const { collapsed, setCollapsed } = useSidebarContext();

  const handleClick = () => {
    if (window.innerWidth < 768) {
      setCollapsed();
    }
  };
  return (
    <NextLink 
      href={""}
      onClick={onClick}
      className="text-default-500 active:bg-none max-w-full"
    >
      <div
        className={clsx(
          "hover:bg-default-100",
          "flex gap-2 w-full min-h-[44px] h-full items-center px-3.5 rounded-xl cursor-pointer transition-all duration-150 active:scale-[0.98]"
        )}
        onClick={handleClick}
      >
        {icon}
        <span className="text-default-500">{title}</span>
      </div>
    </NextLink>
  );
};