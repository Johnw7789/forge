import {
    Button,
    Input,
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
  
  import { accountModalState } from "../state/accounts/atoms";
  import { useRecoilState } from "recoil";
  import { accountsState } from "../state/accounts/atoms";
  import { AddAddresses } from "@/wailsjs/go/main/BackgroundController";
  
  export const ImportAddressesModal = () => {
    const { isOpen, onOpen, onOpenChange } = useDisclosure();
    const [bot, setBot] = React.useState("");
    const [importState, setImportState] = React.useState("");

    const handleUpdateImport = (event) => {
        const { value } = event.target;
        setImportState(value);
      };

    const isValidState = (state) => {
        return state === "Alabama" || state === "Alaska" || state === "Arizona" || state === "Arkansas" || state === "California" || state === "Colorado" || state === "Connecticut" || state === "Delaware" || state === "Florida" || state === "Georgia" || state === "Hawaii" || state === "Idaho" || state === "Illinois" || state === "Indiana" || state === "Iowa" || state === "Kansas" || state === "Kentucky" || state === "Louisiana" || state === "Maine" || state === "Maryland" || state === "Massachusetts" || state === "Michigan" || state === "Minnesota" || state === "Mississippi" || state === "Missouri" || state === "Montana" || state === "Nebraska" || state === "Nevada" || state === "New Hampshire" || state === "New Jersey" || state === "New Mexico" || state === "New York" || state === "North Carolina" || state === "North Dakota" || state === "Ohio" || state === "Oklahoma" || state === "Oregon" || state === "Pennsylvania" || state === "Rhode Island" || state === "South Carolina" || state === "South Dakota" || state === "Tennessee" || state === "Texas" || state === "Utah" || state === "Vermont" || state === "Virginia" || state === "Washington" || state === "West Virginia" || state === "Wisconsin" || state === "Wyoming";
    }
  
    async function importAddresses() {
        let importString = importState;
        let addresses = [] as any; 

        // replace any \r or \t
        importString = importString.replace("\r", "");
        importString = importString.replace("\t", "");

        // split on \n
        let lines = importString.split("\n");
        // for each line, split by comma. format is profile_name,full_name,address_line_1,address_line_2,city,zip,state,phone
        lines.forEach((line, index) => {
            let address = line.split(",");
            if (address.length !== 8) {
                toast.error(`Error on line ${index + 1}: Address must have 8 fields!`)
                return;
            }

            let state = address[6];
            if (!isValidState(state)) {
                toast.error(`Error on line ${index + 1}: Invalid state!`)
                return;
            }

            let addressObj = {
                profileName: address[0],
                name: address[1],
                line1: address[2],
                line2: address[3],
                city: address[4],
                zip: address[5],
                state: address[6],
                phone: address[7]
            }
            addresses.push(addressObj)
        })

        AddAddresses(addresses).then(() => {
            toast.success("Successfully imported addresses!")
            setImportState("");
        })
    }
  
    return (
      <div>
        <>
          <Button onClick={onOpen} color="primary" startContent={<ExportIcon />}>
            Import
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
                    Import Addresses
                  </ModalHeader>
                  <ModalBody>
                <Textarea variant="faded" placeholder="Formatted addresses..." value={importState} onChange={handleUpdateImport}         
                        classNames={{
                        input: "resize-y min-h-[15rem] ",
                        }} 
                    />
                  </ModalBody>
                  <ModalFooter>
                    <Button color="danger" variant="flat" onClick={() => {onClose();}}>
                      Close
                    </Button>
                    <Button color="primary" onPress={importAddresses}>
                      Import
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
  