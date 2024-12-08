import {
    Button,
    Modal,
    ModalBody,
    ModalContent,
    ModalFooter,
    ModalHeader,
    Select,
    SelectItem,
    Textarea,
    useDisclosure,
  } from "@nextui-org/react";
  import React from "react";
  import { ExportIcon } from "@/components/icons/accounts/export-icon";

  import { toast } from "react-toastify";
  import 'react-toastify/dist/ReactToastify.css';
  
  import { useRecoilState } from "recoil";
  import { selectedAccountsState } from "../state/accounts/atoms";
  
  export const ExportAccountsModal = () => {
    const { isOpen, onOpen, onOpenChange } = useDisclosure();
    const [bot, setBot] = React.useState("");
    const [accounts, setAccounts] = useRecoilState(selectedAccountsState);
    const [exportState, setExportState] = React.useState("");

    const handleUpdateExport = (event) => {
        const { value } = event.target;
        setExportState(value);
      };
      
      const handleBot = (event) => {
        const { value } = event.target;
        setBot(value);
      };
      

    // Frozen: user:::pass:::2fa:::proxy
    // Frenzy: user::pass::proxy::2fa
    // Refract: user:pass;proxy;2fa
    // Enven: user:::pass:::proxy:::2fa

  
    async function exportAccounts() {
        if (accounts.length === 0) {
            toast.error("No accounts selected");
            return;
        }

        let exportString = "";

        switch(bot) {
            case "Frozen":
                accounts.forEach((account: any) => {
                    let proxy = account.proxy;
                    if (proxy === "Localhost") {
                        proxy = "";
                    }
                    exportString += `${account.email}:::${account.password}:::${account.key2fa}:::${proxy}\n`;
                })
                break;
            case "Frenzy":
                accounts.forEach((account: any) => {
                    let proxy = account.proxy;
                    if (proxy === "Localhost") {
                        proxy = "";
                    }
                    exportString += `${account.email}::${account.password}::${proxy}::${account.key2fa}\n`;
                })
                break;
            case "Refract":
                accounts.forEach((account: any) => {
                    let proxy = account.proxy;
                    if (proxy === "Localhost") {
                        proxy = "";
                    }
                    exportString += `${account.email}:${account.password};${proxy};${account.key2fa}\n`;
                })
                break;
            case "Enven":
                accounts.forEach((account: any) => {
                    let proxy = account.proxy;
                    if (proxy === "Localhost") {
                        proxy = "";
                    }
                    exportString += `${account.email}:::${account.password}:::${proxy}:::${account.key2fa}\n`;
                })
                break;
            default:
                toast.error("Invalid bot selected");
                return;
        }

        setExportState(exportString);
    }
  
    return (
      <div>
        <>
          <Button onClick={onOpen} color="primary" startContent={<ExportIcon />}>
            Export to Bot
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
                    Export Accounts
                  </ModalHeader>
                  <ModalBody>
                  <Select
                    // value={settings.smsConfig.provider}
                    selectedKeys={[bot]}
                    onChange={handleBot}
                    variant="bordered"
                    label="Select Bot"
                    className="max-w-[180px]"
                    classNames={{ trigger: "data-[open=true]:border-default-400 data-[focus=true]:border-default-400"}}
                >
                {/* {bots.map((p) => (
                    <SelectItem key={p} value={p}>
                        {p}
                    </SelectItem>
                    ))} */}

                    <SelectItem key="Frozen" value="Frozen">
                        Frozen
                    </SelectItem>
                    <SelectItem key="Frenzy" value="Frenzy">
                        Frenzy
                    </SelectItem>
                    <SelectItem key="Refract" value="Refract">
                        Refract
                    </SelectItem>
                    <SelectItem key="Enven" value="Enven">
                        Enven
                    </SelectItem>
                </Select>
                <Textarea variant="faded" placeholder="Formatted accounts..." value={exportState} onChange={handleUpdateExport}         
                        classNames={{
                        input: "resize-y min-h-[15rem] ",
                        // inputWrapper: "bg-default-100 hover:bg-default-100 focus:bg-default-100", 
                        // innerWrapper: "bg-default-100 hover:bg-default-100 focus:bg-default-100"
                        }} 
                    />
                  </ModalBody>
                  <ModalFooter>
                    <Button color="danger" variant="flat" onClick={() => {onClose();}}>
                      Close
                    </Button>
                    <Button color="primary" onPress={exportAccounts}>
                      Format
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
  