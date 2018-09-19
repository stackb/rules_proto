goog.module('example.routeguide.Client');

const Feature = goog.require('proto.routeguide.Feature');
const arrays = goog.require('goog.array');

/**
 * A skeleton client. The point of this exercise is not actually to create a
 * routeguide client, but show how we can use protobufs in closure code.
 */
class Client {

    constructor() {
        /**
         * A list of features
         * @const @private @type {!Array<!Feature>} 
         */
        this.features_ = [];
    }

    /**
     * Add a feature.
     * 
     * @param {!Feature} feature 
     */
    addFeature(feature) {
        this.features_.push(feature);
    }

    /**
     * Get a list of features.
     * 
     * @return {!Array<!Feature>} 
     */
    getFeatures() {
        return arrays.clone(this.features_);
    }

    /**
     * Run the client.  A real implementation might actually do something here.
     */
    run() {
        console.log("Running browser routeguide client...");
    }

}

exports = Client;