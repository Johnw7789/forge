import {
  Button,
  Input,
  Modal,
  ModalBody,
  ModalContent,
  ModalFooter,
  ModalHeader,
  useDisclosure,
  Select,
  SelectItem,
} from "@nextui-org/react";
import React from "react"

import { toast } from "react-toastify";
import 'react-toastify/dist/ReactToastify.css';

import { addressModalState } from "../state/addresses/atoms";
import { useRecoilState } from "recoil";

import { AddAddress } from "../../wailsjs/go/main/BackgroundController";
import { JigAddress } from "../../wailsjs/go/main/BackgroundController";


export const AddAddressModal = () => {
  const { isOpen, onOpen, onOpenChange } = useDisclosure();
  
  const [addressModal, setAddressModal] = useRecoilState(addressModalState)

  const states = ["AL", "AK", "AZ", "AR", "CA", "CO", "CT", "DE", "FL", "GA", "HI", "ID", "IL", "IN", "IA", "KS", "KY", "LA", "ME", "MD", "MA", "MI", "MN", "MS", "MO", "MT", "NE", "NV", "NH", "NJ", "NM", "NY", "NC", "ND", "OH", "OK", "OR", "PA", "RI", "SC", "SD", "TN", "TX", "UT", "VT", "VA", "WA", "WV", "WI", "WY"]

  const handleUpdate = (name) => ({
		target: {
			value
		}
	}) => {
		setAddressModal({
			...addressModal,
			[name]: value
		})
	}

  const resetState = () => {
    let addressCopy = {...addressModal};

    addressCopy.id = "";
    addressCopy.profileName = "";
    addressCopy.name = "";
    addressCopy.line1 = "";
    addressCopy.line2 = "";
    addressCopy.city = "";
    addressCopy.state = "";
    addressCopy.zip = "";
    addressCopy.phone = "";

    setAddressModal(addressCopy)
  }

  // deprecated
  // const isValidState = (state) => {
  //   return state === "Alabama" || state === "Alaska" || state === "Arizona" || state === "Arkansas" || state === "California" || state === "Colorado" || state === "Connecticut" || state === "Delaware" || state === "Florida" || state === "Georgia" || state === "Hawaii" || state === "Idaho" || state === "Illinois" || state === "Indiana" || state === "Iowa" || state === "Kansas" || state === "Kentucky" || state === "Louisiana" || state === "Maine" || state === "Maryland" || state === "Massachusetts" || state === "Michigan" || state === "Minnesota" || state === "Mississippi" || state === "Missouri" || state === "Montana" || state === "Nebraska" || state === "Nevada" || state === "New Hampshire" || state === "New Jersey" || state === "New Mexico" || state === "New York" || state === "North Carolina" || state === "North Dakota" || state === "Ohio" || state === "Oklahoma" || state === "Oregon" || state === "Pennsylvania" || state === "Rhode Island" || state === "South Carolina" || state === "South Dakota" || state === "Tennessee" || state === "Texas" || state === "Utah" || state === "Vermont" || state === "Virginia" || state === "Washington" || state === "West Virginia" || state === "Wisconsin" || state === "Wyoming";
  // }

  async function jigAddress() {
    let jigged = await JigAddress(addressModal.line1);

    setAddressModal({
      ...addressModal,
      line1: jigged
    })
  }

  async function addAddress() {
    let address = addressModal;

    if (address.profileName === "") {
      toast.error("Profile Name is required");
      return;
    }

    if (address.name === "") {
      toast.error("Name is required");
      return;
    }

    if (address.line1 === "") {
      toast.error("Line 1 is required");
      return;
    }

    if (address.city === "") {
      toast.error("City is required");
      return;
    }

    if (address.state === "") {
      toast.error("State is required");
      return;
    }

    if (address.zip === "") {
      toast.error("Zip is required");
      return;
    }

    if (address.phone === "") {
      toast.error("Phone is required");
      return;
    }

    if (address.state === "" || !states.includes(address.state)) {
      toast.error("Select a State");
      return;
    }

    AddAddress(address);

    resetState();
  }

  return (
    <div>
      <>
        <Button onPress={onOpen} color="primary">
          Add Address
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
                  Add Address
                </ModalHeader>
                <ModalBody>
                  <Input spellCheck={false} value={addressModal.profileName} onChange={handleUpdate("profileName")} variant="faded" label="Profile Name"  />
                  <Input spellCheck={false} value={addressModal.name} onChange={handleUpdate("name")} variant="faded" label="Full Name"  />
                  <Input spellCheck={false} value={addressModal.line1} onChange={handleUpdate("line1")} variant="faded" label="Address Line 1"  />
                  <Input spellCheck={false} value={addressModal.line2} onChange={handleUpdate("line2")} variant="faded" label="Address Line 2"  />
                  {/* grid 2 by 2 for last 4 inputs */}
                  <div className="grid grid-cols-2 gap-2">
                    <Input spellCheck={false} value={addressModal.city} onChange={handleUpdate("city")} variant="faded" label="City"  />
                    {/* <Input spellCheck={false} value={addressModal.state} onChange={handleUpdate("state")} variant="faded" label="State"  /> */}
                    {/* select state eg FL, NY, etc */}
                    <Select
                      selectedKeys={[addressModal.state]}
                      onChange={handleUpdate("state")}
                      variant="faded"
                      label="State"
                      // classNames={{ trigger: "data-[open=true]:border-default-400 data-[focus=true]:border-default-400"}}
                      >
                      {states.map((state: any) => (
                      <SelectItem key={state} value={state}>
                          {state}
                      </SelectItem>
                      ))}
                    </Select>

                    <Input spellCheck={false} value={addressModal.zip} onChange={handleUpdate("zip")} variant="faded" label="Zip"  />
                    <Input spellCheck={false} value={addressModal.phone} onChange={handleUpdate("phone")} variant="faded" label="Phone"  />
                  </div>
      
                </ModalBody>
                <ModalFooter>
                  <Button color="danger" variant="flat" onClick={() => {onClose(); resetState();}}>
                    Close
                  </Button>
                  <Button color="primary" onPress={async() => {jigAddress();}}>
                      Jig Address
                    </Button>
                  <Button color="primary" onPress={addAddress}>
                    Add Address
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
