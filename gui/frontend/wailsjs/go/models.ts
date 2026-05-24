export namespace journal {
	
	export class Entry {
	    // Go type: time
	    Date: any;
	    Path: string;
	
	    static createFrom(source: any = {}) {
	        return new Entry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Date = this.convertValues(source["Date"], null);
	        this.Path = source["Path"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

export namespace service {
	
	export class Stats {
	    TotalEntries: number;
	    TotalWords: number;
	    AverageWords: number;
	    CurrentStreak: number;
	    LongestStreak: number;
	    // Go type: time
	    FirstEntry?: any;
	    // Go type: time
	    LastEntry?: any;
	
	    static createFrom(source: any = {}) {
	        return new Stats(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.TotalEntries = source["TotalEntries"];
	        this.TotalWords = source["TotalWords"];
	        this.AverageWords = source["AverageWords"];
	        this.CurrentStreak = source["CurrentStreak"];
	        this.LongestStreak = source["LongestStreak"];
	        this.FirstEntry = this.convertValues(source["FirstEntry"], null);
	        this.LastEntry = this.convertValues(source["LastEntry"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

