export namespace activities {
	
	export class Activity {
	    path: string;
	    name: string;
	    type: string;
	    experience: number;
	    water: number;
	    hunger: number;
	    health: number;
	    mood: number;
	    energy: number;
	    money: number;
	    duration_minute: number;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Activity(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.experience = source["experience"];
	        this.water = source["water"];
	        this.hunger = source["hunger"];
	        this.health = source["health"];
	        this.mood = source["mood"];
	        this.energy = source["energy"];
	        this.money = source["money"];
	        this.duration_minute = source["duration_minute"];
	        this.description = source["description"];
	    }
	}
	export class ActivityWithFrames {
	    path: string;
	    name: string;
	    type: string;
	    experience: number;
	    water: number;
	    hunger: number;
	    health: number;
	    mood: number;
	    energy: number;
	    money: number;
	    duration_minute: number;
	    description: string;
	    Frames: string[];
	
	    static createFrom(source: any = {}) {
	        return new ActivityWithFrames(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.experience = source["experience"];
	        this.water = source["water"];
	        this.hunger = source["hunger"];
	        this.health = source["health"];
	        this.mood = source["mood"];
	        this.energy = source["energy"];
	        this.money = source["money"];
	        this.duration_minute = source["duration_minute"];
	        this.description = source["description"];
	        this.Frames = source["Frames"];
	    }
	}

}

export namespace items {
	
	export class Item {
	    path: string;
	    name: string;
	    type: string;
	    experience: number;
	    water: number;
	    hunger: number;
	    health: number;
	    mood: number;
	    energy: number;
	    moneyCost: number;
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.experience = source["experience"];
	        this.water = source["water"];
	        this.hunger = source["hunger"];
	        this.health = source["health"];
	        this.mood = source["mood"];
	        this.energy = source["energy"];
	        this.moneyCost = source["moneyCost"];
	        this.description = source["description"];
	    }
	}
	export class ItemWithFrame {
	    path: string;
	    name: string;
	    type: string;
	    experience: number;
	    water: number;
	    hunger: number;
	    health: number;
	    mood: number;
	    energy: number;
	    moneyCost: number;
	    description: string;
	    Frame: string;
	
	    static createFrom(source: any = {}) {
	        return new ItemWithFrame(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.experience = source["experience"];
	        this.water = source["water"];
	        this.hunger = source["hunger"];
	        this.health = source["health"];
	        this.mood = source["mood"];
	        this.energy = source["energy"];
	        this.moneyCost = source["moneyCost"];
	        this.description = source["description"];
	        this.Frame = source["Frame"];
	    }
	}

}

