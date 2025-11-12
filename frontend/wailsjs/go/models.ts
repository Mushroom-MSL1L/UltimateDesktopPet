export namespace window {
	
	export class PetHitbox {
	    left: number;
	    top: number;
	    width: number;
	    height: number;
	    devicePixelRatio: number;
	
	    static createFrom(source: any = {}) {
	        return new PetHitbox(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.left = source["left"];
	        this.top = source["top"];
	        this.width = source["width"];
	        this.height = source["height"];
	        this.devicePixelRatio = source["devicePixelRatio"];
	    }
	}

}

