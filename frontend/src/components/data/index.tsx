"use client";
import { Button, Textarea } from "@nextui-org/react";
import React from "react";

import {Card, CardHeader, CardBody} from "@nextui-org/react";

import { useRecoilState } from "recoil";
import { emailsInputState } from "@/components/state/emails/atoms";
import { proxiesInputState } from "@/components/state/proxies/atoms";

import { SaveData } from "@/wailsjs/go/main/BackgroundController";

export const Data = () => {
  const [emails, setEmails] = useRecoilState(emailsInputState);
  const [proxies, setProxies] = useRecoilState(proxiesInputState);

  async function saveInfo() {
    SaveData(emails.emailsString, proxies.proxiesString)
  }

  const handleUpdateEmails = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setEmails((e) => ({
      ...e,
      [name]: value,
    }));
  };

  const handleUpdateProxies = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setProxies((p) => ({
      ...p,
      [name]: value,
    }));
  };

  return (
    <div className="my-14 lg:px-6 max-w-[95rem] mx-auto w-full flex flex-col gap-4">
      <h3 className="text-xl font-semibold">Data</h3>

      <div>
        <Button onClick={saveInfo} size="md" color="primary" >
            Save All
        </Button>
      </div>

      <div className="mx-auto w-full flex flex-row gap-4">

      <Card className="w-full p-2" >
      <CardHeader className="flex flex-col items-start px-4 pb-0 pt-4">
        <p className="text-large">Emails</p>
      </CardHeader>
      <CardBody >
        <Textarea variant="faded" placeholder="example@icloud.com" value={emails.emailsString} onChange={handleUpdateEmails("emailsString")} 
            classNames={{
              input: "resize-y min-h-[40rem] ",
            }} 
        />       
      </CardBody>
    </Card>

    <Card className="w-full p-2" >
      <CardHeader className="flex flex-col items-start px-4 pb-0 pt-4">
        <p className="text-large">Proxies</p>
      </CardHeader>
      <CardBody >
        <Textarea variant="faded" placeholder="ip:port:user:pass" value={proxies.proxiesString} onChange={handleUpdateProxies("proxiesString")}         
            classNames={{
              input: "resize-y min-h-[40rem] ",
            }} 
        />
      </CardBody>
    </Card>
    </div>
    </div>

  );
};