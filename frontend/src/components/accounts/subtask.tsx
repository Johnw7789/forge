import {
  Button,
  Tooltip,
  Input,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
  DropdownTrigger,
  Dropdown,
  DropdownMenu,
  Select,
  SelectItem,
  DropdownItem,
} from "@nextui-org/react";
import React from "react";

import { toast } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

import { useRecoilState } from "recoil";
import { cardsState } from "../state/cards/atoms";
import { addressesState } from "../state/addresses/atoms";

import { Account } from "./data";

import { CreatePrimeTask } from "@/wailsjs/go/main/BackgroundController";
import { CreateInfoTask } from "@/wailsjs/go/main/BackgroundController";

interface EditProps {
  account: Account,
}

const VerticalDotsIcon = () => (
  <svg
    width="24"
    height="24"
    fill="none"
    viewBox="0 0 24 24"
    xmlns="http://www.w3.org/2000/svg"
  >
    <path
      className="fill-default-400"
      fillRule="evenodd"
      clipRule="evenodd"
      d="M12 10c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0-6c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2zm0 12c-1.1 0-2 .9-2 2s.9 2 2 2 2-.9 2-2-.9-2-2-2z"
      fill="#969696"
    />
  </svg>
);

export const SubtaskModal = ({ account }: EditProps) => {
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  const [cardProfile, setCardProfile] = React.useState(""); 
  const [addressProfile, setAddressProfile] = React.useState("");
  const [mode, setMode] = React.useState("");
  const [cards, setCards] = useRecoilState(cardsState);
  const [addresses, setAddresses] = useRecoilState(addressesState);

  const handleCardProfile = (event) => {
    const { value } = event.target;
    setCardProfile(value);
  };

  const handleAddressProfile = (event) => {
    const { value } = event.target;
    setAddressProfile(value);
  };

  async function startTask() {  
    if ((cardProfile === "" || addressProfile === "") && (mode === "info")) {
      toast.error("Please select a card and address profile");
      return;
    }

    if (account.cookies === "") {
      toast.error("Account has no stored cookies. Likely created from a previous update and this feature is not supported.");
      return;
    }

    if (mode === "prime") {
      CreatePrimeTask(account.cookies, account.proxy, account.id);
    } else {
      CreateInfoTask(addressProfile, cardProfile, account.cookies, account.proxy, account.id);
    }
  }

  return (
    <Tooltip content="Start Subtask" color="secondary">
      <div>
        <Dropdown>
          <DropdownTrigger>
            <button>
              <VerticalDotsIcon />
            </button>
          </DropdownTrigger>
          <DropdownMenu>
            <DropdownItem onClick={() => {setMode("info"); onOpen()}}>Add Info</DropdownItem>
            <DropdownItem onClick={() => {setMode("prime"); onOpen()}}>Add Prime</DropdownItem>
          </DropdownMenu>
        </Dropdown>

        <Modal
          isOpen={isOpen}
          onOpenChange={onOpenChange}
          placement="top-center"
          isDismissable={false}
        >
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  {mode === "prime" ? "Prime Task" : "Info Task"}
                </ModalHeader>
                <ModalBody>
                  {mode === "info" ? (
                    <>
                      {/* grid of 2 */}
                      <div className="grid grid-cols-2 gap-4">
                        <Select
                          selectedKeys={[addressProfile]}
                          onChange={handleAddressProfile}
                          variant="bordered"
                          label="Select Address"
                          className="max-w-[180px]"
                          classNames={{ trigger: "data-[open=true]:border-default-400 data-[focus=true]:border-default-400"}}
                        >
                          {addresses.map((p: any) => (
                            <SelectItem key={p.id} value={p.id}>
                              {p.profileName}
                            </SelectItem>
                          ))}
                        </Select>

                        <Select
                          selectedKeys={[cardProfile]}
                          onChange={handleCardProfile}
                          variant="bordered"
                          label="Select Card"
                          className="max-w-[180px]"
                          classNames={{ trigger: "data-[open=true]:border-default-400 data-[focus=true]:border-default-400"}}
                        >
                          {cards.map((p: any) => (
                            <SelectItem key={p.id} value={p.id}>
                              {p.profileName}
                            </SelectItem>
                          ))}
                        </Select>
                      </div>
                    </>
                  ) : (
                    <> 
                    </> 
                  )}
                </ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="flat" onClick={() => {onClose();}}>
                    Close
                  </Button>
                  <Button color="primary" onPress={() => {startTask(); onClose();}}>
                    Start Task
                  </Button>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>
      </div>
    </Tooltip>
  );
};
