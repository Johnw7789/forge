import React from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure, Checkbox, Input, Link} from "@nextui-org/react";
import { Icon } from "@iconify/react";

import { useRecoilState } from "recoil";
import { settingsState } from "../state/settings/atoms";
import { TestSms } from "@/wailsjs/go/main/BackgroundController";

export default function SmsModal() {
  const {isOpen, onOpen, onOpenChange} = useDisclosure();
  const [settings, setSettings] = useRecoilState(settingsState);
  const [loading, setLoading] = React.useState(false);

  const handleSmsCfgUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      smsConfig: {
        ...settings.smsConfig,
        [name]: value,
      }
    }));
  }

  const testSms = async () => {
    setLoading(true);
    let sms = settings.smsConfig;
    let result = await TestSms(sms.provider, sms.username, sms.apiKey);
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
              <ModalHeader className="flex flex-col gap-1">SMS Credentials</ModalHeader>
              <ModalBody>
                {/* <Input
                  disabled={settings.smsConfig.provider === "SMS Pool" ? true : false}
                  value={settings.smsConfig.username}
                  onChange={handleSmsCfgUpdate("username")}
                  maxLength={40}
                  label="Username"
                  variant="faded"
                  spellCheck={false}
                /> */}
                <Input
                  value={settings.smsConfig.apiKey}
                  onChange={handleSmsCfgUpdate("apiKey")}
                  maxLength={70}
                  label="API Key"
                  variant="faded"
                  spellCheck={false}
                />
              </ModalBody>
              <ModalFooter>
                <Button color="danger" variant="flat" onPress={onClose}>
                  Close
                </Button>
                <Button isLoading={loading ? true : false} color="primary" onPress={testSms}>
                  Test SMS API
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </>
  );
}
