import {
  Avatar,
} from "@nextui-org/react";

import React from "react";

import { useRecoilState } from "recoil";
import { authState } from "../state/auth/atoms";

export const UserDropdown = () => {
  const [auth, setAuth] = useRecoilState(authState);

  return (
        <Avatar
        className="mb-3"
        as="button"
        // color="secondary"
        color="warning"
        size="md"
        src={auth.discordImage}
      />
  );
};
