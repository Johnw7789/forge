import React from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure, Checkbox, Input, Link} from "@nextui-org/react";
import { Icon } from "@iconify/react";

import { useRecoilState } from "recoil";
import { settingsState } from "../state/settings/atoms";
import { ImapLogin } from "@/wailsjs/go/main/BackgroundController";

export default function ImapModal() {
  const {isOpen, onOpen, onOpenChange} = useDisclosure();
  const [settings, setSettings] = useRecoilState(settingsState);
  const [loading, setLoading] = React.useState(false);

  const handleImapCfgUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      imapConfig: {
        ...settings.imapConfig,
        [name]: value,
      }
    }));
  }

  const imapLogin = async () => {
    setLoading(true);
    let imap = settings.imapConfig;
    let result = await ImapLogin(imap.username, imap.password);
    setLoading(false);
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
              <ModalHeader className="flex flex-col gap-1">Imap Credentials</ModalHeader>
              <ModalBody>
                <Input
                  value={settings.imapConfig.username}
                  onChange={handleImapCfgUpdate("username")}
                  maxLength={40}
                  label="Username"
                  variant="faded"
                  spellCheck={false}
                />
                <Input
                  value={settings.imapConfig.password}
                  onChange={handleImapCfgUpdate("password")}
                  maxLength={40}
                  label="Password"
                  variant="faded"
                  spellCheck={false}
                />
              </ModalBody>
              <ModalFooter>
                <Button color="danger" variant="flat" onPress={onClose}>
                  Close
                </Button>
                <Button isLoading={loading ? true : false} color="primary" onPress={imapLogin}>
                  Login
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </>
  );
}
