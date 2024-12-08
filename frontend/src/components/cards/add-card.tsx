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

import { cardModalState } from "../state/cards/atoms";
import { useRecoilState } from "recoil";

import { AddCard } from "../../wailsjs/go/main/BackgroundController";

export const AddCardModal = () => {
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  
  const [cardModal, setCardModal] = useRecoilState(cardModalState)

  const handleUpdate = (name) => ({
		target: {
			value
		}
	}) => {
		setCardModal({
			...cardModal,
			[name]: value
		})
	}

  const resetState = () => {
    let cardCopy = {...cardModal};

    cardCopy.id = "";
    cardCopy.profileName = "";
    cardCopy.name = "";
    cardCopy.number = "";
    cardCopy.expMonth = "";
    cardCopy.expYear = "";
    cardCopy.cvv = "";

    setCardModal(cardCopy)
  }

  async function addCard() {
    let card = cardModal;

    if (card.profileName === "") {
      toast.error("Profile Name is required");
      return;
    }

    if (card.name === "") {
      toast.error("Card Holder is required");
      return;
    }

    if (card.number === "") {
      toast.error("Number is required");
      return;
    }

    if (card.expMonth === "") {
      toast.error("Exp Month is required");
      return;
    }

    if (card.expYear === "") {
      toast.error("Exp Year is required");
      return;
    }

    if (card.cvv === "") {
      toast.error("CVV is required");
      return;
    }

    AddCard(card);

    resetState();
  }

  return (
    <div>
      <>
        <Button onPress={onOpen} color="primary">
          Add Card
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
                  Add Card
                </ModalHeader>
                <ModalBody>
                  <Input
                    variant="faded"
                    label="Profile Name"
                    placeholder="Profile 1"
                    value={cardModal.profileName}
                    onChange={handleUpdate("profileName")}
                  />
                  <Input
                    variant="faded"
                    label="Card Holder"
                    placeholder="John Doe"
                    value={cardModal.name}
                    onChange={handleUpdate("name")}
                  />
                  <Input
                    variant="faded"
                    label="Number"
                    placeholder="4242424242424242"
                    value={cardModal.number}
                    onChange={handleUpdate("number")}
                  />
                  {/* row of 3 inputs  */}
                  <div className="grid grid-cols-3 gap-2">
                    <Input
                      variant="faded"
                      label="Exp Month"
                      placeholder="01"
                      value={cardModal.expMonth}
                      onChange={handleUpdate("expMonth")}
                    />
                    <Input
                      variant="faded"
                      label="Exp Year"
                      placeholder="25"
                      value={cardModal.expYear}
                      onChange={handleUpdate("expYear")}
                    />
                    <Input
                      variant="faded"
                      label="CVV"
                      placeholder="123"
                      value={cardModal.cvv}
                      onChange={handleUpdate("cvv")}
                    />
                  </div>
                </ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="flat" onClick={() => {onClose(); resetState();}}>
                    Close
                  </Button>
                  <Button color="primary" onPress={addCard}>
                    Add Card
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
