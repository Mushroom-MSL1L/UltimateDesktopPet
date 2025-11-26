export namespace activities {
	
	export class Activity {
	    id: number;
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
	        this.id = source["id"];
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

}

export namespace attributes {
	
	export class Attributes {
	    experience: number;
	    water: number;
	    hunger: number;
	    health: number;
	    mood: number;
	    energy: number;
	    money: number;
	
	    static createFrom(source: any = {}) {
	        return new Attributes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.experience = source["experience"];
	        this.water = source["water"];
	        this.hunger = source["hunger"];
	        this.health = source["health"];
	        this.mood = source["mood"];
	        this.energy = source["energy"];
	        this.money = source["money"];
	    }
	}

}

export namespace items {
	
	export class Item {
	    id: number;
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
	    description: string;
	
	    static createFrom(source: any = {}) {
	        return new Item(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
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
	        this.description = source["description"];
	    }
	}

}

export namespace pet {
	
	export class Pet {
	    id: number;
	    experience: number;
	    water: number;
	    hunger: number;
	    health: number;
	    mood: number;
	    energy: number;
	    money: number;
	
	    static createFrom(source: any = {}) {
	        return new Pet(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.experience = source["experience"];
	        this.water = source["water"];
	        this.hunger = source["hunger"];
	        this.health = source["health"];
	        this.mood = source["mood"];
	        this.energy = source["energy"];
	        this.money = source["money"];
	    }
	}

}

