import { User, Tooltip, Chip } from "@nextui-org/react";
import React from "react";
import { DeleteIcon } from "../icons/table/delete-icon";
import { Account } from "./data";
import { DeleteAccount } from "../../wailsjs/go/main/BackgroundController";
import { EditAccountModal } from "./edit-account";
import  { SubtaskModal } from "./subtask";


interface Props {
  account: Account;
  columnKey: string | React.Key;
}

export const RenderCell = ({ account, columnKey }: Props) => {
  const deleteAccount = () => {
    DeleteAccount(account as any)
  }

  const cellValue = account[columnKey];
  switch (columnKey) {
    case "name":
      return (
        <User
          avatarProps={{
            src: "https://cdn.icon-icons.com/icons2/2699/PNG/512/amazon_tile_logo_icon_170594.png",
            size: "md",
            isBordered: false,
          }}
          
          name={cellValue ? cellValue : "N/A"}
        >
        </User>
      );
    case "email":
      return (
        <div>
          <div>
            <span>{account.email ? account.email : "N/A"}</span>
          </div>
        </div>
      );
    case "password":
      return (
        <div>
          <div>
            <span>{account.password ? account.password : "N/A"}</span>
          </div>
        </div>
      );
      case "prime":
        return (
          <div>
            <div>
            <span>{account.prime ? "Yes" : "No"}</span>
            </div>
          </div>
        );
    case "status":
      return (
        <div>
          <div>
            <span>
              {(account.status === "Prime Activated" || account.status === "Submitted Profile") ? (
                <Chip color="success" variant="faded">{account.status}</Chip>
              ) : account.status?.includes("Error") ? (
                <Chip color="warning" variant="faded">{account?.status?.length > 80 ? account?.status?.substring(0, 80) + "..." : account.status}</Chip>
              ) : account.status === "" ? (
                <Chip variant="faded">Idle</Chip>
              ) : (
                  <Chip variant="faded">{account?.status?.length > 80 ? account?.status?.substring(0, 80) + "..." : account.status}</Chip>
              )}
              </span>
          </div>
        </div>
      );
  

    case "actions":
      return (
        <div className="flex gap-4 ">
          <div>
            <EditAccountModal account={account} />
          </div>
          <div>
            <Tooltip
              content="Delete account"
              color="danger"
              onClick={() => console.log("Delete account", account.id)}
            >
              <button onClick={deleteAccount}>
                <DeleteIcon size={20} fill="#FF0080" />
              </button>
            </Tooltip>
          </div>
          <div>
            <SubtaskModal account={account} />
          </div>
        </div>
      );
    default:
      return cellValue;
  }
};
