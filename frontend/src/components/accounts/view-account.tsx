import {
  Button,
  Tooltip,
  Snippet,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@nextui-org/react";
import React from "react";

import 'react-toastify/dist/ReactToastify.css';

import { EyeIcon } from "../icons/table/eye-icon";
import { Account } from "./data";

interface ViewProps {
  account: Account,
}

export const ViewAccountModal = ({ account }: ViewProps) => {
  const { isOpen, onOpen, onOpenChange } = useDisclosure();

  const loadAccount = () => {
    if (isOpen) {
      return;
    }
    // setAccountModal(account);
    onOpen();
  }

  return (
    <Tooltip content="Edit account" color="secondary">
      <button onClick={loadAccount}>
        <EyeIcon size={20} fill="#979797" />

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
                  View Account
                </ModalHeader>
                <ModalBody>
                  <Snippet spellCheck={false} variant="bordered" symbol="Name:">{account.name}</Snippet>
                  <Snippet spellCheck={false} variant="bordered" symbol="Email:">{account.email}</Snippet>
                  <Snippet spellCheck={false} variant="bordered" symbol="Password:">{account.password}</Snippet>
                  <Snippet spellCheck={false} variant="bordered" symbol="Phone:">{account.phone}</Snippet>
                  <Snippet spellCheck={false} variant="bordered" symbol="Proxy:">{account.proxy}</Snippet>
                  <Snippet spellCheck={false} variant="bordered" symbol="2FA Key:">{account.key2fa}</Snippet>
                </ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="flat" onClick={() => {onClose();}}>
                    Close
                  </Button>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>
      </button>
    </Tooltip>
  );
};
