import type {InputProps} from "@nextui-org/react";

import React from "react";
import {Button, Input, Checkbox, Link, Divider} from "@nextui-org/react";
import { authState } from "@/components/state/auth/atoms";
import { useRecoilState } from 'recoil'
// import { Authenticate } from "../../wailsjs/go/main/SettingsController";
import { EventsOn, EventsEmit } from "../../wailsjs/runtime/runtime"
import { WindowMinimise } from "../../wailsjs/runtime/runtime"
import { Quit } from "../../wailsjs/runtime/runtime"
// import "./login.css";

export const Login = () => {
  const [isVisible, setIsVisible] = React.useState(false);
  const [auth, setAuth] = useRecoilState(authState)

  const toggleVisibility = () => setIsVisible((prev) => !prev);

  const inputClasses: InputProps["classNames"] = {
    inputWrapper:
      "border-transparent bg-default-50/40 dark:bg-default-50/20 group-data-[focus=true]:border-foreground/20 data-[hover=true]:border-foreground/20",
  };

  const buttonClasses = "bg-foreground/10 dark:bg-foreground/20";

  const handleUpdate = (name: keyof any) => ({ target: { value } }: { target: { value: any } }) => {
    setAuth((auth) => ({
      ...auth,
      [name]: value,
    }));
  };

  const keyEvent = (key) => {
    setAuth({
      ...auth,
      licenseKey: key,
      loading: true
    })
    // Authenticate(key);
  }

  React.useEffect(() => {
    EventsEmit("frontend:auth")

    EventsOn("backend:key", keyEvent)

    EventsOn("backend:authloading", () =>{
      console.log("Auth Loading")
      setAuth({
        ...auth,
        loading: true
      })
      console.log(auth)
    })

    EventsOn("backend:authfail", () =>{
      console.log("Auth Fail")
      setAuth({
        ...auth,
        loading: false,
        authenticated: false
      })
    })

    EventsOn("backend:authsuccess", (userData) =>{
      console.log("Auth Success")
      setAuth({
        ...auth,
        loading: false,
        authenticated: true,
        discordImage: userData.discordImage,
        discordUser: userData.discordUser
      })
      console.log(auth)
    })
  }, []);

  
	async function submitLogin() {
    if (auth.licenseKey !== '') {
      // await Authenticate(auth.licenseKey);
    }
  }
	

  return (
    <div style={{"--wails-draggable":"drag"} as React.CSSProperties} className="bg-gradient-to-br from-red-700 to-orange-400">
     <div className="flex justify-end mr-4 mb-0" >
     <div className="mr-3" onClick={() => WindowMinimise()} id="minimize">–</div>
        <div onClick={() => Quit()} id="close" >✕</div>
    </div>
    <div className="flex h-screen w-screen items-center justify-center  p-2 sm:p-4 lg:p-8">
      <div className="flex w-full max-w-sm flex-col gap-4 rounded-large bg-background/60 px-8 pb-10 pt-6 shadow-small backdrop-blur-md backdrop-saturate-150 dark:bg-default-100/50">
        <p className="pb-2 text-xl font-medium">Log In</p>
        <form className="flex flex-col gap-3" onSubmit={(e) => e.preventDefault()}>
          <Input
            classNames={inputClasses}
            value={auth.licenseKey}
            onChange={handleUpdate("licenseKey")}
            label="License Key"
            name="license"
            placeholder="••••• ••••• •••••"
            type={isVisible ? "text" : "password"}
            variant="bordered"
          />
          <Button type="button" isLoading={auth.loading ? true : false} onClick={submitLogin} className={buttonClasses}>
            Log In
          </Button>
        </form>
      </div> 
    </div>
    </div>
  );
}