import {
  Button,
  Input,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
} from "@nextui-org/react";
import React from "react";

import { toast } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

import { accountModalState } from "../state/accounts/atoms";
import { useRecoilState } from "recoil";

import { AddAccount } from "../../wailsjs/go/main/BackgroundController";

export const AddAccountModal = () => {
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  
  const [accountModal, setAccountModal] = useRecoilState(accountModalState)

  const handleUpdate = (name) => ({
		target: {
			value
		}
	}) => {
		setAccountModal({
			...accountModal,
			[name]: value
		})
	}

  const resetState = () => {
    let accCopy = {...accountModal};
    accCopy.id = "";
    accCopy.name = "";
    accCopy.email = "";
    accCopy.password = "";
    accCopy.phone = "";
    accCopy.proxy = "";
    accCopy.key2fa = "";

    setAccountModal(accCopy)
  }

  async function addAccount() {
    let account = accountModal;

    // id: "",
    // name: "",
    // line1: "",
    // line2: "",
    // city: "",
    // state: "",
    // zip: "",
    // phone: "",

    AddAccount(account as any);

    resetState();
  }

  return (
    <div>
      <>
        <Button onPress={onOpen} color="primary">
          Add Account
        </Button>
        <Modal
          isOpen={isOpen}
          onOpenChange={onOpenChange}
          placement="top-center"
        >
          <ModalContent>
            {(onClose) => (
              <>
                <ModalHeader className="flex flex-col gap-1">
                  Add Account
                </ModalHeader>
                <ModalBody>
                  <Input spellCheck={false} value={accountModal.name} onChange={handleUpdate("name")} variant="faded" label="Name"  />
                  <Input spellCheck={false} value={accountModal.email} onChange={handleUpdate("email")} variant="faded" label="Email" />
                  <Input spellCheck={false} value={accountModal.password} onChange={handleUpdate("password")} variant="faded" label="Password" />
                  <Input spellCheck={false} value={accountModal.phone} onChange={handleUpdate("phone")} variant="faded" label="Phone" />
                  <Input spellCheck={false} value={accountModal.proxy} onChange={handleUpdate("proxy")} variant="faded" label="Proxy" />
                  <Input spellCheck={false} value={accountModal.key2fa} onChange={handleUpdate("key2fa")} variant="faded" label="2FA Key" />
                </ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="flat" onClick={() => {onClose(); resetState();}}>
                    Close
                  </Button>
                  <Button color="primary" onPress={addAccount}>
                    Add Account
                  </Button>
                </ModalFooter>
              </>
            )}
          </ModalContent>
        </Modal>
      </>
    </div>
  );
};
