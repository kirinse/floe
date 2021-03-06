"use strict";

export function store(initial, restReq) {
    this.changed = true;
    this.data = initial;
    
    // Update updates the data at the given key and marks the Store as having a change.
    this.Update = function(key, val) {
        if (this.data == null) {
            this.data = {};
        }
        this.data[key] = val;
        this.changed = true;
    }

    // Get returns the data at the given key. If the data is unchanged then return null, 
    // unless force is true, then return the data in any case.
    this.Get = function(force) {
        if (!this.changed && !force) {
            return null;
        }
        this.changed = false;
        return this.data;
    }

    this.TrashAll = function(force) {
        this.changed = true;
        this.data = initial;
    }

    this.IsEmpty = function(){
        return this.data == null;
    }

    this.Reset = function() {
        this.changed = true;
        this.data = initial;
    }
}
