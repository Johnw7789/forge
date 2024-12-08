import { User, Tooltip, Chip } from "@nextui-org/react";
import React from "react";
import { DeleteIcon } from "../icons/table/delete-icon";
import { EditIcon } from "../icons/table/edit-icon";
import { EyeIcon } from "../icons/table/eye-icon";
import { Card } from "./data";
import { DeleteCard } from "../../wailsjs/go/main/BackgroundController";
import { EditCardModal } from "./edit-card";

interface Props {
  card: Card;
  columnKey: string | React.Key;
}

export const RenderCell = ({ card, columnKey }: Props) => {
  const deleteCard = () => {
    DeleteCard(card)
  }

  const formatCard = (card) => {
    // only show last 4 digits, and replace the ones before with a *. handle cases of card length 15 and 16
    // const cardLength = card.length;
    const lastFourDigits = card.slice(-4);
    // const firstDigits = card.slice(0, cardLength - 4);
    // const asterisks = "*".repeat(firstDigits.length); 
    return "Ending in " + lastFourDigits;
  }

  // @ts-ignore
  const cellValue = card[columnKey];
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
          <span>{card.profileName}</span>
        </div>
      </div>
      );
    case "name":
        return (
          <div>
            <div>
              <span>{card.name}</span>
            </div>
          </div>
      );
    case "number":
      return (
        <div>
          <div>
            <span>{formatCard(card.number)}</span>
          </div>
        </div>
      );
      case "expiration":
        return (
          <div>
            <div>
              <span>{card.expMonth + "/" + card.expYear}</span>
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
            <EditCardModal card={card} />
          </div>
          <div>
            <Tooltip
              content="Delete card"
              color="danger"
              onClick={() => console.log("Delete card", card.id)}
            >
              <button onClick={deleteCard}>
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
