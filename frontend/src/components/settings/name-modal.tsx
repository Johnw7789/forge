import React from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure, Checkbox, Input, Link} from "@nextui-org/react";
import { Icon } from "@iconify/react";

import { useRecoilState } from "recoil";
import { settingsState } from "../state/settings/atoms";

export default function NameModal() {
  const {isOpen, onOpen, onOpenChange} = useDisclosure();
  const [settings, setSettings] = useRecoilState(settingsState);

  const handleNameUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      [name]: value,
    }));
  }

  return (
    <>
      <Button
        onPress={onOpen}
        endContent={<Icon icon="solar:pen-2-linear" />}
        radius="full"
        variant="bordered"
      >
        Change
      </Button>      
      <Modal 
        isOpen={isOpen} 
        onOpenChange={onOpenChange}
        placement="top-center"
      >
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1">Global Name Override</ModalHeader>
              <ModalBody>
                <Input
                  value={settings.nameOverride}
                  onChange={handleNameUpdate("nameOverride")}
                  maxLength={70}
                  label="Full Name"
                  variant="faded"
                  spellCheck={false}
                />
              </ModalBody>
              <ModalFooter>
                <Button color="danger" variant="flat" onPress={onClose}>
                  Close
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </>
  );
}
