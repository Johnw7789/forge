import { User, Tooltip, Chip } from "@nextui-org/react";
import React from "react";
import { DeleteIcon } from "../icons/table/delete-icon";
import { EditIcon } from "../icons/table/edit-icon";
import { EyeIcon } from "../icons/table/eye-icon";
import { Address } from "./data";
import { DeleteAddress } from "../../wailsjs/go/main/BackgroundController";
import { EditAddressModal } from "./edit-address";

interface Props {
  address: Address;
  columnKey: string | React.Key;
}

export const RenderCell = ({ address, columnKey }: Props) => {
  const deleteAddress = () => {
    DeleteAddress(address)
  }

  // @ts-ignore
  const cellValue = address[columnKey];
  switch (columnKey) {
    case "profile":
      return (
        // <User
        //   avatarProps={{
        //     src: "https://cdn.icon-icons.com/icons2/2699/PNG/512/amazon_tile_logo_icon_170594.png",
        //     size: "md",
        //     isBordered: false,
        //   }}
          
        //   name={cellValue}
        // >
        // </User>

        <div>
        <div>
          <span>{address.profileName}</span>
        </div>
      </div>
      );
    case "name":
        return (
          <div>
            <div>
              <span>{address.name}</span>
            </div>
          </div>
      );
    case "street":
      return (
        <div>
          <div>
            <span>{address.line1}</span>
          </div>
        </div>
      );
    case "phone":
      return (
        <div>
          <div>
          <span>{address.phone}</span>
          </div>
        </div>
      );
    case "actions":
      return (
        <div className="flex gap-4 ">
          {/* <div>
            <Tooltip content="Details">
              <button onClick={() => console.log("View account", account.id)}>
                <EyeIcon size={20} fill="#979797" />
              </button>
            </Tooltip>
          </div> */}
          <div>
            <EditAddressModal address={address} />
          </div>
          <div>
            <Tooltip
              content="Delete address"
              color="danger"
              onClick={() => console.log("Delete address", address.id)}
            >
              <button onClick={deleteAddress}>
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
