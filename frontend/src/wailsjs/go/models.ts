export namespace main {
	
	export class Account {
	    id: string;
	    name: string;
	    email: string;
	    password: string;
	    phone: string;
	    proxy: string;
	    key2fa: string;
	    cookies: string;
	    prime: boolean;
	    status: string;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.email = source["email"];
	        this.password = source["password"];
	        this.phone = source["phone"];
	        this.proxy = source["proxy"];
	        this.key2fa = source["key2fa"];
	        this.cookies = source["cookies"];
	        this.prime = source["prime"];
	        this.status = source["status"];
	    }
	}
	export class Address {
	    id: string;
	    profileName: string;
	    name: string;
	    line1: string;
	    line2: string;
	    city: string;
	    state: string;
	    zip: string;
	    phone: string;
	
	    static createFrom(source: any = {}) {
	        return new Address(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profileName = source["profileName"];
	        this.name = source["name"];
	        this.line1 = source["line1"];
	        this.line2 = source["line2"];
	        this.city = source["city"];
	        this.state = source["state"];
	        this.zip = source["zip"];
	        this.phone = source["phone"];
	    }
	}
	export class Card {
	    id: string;
	    profileName: string;
	    name: string;
	    number: string;
	    expMonth: string;
	    expYear: string;
	    cvv: string;
	
	    static createFrom(source: any = {}) {
	        return new Card(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.profileName = source["profileName"];
	        this.name = source["name"];
	        this.number = source["number"];
	        this.expMonth = source["expMonth"];
	        this.expYear = source["expYear"];
	        this.cvv = source["cvv"];
	    }
	}
	export class IcloudConfig {
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new IcloudConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class ImapConfig {
	    uniqueTaskClient: boolean;
	    username: string;
	    password: string;
	
	    static createFrom(source: any = {}) {
	        return new ImapConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uniqueTaskClient = source["uniqueTaskClient"];
	        this.username = source["username"];
	        this.password = source["password"];
	    }
	}
	export class SmsConfig {
	    maxTries: number;
	    provider: string;
	    username: string;
	    apiKey: string;
	
	    static createFrom(source: any = {}) {
	        return new SmsConfig(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.maxTries = source["maxTries"];
	        this.provider = source["provider"];
	        this.username = source["username"];
	        this.apiKey = source["apiKey"];
	    }
	}
	export class Webhooks {
	    success: string;
	    fail: string;
	
	    static createFrom(source: any = {}) {
	        return new Webhooks(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.fail = source["fail"];
	    }
	}
	export class Settings {
	    licenseKey: string;
	    maxTasks: number;
	    limitProxyUse: boolean;
	    persistState: boolean;
	    nameOverride: string;
	    webhooks: Webhooks;
	    imapConfig: ImapConfig;
	    smsConfig: SmsConfig;
	    captchaKey: string;
	    captchaMaxTries: number;
	    icloudConfig: IcloudConfig;
	    appleCookies: string;
	    localHost: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Settings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.licenseKey = source["licenseKey"];
	        this.maxTasks = source["maxTasks"];
	        this.limitProxyUse = source["limitProxyUse"];
	        this.persistState = source["persistState"];
	        this.nameOverride = source["nameOverride"];
	        this.webhooks = this.convertValues(source["webhooks"], Webhooks);
	        this.imapConfig = this.convertValues(source["imapConfig"], ImapConfig);
	        this.smsConfig = this.convertValues(source["smsConfig"], SmsConfig);
	        this.captchaKey = source["captchaKey"];
	        this.captchaMaxTries = source["captchaMaxTries"];
	        this.icloudConfig = this.convertValues(source["icloudConfig"], IcloudConfig);
	        this.appleCookies = source["appleCookies"];
	        this.localHost = source["localHost"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	

}

