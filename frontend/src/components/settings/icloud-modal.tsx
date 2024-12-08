import React from "react";
import {Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure, Checkbox, Input, Link} from "@nextui-org/react";
import { Icon } from "@iconify/react";

import { useRecoilState } from "recoil";
import { settingsState } from "../state/settings/atoms";
import { IcloudLogin, SubmitOTP, GenerateHMELoop } from "@/wailsjs/go/main/BackgroundController";

export default function IcloudModal() {
  const {isOpen, onOpen, onOpenChange} = useDisclosure();
  const [settings, setSettings] = useRecoilState(settingsState);
  const [loading, setLoading] = React.useState(false);
  const [loadingOtp, setLoadingOtp] = React.useState(false);
  const [otp, setOtp] = React.useState("");

  const handleIcloudCfgUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      icloudConfig: {
        ...settings.icloudConfig,
        [name]: value,
      }
    }));
  }

  const icloudGen = async () => {
    GenerateHMELoop();
  }

  const icloudLogin = async () => {
    setLoading(true);
    let ic = settings.icloudConfig;
    let result = await IcloudLogin(ic.username, ic.password);
    setOtp("")
    setLoading(false);
    setLoadingOtp(false);
  }

  const submitOtp = async () => {
    setLoadingOtp(true);
    SubmitOTP(otp);
  }

  return (
    <>
      <Button
        onPress={onOpen}
        endContent={<Icon icon="solar:pen-2-linear" />}
        radius="full"
        variant="bordered"
      >
        Open
      </Button>
      <Modal 
        isOpen={isOpen} 
        onOpenChange={onOpenChange}
        placement="top-center"
      >
        <ModalContent>
          {(onClose) => (
            <>
              <ModalHeader className="flex flex-col gap-1">iCloud Setup</ModalHeader>
              <ModalBody>
                <Input
                  value={settings.icloudConfig.username}
                  onChange={handleIcloudCfgUpdate("username")}
                  maxLength={40}
                  label="Username"
                  variant="faded"
                  spellCheck={false}
                />
                <Input
                  value={settings.icloudConfig.password}
                  onChange={handleIcloudCfgUpdate("password")}
                  maxLength={40}
                  label="Password"
                  variant="faded"
                  spellCheck={false}
                />

                <Input
                  value={otp}
                  disabled={!loading}
                  onChange={(e) => setOtp(e.target.value)}
                  maxLength={6}
                  label="OTP"
                  variant="faded"
                  spellCheck={false}
                />

              </ModalBody>
              <ModalFooter>
                <Button color="danger" variant="flat" onPress={onClose}>
                  Close
                </Button>
                <Button color="primary" onPress={icloudGen}>
                  Start Gen
                </Button>
                <Button isLoading={loading ? true : false} color="primary" onPress={icloudLogin}>
                  Login
                </Button>
                <Button isLoading={loadingOtp} disabled={!loading} color="primary" onPress={submitOtp}>
                  Submit OTP
                </Button>
              </ModalFooter>
            </>
          )}
        </ModalContent>
      </Modal>
    </>
  );
}
