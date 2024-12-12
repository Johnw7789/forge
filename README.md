# Disclaimer
This project was published for educational purposes ONLY. I am not responsible for any individual who may modify it for the intent of malicious or harmful purposes. 

# About this project
Forge is an advanced account generator and toolkit that uses various bot evasion techniques to bypass restrictions on making accounts. Not a single automated browser is used, all automation is executed through http requests using a matching user-agent and tls fingerprint combination along with proper header order. 

Metadata1 (antibot by Amazon) payload generation has been excluded, but it is not fully needed in order to properly generate an account, however it does help in terms of success rate.

Additionally, all app authentication has been removed in addition to the login page.

# Technologies used
Forge was built with [Wails](https://wails.io/), using Go for the backend and React for the frontend. [NextUI](https://nextui.org/) (components) and [Recoiljs](https://recoiljs.org/) (state management) were also frequently used.

# Building
Fist you must install the Wails CLI:
``go install github.com/wailsapp/wails/v2/cmd/wails@latest``

Then to build:
``wails build``

# Accounts page
![alt text](https://github.com/Johnw7789/forge/blob/main/frontend/images/accounts.png)

# Tasks page
![alt text](https://github.com/Johnw7789/forge/blob/main/frontend/images/tasks.png)

# Settings page
![alt text](https://github.com/Johnw7789/forge/blob/main/frontend/images/settings.png)
