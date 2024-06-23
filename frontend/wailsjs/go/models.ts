export namespace main {
	
	export class GameInfo {
	    name: string;
	    url: string;
	    img: string;
	
	    static createFrom(source: any = {}) {
	        return new GameInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.url = source["url"];
	        this.img = source["img"];
	    }
	}
	export class localGame {
	    name: string;
	    path: string;
	
	    static createFrom(source: any = {}) {
	        return new localGame(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	    }
	}

}

