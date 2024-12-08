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
  import { AddCards } from "../../wailsjs/go/main/BackgroundController";
  
  export const ImportCardsModal = () => {
    const { isOpen, onOpen, onOpenChange } = useDisclosure();
    const [bot, setBot] = React.useState("");
    const [importState, setImportState] = React.useState("");

    const handleUpdateImport = (event) => {
        const { value } = event.target;
        setImportState(value);
      };
  
    async function importCards() {
        let importString = importState;
        let cards = [] as any; 

        // replace any \r or \t
        importString = importString.replace("\r", "");
        importString = importString.replace("\t", "");

        // split on \n
        let lines = importString.split("\n");
        // for each line, split by comma. format is profile_name,card_holder,card_number,exp_month,exp_year,cvv
        lines.forEach((line, index) => {
            let card = line.split(",");
            if (card.length !== 6) {
                toast.error(`Error on line ${index + 1}: Card must have 6 fields!`)
                return;
            }

            let cardObj = {
                profileName: card[0],
                name: card[1],
                number: card[2],
                expMonth: card[3],
                expYear: card[4],
                cvv: card[5]
            }
            cards.push(cardObj)
        })

        AddCards(cards).then(() => {
            toast.success("Successfully imported cards!")
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
                    Import Cards
                  </ModalHeader>
                  <ModalBody>
                <Textarea variant="faded" placeholder="Formatted cards..." value={importState} onChange={handleUpdateImport}         
                        classNames={{
                        input: "resize-y min-h-[15rem] ",
                        }} 
                    />
                  </ModalBody>
                  <ModalFooter>
                    <Button color="danger" variant="flat" onClick={() => {onClose();}}>
                      Close
                    </Button>
                    <Button color="primary" onPress={importCards}>
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
  