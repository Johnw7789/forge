import React from "react";
import { useLockedBody } from "../hooks/useBodyLock";
import { NavbarWrapper } from "../navbar/navbar";
import { SidebarWrapper } from "../sidebar/sidebar";
import { SidebarContext } from "./layout-context";
import { Login } from "../login";
import { authState } from "@/components/state/auth/atoms";
import { useRecoilState } from 'recoil'
import { ToastContainer, toast } from 'react-toastify'
import 'react-toastify/dist/ReactToastify.css';

import { emailsInputState } from "../state/emails/atoms";
import { proxiesInputState } from "../state/proxies/atoms";
import { settingsState } from "../state/settings/atoms";
import { accountsState } from "../state/accounts/atoms";
import { addressesState } from "../state/addresses/atoms";
import { cardsState } from "../state/cards/atoms";
import { tasksState } from "../state/tasks/atoms";

import { EventsOn } from "../../wailsjs/runtime/runtime"
import { EventsEmit } from "../../wailsjs/runtime/runtime"

interface Props {
  children: React.ReactNode;
}

export const Layout = ({ children }: Props) => {
  const [sidebarOpen, setSidebarOpen] = React.useState(false);
  const [_, setLocked] = useLockedBody(false);
  const handleToggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
    setLocked(!sidebarOpen);
  };

  const [auth, setAuth] = useRecoilState(authState)
  const [tasks, setTasks] = useRecoilState(tasksState);
  const [emails, setEmails] = useRecoilState(emailsInputState);
  const [proxies, setProxies] = useRecoilState(proxiesInputState);
  const [settings, setSettings] = useRecoilState(settingsState);
  const [accounts, setAccounts] = useRecoilState(accountsState)
  const [addresses, setAddresses] = useRecoilState(addressesState)
  const [cards, setCards] = useRecoilState(cardsState)

  const emailsEvent = (emailStr) => {
    setEmails({emailsString: emailStr});
  }

  const proxiesEvent = (proxyStr) => {
    setProxies({proxiesString: proxyStr});
  }

  const settingsEvent = (settingsObj) => {
    setSettings(settingsObj);
  }

  const accountsEvent = (accounts) => {
    setAccounts(accounts);
  }

  const addressesEvent = (addresses) => {
    setAddresses(addresses);
  }

  const cardsEvent = (cards) => {
    setCards(cards);
  }

  const tasksEvent = (tasks) => {
    setTasks(tasks);
  }

  const errorEvent = (error: string) => {
    toast.error(error, {
      position: "top-right",
      autoClose: 4000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      // theme: "colored",
      });
  }

  const successEvent = (msg: string) => {
    toast.success(msg, {
      position: "top-right",
      autoClose: 4000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      progress: undefined,
      });
  }

  React.useEffect(() => {
    EventsEmit("frontend:auth")
  }, []);

  React.useEffect(() => {
    EventsOn("tasks", tasksEvent)
    EventsOn("emails", emailsEvent)
    EventsOn("proxies", proxiesEvent)
    EventsOn("settings", settingsEvent)
    EventsOn("accounts", accountsEvent)
    EventsOn("addresses", addressesEvent)
    EventsOn("cards", cardsEvent)

    EventsOn("error", errorEvent)
    EventsOn("success", successEvent)
  }, []);

  return (
    <div>
        {/* {auth.authenticated ? (  */}
          <SidebarContext.Provider
          value={{
            collapsed: sidebarOpen,
            setCollapsed: handleToggleSidebar,
          }}
        >
          <section className="flex">
            <SidebarWrapper  />
            <NavbarWrapper >{children}</NavbarWrapper>
            <ToastContainer
              position="top-right"
              autoClose={5000}
              hideProgressBar={false}
              newestOnTop={false}
              closeOnClick
              rtl={false}
              pauseOnFocusLoss
              draggable
              pauseOnHover
              toastStyle={{ backgroundColor: "bg-primary-200" }}
              theme="dark"
            />
          </section>
        </SidebarContext.Provider>
        {/* ) : ( 
          <section className="flex">
          <Login></Login>
        </section>
        ) */}
      {/* } */}
    </div>
  );
};
