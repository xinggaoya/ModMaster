export namespace model {
	
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
	export class LocalGame {
	    name: string;
	    path: string;
	    img: string;
	
	    static createFrom(source: any = {}) {
	        return new LocalGame(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.path = source["path"];
	        this.img = source["img"];
	    }
	}

}

