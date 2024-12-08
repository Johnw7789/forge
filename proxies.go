package main

import (
	"errors"
	"os"
	"strings"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func isValidProxy(proxy string) bool {
	proxySpl := strings.Split(proxy, ":")
	if len(proxySpl) < 2 || len(proxySpl) > 4 {
		return false
	}

	return true
}

func (bc *BackgroundController) EditProxies(proxies string) {
	proxiesSpl := strings.Split(proxies, "\n")
	for _, proxy := range proxiesSpl {
		if !isValidProxy(proxy) {
			// emit err
			return
		}
	}

	bc.dataMu.Lock()
	defer bc.dataMu.Unlock()

	for _, proxy := range proxiesSpl {
		// check to see if exists in bc.proxies and if it does check if in use, if not in use then delete
		if _, ok := bc.proxies[proxy]; ok {
			if !bc.proxies[proxy] {
				delete(bc.proxies, proxy)
			}
		} else {
			bc.proxies[proxy] = false
		}
	}

	bc.SaveProxies(proxies)
}

func (bc *BackgroundController) LoadProxies() {
	bc.proxiesMu.Lock()

	proxies := ""

	proxiesFile, err := GetFilePath("config", "proxies.txt")
	if err != nil {
		bc.proxiesMu.Unlock()
		runtime.EventsEmit(bc.ctx, "error", "Error getting proxies file path")
	}

	_, err = os.Stat(proxiesFile)
	if os.IsNotExist(err) {
		bc.proxiesMu.Unlock()
		// runtime.EventsEmit(bc.ctx, "error", "Failed to load proxies")
		bc.SaveProxies(proxies)
		return
	} else if err != nil {
		bc.proxiesMu.Unlock()
		runtime.EventsEmit(bc.ctx, "error", "Failed to load proxies 2")
		return
	}

	file, err := os.ReadFile(proxiesFile)
	if err != nil {
		bc.proxiesMu.Unlock()
		runtime.EventsEmit(bc.ctx, "error", "Error reading proxies file")
		return
	}

	proxies = string(file)

	proxies = strings.ReplaceAll(proxies, "\r", "")
	proxies = strings.ReplaceAll(proxies, "\t", "")

	proxies = strings.TrimSuffix(proxies, "\n")

	bc.proxies = make(map[string]bool)
	proxiesSpl := strings.Split(proxies, "\n")
	for _, proxy := range proxiesSpl {
		if !isValidProxy(proxy) {
			proxiesSpl = append(proxiesSpl[:0], proxiesSpl[1:]...)
			// return errors.New("invalid proxy found in proxies.txt")
			// emit err
		}
	}

	for _, proxy := range proxiesSpl {
		bc.proxies[proxy] = false
	}

	runtime.EventsEmit(bc.ctx, "proxies", proxies)

	bc.proxiesMu.Unlock()
}

func (bc *BackgroundController) SaveProxies(proxies string) error {
	// bc.proxiesMu.Lock()
	// defer bc.proxiesMu.Unlock()

	// Clean up any /r or /t
	proxies = strings.ReplaceAll(proxies, "\r", "")
	proxies = strings.ReplaceAll(proxies, "\t", "")

	proxies = strings.TrimSuffix(proxies, "\n")

	proxiesFile, err := GetFilePath("config", "proxies.txt")
	if err != nil {
		runtime.EventsEmit(bc.ctx, "error", "Error getting proxies file path")
		return errors.New("Error getting proxies file path")
	}

	err = os.WriteFile(proxiesFile, []byte(proxies), 0644)
	if err != nil {
		runtime.EventsEmit(bc.ctx, "error", "Error writing proxies file")
		return errors.New("Error writing proxies file")
	}

	return nil
}
