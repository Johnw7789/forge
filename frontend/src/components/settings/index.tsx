"use client";
import { Button, Select, SelectItem, Input } from "@nextui-org/react";
import React from "react";

import {Card, CardHeader, CardBody} from "@nextui-org/react";

import SwitchCell from "./switch-cell";
import CellWrapper from "./cell-wrapper";

import { settingsState } from "../state/settings/atoms";
import { useRecoilState } from 'recoil'
import { SaveSettings } from "@/wailsjs/go/main/SettingsController";

import ImapModal from "./imap-modal";
import IcloudModal from "./icloud-modal";
import SmsModal from "./sms-modal";
import NameModal from "./name-modal";

export const Settings = () => {
  const [settings, setSettings] = useRecoilState(settingsState)
  let smsProviders = ["SMS Man", "SMS Pool", "Daisy SMS"]

  const handleSwitch = (name: keyof any) => ({ target: { checked } }: { target: { checked: any } }) => {
    setSettings((settings) => ({
      ...settings,
      [name]: checked,
    }));
  }

  const handleSwitchSms = (name: keyof any) => ({ target: { checked } }: { target: { checked: any } }) => {
    setSettings((settings) => ({
      ...settings,
      smsConfig: {
        ...settings.smsConfig,
        [name]: checked,
      }
    }));
  }

  const handleSwitchImap = (name: keyof any) => ({ target: { checked }
  }: { target: { checked: any } }) => {
    setSettings((settings) => ({
      ...settings,
      imapConfig: {
        ...settings.imapConfig,
        [name]: checked,
      }
    }));
  }

  const handleUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      [name]: value,
    }));
  }

  const handleMaxTasksUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      maxTasks: Number(value),
    }));
  }

  const handleSmsMaxTriesUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      smsConfig: {
        ...settings.smsConfig,
        [name]: Number(value),
      }
    }));
  }

  const handleCaptchaMaxTriesUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      [name]: Number(value),
    }));
  }

  const handleWebhookUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      webhooks: {
        ...settings.webhooks,
        [name]: value,
      }
    }));
  }

  const handleImapCfgUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      imapConfig: {
        ...settings.imapConfig,
        [name]: value,
      }
    }));
  }

  const handleSmsCfgUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      smsConfig: {
        ...settings.smsConfig,
        [name]: value,
      }
    }));
  }

  const handleSmsCfgUpdateProvider = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setSettings((settings) => ({
      ...settings,
      smsConfig: {
        ...settings.smsConfig,
        username: "",
        [name]: value,
      }
    }));
  }


  async function saveSettings () {
    SaveSettings(settings as any, true)
  }

  return (
    <div className="my-14 lg:px-6 max-w-[95rem] mx-auto w-full flex flex-col gap-4">
      <h3 className="text-xl font-semibold">Settings</h3>

      <div>
        <Button onClick={saveSettings} size="md" color="primary" >
            Save All
        </Button>
      </div>

      <div className="mx-auto w-full flex flex-row gap-4">

      <Card className="w-full p-2" >
      <CardHeader className="flex flex-col items-start px-4 pb-0 pt-4">
        <p className="text-large">User Data</p>
        <p className="text-small text-default-500">Manage your user info</p>
      </CardHeader>
      <CardBody className="space-y-2">
        <CellWrapper>
          <div>
            <p>IMAP Credentials</p>
            <p className="text-small text-default-500">
              Save your IMAP credentials used to login to your iCloud account. Used for fetching OTP codes.
            </p>
          </div>
          <div className="flex w-full flex-wrap items-center justify-end gap-6 sm:w-auto sm:flex-nowrap">
            <ImapModal />
          </div>
        </CellWrapper>

        <CellWrapper>
          <div>
            <p>iCloud Gen</p>
            <p className="text-small text-default-500">
              Login to your iCloud account to generate emails for account creation. 5 emails will be generated each hour the gen is running.
            </p>
          </div>
          <div className="flex w-full flex-wrap items-center justify-end gap-6 sm:w-auto sm:flex-nowrap">
            <IcloudModal />
          </div>
        </CellWrapper>
{/* 
        <SwitchCell
          isSelected={settings.limitProxyUse}
          onChange={handleSwitch("limitProxyUse")}
          description="Prevents accounts from using the same proxy more than once if used during a successful creation before."
          label="Limit Proxy on Creation"
        /> */}

        <CellWrapper>
          <div>
            <p>SMS Provider</p>
            <p className="text-small text-default-500">
              Choose your SMS provider for retrieving SMS codes.
            </p>
          </div>

          <Select
            // value={settings.smsConfig.provider}
            selectedKeys={[settings.smsConfig.provider]}
            onChange={handleSmsCfgUpdateProvider("provider")}
            variant="bordered"
            label="Select Provider"
            className="max-w-[180px]"
            classNames={{ trigger: "data-[open=true]:border-default-400 data-[focus=true]:border-default-400"}}
          >
           {smsProviders.map((p) => (
              <SelectItem key={p} value={p}>
                {p}
              </SelectItem>
            ))}
          </Select>
        </CellWrapper>

        <CellWrapper>
          <div>
            <p>SMS Provider Credentials</p>
            <p className="text-small text-default-500">
              Save your credentials so that we can retrieve SMS codes for you.
            </p>
          </div>
          <div className="flex w-full flex-wrap items-center justify-end gap-6 sm:w-auto sm:flex-nowrap">
            <SmsModal />
          </div>
        </CellWrapper>

        <CellWrapper>
          <Input spellCheck={false} type="number" value={String(settings.smsConfig.maxTries)} onChange={handleSmsMaxTriesUpdate("maxTries")} classNames={{ inputWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", innerWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", input: "bg-default-200 group-data-[focus=true]:bg-default-200"}} label="Max SMS Attempts" labelPlacement="outside" placeholder="Enter max sms verification attempts" />
        </CellWrapper>

        <CellWrapper>
          <div>
            <p>Global Name Override</p>
            <p className="text-small text-default-500">
              Save your name so that we can use it for all tasks.
            </p>
          </div>
          <div className="flex w-full flex-wrap items-center justify-end gap-6 sm:w-auto sm:flex-nowrap">
            <NameModal />
          </div>
        </CellWrapper>
      </CardBody>
    </Card>

    {/* placeholder discord webhook: https://api.discord.com/webhooks/1234567890/ABCDEFGHIJKL */}
    <Card className="w-full p-2" >
      <CardHeader className="flex flex-col items-start px-4 pb-0 pt-4">
        <p className="text-large">Misc Settings</p>
        <p className="text-small text-default-500">Manage your general settings</p>
      </CardHeader>
      <CardBody className="space-y-2">
        <SwitchCell
          isSelected={settings.localHost}
          onChange={handleSwitch("localHost")}
          description="Skips proxies for any task started, uses your local IP."
          label="Use Local Host"
        />
        
      <SwitchCell
          isSelected={settings.limitProxyUse}
          onChange={handleSwitch("limitProxyUse")}
          description="Prevents accounts from using the same proxy more than once if used during a successful creation before."
          label="Limit Proxy on Creation"
        />

        <CellWrapper>
        <Input spellCheck={false} value={settings.captchaKey} onChange={handleUpdate("captchaKey")} classNames={{ inputWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", innerWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", input: "bg-default-200 group-data-[focus=true]:bg-default-200"}} label="Captcha Key" labelPlacement="outside" placeholder="Enter your Captcha API key" />
        </CellWrapper>

        <CellWrapper>
          <Input spellCheck={false} type="number" value={String(settings.captchaMaxTries)} onChange={handleCaptchaMaxTriesUpdate("captchaMaxTries")} classNames={{ inputWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", innerWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", input: "bg-default-200 group-data-[focus=true]:bg-default-200"}} label="Captcha Max Attempts" labelPlacement="outside" placeholder="Enter max captcha solve attempts" />
        </CellWrapper>

        <CellWrapper>
          <Input spellCheck={false} type="number" value={String(settings.maxTasks)} onChange={handleMaxTasksUpdate("maxTasks")} classNames={{ inputWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", innerWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", input: "bg-default-200 group-data-[focus=true]:bg-default-200"}} label="Max Tasks" labelPlacement="outside" placeholder="Enter max concurrent tasks" />
        </CellWrapper>

        <CellWrapper>
        <Input spellCheck={false} value={settings.webhooks.success} onChange={handleWebhookUpdate("success")} classNames={{ inputWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", innerWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", input: "bg-default-200 group-data-[focus=true]:bg-default-200"}} label="Success Webhook" labelPlacement="outside" placeholder="https://www.discord.com/api/webhooks/xxxxxxxx" />
        </CellWrapper>

        {/* <CellWrapper>
        <Input spellCheck={false} value={settings.webhooks.fail} onChange={handleWebhookUpdate("fail")} classNames={{ inputWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", innerWrapper: "bg-default-200 group-data-[focus=true]:bg-default-200", input: "bg-default-200 group-data-[focus=true]:bg-default-200"}} label="Fail Webhook" labelPlacement="outside" placeholder="https://www.discord.com/api/webhooks/xxxxxxxx" />
        </CellWrapper> */}
      </CardBody>
    </Card>
    </div>
    </div>
  );
};